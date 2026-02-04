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

	// Get installment payment transactions
	installmentTxs, err := s.repos.Transaction.GetByUserIDAndDateRangeAndReferenceType(
		userID,
		startDate,
		endDate,
		"installment_payment",
	)
	if err == nil {
		// Group by reference_id to get unique installments
		installmentMap := make(map[uuid.UUID]*ActualInstallmentPayment)

		for _, tx := range installmentTxs {
			if tx.ReferenceID == nil {
				continue
			}

			// Get installment details
			installment, err := s.repos.Installment.GetByID(*tx.ReferenceID)
			if err != nil {
				continue
			}

			// Calculate amount from transaction entries
			var amount int64
			for _, entry := range tx.Entries {
				if entry.Account.AccountType == "LIABILITY" {
					amount += entry.Debit
				}
			}

			// If this installment already exists, add to the amount
			if existing, exists := installmentMap[*tx.ReferenceID]; exists {
				existing.Amount += amount
			} else {
				installmentMap[*tx.ReferenceID] = &ActualInstallmentPayment{
					InstallmentID:   *tx.ReferenceID,
					Name:            installment.Name,
					Amount:          amount,
					TransactionDate: tx.TransactionDate.Format("2006-01-02"),
					Description:     tx.Description,
				}
			}

			report.TotalInstallment += amount
		}

		for _, payment := range installmentMap {
			report.Installments = append(report.Installments, *payment)
		}
	}

	// Get debt payment transactions
	debtTxs, err := s.repos.Transaction.GetByUserIDAndDateRangeAndReferenceType(
		userID,
		startDate,
		endDate,
		"debt_payment",
	)
	if err == nil {
		// Group by reference_id to get unique debts
		debtMap := make(map[uuid.UUID]*ActualDebtPayment)

		for _, tx := range debtTxs {
			if tx.ReferenceID == nil {
				continue
			}

			// Get debt details
			debt, err := s.repos.Debt.GetByID(*tx.ReferenceID)
			if err != nil {
				continue
			}

			// Calculate amount from transaction entries
			var amount int64
			for _, entry := range tx.Entries {
				if entry.Account.AccountType == "LIABILITY" {
					amount += entry.Debit
				}
			}

			// If this debt already exists, add to the amount
			if existing, exists := debtMap[*tx.ReferenceID]; exists {
				existing.Amount += amount
			} else {
				debtMap[*tx.ReferenceID] = &ActualDebtPayment{
					DebtID:          *tx.ReferenceID,
					PersonName:      debt.PersonName,
					Amount:          amount,
					TransactionDate: tx.TransactionDate.Format("2006-01-02"),
					Description:     tx.Description,
				}
			}

			report.TotalDebt += amount
		}

		for _, payment := range debtMap {
			report.Debts = append(report.Debts, *payment)
		}
	}

	report.TotalPayments = report.TotalInstallment + report.TotalDebt

	return report, nil
}
