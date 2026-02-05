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
	accountRepo     repository.AccountRepository
	accountService  *AccountService
	ledgerService   *LedgerService
}

func NewInstallmentService(
	installmentRepo repository.InstallmentRepository,
	paymentRepo repository.InstallmentPaymentRepository,
	accountRepo repository.AccountRepository,
	accountService *AccountService,
	ledgerService *LedgerService,
) *InstallmentService {
	return &InstallmentService{
		installmentRepo: installmentRepo,
		paymentRepo:     paymentRepo,
		accountRepo:     accountRepo,
		accountService:  accountService,
		ledgerService:   ledgerService,
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

	// Create linked LIABILITY account for this installment
	if _, err := s.accountService.CreateLinkedAccount(userID, input.Name, models.AccountTypeLiability, installment.ID, "installment"); err != nil {
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
	// Get installment with payments to cleanup transactions
	installment, err := s.installmentRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete all payment transactions first (before CASCADE deletes payments)
	for _, payment := range installment.Payments {
		_ = s.ledgerService.DeleteByReference(payment.ID, "installment_payment")
	}

	// Delete linked account
	if err := s.accountService.DeleteAccountByReference(id, "installment"); err != nil {
		return err
	}

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

	// Create ledger entry: DEBIT Liability Account, CREDIT Cash Account
	if err := s.createPaymentLedgerEntry(installment.UserID, installment, payment); err != nil {
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

func (s *InstallmentService) MarkComplete(id uuid.UUID) (*models.Installment, error) {
	installment, err := s.installmentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	installment.Status = models.InstallmentStatusCompleted
	if err := s.installmentRepo.Update(installment); err != nil {
		return nil, err
	}

	return s.installmentRepo.GetByID(id)
}

func (s *InstallmentService) createPaymentLedgerEntry(userID uuid.UUID, installment *models.Installment, payment *models.InstallmentPayment) error {
	// Get liability account (linked to installment)
	liabilityAccount, err := s.accountRepo.GetByReference(installment.ID, "installment")
	if err != nil {
		return err
	}

	// Get default cash account
	cashAccount, err := s.accountRepo.GetDefaultByUserID(userID)
	if err != nil {
		return err
	}

	entries := []LedgerEntry{
		{AccountID: liabilityAccount.ID, Debit: payment.Amount, Credit: 0},
		{AccountID: cashAccount.ID, Debit: 0, Credit: payment.Amount},
	}

	// Calculate period date based on installment start_date + (payment_number - 1) months
	periodDate := installment.StartDate.AddDate(0, payment.PaymentNumber-1, 0)

	_, err = s.ledgerService.CreateJournalEntry(
		userID,
		periodDate,
		"Installment Payment: "+installment.Name,
		entries,
		&payment.ID,
		"installment_payment",
	)
	return err
}
