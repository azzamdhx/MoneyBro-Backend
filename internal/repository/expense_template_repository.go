package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type expenseTemplateGroupRepository struct {
	db *gorm.DB
}

func NewExpenseTemplateGroupRepository(db *gorm.DB) ExpenseTemplateGroupRepository {
	return &expenseTemplateGroupRepository{db: db}
}

func (r *expenseTemplateGroupRepository) Create(group *models.ExpenseTemplateGroup) error {
	return r.db.Create(group).Error
}

func (r *expenseTemplateGroupRepository) GetByID(id uuid.UUID) (*models.ExpenseTemplateGroup, error) {
	var group models.ExpenseTemplateGroup
	err := r.db.Preload("Items").Preload("Items.Category").First(&group, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *expenseTemplateGroupRepository) GetByUserID(userID uuid.UUID) ([]models.ExpenseTemplateGroup, error) {
	var groups []models.ExpenseTemplateGroup
	err := r.db.Preload("Items").Preload("Items.Category").Where("user_id = ?", userID).Order("name ASC").Find(&groups).Error
	return groups, err
}

func (r *expenseTemplateGroupRepository) Update(group *models.ExpenseTemplateGroup) error {
	return r.db.Save(group).Error
}

func (r *expenseTemplateGroupRepository) Delete(id uuid.UUID) error {
	// Delete items first, then delete group
	if err := r.db.Delete(&models.ExpenseTemplateItem{}, "group_id = ?", id).Error; err != nil {
		return err
	}
	return r.db.Delete(&models.ExpenseTemplateGroup{}, "id = ?", id).Error
}

func (r *expenseTemplateGroupRepository) AddItem(item *models.ExpenseTemplateItem) error {
	return r.db.Create(item).Error
}

func (r *expenseTemplateGroupRepository) UpdateItem(item *models.ExpenseTemplateItem) error {
	return r.db.Save(item).Error
}

func (r *expenseTemplateGroupRepository) DeleteItem(itemID uuid.UUID) error {
	return r.db.Delete(&models.ExpenseTemplateItem{}, "id = ?", itemID).Error
}

func (r *expenseTemplateGroupRepository) GetItemByID(itemID uuid.UUID) (*models.ExpenseTemplateItem, error) {
	var item models.ExpenseTemplateItem
	err := r.db.Preload("Category").First(&item, "id = ?", itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
