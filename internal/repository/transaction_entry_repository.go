package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type transactionEntryRepository struct {
	db *gorm.DB
}

func NewTransactionEntryRepository(db *gorm.DB) TransactionEntryRepository {
	return &transactionEntryRepository{db: db}
}

func (r *transactionEntryRepository) Create(entry *models.TransactionEntry) error {
	return r.db.Create(entry).Error
}

func (r *transactionEntryRepository) CreateBatch(entries []models.TransactionEntry) error {
	if len(entries) == 0 {
		return nil
	}
	return r.db.Create(&entries).Error
}

func (r *transactionEntryRepository) GetByTransactionID(transactionID uuid.UUID) ([]models.TransactionEntry, error) {
	var entries []models.TransactionEntry
	err := r.db.Preload("Account").Where("transaction_id = ?", transactionID).Find(&entries).Error
	return entries, err
}

func (r *transactionEntryRepository) GetByAccountID(accountID uuid.UUID) ([]models.TransactionEntry, error) {
	var entries []models.TransactionEntry
	err := r.db.Preload("Transaction").Where("account_id = ?", accountID).
		Order("created_at DESC").Find(&entries).Error
	return entries, err
}

func (r *transactionEntryRepository) GetByAccountIDAndDateRange(accountID uuid.UUID, startDate, endDate string) ([]models.TransactionEntry, error) {
	var entries []models.TransactionEntry
	err := r.db.Preload("Transaction").
		Joins("JOIN transactions ON transactions.id = transaction_entries.transaction_id").
		Where("transaction_entries.account_id = ? AND transactions.transaction_date BETWEEN ? AND ?", accountID, startDate, endDate).
		Order("transactions.transaction_date DESC").
		Find(&entries).Error
	return entries, err
}

func (r *transactionEntryRepository) GetByUserIDAndDateRange(userID uuid.UUID, startDate, endDate string) ([]models.TransactionEntry, error) {
	var entries []models.TransactionEntry
	err := r.db.Preload("Transaction").Preload("Account").
		Joins("JOIN transactions ON transactions.id = transaction_entries.transaction_id").
		Where("transactions.user_id = ? AND transactions.transaction_date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("transactions.transaction_date DESC").
		Find(&entries).Error
	return entries, err
}

func (r *transactionEntryRepository) DeleteByTransactionID(transactionID uuid.UUID) error {
	return r.db.Delete(&models.TransactionEntry{}, "transaction_id = ?", transactionID).Error
}

func (r *transactionEntryRepository) SumByAccountID(accountID uuid.UUID) (debit int64, credit int64, err error) {
	var result struct {
		TotalDebit  int64
		TotalCredit int64
	}
	err = r.db.Model(&models.TransactionEntry{}).
		Select("COALESCE(SUM(debit), 0) as total_debit, COALESCE(SUM(credit), 0) as total_credit").
		Where("account_id = ?", accountID).
		Scan(&result).Error
	return result.TotalDebit, result.TotalCredit, err
}

func (r *transactionEntryRepository) SumByAccountIDAndDateRange(accountID uuid.UUID, startDate, endDate string) (debit int64, credit int64, err error) {
	var result struct {
		TotalDebit  int64
		TotalCredit int64
	}
	err = r.db.Model(&models.TransactionEntry{}).
		Select("COALESCE(SUM(transaction_entries.debit), 0) as total_debit, COALESCE(SUM(transaction_entries.credit), 0) as total_credit").
		Joins("JOIN transactions ON transactions.id = transaction_entries.transaction_id").
		Where("transaction_entries.account_id = ? AND transactions.transaction_date BETWEEN ? AND ?", accountID, startDate, endDate).
		Scan(&result).Error
	return result.TotalDebit, result.TotalCredit, err
}
