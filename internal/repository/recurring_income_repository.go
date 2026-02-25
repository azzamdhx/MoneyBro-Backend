package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type recurringIncomeGroupRepository struct {
	db *gorm.DB
}

func NewRecurringIncomeGroupRepository(db *gorm.DB) RecurringIncomeGroupRepository {
	return &recurringIncomeGroupRepository{db: db}
}

func (r *recurringIncomeGroupRepository) Create(group *models.RecurringIncomeGroup) error {
	return r.db.Create(group).Error
}

func (r *recurringIncomeGroupRepository) GetByID(id uuid.UUID) (*models.RecurringIncomeGroup, error) {
	var group models.RecurringIncomeGroup
	err := r.db.Preload("Items").Preload("Items.Category").First(&group, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *recurringIncomeGroupRepository) GetByUserID(userID uuid.UUID, isActive *bool) ([]models.RecurringIncomeGroup, error) {
	var groups []models.RecurringIncomeGroup
	query := r.db.Preload("Items").Preload("Items.Category").Where("user_id = ?", userID)

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("name ASC").Find(&groups).Error
	return groups, err
}

func (r *recurringIncomeGroupRepository) Update(group *models.RecurringIncomeGroup) error {
	return r.db.Save(group).Error
}

func (r *recurringIncomeGroupRepository) Delete(id uuid.UUID) error {
	// Delete items first, then delete group
	if err := r.db.Delete(&models.RecurringIncomeItem{}, "group_id = ?", id).Error; err != nil {
		return err
	}
	return r.db.Delete(&models.RecurringIncomeGroup{}, "id = ?", id).Error
}

func (r *recurringIncomeGroupRepository) AddItem(item *models.RecurringIncomeItem) error {
	return r.db.Create(item).Error
}

func (r *recurringIncomeGroupRepository) UpdateItem(item *models.RecurringIncomeItem) error {
	return r.db.Save(item).Error
}

func (r *recurringIncomeGroupRepository) DeleteItem(itemID uuid.UUID) error {
	return r.db.Delete(&models.RecurringIncomeItem{}, "id = ?", itemID).Error
}

func (r *recurringIncomeGroupRepository) GetItemByID(itemID uuid.UUID) (*models.RecurringIncomeItem, error) {
	var item models.RecurringIncomeItem
	err := r.db.Preload("Category").First(&item, "id = ?", itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
