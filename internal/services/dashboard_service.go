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
	BalanceSummary         BalanceSummary
	UpcomingInstallments   []models.Installment
	UpcomingDebts          []models.Debt
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

	installmentActiveStatus := models.InstallmentStatusActive
	installments, err := s.repos.Installment.GetByUserID(userID, &installmentActiveStatus)
	if err != nil {
		return nil, err
	}
	for _, inst := range installments {
		dashboard.TotalActiveInstallment += inst.RemainingAmount()
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

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

	// Get actual installment payments this month from ledger (only installment reference_type)
	actualInstallmentPayment, err := s.ledgerService.GetActualPaymentsByDateRangeAndReferenceType(
		userID,
		startOfMonth.Format("2006-01-02"),
		endOfMonth.Format("2006-01-02"),
		models.AccountTypeLiability,
		"installment",
	)
	if err != nil {
		actualInstallmentPayment = 0
	}

	// Get actual debt payments this month from ledger (only debt reference_type)
	actualDebtPayment, err := s.ledgerService.GetActualPaymentsByDateRangeAndReferenceType(
		userID,
		startOfMonth.Format("2006-01-02"),
		endOfMonth.Format("2006-01-02"),
		models.AccountTypeLiability,
		"debt",
	)
	if err != nil {
		actualDebtPayment = 0
	}

	// Calculate balance summary using actual payments
	netBalance := dashboard.TotalIncomeThisMonth - dashboard.TotalExpenseThisMonth - actualInstallmentPayment - actualDebtPayment
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
		TotalInstallmentPayment: actualInstallmentPayment,
		TotalDebtPayment:        actualDebtPayment,
		NetBalance:              netBalance,
		Status:                  status,
	}

	today := now.Day()
	for _, inst := range installments {
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

	allExpenses, err := s.repos.Expense.GetByUserID(userID, nil)
	if err != nil {
		return nil, err
	}
	limit := 10
	if len(allExpenses) < limit {
		limit = len(allExpenses)
	}
	dashboard.RecentExpenses = allExpenses[:limit]

	return dashboard, nil
}
