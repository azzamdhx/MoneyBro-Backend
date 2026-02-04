package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type LedgerService struct {
	db              *gorm.DB
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	entryRepo       repository.TransactionEntryRepository
}

func NewLedgerService(
	db *gorm.DB,
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	entryRepo repository.TransactionEntryRepository,
) *LedgerService {
	return &LedgerService{
		db:              db,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		entryRepo:       entryRepo,
	}
}

type LedgerEntry struct {
	AccountID uuid.UUID
	Debit     int64
	Credit    int64
}

func (s *LedgerService) CreateJournalEntry(
	userID uuid.UUID,
	date time.Time,
	description string,
	entries []LedgerEntry,
	referenceID *uuid.UUID,
	referenceType string,
) (*models.Transaction, error) {
	if err := s.validateEntries(entries); err != nil {
		return nil, err
	}

	var transaction *models.Transaction
	var refType *string
	if referenceType != "" {
		refType = &referenceType
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		transaction = &models.Transaction{
			UserID:          userID,
			TransactionDate: date,
			Description:     description,
			ReferenceID:     referenceID,
			ReferenceType:   refType,
		}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		for _, entry := range entries {
			txEntry := &models.TransactionEntry{
				TransactionID: transaction.ID,
				AccountID:     entry.AccountID,
				Debit:         entry.Debit,
				Credit:        entry.Credit,
			}
			if err := tx.Create(txEntry).Error; err != nil {
				return err
			}

			balanceChange := s.calculateBalanceChange(entry)
			if err := tx.Model(&models.Account{}).Where("id = ?", entry.AccountID).
				Update("current_balance", gorm.Expr("current_balance + ?", balanceChange)).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.transactionRepo.GetByID(transaction.ID)
}

func (s *LedgerService) UpdateJournalEntry(
	transactionID uuid.UUID,
	date time.Time,
	description string,
	entries []LedgerEntry,
) (*models.Transaction, error) {
	if err := s.validateEntries(entries); err != nil {
		return nil, err
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		oldEntries, err := s.entryRepo.GetByTransactionID(transactionID)
		if err != nil {
			return err
		}

		for _, entry := range oldEntries {
			oldLedgerEntry := LedgerEntry{
				AccountID: entry.AccountID,
				Debit:     entry.Debit,
				Credit:    entry.Credit,
			}
			reverseChange := -s.calculateBalanceChange(oldLedgerEntry)
			if err := tx.Model(&models.Account{}).Where("id = ?", entry.AccountID).
				Update("current_balance", gorm.Expr("current_balance + ?", reverseChange)).Error; err != nil {
				return err
			}
		}

		if err := tx.Delete(&models.TransactionEntry{}, "transaction_id = ?", transactionID).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Transaction{}).Where("id = ?", transactionID).Updates(map[string]interface{}{
			"transaction_date": date,
			"description":      description,
		}).Error; err != nil {
			return err
		}

		for _, entry := range entries {
			txEntry := &models.TransactionEntry{
				TransactionID: transactionID,
				AccountID:     entry.AccountID,
				Debit:         entry.Debit,
				Credit:        entry.Credit,
			}
			if err := tx.Create(txEntry).Error; err != nil {
				return err
			}

			balanceChange := s.calculateBalanceChange(entry)
			if err := tx.Model(&models.Account{}).Where("id = ?", entry.AccountID).
				Update("current_balance", gorm.Expr("current_balance + ?", balanceChange)).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.transactionRepo.GetByID(transactionID)
}

func (s *LedgerService) DeleteJournalEntry(transactionID uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		entries, err := s.entryRepo.GetByTransactionID(transactionID)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			ledgerEntry := LedgerEntry{
				AccountID: entry.AccountID,
				Debit:     entry.Debit,
				Credit:    entry.Credit,
			}
			reverseChange := -s.calculateBalanceChange(ledgerEntry)
			if err := tx.Model(&models.Account{}).Where("id = ?", entry.AccountID).
				Update("current_balance", gorm.Expr("current_balance + ?", reverseChange)).Error; err != nil {
				return err
			}
		}

		if err := tx.Delete(&models.TransactionEntry{}, "transaction_id = ?", transactionID).Error; err != nil {
			return err
		}

		return tx.Delete(&models.Transaction{}, "id = ?", transactionID).Error
	})
}

func (s *LedgerService) DeleteByReference(referenceID uuid.UUID, referenceType string) error {
	transaction, err := s.transactionRepo.GetByReference(referenceID, referenceType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return s.DeleteJournalEntry(transaction.ID)
}

func (s *LedgerService) GetTransaction(id uuid.UUID) (*models.Transaction, error) {
	return s.transactionRepo.GetByID(id)
}

func (s *LedgerService) GetTransactionByReference(referenceID uuid.UUID, referenceType string) (*models.Transaction, error) {
	return s.transactionRepo.GetByReference(referenceID, referenceType)
}

func (s *LedgerService) GetTransactionsByDateRange(userID uuid.UUID, startDate, endDate string) ([]models.Transaction, error) {
	return s.transactionRepo.GetByUserIDAndDateRange(userID, startDate, endDate)
}

func (s *LedgerService) validateEntries(entries []LedgerEntry) error {
	if len(entries) < 2 {
		return errors.New("journal entry must have at least 2 entries")
	}

	var totalDebit, totalCredit int64
	for _, entry := range entries {
		if entry.Debit > 0 && entry.Credit > 0 {
			return errors.New("entry cannot have both debit and credit")
		}
		if entry.Debit == 0 && entry.Credit == 0 {
			return errors.New("entry must have either debit or credit")
		}
		totalDebit += entry.Debit
		totalCredit += entry.Credit
	}

	if totalDebit != totalCredit {
		return errors.New("total debit must equal total credit")
	}

	return nil
}

func (s *LedgerService) calculateBalanceChange(entry LedgerEntry) int64 {
	account, err := s.accountRepo.GetByID(entry.AccountID)
	if err != nil {
		return 0
	}

	switch account.AccountType {
	case models.AccountTypeAsset, models.AccountTypeExpense:
		return entry.Debit - entry.Credit
	case models.AccountTypeLiability, models.AccountTypeIncome:
		return entry.Credit - entry.Debit
	}
	return 0
}

func (s *LedgerService) GetActualPaymentsByDateRange(userID uuid.UUID, startDate, endDate string, accountType models.AccountType) (int64, error) {
	accounts, err := s.accountRepo.GetByUserIDAndType(userID, accountType)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, account := range accounts {
		debit, credit, err := s.entryRepo.SumByAccountIDAndDateRange(account.ID, startDate, endDate)
		if err != nil {
			return 0, err
		}
		switch accountType {
		case models.AccountTypeLiability:
			total += debit
		case models.AccountTypeExpense:
			total += debit
		case models.AccountTypeIncome:
			total += credit
		case models.AccountTypeAsset:
			total += debit - credit
		}
	}
	return total, nil
}
