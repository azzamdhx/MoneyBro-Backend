package services

import (
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type UpcomingPaymentsService struct {
	repos *repository.Repositories
}

func NewUpcomingPaymentsService(repos *repository.Repositories) *UpcomingPaymentsService {
	return &UpcomingPaymentsService{
		repos: repos,
	}
}

type UpcomingInstallmentPayment struct {
	InstallmentID     uuid.UUID
	Name              string
	MonthlyPayment    int64
	DueDay            int
	DueDate           time.Time
	RemainingAmount   int64
	RemainingPayments int
}

type UpcomingDebtPayment struct {
	DebtID          uuid.UUID
	PersonName      string
	MonthlyPayment  int64
	DueDate         time.Time
	RemainingAmount int64
	PaymentType     string
}

type UpcomingPaymentsReport struct {
	Installments     []UpcomingInstallmentPayment
	Debts            []UpcomingDebtPayment
	TotalInstallment int64
	TotalDebt        int64
	TotalPayments    int64
}

func (s *UpcomingPaymentsService) GetUpcomingPayments(userID uuid.UUID, month, year int) (*UpcomingPaymentsReport, error) {
	report := &UpcomingPaymentsReport{
		Installments: []UpcomingInstallmentPayment{},
		Debts:        []UpcomingDebtPayment{},
	}

	// Get all active installments
	activeStatus := models.InstallmentStatusActive
	installments, err := s.repos.Installment.GetByUserID(userID, &activeStatus)
	if err != nil {
		return nil, err
	}

	// Filter installments that have payments due in the specified month
	for _, inst := range installments {
		if isDueInMonth(inst.DueDay, month, year, inst.StartDate, inst.Tenor) {
			dueDate := calculateDueDate(inst.DueDay, month, year)

			payment := UpcomingInstallmentPayment{
				InstallmentID:     inst.ID,
				Name:              inst.Name,
				MonthlyPayment:    inst.MonthlyPayment,
				DueDay:            inst.DueDay,
				DueDate:           dueDate,
				RemainingAmount:   inst.RemainingAmount(),
				RemainingPayments: inst.RemainingPayments(),
			}

			report.Installments = append(report.Installments, payment)
			report.TotalInstallment += inst.MonthlyPayment
			report.TotalPayments += inst.MonthlyPayment
		}
	}

	// Get all active debts
	debtActiveStatus := models.DebtStatusActive
	debts, err := s.repos.Debt.GetByUserID(userID, &debtActiveStatus)
	if err != nil {
		return nil, err
	}

	// Filter debts that have payments due in the specified month
	for _, debt := range debts {
		if debt.DueDate != nil && isDebtDueInMonth(*debt.DueDate, month, year) {
			var monthlyPayment int64
			if debt.PaymentType == models.DebtPaymentTypeInstallment && debt.MonthlyPayment != nil {
				monthlyPayment = *debt.MonthlyPayment
			} else {
				monthlyPayment = debt.RemainingAmount()
			}

			payment := UpcomingDebtPayment{
				DebtID:          debt.ID,
				PersonName:      debt.PersonName,
				MonthlyPayment:  monthlyPayment,
				DueDate:         *debt.DueDate,
				RemainingAmount: debt.RemainingAmount(),
				PaymentType:     string(debt.PaymentType),
			}

			report.Debts = append(report.Debts, payment)
			report.TotalDebt += monthlyPayment
			report.TotalPayments += monthlyPayment
		}
	}

	return report, nil
}

// isDueInMonth checks if an installment is due in the specified month
func isDueInMonth(dueDay, month, year int, startDate time.Time, tenor int) bool {
	targetMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// Check if the installment has started by the target month
	if startDate.After(targetMonth.AddDate(0, 1, -1)) {
		return false
	}

	// If startDate is in the future relative to target month, not due yet
	if startDate.Year() > year || (startDate.Year() == year && int(startDate.Month()) > month) {
		return false
	}

	// Check if the installment has ended before the target month
	// End month is startDate + (tenor - 1) months (since first payment is in start month)
	endDate := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, tenor, 0)
	if targetMonth.After(endDate) || targetMonth.Equal(endDate) {
		return false
	}

	return true
}

// isDebtDueInMonth checks if a debt payment is due in the specified month
func isDebtDueInMonth(dueDate time.Time, month, year int) bool {
	return dueDate.Year() == year && int(dueDate.Month()) == month
}

// calculateDueDate creates a date from dueDay, month, and year
func calculateDueDate(dueDay, month, year int) time.Time {
	// Get the last day of the month
	lastDay := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()

	// If dueDay is greater than the last day of the month, use the last day
	if dueDay > lastDay {
		dueDay = lastDay
	}

	return time.Date(year, time.Month(month), dueDay, 0, 0, 0, 0, time.UTC)
}
