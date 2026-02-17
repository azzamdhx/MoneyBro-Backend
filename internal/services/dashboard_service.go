package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type DashboardService struct {
	repos         *repository.Repositories
	redis         *redis.Client
	ledgerService *LedgerService
}

func NewDashboardService(repos *repository.Repositories, redis *redis.Client, ledgerService *LedgerService) *DashboardService {
	return &DashboardService{
		repos:         repos,
		redis:         redis,
		ledgerService: ledgerService,
	}
}

type BalanceStatus string

const (
	BalanceStatusSurplus  BalanceStatus = "SURPLUS"
	BalanceStatusDeficit  BalanceStatus = "DEFICIT"
	BalanceStatusBalanced BalanceStatus = "BALANCED"
)

type BalanceSummary struct {
	TotalIncome             int64
	TotalExpense            int64
	TotalInstallmentPayment int64
	TotalDebtPayment        int64
	NetBalance              int64
	Status                  BalanceStatus
}

type Dashboard struct {
	TotalActiveDebt        int64
	TotalActiveInstallment int64
	TotalExpenseThisMonth  int64
	TotalIncomeThisMonth   int64
	TotalSavingsGoal       int64
	BalanceSummary         BalanceSummary
	UpcomingInstallments   []models.Installment
	UpcomingDebts          []models.Debt
	ActiveSavingsGoals     []models.SavingsGoal
	ExpensesByCategory     []CategorySummary
	RecentExpenses         []models.Expense
}

type CategorySummary struct {
	Category     models.Category
	TotalAmount  int64
	ExpenseCount int
}

func (s *DashboardService) GetDashboard(userID uuid.UUID) (*Dashboard, error) {
	dashboard := &Dashboard{}

	activeStatus := models.DebtStatusActive
	debts, err := s.repos.Debt.GetByUserID(userID, &activeStatus)
	if err != nil {
		return nil, err
	}
	for _, debt := range debts {
		dashboard.TotalActiveDebt += debt.RemainingAmount()
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	installmentActiveStatus := models.InstallmentStatusActive
	installments, err := s.repos.Installment.GetByUserID(userID, &installmentActiveStatus)
	if err != nil {
		return nil, err
	}
	for _, inst := range installments {
		dashboard.TotalActiveInstallment += inst.RemainingAmount()
	}

	// Get monthly obligation from transactions (by period date, includes completed installments)
	totalMonthlyInstallment, err := s.ledgerService.GetMonthlyObligationByReferenceType(
		userID,
		startOfMonth.Format("2006-01-02"),
		endOfMonth.Format("2006-01-02"),
		"installment_payment",
	)
	if err != nil {
		totalMonthlyInstallment = 0
	}

	totalMonthlyDebt, err := s.ledgerService.GetMonthlyObligationByReferenceType(
		userID,
		startOfMonth.Format("2006-01-02"),
		endOfMonth.Format("2006-01-02"),
		"debt_payment",
	)
	if err != nil {
		totalMonthlyDebt = 0
	}

	expenses, err := s.repos.Expense.GetByUserIDAndDateRange(
		userID,
		startOfMonth.Format("2006-01-02"),
		endOfMonth.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[uuid.UUID]*CategorySummary)
	for _, exp := range expenses {
		dashboard.TotalExpenseThisMonth += exp.Total()

		if _, exists := categoryMap[exp.CategoryID]; !exists {
			categoryMap[exp.CategoryID] = &CategorySummary{
				Category: *exp.Category,
			}
		}
		categoryMap[exp.CategoryID].TotalAmount += exp.Total()
		categoryMap[exp.CategoryID].ExpenseCount++
	}

	for _, summary := range categoryMap {
		dashboard.ExpensesByCategory = append(dashboard.ExpensesByCategory, *summary)
	}

	// Get income for this month
	incomes, err := s.repos.Income.GetByUserIDAndDateRange(
		userID,
		startOfMonth.Format("2006-01-02"),
		endOfMonth.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	for _, inc := range incomes {
		dashboard.TotalIncomeThisMonth += inc.Amount
	}

	// Calculate balance summary using total monthly obligations from active installments/debts
	netBalance := dashboard.TotalIncomeThisMonth - dashboard.TotalExpenseThisMonth - totalMonthlyInstallment - totalMonthlyDebt
	var status BalanceStatus
	if netBalance > 0 {
		status = BalanceStatusSurplus
	} else if netBalance < 0 {
		status = BalanceStatusDeficit
	} else {
		status = BalanceStatusBalanced
	}

	dashboard.BalanceSummary = BalanceSummary{
		TotalIncome:             dashboard.TotalIncomeThisMonth,
		TotalExpense:            dashboard.TotalExpenseThisMonth,
		TotalInstallmentPayment: totalMonthlyInstallment,
		TotalDebtPayment:        totalMonthlyDebt,
		NetBalance:              netBalance,
		Status:                  status,
	}

	today := now.Day()
	currentMonth := int(now.Month())
	currentYear := now.Year()
	for _, inst := range installments {
		// Check if the installment is still within its payment period
		if !isInstallmentDueThisMonth(inst.StartDate, inst.Tenor, currentMonth, currentYear) {
			continue
		}
		daysUntilDue := inst.DueDay - today
		if daysUntilDue < 0 {
			daysUntilDue += 30
		}
		if daysUntilDue <= 7 {
			dashboard.UpcomingInstallments = append(dashboard.UpcomingInstallments, inst)
		}
	}

	for _, debt := range debts {
		if debt.DueDate != nil {
			daysUntilDue := int(debt.DueDate.Sub(now).Hours() / 24)
			if daysUntilDue >= 0 && daysUntilDue <= 7 {
				dashboard.UpcomingDebts = append(dashboard.UpcomingDebts, debt)
			}
		}
	}

	// Get active savings goals
	activeSavingsGoals, err := s.repos.SavingsGoal.GetActiveByUserID(userID)
	if err != nil {
		activeSavingsGoals = nil
	}
	for _, goal := range activeSavingsGoals {
		dashboard.TotalSavingsGoal += goal.CurrentAmount
	}
	dashboard.ActiveSavingsGoals = activeSavingsGoals

	recentExpenses, err := s.repos.Expense.GetRecentByUserID(userID, 10)
	if err != nil {
		return nil, err
	}
	dashboard.RecentExpenses = recentExpenses

	return dashboard, nil
}

// isInstallmentDueThisMonth checks if an installment has a payment due in the specified month
func isInstallmentDueThisMonth(startDate time.Time, tenor int, month int, year int) bool {
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
