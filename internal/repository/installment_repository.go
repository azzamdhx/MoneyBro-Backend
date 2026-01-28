package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type installmentRepository struct {
	db *gorm.DB
}

func NewInstallmentRepository(db *gorm.DB) InstallmentRepository {
	return &installmentRepository{db: db}
}

func (r *installmentRepository) Create(installment *models.Installment) error {
	return r.db.Create(installment).Error
}

func (r *installmentRepository) GetByID(id uuid.UUID) (*models.Installment, error) {
	var installment models.Installment
	err := r.db.Preload("Payments").First(&installment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &installment, nil
}

func (r *installmentRepository) GetByUserID(userID uuid.UUID, status *models.InstallmentStatus) ([]models.Installment, error) {
	var installments []models.Installment
	query := r.db.Preload("Payments").Where("user_id = ?", userID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("created_at DESC").Find(&installments).Error
	return installments, err
}

func (r *installmentRepository) GetByDueDay(dueDay int, status models.InstallmentStatus) ([]models.Installment, error) {
	var installments []models.Installment
	err := r.db.Preload("Payments").Preload("User").
		Where("due_day = ? AND status = ?", dueDay, status).
		Find(&installments).Error
	return installments, err
}

func (r *installmentRepository) Update(installment *models.Installment) error {
	return r.db.Save(installment).Error
}

func (r *installmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Installment{}, "id = ?", id).Error
}
