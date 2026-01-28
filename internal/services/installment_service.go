package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type InstallmentService struct {
	installmentRepo repository.InstallmentRepository
	paymentRepo     repository.InstallmentPaymentRepository
}

func NewInstallmentService(installmentRepo repository.InstallmentRepository, paymentRepo repository.InstallmentPaymentRepository) *InstallmentService {
	return &InstallmentService{
		installmentRepo: installmentRepo,
		paymentRepo:     paymentRepo,
	}
}

type CreateInstallmentInput struct {
	Name           string
	ActualAmount   int64
	LoanAmount     int64
	MonthlyPayment int64
	Tenor          int
	StartDate      time.Time
	DueDay         int
	Notes          *string
}

func (s *InstallmentService) Create(userID uuid.UUID, input CreateInstallmentInput) (*models.Installment, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	if input.ActualAmount <= 0 {
		return nil, errors.New("actual amount must be positive")
	}
	if input.LoanAmount <= 0 {
		return nil, errors.New("loan amount must be positive")
	}
	if input.MonthlyPayment <= 0 {
		return nil, errors.New("monthly payment must be positive")
	}
	if input.Tenor <= 0 {
		return nil, errors.New("tenor must be positive")
	}
	if input.DueDay < 1 || input.DueDay > 31 {
		return nil, errors.New("due day must be between 1 and 31")
	}

	installment := &models.Installment{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           input.Name,
		ActualAmount:   input.ActualAmount,
		LoanAmount:     input.LoanAmount,
		MonthlyPayment: input.MonthlyPayment,
		Tenor:          input.Tenor,
		StartDate:      input.StartDate,
		DueDay:         input.DueDay,
		Status:         models.InstallmentStatusActive,
		Notes:          input.Notes,
	}

	if err := s.installmentRepo.Create(installment); err != nil {
		return nil, err
	}

	return s.installmentRepo.GetByID(installment.ID)
}

func (s *InstallmentService) GetByID(id uuid.UUID) (*models.Installment, error) {
	return s.installmentRepo.GetByID(id)
}

func (s *InstallmentService) GetByUserID(userID uuid.UUID, status *models.InstallmentStatus) ([]models.Installment, error) {
	return s.installmentRepo.GetByUserID(userID, status)
}

func (s *InstallmentService) Update(id uuid.UUID, input CreateInstallmentInput, status *models.InstallmentStatus) (*models.Installment, error) {
	installment, err := s.installmentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	installment.Name = input.Name
	installment.ActualAmount = input.ActualAmount
	installment.LoanAmount = input.LoanAmount
	installment.MonthlyPayment = input.MonthlyPayment
	installment.Tenor = input.Tenor
	installment.StartDate = input.StartDate
	installment.DueDay = input.DueDay
	installment.Notes = input.Notes

	if status != nil {
		installment.Status = *status
	}

	if err := s.installmentRepo.Update(installment); err != nil {
		return nil, err
	}

	return s.installmentRepo.GetByID(installment.ID)
}

func (s *InstallmentService) Delete(id uuid.UUID) error {
	return s.installmentRepo.Delete(id)
}

func (s *InstallmentService) RecordPayment(installmentID uuid.UUID, amount int64, paidAt time.Time) (*models.InstallmentPayment, error) {
	installment, err := s.installmentRepo.GetByID(installmentID)
	if err != nil {
		return nil, err
	}

	lastNumber, err := s.paymentRepo.GetLastPaymentNumber(installmentID)
	if err != nil {
		return nil, err
	}

	payment := &models.InstallmentPayment{
		ID:            uuid.New(),
		InstallmentID: installmentID,
		PaymentNumber: lastNumber + 1,
		Amount:        amount,
		PaidAt:        paidAt,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	if payment.PaymentNumber >= installment.Tenor {
		installment.Status = models.InstallmentStatusCompleted
		if err := s.installmentRepo.Update(installment); err != nil {
			return nil, err
		}
	}

	return payment, nil
}
