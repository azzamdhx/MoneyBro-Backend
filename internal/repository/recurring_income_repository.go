package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type recurringIncomeRepository struct {
	db *gorm.DB
}

func NewRecurringIncomeRepository(db *gorm.DB) RecurringIncomeRepository {
	return &recurringIncomeRepository{db: db}
}

func (r *recurringIncomeRepository) Create(recurringIncome *models.RecurringIncome) error {
	return r.db.Create(recurringIncome).Error
}

func (r *recurringIncomeRepository) GetByID(id uuid.UUID) (*models.RecurringIncome, error) {
	var recurringIncome models.RecurringIncome
	err := r.db.Preload("Category").First(&recurringIncome, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &recurringIncome, nil
}

func (r *recurringIncomeRepository) GetByUserID(userID uuid.UUID, isActive *bool) ([]models.RecurringIncome, error) {
	var recurringIncomes []models.RecurringIncome
	query := r.db.Preload("Category").Where("user_id = ?", userID)

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("recurring_day ASC, source_name ASC").Find(&recurringIncomes).Error
	return recurringIncomes, err
}

func (r *recurringIncomeRepository) GetByRecurringDay(day int, isActive bool) ([]models.RecurringIncome, error) {
	var recurringIncomes []models.RecurringIncome
	err := r.db.Preload("Category").Preload("User").
		Where("recurring_day = ? AND is_active = ?", day, isActive).
		Find(&recurringIncomes).Error
	return recurringIncomes, err
}

func (r *recurringIncomeRepository) Update(recurringIncome *models.RecurringIncome) error {
	return r.db.Save(recurringIncome).Error
}

func (r *recurringIncomeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.RecurringIncome{}, "id = ?", id).Error
}
