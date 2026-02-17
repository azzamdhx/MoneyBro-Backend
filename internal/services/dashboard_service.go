package services

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

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
	TotalActiveDebt                   int64
	TotalActiveInstallment            int64
	TotalExpenseThisMonth             int64
	TotalIncomeThisMonth              int64
	TotalSavingsContributionThisMonth int64
	BalanceSummary                    BalanceSummary
	ActiveSavingsGoals                []models.SavingsGoal
	ExpensesByCategory                []CategorySummary
	RecentExpenses                    []models.Expense
}

type CategorySummary struct {
	Category     models.Category
	TotalAmount  int64
	ExpenseCount int
}

func (s *DashboardService) GetDashboard(userID uuid.UUID) (*Dashboard, error) {
	dashboard := &Dashboard{}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	startStr := startOfMonth.Format("2006-01-02")
	endStr := endOfMonth.Format("2006-01-02")

	var (
		mu                      sync.Mutex
		totalMonthlyInstallment int64
		totalMonthlyDebt        int64
	)

	var g errgroup.Group

	// 1. Active debts
	g.Go(func() error {
		activeStatus := models.DebtStatusActive
		debts, err := s.repos.Debt.GetByUserID(userID, &activeStatus)
		if err != nil {
			return err
		}
		var total int64
		for _, debt := range debts {
			total += debt.RemainingAmount()
		}
		mu.Lock()
		dashboard.TotalActiveDebt = total
		mu.Unlock()
		return nil
	})

	// 2. Active installments
	g.Go(func() error {
		activeStatus := models.InstallmentStatusActive
		installments, err := s.repos.Installment.GetByUserID(userID, &activeStatus)
		if err != nil {
			return err
		}
		var total int64
		for _, inst := range installments {
			total += inst.RemainingAmount()
		}
		mu.Lock()
		dashboard.TotalActiveInstallment = total
		mu.Unlock()
		return nil
	})

	// 3. Monthly installment obligation
	g.Go(func() error {
		total, err := s.ledgerService.GetMonthlyObligationByReferenceType(userID, startStr, endStr, "installment_payment")
		if err != nil {
			total = 0
		}
		mu.Lock()
		totalMonthlyInstallment = total
		mu.Unlock()
		return nil
	})

	// 4. Monthly debt obligation
	g.Go(func() error {
		total, err := s.ledgerService.GetMonthlyObligationByReferenceType(userID, startStr, endStr, "debt_payment")
		if err != nil {
			total = 0
		}
		mu.Lock()
		totalMonthlyDebt = total
		mu.Unlock()
		return nil
	})

	// 5. Expenses this month + by category
	g.Go(func() error {
		expenses, err := s.repos.Expense.GetByUserIDAndDateRange(userID, startStr, endStr)
		if err != nil {
			return err
		}
		var totalExpense int64
		categoryMap := make(map[uuid.UUID]*CategorySummary)
		for _, exp := range expenses {
			totalExpense += exp.Total()
			if exp.Category != nil {
				if _, exists := categoryMap[exp.CategoryID]; !exists {
					categoryMap[exp.CategoryID] = &CategorySummary{Category: *exp.Category}
				}
				categoryMap[exp.CategoryID].TotalAmount += exp.Total()
				categoryMap[exp.CategoryID].ExpenseCount++
			}
		}
		var byCategory []CategorySummary
		for _, cs := range categoryMap {
			byCategory = append(byCategory, *cs)
		}
		mu.Lock()
		dashboard.TotalExpenseThisMonth = totalExpense
		dashboard.ExpensesByCategory = byCategory
		mu.Unlock()
		return nil
	})

	// 6. Income this month
	g.Go(func() error {
		incomes, err := s.repos.Income.GetByUserIDAndDateRange(userID, startStr, endStr)
		if err != nil {
			return err
		}
		var total int64
		for _, inc := range incomes {
			total += inc.Amount
		}
		mu.Lock()
		dashboard.TotalIncomeThisMonth = total
		mu.Unlock()
		return nil
	})

	// 7. Active savings goals
	g.Go(func() error {
		goals, err := s.repos.SavingsGoal.GetActiveByUserID(userID)
		if err != nil {
			goals = nil
		}
		mu.Lock()
		dashboard.ActiveSavingsGoals = goals
		mu.Unlock()
		return nil
	})

	// 8. Savings contributions this month
	g.Go(func() error {
		total, err := s.repos.SavingsContribution.GetTotalByUserIDAndDateRange(userID, startStr, endStr)
		if err != nil {
			total = 0
		}
		mu.Lock()
		dashboard.TotalSavingsContributionThisMonth = total
		mu.Unlock()
		return nil
	})

	// 9. Recent expenses
	g.Go(func() error {
		recent, err := s.repos.Expense.GetRecentByUserID(userID, 10)
		if err != nil {
			return err
		}
		mu.Lock()
		dashboard.RecentExpenses = recent
		mu.Unlock()
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Calculate balance summary after all parallel queries complete
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

	return dashboard, nil
}
