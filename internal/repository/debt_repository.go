package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type debtRepository struct {
	db *gorm.DB
}

func NewDebtRepository(db *gorm.DB) DebtRepository {
	return &debtRepository{db: db}
}

func (r *debtRepository) Create(debt *models.Debt) error {
	return r.db.Create(debt).Error
}

func (r *debtRepository) GetByID(id uuid.UUID) (*models.Debt, error) {
	var debt models.Debt
	err := r.db.Preload("Payments").First(&debt, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &debt, nil
}

func (r *debtRepository) GetByUserID(userID uuid.UUID, status *models.DebtStatus) ([]models.Debt, error) {
	var debts []models.Debt
	query := r.db.Preload("Payments").Where("user_id = ?", userID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("created_at DESC").Find(&debts).Error
	return debts, err
}

func (r *debtRepository) GetByDueDateRange(startDate, endDate string, status models.DebtStatus) ([]models.Debt, error) {
	var debts []models.Debt
	err := r.db.Preload("Payments").Preload("User").
		Where("due_date >= ? AND due_date <= ? AND status = ?", startDate, endDate, status).
		Find(&debts).Error
	return debts, err
}

func (r *debtRepository) Update(debt *models.Debt) error {
	return r.db.Save(debt).Error
}

func (r *debtRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Debt{}, "id = ?", id).Error
}
