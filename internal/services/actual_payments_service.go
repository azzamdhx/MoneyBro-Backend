package services

import (
	"github.com/azzamdhx/moneybro/backend/internal/models"
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
		userID, startDate, endDate, "installment_payment",
	)
	if err == nil && len(installmentTxs) > 0 {
		// Collect all reference IDs for batch fetch
		refIDs := make([]uuid.UUID, 0, len(installmentTxs))
		for _, tx := range installmentTxs {
			if tx.ReferenceID != nil {
				refIDs = append(refIDs, *tx.ReferenceID)
			}
		}

		// Batch fetch all installment payments
		instPayments, _ := s.repos.InstallmentPayment.GetByIDs(refIDs)
		paymentMap := make(map[uuid.UUID]*models.InstallmentPayment, len(instPayments))
		for i := range instPayments {
			paymentMap[instPayments[i].ID] = &instPayments[i]
		}

		for _, tx := range installmentTxs {
			if tx.ReferenceID == nil {
				continue
			}
			instPmt, ok := paymentMap[*tx.ReferenceID]
			if !ok || instPmt.Installment == nil {
				continue
			}

			var amount int64
			for _, entry := range tx.Entries {
				if entry.Account.AccountType == "LIABILITY" {
					amount += entry.Debit
				}
			}

			report.Installments = append(report.Installments, ActualInstallmentPayment{
				InstallmentID:   instPmt.InstallmentID,
				Name:            instPmt.Installment.Name,
				Amount:          amount,
				TransactionDate: tx.TransactionDate.Format("2006-01-02"),
				Description:     tx.Description,
			})
			report.TotalInstallment += amount
		}
	}

	// Get debt payment transactions from transactions table only
	debtTxs, err := s.repos.Transaction.GetByUserIDAndDateRangeAndReferenceType(
		userID, startDate, endDate, "debt_payment",
	)
	if err == nil && len(debtTxs) > 0 {
		// Collect all reference IDs for batch fetch
		refIDs := make([]uuid.UUID, 0, len(debtTxs))
		for _, tx := range debtTxs {
			if tx.ReferenceID != nil {
				refIDs = append(refIDs, *tx.ReferenceID)
			}
		}

		// Batch fetch all debt payments
		debtPayments, _ := s.repos.DebtPayment.GetByIDs(refIDs)
		paymentMap := make(map[uuid.UUID]*models.DebtPayment, len(debtPayments))
		for i := range debtPayments {
			paymentMap[debtPayments[i].ID] = &debtPayments[i]
		}

		for _, tx := range debtTxs {
			if tx.ReferenceID == nil {
				continue
			}
			debtPmt, ok := paymentMap[*tx.ReferenceID]
			if !ok || debtPmt.Debt == nil {
				continue
			}

			var amount int64
			for _, entry := range tx.Entries {
				if entry.Account.AccountType == "LIABILITY" {
					amount += entry.Debit
				}
			}

			report.Debts = append(report.Debts, ActualDebtPayment{
				DebtID:          debtPmt.DebtID,
				PersonName:      debtPmt.Debt.PersonName,
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
