package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type DebtService struct {
	debtRepo       repository.DebtRepository
	paymentRepo    repository.DebtPaymentRepository
	accountRepo    repository.AccountRepository
	accountService *AccountService
	ledgerService  *LedgerService
}

func NewDebtService(
	debtRepo repository.DebtRepository,
	paymentRepo repository.DebtPaymentRepository,
	accountRepo repository.AccountRepository,
	accountService *AccountService,
	ledgerService *LedgerService,
) *DebtService {
	return &DebtService{
		debtRepo:       debtRepo,
		paymentRepo:    paymentRepo,
		accountRepo:    accountRepo,
		accountService: accountService,
		ledgerService:  ledgerService,
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

	// Validate and calculate for INSTALLMENT payment type
	if input.PaymentType == models.DebtPaymentTypeInstallment {
		if input.Tenor == nil || *input.Tenor <= 0 {
			return nil, errors.New("tenor is required for installment debt")
		}

		// Auto-calculate monthly payment
		calculatedMonthly := input.ActualAmount / int64(*input.Tenor)
		if input.MonthlyPayment != nil && *input.MonthlyPayment != calculatedMonthly {
			return nil, errors.New("monthly payment doesn't match calculated value")
		}
		input.MonthlyPayment = &calculatedMonthly
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

	// Create linked LIABILITY account for this debt
	accountName := "Hutang: " + input.PersonName
	if _, err := s.accountService.CreateLinkedAccount(userID, accountName, models.AccountTypeLiability, debt.ID, "debt"); err != nil {
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

	// Validate and calculate for INSTALLMENT payment type
	if input.PaymentType == models.DebtPaymentTypeInstallment {
		if input.Tenor == nil || *input.Tenor <= 0 {
			return nil, errors.New("tenor is required for installment debt")
		}

		// Auto-calculate monthly payment
		calculatedMonthly := input.ActualAmount / int64(*input.Tenor)
		if input.MonthlyPayment != nil && *input.MonthlyPayment != calculatedMonthly {
			return nil, errors.New("monthly payment doesn't match calculated value")
		}
		input.MonthlyPayment = &calculatedMonthly
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
	// Get debt with payments to cleanup transactions
	debt, err := s.debtRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete all payment transactions first (before CASCADE deletes payments)
	for _, payment := range debt.Payments {
		_ = s.ledgerService.DeleteByReference(payment.ID, "debt_payment")
	}

	// Delete linked account
	if err := s.accountService.DeleteAccountByReference(id, "debt"); err != nil {
		return err
	}

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

	// Create ledger entry: DEBIT Liability Account, CREDIT Cash Account
	if err := s.createPaymentLedgerEntry(debt.UserID, debt, payment); err != nil {
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

func (s *DebtService) createPaymentLedgerEntry(userID uuid.UUID, debt *models.Debt, payment *models.DebtPayment) error {
	// Get liability account (linked to debt)
	liabilityAccount, err := s.accountRepo.GetByReference(debt.ID, "debt")
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

	_, err = s.ledgerService.CreateJournalEntry(
		userID,
		payment.PaidAt,
		"Debt Payment: "+debt.PersonName,
		entries,
		&payment.ID,
		"debt_payment",
	)
	return err
}
