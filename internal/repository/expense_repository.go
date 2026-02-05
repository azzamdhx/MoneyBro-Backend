package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type expenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

func (r *expenseRepository) GetByID(id uuid.UUID) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.Preload("Category").First(&expense, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepository) GetByUserID(userID uuid.UUID, filter *ExpenseFilter) ([]models.Expense, error) {
	var expenses []models.Expense
	query := r.db.Preload("Category").Where("user_id = ?", userID)

	if filter != nil {
		if filter.CategoryID != nil {
			query = query.Where("category_id = ?", *filter.CategoryID)
		}
		if filter.StartDate != nil {
			query = query.Where("expense_date >= ?", *filter.StartDate)
		}
		if filter.EndDate != nil {
			query = query.Where("expense_date <= ?", *filter.EndDate)
		}
	}

	err := query.Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepository) GetByUserIDAndDateRange(userID uuid.UUID, startDate, endDate string) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Preload("Category").
		Where("user_id = ? AND expense_date >= ? AND expense_date <= ?", userID, startDate, endDate).
		Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepository) GetRecentByUserID(userID uuid.UUID, limit int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Preload("Category").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepository) Update(expense *models.Expense) error {
	return r.db.Save(expense).Error
}

func (r *expenseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Expense{}, "id = ?", id).Error
}
