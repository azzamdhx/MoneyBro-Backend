package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type BalanceService struct {
	repos         *repository.Repositories
	ledgerService *LedgerService
}

func NewBalanceService(repos *repository.Repositories, ledgerService *LedgerService) *BalanceService {
	return &BalanceService{
		repos:         repos,
		ledgerService: ledgerService,
	}
}

type BalancePeriod string

const (
	BalancePeriodThisMonth BalancePeriod = "THIS_MONTH"
	BalancePeriodLastMonth BalancePeriod = "LAST_MONTH"
	BalancePeriodThisYear  BalancePeriod = "THIS_YEAR"
	BalancePeriodCustom    BalancePeriod = "CUSTOM"
)

type BalanceFilterInput struct {
	Period    BalancePeriod
	StartDate *time.Time
	EndDate   *time.Time
}

type IncomeCategorySummary struct {
	Category    models.IncomeCategory
	TotalAmount int64
	IncomeCount int
}

type IncomeTypeSummary struct {
	IncomeType  models.IncomeType
	TotalAmount int64
	IncomeCount int
}

type IncomeBreakdown struct {
	Total      int64
	Count      int
	ByCategory []IncomeCategorySummary
	ByType     []IncomeTypeSummary
}

type ExpenseBreakdown struct {
	Total      int64
	Count      int
	ByCategory []CategorySummary
}

type BalanceBreakdown struct {
	Total int64
	Count int
}

type BalanceReport struct {
	PeriodLabel string
	StartDate   time.Time
	EndDate     time.Time

	Income      IncomeBreakdown
	Expense     ExpenseBreakdown
	Installment BalanceBreakdown
	Debt        BalanceBreakdown

	NetBalance int64
	Status     BalanceStatus
}

func (s *BalanceService) GetBalance(userID uuid.UUID, filter BalanceFilterInput) (*BalanceReport, error) {
	startDate, endDate, periodLabel := s.calculateDateRange(filter)

	report := &BalanceReport{
		PeriodLabel: periodLabel,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	// Get incomes for the period
	incomes, err := s.repos.Income.GetByUserIDAndDateRange(
		userID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}

	incomeCategoryMap := make(map[uuid.UUID]*IncomeCategorySummary)
	incomeTypeMap := make(map[models.IncomeType]*IncomeTypeSummary)

	for _, inc := range incomes {
		report.Income.Total += inc.Amount
		report.Income.Count++

		// By category
		if _, exists := incomeCategoryMap[inc.CategoryID]; !exists {
			incomeCategoryMap[inc.CategoryID] = &IncomeCategorySummary{
				Category: *inc.Category,
			}
		}
		incomeCategoryMap[inc.CategoryID].TotalAmount += inc.Amount
		incomeCategoryMap[inc.CategoryID].IncomeCount++

		// By type
		if _, exists := incomeTypeMap[inc.IncomeType]; !exists {
			incomeTypeMap[inc.IncomeType] = &IncomeTypeSummary{
				IncomeType: inc.IncomeType,
			}
		}
		incomeTypeMap[inc.IncomeType].TotalAmount += inc.Amount
		incomeTypeMap[inc.IncomeType].IncomeCount++
	}

	for _, summary := range incomeCategoryMap {
		report.Income.ByCategory = append(report.Income.ByCategory, *summary)
	}
	for _, summary := range incomeTypeMap {
		report.Income.ByType = append(report.Income.ByType, *summary)
	}

	// Get expenses for the period
	expenses, err := s.repos.Expense.GetByUserIDAndDateRange(
		userID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}

	expenseCategoryMap := make(map[uuid.UUID]*CategorySummary)
	for _, exp := range expenses {
		report.Expense.Total += exp.Total()
		report.Expense.Count++

		if _, exists := expenseCategoryMap[exp.CategoryID]; !exists {
			expenseCategoryMap[exp.CategoryID] = &CategorySummary{
				Category: *exp.Category,
			}
		}
		expenseCategoryMap[exp.CategoryID].TotalAmount += exp.Total()
		expenseCategoryMap[exp.CategoryID].ExpenseCount++
	}

	for _, summary := range expenseCategoryMap {
		report.Expense.ByCategory = append(report.Expense.ByCategory, *summary)
	}

	// Get actual liability payments from ledger for the period (for installments)
	actualLiabilityPayments, err := s.ledgerService.GetActualPaymentsByDateRange(
		userID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
		models.AccountTypeLiability,
	)
	if err != nil {
		actualLiabilityPayments = 0
	}
	report.Installment.Total = actualLiabilityPayments

	// Get debt payments from ledger for the period
	debtPayments, err := s.ledgerService.GetActualPaymentsByDateRangeAndReferenceType(
		userID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
		models.AccountTypeLiability,
		"debt_payment",
	)
	if err != nil {
		debtPayments = 0
	}
	report.Debt.Total = debtPayments

	// Get counts from active installments and debts
	installmentActiveStatus := models.InstallmentStatusActive
	installments, _ := s.repos.Installment.GetByUserID(userID, &installmentActiveStatus)
	report.Installment.Count = len(installments)

	debtActiveStatus := models.DebtStatusActive
	debts, _ := s.repos.Debt.GetByUserID(userID, &debtActiveStatus)
	report.Debt.Count = len(debts)

	// Calculate net balance
	report.NetBalance = report.Income.Total - report.Expense.Total - report.Installment.Total - report.Debt.Total

	if report.NetBalance > 0 {
		report.Status = BalanceStatusSurplus
	} else if report.NetBalance < 0 {
		report.Status = BalanceStatusDeficit
	} else {
		report.Status = BalanceStatusBalanced
	}

	return report, nil
}

func (s *BalanceService) calculateDateRange(filter BalanceFilterInput) (time.Time, time.Time, string) {
	now := time.Now()

	switch filter.Period {
	case BalancePeriodThisMonth:
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, -1)
		return start, end, fmt.Sprintf("%s %d", now.Month().String(), now.Year())

	case BalancePeriodLastMonth:
		lastMonth := now.AddDate(0, -1, 0)
		start := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, -1)
		return start, end, fmt.Sprintf("%s %d", lastMonth.Month().String(), lastMonth.Year())

	case BalancePeriodThisYear:
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end := time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
		return start, end, fmt.Sprintf("Year %d", now.Year())

	case BalancePeriodCustom:
		if filter.StartDate != nil && filter.EndDate != nil {
			return *filter.StartDate, *filter.EndDate, fmt.Sprintf("%s - %s",
				filter.StartDate.Format("02 Jan 2006"),
				filter.EndDate.Format("02 Jan 2006"))
		}
		// Fallback to this month if custom dates not provided
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, -1)
		return start, end, fmt.Sprintf("%s %d", now.Month().String(), now.Year())

	default:
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, -1)
		return start, end, fmt.Sprintf("%s %d", now.Month().String(), now.Year())
	}
}

func (s *BalanceService) calculateMonths(start, end time.Time) int {
	months := (end.Year()-start.Year())*12 + int(end.Month()) - int(start.Month()) + 1
	if months < 1 {
		months = 1
	}
	return months
}
