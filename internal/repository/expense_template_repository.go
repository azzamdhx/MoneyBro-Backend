package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type expenseTemplateRepository struct {
	db *gorm.DB
}

func NewExpenseTemplateRepository(db *gorm.DB) ExpenseTemplateRepository {
	return &expenseTemplateRepository{db: db}
}

func (r *expenseTemplateRepository) Create(template *models.ExpenseTemplate) error {
	return r.db.Create(template).Error
}

func (r *expenseTemplateRepository) GetByID(id uuid.UUID) (*models.ExpenseTemplate, error) {
	var template models.ExpenseTemplate
	err := r.db.Preload("Category").First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *expenseTemplateRepository) GetByUserID(userID uuid.UUID) ([]models.ExpenseTemplate, error) {
	var templates []models.ExpenseTemplate
	err := r.db.Preload("Category").Where("user_id = ?", userID).Order("item_name ASC").Find(&templates).Error
	return templates, err
}

func (r *expenseTemplateRepository) Update(template *models.ExpenseTemplate) error {
	return r.db.Save(template).Error
}

func (r *expenseTemplateRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ExpenseTemplate{}, "id = ?", id).Error
}
