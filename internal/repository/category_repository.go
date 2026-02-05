package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetByUserID(userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("user_id = ?", userID).Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetByUserIDWithStats(userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Raw(`
		SELECT c.*, 
			   COUNT(e.id) as expense_count, 
			   COALESCE(SUM(e.unit_price * e.quantity), 0) as total_spent
		FROM categories c
		LEFT JOIN expenses e ON e.category_id = c.id
		WHERE c.user_id = ?
		GROUP BY c.id
		ORDER BY c.name ASC
	`, userID).Scan(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}
