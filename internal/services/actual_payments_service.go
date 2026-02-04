package services

import (
	"github.com/azzamdhx/moneybro/backend/internal/repository"
	"github.com/google/uuid"
)

type ActualPaymentsService struct {
	repos *repository.Repositories
}

func NewActualPaymentsService(repos *repository.Repositories) *ActualPaymentsService {
	return &ActualPaymentsService{
		repos: repos,
	}
}

type ActualInstallmentPayment struct {
	InstallmentID   uuid.UUID
	Name            string
	Amount          int64
	TransactionDate string
	Description     string
}

type ActualDebtPayment struct {
	DebtID          uuid.UUID
	PersonName      string
	Amount          int64
	TransactionDate string
	Description     string
}

type ActualPaymentsReport struct {
	Installments     []ActualInstallmentPayment
	Debts            []ActualDebtPayment
	TotalInstallment int64
	TotalDebt        int64
	TotalPayments    int64
}

func (s *ActualPaymentsService) GetActualPayments(userID uuid.UUID, startDate, endDate string) (*ActualPaymentsReport, error) {
	report := &ActualPaymentsReport{
		Installments: []ActualInstallmentPayment{},
		Debts:        []ActualDebtPayment{},
	}

	// Get installment payment transactions from transactions table only
	installmentTxs, err := s.repos.Transaction.GetByUserIDAndDateRangeAndReferenceType(
		userID,
		startDate,
		endDate,
		"installment_payment",
	)
	if err == nil {
		// Each transaction record represents an actual payment
		for _, tx := range installmentTxs {
			if tx.ReferenceID == nil {
				continue
			}

			// Get installment_payment record (reference_id points to installment_payments.id)
			installmentPayment, err := s.repos.InstallmentPayment.GetByID(*tx.ReferenceID)
			if err != nil || installmentPayment.Installment == nil {
				continue
			}

			// Calculate amount from transaction entries
			var amount int64
			for _, entry := range tx.Entries {
				if entry.Account.AccountType == "LIABILITY" {
					amount += entry.Debit
				}
			}

			// Add each transaction as a separate payment
			report.Installments = append(report.Installments, ActualInstallmentPayment{
				InstallmentID:   installmentPayment.InstallmentID,
				Name:            installmentPayment.Installment.Name,
				Amount:          amount,
				TransactionDate: tx.TransactionDate.Format("2006-01-02"),
				Description:     tx.Description,
			})

			report.TotalInstallment += amount
		}
	}

	// Get debt payment transactions from transactions table only
	debtTxs, err := s.repos.Transaction.GetByUserIDAndDateRangeAndReferenceType(
		userID,
		startDate,
		endDate,
		"debt_payment",
	)
	if err == nil {
		// Each transaction record represents an actual payment
		for _, tx := range debtTxs {
			if tx.ReferenceID == nil {
				continue
			}

			// Get debt_payment record (reference_id points to debt_payments.id)
			debtPayment, err := s.repos.DebtPayment.GetByID(*tx.ReferenceID)
			if err != nil || debtPayment.Debt == nil {
				continue
			}

			// Calculate amount from transaction entries
			var amount int64
			for _, entry := range tx.Entries {
				if entry.Account.AccountType == "LIABILITY" {
					amount += entry.Debit
				}
			}

			// Add each transaction as a separate payment
			report.Debts = append(report.Debts, ActualDebtPayment{
				DebtID:          debtPayment.DebtID,
				PersonName:      debtPayment.Debt.PersonName,
				Amount:          amount,
				TransactionDate: tx.TransactionDate.Format("2006-01-02"),
				Description:     tx.Description,
			})

			report.TotalDebt += amount
		}
	}

	report.TotalPayments = report.TotalInstallment + report.TotalDebt

	return report, nil
}
