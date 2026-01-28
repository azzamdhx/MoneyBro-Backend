package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type incomeCategoryRepository struct {
	db *gorm.DB
}

func NewIncomeCategoryRepository(db *gorm.DB) IncomeCategoryRepository {
	return &incomeCategoryRepository{db: db}
}

func (r *incomeCategoryRepository) Create(category *models.IncomeCategory) error {
	return r.db.Create(category).Error
}

func (r *incomeCategoryRepository) GetByID(id uuid.UUID) (*models.IncomeCategory, error) {
	var category models.IncomeCategory
	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *incomeCategoryRepository) GetByUserID(userID uuid.UUID) ([]models.IncomeCategory, error) {
	var categories []models.IncomeCategory
	err := r.db.Where("user_id = ?", userID).Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *incomeCategoryRepository) Update(category *models.IncomeCategory) error {
	return r.db.Save(category).Error
}

func (r *incomeCategoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.IncomeCategory{}, "id = ?", id).Error
}
