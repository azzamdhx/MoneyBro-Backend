package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) GetByID(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("id = ?", id).First(&account).Error
	return &account, err
}

func (r *accountRepository) GetByUserID(userID uuid.UUID) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("user_id = ?", userID).Order("account_type, name").Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) GetByUserIDAndType(userID uuid.UUID, accountType models.AccountType) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("user_id = ? AND account_type = ?", userID, accountType).Order("name").Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) GetByUserIDAndTypeAndReferenceType(userID uuid.UUID, accountType models.AccountType, referenceType string) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("user_id = ? AND account_type = ? AND reference_type = ?", userID, accountType, referenceType).Order("name").Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) GetDefaultByUserID(userID uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&account).Error
	return &account, err
}

func (r *accountRepository) GetByReference(referenceID uuid.UUID, referenceType string) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("reference_id = ? AND reference_type = ?", referenceID, referenceType).First(&account).Error
	return &account, err
}

func (r *accountRepository) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

func (r *accountRepository) UpdateBalance(id uuid.UUID, balance int64) error {
	return r.db.Model(&models.Account{}).Where("id = ?", id).Update("current_balance", balance).Error
}

func (r *accountRepository) AddToBalance(id uuid.UUID, amount int64) error {
	return r.db.Model(&models.Account{}).Where("id = ?", id).
		Update("current_balance", gorm.Expr("current_balance + ?", amount)).Error
}

func (r *accountRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Account{}, "id = ?", id).Error
}

func (r *accountRepository) DeleteByReference(referenceID uuid.UUID, referenceType string) error {
	return r.db.Delete(&models.Account{}, "reference_id = ? AND reference_type = ?", referenceID, referenceType).Error
}
