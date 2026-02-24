package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type IncomeService struct {
	incomeRepo         repository.IncomeRepository
	incomeCategoryRepo repository.IncomeCategoryRepository
	accountRepo        repository.AccountRepository
	ledgerService      *LedgerService
}

func NewIncomeService(
	incomeRepo repository.IncomeRepository,
	incomeCategoryRepo repository.IncomeCategoryRepository,
	accountRepo repository.AccountRepository,
	ledgerService *LedgerService,
) *IncomeService {
	return &IncomeService{
		incomeRepo:         incomeRepo,
		incomeCategoryRepo: incomeCategoryRepo,
		accountRepo:        accountRepo,
		ledgerService:      ledgerService,
	}
}

type CreateIncomeInput struct {
	CategoryID  uuid.UUID
	SourceName  string
	Amount      int64
	IncomeDate  *time.Time
	IsRecurring bool
	Notes       *string
	PocketID    *uuid.UUID
}

type UpdateIncomeInput struct {
	CategoryID  *uuid.UUID
	SourceName  *string
	Amount      *int64
	IncomeDate  *time.Time
	IsRecurring *bool
	Notes       *string
	PocketID    *uuid.UUID
}

func (s *IncomeService) Create(userID uuid.UUID, input CreateIncomeInput) (*models.Income, error) {
	if input.SourceName == "" {
		return nil, errors.New("source name is required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	_, err := s.incomeCategoryRepo.GetByID(input.CategoryID)
	if err != nil {
		return nil, errors.New("invalid income category")
	}

	incomeDate := time.Now()
	if input.IncomeDate != nil {
		incomeDate = *input.IncomeDate
	}

	// Resolve pocket: use provided or fall back to default
	pocketID := input.PocketID
	if pocketID == nil {
		defaultAccount, defErr := s.accountRepo.GetDefaultByUserID(userID)
		if defErr != nil {
			return nil, errors.New("no default pocket found")
		}
		pocketID = &defaultAccount.ID
	}

	income := &models.Income{
		ID:          uuid.New(),
		UserID:      userID,
		CategoryID:  input.CategoryID,
		SourceName:  input.SourceName,
		Amount:      input.Amount,
		IncomeDate:  incomeDate,
		IsRecurring: input.IsRecurring,
		Notes:       input.Notes,
		PocketID:    pocketID,
	}

	if err := s.incomeRepo.Create(income); err != nil {
		return nil, err
	}

	// Create ledger entry: DEBIT Cash Account, CREDIT Income Account
	if err := s.createLedgerEntry(userID, income); err != nil {
		return nil, err
	}

	return s.incomeRepo.GetByID(income.ID)
}

func (s *IncomeService) GetByID(id uuid.UUID) (*models.Income, error) {
	return s.incomeRepo.GetByID(id)
}

func (s *IncomeService) GetByUserID(userID uuid.UUID, filter *repository.IncomeFilter) ([]models.Income, error) {
	return s.incomeRepo.GetByUserID(userID, filter)
}

func (s *IncomeService) Update(id uuid.UUID, input UpdateIncomeInput) (*models.Income, error) {
	income, err := s.incomeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.CategoryID != nil {
		_, err := s.incomeCategoryRepo.GetByID(*input.CategoryID)
		if err != nil {
			return nil, errors.New("invalid income category")
		}
		income.CategoryID = *input.CategoryID
	}

	if input.SourceName != nil {
		if *input.SourceName == "" {
			return nil, errors.New("source name is required")
		}
		income.SourceName = *input.SourceName
	}

	if input.Amount != nil {
		if *input.Amount <= 0 {
			return nil, errors.New("amount must be greater than 0")
		}
		income.Amount = *input.Amount
	}

	if input.IncomeDate != nil {
		income.IncomeDate = *input.IncomeDate
	}

	if input.IsRecurring != nil {
		income.IsRecurring = *input.IsRecurring
	}

	if input.Notes != nil {
		income.Notes = input.Notes
	}
	if input.PocketID != nil {
		income.PocketID = input.PocketID
	}

	if err := s.incomeRepo.Update(income); err != nil {
		return nil, err
	}

	// Update ledger entry
	if err := s.updateLedgerEntry(income); err != nil {
		return nil, err
	}

	return s.incomeRepo.GetByID(id)
}

func (s *IncomeService) Delete(id uuid.UUID) error {
	income, err := s.incomeRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete ledger entry first
	if err := s.ledgerService.DeleteByReference(income.ID, "income"); err != nil {
		return err
	}

	return s.incomeRepo.Delete(id)
}

func (s *IncomeService) createLedgerEntry(userID uuid.UUID, income *models.Income) error {
	// Get income account (linked to income category)
	incomeAccount, err := s.accountRepo.GetByReference(income.CategoryID, "income_category")
	if err != nil {
		return err
	}

	// Get pocket account
	var pocketAccount *models.Account
	if income.PocketID != nil {
		pocketAccount, err = s.accountRepo.GetByID(*income.PocketID)
	} else {
		pocketAccount, err = s.accountRepo.GetDefaultByUserID(userID)
	}
	if err != nil {
		return err
	}

	entries := []LedgerEntry{
		{AccountID: pocketAccount.ID, Debit: income.Amount, Credit: 0},
		{AccountID: incomeAccount.ID, Debit: 0, Credit: income.Amount},
	}

	_, err = s.ledgerService.CreateJournalEntry(
		userID,
		income.IncomeDate,
		"Income: "+income.SourceName,
		entries,
		&income.ID,
		"income",
	)
	return err
}

func (s *IncomeService) updateLedgerEntry(income *models.Income) error {
	tx, err := s.ledgerService.GetTransactionByReference(income.ID, "income")
	if err != nil {
		return err
	}

	incomeAccount, err := s.accountRepo.GetByReference(income.CategoryID, "income_category")
	if err != nil {
		return err
	}

	// Get pocket account
	var pocketAccount *models.Account
	if income.PocketID != nil {
		pocketAccount, err = s.accountRepo.GetByID(*income.PocketID)
	} else {
		pocketAccount, err = s.accountRepo.GetDefaultByUserID(income.UserID)
	}
	if err != nil {
		return err
	}

	entries := []LedgerEntry{
		{AccountID: pocketAccount.ID, Debit: income.Amount, Credit: 0},
		{AccountID: incomeAccount.ID, Debit: 0, Credit: income.Amount},
	}

	_, err = s.ledgerService.UpdateJournalEntry(tx.ID, income.IncomeDate, "Income: "+income.SourceName, entries)
	return err
}
