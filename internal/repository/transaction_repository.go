package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Entries").Preload("Entries.Account").Where("id = ?", id).First(&transaction).Error
	return &transaction, err
}

func (r *transactionRepository) GetByUserID(userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Entries").Preload("Entries.Account").
		Where("user_id = ?", userID).
		Order("transaction_date DESC, created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) GetByUserIDAndDateRange(userID uuid.UUID, startDate, endDate string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Entries").Preload("Entries.Account").
		Where("user_id = ? AND transaction_date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("transaction_date DESC, created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) GetByUserIDAndDateRangeAndReferenceType(userID uuid.UUID, startDate, endDate, referenceType string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Entries").Preload("Entries.Account").
		Where("user_id = ? AND transaction_date BETWEEN ? AND ? AND reference_type = ?", userID, startDate, endDate, referenceType).
		Order("transaction_date DESC, created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) GetByReference(referenceID uuid.UUID, referenceType string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Entries").Preload("Entries.Account").
		Where("reference_id = ? AND reference_type = ?", referenceID, referenceType).
		First(&transaction).Error
	return &transaction, err
}

func (r *transactionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Transaction{}, "id = ?", id).Error
}

func (r *transactionRepository) DeleteByReference(referenceID uuid.UUID, referenceType string) error {
	return r.db.Delete(&models.Transaction{}, "reference_id = ? AND reference_type = ?", referenceID, referenceType).Error
}
