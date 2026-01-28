package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type DebtService struct {
	debtRepo    repository.DebtRepository
	paymentRepo repository.DebtPaymentRepository
}

func NewDebtService(debtRepo repository.DebtRepository, paymentRepo repository.DebtPaymentRepository) *DebtService {
	return &DebtService{
		debtRepo:    debtRepo,
		paymentRepo: paymentRepo,
	}
}

type CreateDebtInput struct {
	PersonName     string
	ActualAmount   int64
	LoanAmount     *int64
	PaymentType    models.DebtPaymentType
	MonthlyPayment *int64
	Tenor          *int
	DueDate        *time.Time
	Notes          *string
}

func (s *DebtService) Create(userID uuid.UUID, input CreateDebtInput) (*models.Debt, error) {
	if input.PersonName == "" {
		return nil, errors.New("person name is required")
	}
	if input.ActualAmount <= 0 {
		return nil, errors.New("actual amount must be positive")
	}

	debt := &models.Debt{
		ID:             uuid.New(),
		UserID:         userID,
		PersonName:     input.PersonName,
		ActualAmount:   input.ActualAmount,
		LoanAmount:     input.LoanAmount,
		PaymentType:    input.PaymentType,
		MonthlyPayment: input.MonthlyPayment,
		Tenor:          input.Tenor,
		DueDate:        input.DueDate,
		Status:         models.DebtStatusActive,
		Notes:          input.Notes,
	}

	if err := s.debtRepo.Create(debt); err != nil {
		return nil, err
	}

	return s.debtRepo.GetByID(debt.ID)
}

func (s *DebtService) GetByID(id uuid.UUID) (*models.Debt, error) {
	return s.debtRepo.GetByID(id)
}

func (s *DebtService) GetByUserID(userID uuid.UUID, status *models.DebtStatus) ([]models.Debt, error) {
	return s.debtRepo.GetByUserID(userID, status)
}

func (s *DebtService) Update(id uuid.UUID, input CreateDebtInput, status *models.DebtStatus) (*models.Debt, error) {
	debt, err := s.debtRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	debt.PersonName = input.PersonName
	debt.ActualAmount = input.ActualAmount
	debt.LoanAmount = input.LoanAmount
	debt.PaymentType = input.PaymentType
	debt.MonthlyPayment = input.MonthlyPayment
	debt.Tenor = input.Tenor
	debt.DueDate = input.DueDate
	debt.Notes = input.Notes

	if status != nil {
		debt.Status = *status
	}

	if err := s.debtRepo.Update(debt); err != nil {
		return nil, err
	}

	return s.debtRepo.GetByID(debt.ID)
}

func (s *DebtService) Delete(id uuid.UUID) error {
	return s.debtRepo.Delete(id)
}

func (s *DebtService) RecordPayment(debtID uuid.UUID, amount int64, paidAt time.Time) (*models.DebtPayment, error) {
	debt, err := s.debtRepo.GetByID(debtID)
	if err != nil {
		return nil, err
	}

	lastNumber, err := s.paymentRepo.GetLastPaymentNumber(debtID)
	if err != nil {
		return nil, err
	}

	payment := &models.DebtPayment{
		ID:            uuid.New(),
		DebtID:        debtID,
		PaymentNumber: lastNumber + 1,
		Amount:        amount,
		PaidAt:        paidAt,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	debt.Payments = append(debt.Payments, *payment)
	if debt.RemainingAmount() <= 0 {
		debt.Status = models.DebtStatusCompleted
		if err := s.debtRepo.Update(debt); err != nil {
			return nil, err
		}
	}

	return payment, nil
}

func (s *DebtService) MarkComplete(id uuid.UUID) (*models.Debt, error) {
	debt, err := s.debtRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	debt.Status = models.DebtStatusCompleted
	if err := s.debtRepo.Update(debt); err != nil {
		return nil, err
	}

	return s.debtRepo.GetByID(id)
}
