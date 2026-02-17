package services

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type MonthlySummaryService struct {
	repos            *repository.Repositories
	upcomingPayments *UpcomingPaymentsService
	actualPayments   *ActualPaymentsService
}

func NewMonthlySummaryService(repos *repository.Repositories, upcomingPayments *UpcomingPaymentsService, actualPayments *ActualPaymentsService) *MonthlySummaryService {
	return &MonthlySummaryService{
		repos:            repos,
		upcomingPayments: upcomingPayments,
		actualPayments:   actualPayments,
	}
}

type HistorySummaryResult struct {
	AvailableMonths          []string
	SelectedMonth            string
	IncomeSummary            IncomeBreakdown
	ExpenseSummary           ExpenseBreakdown
	Payments                 *ActualPaymentsReport
	TotalSavingsContribution int64
}

type ForecastSummaryResult struct {
	AvailableMonths          []string
	SelectedMonth            string
	IncomeSummary            IncomeBreakdown
	ExpenseSummary           ExpenseBreakdown
	Payments                 *UpcomingPaymentsReport
	TotalSavingsContribution int64
}

func monthKey(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), int(t.Month()))
}

func monthDateRange(month, year int) (string, string) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	return start.Format("2006-01-02"), end.Format("2006-01-02")
}

func (s *MonthlySummaryService) getAvailablePastMonths(userID uuid.UUID) ([]string, error) {
	current := monthKey(time.Now())
	query := `
		SELECT DISTINCT month_key FROM (
			SELECT TO_CHAR(income_date, 'YYYY-MM') AS month_key FROM incomes WHERE user_id = ? AND TO_CHAR(income_date, 'YYYY-MM') < ?
			UNION
			SELECT TO_CHAR(expense_date, 'YYYY-MM') AS month_key FROM expenses WHERE user_id = ? AND expense_date IS NOT NULL AND TO_CHAR(expense_date, 'YYYY-MM') < ?
			UNION
			SELECT TO_CHAR(dp.paid_at, 'YYYY-MM') AS month_key FROM debt_payments dp JOIN debts d ON dp.debt_id = d.id WHERE d.user_id = ? AND TO_CHAR(dp.paid_at, 'YYYY-MM') < ?
			UNION
			SELECT TO_CHAR(ip.paid_at, 'YYYY-MM') AS month_key FROM installment_payments ip JOIN installments i ON ip.installment_id = i.id WHERE i.user_id = ? AND TO_CHAR(ip.paid_at, 'YYYY-MM') < ?
		) sub
		ORDER BY month_key DESC
	`
	var months []string
	err := s.repos.DB.Raw(query, userID, current, userID, current, userID, current, userID, current).Scan(&months).Error
	if err != nil {
		return nil, err
	}
	return months, nil
}

func (s *MonthlySummaryService) getAvailableFutureMonths(userID uuid.UUID) ([]string, error) {
	current := monthKey(time.Now())
	monthSet := make(map[string]bool)

	// SQL DISTINCT for income/expense/debt future months
	query := `
		SELECT DISTINCT month_key FROM (
			SELECT TO_CHAR(income_date, 'YYYY-MM') AS month_key FROM incomes WHERE user_id = ? AND TO_CHAR(income_date, 'YYYY-MM') > ?
			UNION
			SELECT TO_CHAR(expense_date, 'YYYY-MM') AS month_key FROM expenses WHERE user_id = ? AND expense_date IS NOT NULL AND TO_CHAR(expense_date, 'YYYY-MM') > ?
			UNION
			SELECT TO_CHAR(due_date, 'YYYY-MM') AS month_key FROM debts WHERE user_id = ? AND due_date IS NOT NULL AND TO_CHAR(due_date, 'YYYY-MM') > ?
		) sub
	`
	var sqlMonths []string
	if err := s.repos.DB.Raw(query, userID, current, userID, current, userID, current).Scan(&sqlMonths).Error; err == nil {
		for _, m := range sqlMonths {
			monthSet[m] = true
		}
	}

	// Active installments: compute future payment months in Go (needs PaidCount from Payments)
	activeStatus := models.InstallmentStatusActive
	installments, err := s.repos.Installment.GetByUserID(userID, &activeStatus)
	if err == nil {
		for _, inst := range installments {
			remaining := inst.Tenor - inst.PaidCount()
			for i := 0; i < remaining; i++ {
				paymentDate := time.Date(inst.StartDate.Year(), inst.StartDate.Month(), 1, 0, 0, 0, 0, time.UTC)
				paymentDate = paymentDate.AddDate(0, inst.PaidCount()+i, 0)
				mk := monthKey(paymentDate)
				if mk > current {
					monthSet[mk] = true
				}
			}
		}
	}

	months := make([]string, 0, len(monthSet))
	for m := range monthSet {
		months = append(months, m)
	}
	sort.Strings(months)
	return months, nil
}

func (s *MonthlySummaryService) calculateIncomeSummary(userID uuid.UUID, startDate, endDate string) IncomeBreakdown {
	summary := IncomeBreakdown{}
	incomes, err := s.repos.Income.GetByUserIDAndDateRange(userID, startDate, endDate)
	if err != nil {
		return summary
	}

	summary.Count = len(incomes)
	categoryMap := make(map[uuid.UUID]*IncomeCategorySummary)
	typeMap := make(map[models.IncomeType]*IncomeTypeSummary)

	for _, inc := range incomes {
		summary.Total += inc.Amount

		if inc.Category != nil {
			if cs, exists := categoryMap[inc.CategoryID]; exists {
				cs.TotalAmount += inc.Amount
				cs.IncomeCount++
			} else {
				categoryMap[inc.CategoryID] = &IncomeCategorySummary{
					Category:    *inc.Category,
					TotalAmount: inc.Amount,
					IncomeCount: 1,
				}
			}
		}

		if ts, exists := typeMap[inc.IncomeType]; exists {
			ts.TotalAmount += inc.Amount
			ts.IncomeCount++
		} else {
			typeMap[inc.IncomeType] = &IncomeTypeSummary{
				IncomeType:  inc.IncomeType,
				TotalAmount: inc.Amount,
				IncomeCount: 1,
			}
		}
	}

	for _, cs := range categoryMap {
		summary.ByCategory = append(summary.ByCategory, *cs)
	}
	for _, ts := range typeMap {
		summary.ByType = append(summary.ByType, *ts)
	}

	return summary
}

func (s *MonthlySummaryService) calculateExpenseSummary(userID uuid.UUID, startDate, endDate string) ExpenseBreakdown {
	summary := ExpenseBreakdown{}
	expenses, err := s.repos.Expense.GetByUserIDAndDateRange(userID, startDate, endDate)
	if err != nil {
		return summary
	}

	summary.Count = len(expenses)
	categoryMap := make(map[uuid.UUID]*CategorySummary)

	for _, exp := range expenses {
		summary.Total += exp.Total()

		if exp.Category != nil {
			if cs, exists := categoryMap[exp.CategoryID]; exists {
				cs.TotalAmount += exp.Total()
				cs.ExpenseCount++
			} else {
				categoryMap[exp.CategoryID] = &CategorySummary{
					Category:     *exp.Category,
					TotalAmount:  exp.Total(),
					ExpenseCount: 1,
				}
			}
		}
	}

	for _, cs := range categoryMap {
		summary.ByCategory = append(summary.ByCategory, *cs)
	}

	return summary
}

func (s *MonthlySummaryService) GetHistorySummary(userID uuid.UUID, month, year *int) (*HistorySummaryResult, error) {
	availableMonths, err := s.getAvailablePastMonths(userID)
	if err != nil {
		return nil, err
	}

	result := &HistorySummaryResult{
		AvailableMonths: availableMonths,
	}

	var selectedMonth, selectedYear int
	if month != nil && year != nil {
		selectedMonth = *month
		selectedYear = *year
	} else if len(availableMonths) > 0 {
		fmt.Sscanf(availableMonths[0], "%d-%d", &selectedYear, &selectedMonth)
	} else {
		result.SelectedMonth = ""
		result.Payments = &ActualPaymentsReport{
			Installments: []ActualInstallmentPayment{},
			Debts:        []ActualDebtPayment{},
		}
		return result, nil
	}

	result.SelectedMonth = fmt.Sprintf("%04d-%02d", selectedYear, selectedMonth)
	startDate, endDate := monthDateRange(selectedMonth, selectedYear)

	var g errgroup.Group

	g.Go(func() error {
		result.IncomeSummary = s.calculateIncomeSummary(userID, startDate, endDate)
		return nil
	})

	g.Go(func() error {
		result.ExpenseSummary = s.calculateExpenseSummary(userID, startDate, endDate)
		return nil
	})

	g.Go(func() error {
		payments, err := s.actualPayments.GetActualPayments(userID, startDate, endDate)
		if err != nil {
			payments = &ActualPaymentsReport{
				Installments: []ActualInstallmentPayment{},
				Debts:        []ActualDebtPayment{},
			}
		}
		result.Payments = payments
		return nil
	})

	g.Go(func() error {
		totalSavings, err := s.repos.SavingsContribution.GetTotalByUserIDAndDateRange(userID, startDate, endDate)
		if err != nil {
			totalSavings = 0
		}
		result.TotalSavingsContribution = totalSavings
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MonthlySummaryService) GetForecastSummary(userID uuid.UUID, month, year *int) (*ForecastSummaryResult, error) {
	availableMonths, err := s.getAvailableFutureMonths(userID)
	if err != nil {
		return nil, err
	}

	result := &ForecastSummaryResult{
		AvailableMonths: availableMonths,
	}

	var selectedMonth, selectedYear int
	if month != nil && year != nil {
		selectedMonth = *month
		selectedYear = *year
	} else if len(availableMonths) > 0 {
		fmt.Sscanf(availableMonths[0], "%d-%d", &selectedYear, &selectedMonth)
	} else {
		result.SelectedMonth = ""
		result.Payments = &UpcomingPaymentsReport{
			Installments: []UpcomingInstallmentPayment{},
			Debts:        []UpcomingDebtPayment{},
		}
		return result, nil
	}

	result.SelectedMonth = fmt.Sprintf("%04d-%02d", selectedYear, selectedMonth)
	startDate, endDate := monthDateRange(selectedMonth, selectedYear)

	var g errgroup.Group

	g.Go(func() error {
		result.IncomeSummary = s.calculateIncomeSummary(userID, startDate, endDate)
		return nil
	})

	g.Go(func() error {
		result.ExpenseSummary = s.calculateExpenseSummary(userID, startDate, endDate)
		return nil
	})

	g.Go(func() error {
		payments, err := s.upcomingPayments.GetUpcomingPayments(userID, selectedMonth, selectedYear)
		if err != nil {
			payments = &UpcomingPaymentsReport{
				Installments: []UpcomingInstallmentPayment{},
				Debts:        []UpcomingDebtPayment{},
			}
		}
		result.Payments = payments
		return nil
	})

	g.Go(func() error {
		totalSavings, err := s.repos.SavingsContribution.GetTotalByUserIDAndDateRange(userID, startDate, endDate)
		if err != nil {
			totalSavings = 0
		}
		result.TotalSavingsContribution = totalSavings
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}
