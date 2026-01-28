package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type incomeRepository struct {
	db *gorm.DB
}

func NewIncomeRepository(db *gorm.DB) IncomeRepository {
	return &incomeRepository{db: db}
}

func (r *incomeRepository) Create(income *models.Income) error {
	return r.db.Create(income).Error
}

func (r *incomeRepository) GetByID(id uuid.UUID) (*models.Income, error) {
	var income models.Income
	err := r.db.Preload("Category").First(&income, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &income, nil
}

func (r *incomeRepository) GetByUserID(userID uuid.UUID, filter *IncomeFilter) ([]models.Income, error) {
	var incomes []models.Income
	query := r.db.Preload("Category").Where("user_id = ?", userID)

	if filter != nil {
		if filter.CategoryID != nil {
			query = query.Where("category_id = ?", *filter.CategoryID)
		}
		if filter.IncomeType != nil {
			query = query.Where("income_type = ?", *filter.IncomeType)
		}
		if filter.StartDate != nil {
			query = query.Where("income_date >= ?", *filter.StartDate)
		}
		if filter.EndDate != nil {
			query = query.Where("income_date <= ?", *filter.EndDate)
		}
	}

	err := query.Order("income_date DESC, created_at DESC").Find(&incomes).Error
	return incomes, err
}

func (r *incomeRepository) GetByUserIDAndDateRange(userID uuid.UUID, startDate, endDate string) ([]models.Income, error) {
	var incomes []models.Income
	err := r.db.Preload("Category").
		Where("user_id = ? AND income_date >= ? AND income_date <= ?", userID, startDate, endDate).
		Order("income_date DESC").
		Find(&incomes).Error
	return incomes, err
}

func (r *incomeRepository) Update(income *models.Income) error {
	return r.db.Save(income).Error
}

func (r *incomeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Income{}, "id = ?", id).Error
}
