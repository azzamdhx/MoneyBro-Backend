package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type installmentPaymentRepository struct {
	db *gorm.DB
}

func NewInstallmentPaymentRepository(db *gorm.DB) InstallmentPaymentRepository {
	return &installmentPaymentRepository{db: db}
}

func (r *installmentPaymentRepository) Create(payment *models.InstallmentPayment) error {
	return r.db.Create(payment).Error
}

func (r *installmentPaymentRepository) GetByInstallmentID(installmentID uuid.UUID) ([]models.InstallmentPayment, error) {
	var payments []models.InstallmentPayment
	err := r.db.Where("installment_id = ?", installmentID).Order("payment_number ASC").Find(&payments).Error
	return payments, err
}

func (r *installmentPaymentRepository) GetByID(id uuid.UUID) (*models.InstallmentPayment, error) {
	var payment models.InstallmentPayment
	err := r.db.Preload("Installment").Where("id = ?", id).First(&payment).Error
	return &payment, err
}

func (r *installmentPaymentRepository) GetLastPaymentNumber(installmentID uuid.UUID) (int, error) {
	var payment models.InstallmentPayment
	err := r.db.Where("installment_id = ?", installmentID).Order("payment_number DESC").First(&payment).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return payment.PaymentNumber, nil
}
