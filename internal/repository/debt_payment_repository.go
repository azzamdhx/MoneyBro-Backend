package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type debtPaymentRepository struct {
	db *gorm.DB
}

func NewDebtPaymentRepository(db *gorm.DB) DebtPaymentRepository {
	return &debtPaymentRepository{db: db}
}

func (r *debtPaymentRepository) Create(payment *models.DebtPayment) error {
	return r.db.Create(payment).Error
}

func (r *debtPaymentRepository) GetByDebtID(debtID uuid.UUID) ([]models.DebtPayment, error) {
	var payments []models.DebtPayment
	err := r.db.Where("debt_id = ?", debtID).Order("payment_number ASC").Find(&payments).Error
	return payments, err
}

func (r *debtPaymentRepository) GetLastPaymentNumber(debtID uuid.UUID) (int, error) {
	var payment models.DebtPayment
	err := r.db.Where("debt_id = ?", debtID).Order("payment_number DESC").First(&payment).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return payment.PaymentNumber, nil
}
