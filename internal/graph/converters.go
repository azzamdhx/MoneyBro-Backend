package graph

import (
	"github.com/azzamdhx/moneybro/backend/internal/graph/model"
	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/services"
)

func userToModel(u *models.User) *model.User {
	return &model.User{
		ID:                u.ID,
		Email:             u.Email,
		Name:              u.Name,
		ProfileImage:      u.ProfileImage,
		TwoFAEnabled:      u.TwoFAEnabled,
		NotifyInstallment: u.NotifyInstallment,
		NotifyDebt:        u.NotifyDebt,
		NotifySavingsGoal: u.NotifySavingsGoal,
		NotifyDaysBefore:  u.NotifyDaysBefore,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
	}
}

func categoryToModel(c *models.Category) *model.Category {
	return &model.Category{
		ID:           c.ID,
		Name:         c.Name,
		CreatedAt:    c.CreatedAt,
		ExpenseCount: c.ExpenseCount,
		TotalSpent:   int(c.TotalSpent),
	}
}

func expenseToModel(e *models.Expense) *model.Expense {
	exp := &model.Expense{
		ID:          e.ID,
		ItemName:    e.ItemName,
		UnitPrice:   int(e.UnitPrice),
		Quantity:    e.Quantity,
		Total:       int(e.Total()),
		Notes:       e.Notes,
		ExpenseDate: e.ExpenseDate,
		CreatedAt:   e.CreatedAt,
	}
	if e.Category != nil {
		exp.Category = categoryToModel(e.Category)
	}
	return exp
}

func expenseTemplateGroupToModel(g *models.ExpenseTemplateGroup) *model.ExpenseTemplateGroup {
	group := &model.ExpenseTemplateGroup{
		ID:           g.ID,
		Name:         g.Name,
		RecurringDay: g.RecurringDay,
		Notes:        g.Notes,
		Total:        int(g.Total()),
		CreatedAt:    g.CreatedAt,
	}
	if len(g.Items) > 0 {
		items := make([]*model.ExpenseTemplateItem, len(g.Items))
		for i, item := range g.Items {
			items[i] = expenseTemplateItemToModel(&item)
		}
		group.Items = items
	}
	return group
}

func expenseTemplateItemToModel(i *models.ExpenseTemplateItem) *model.ExpenseTemplateItem {
	item := &model.ExpenseTemplateItem{
		ID:        i.ID,
		ItemName:  i.ItemName,
		UnitPrice: int(i.UnitPrice),
		Quantity:  i.Quantity,
		Total:     int(i.Total()),
		CreatedAt: i.CreatedAt,
	}
	if i.Category != nil {
		item.Category = categoryToModel(i.Category)
	}
	return item
}

func installmentToModel(i *models.Installment) *model.Installment {
	inst := &model.Installment{
		ID:                 i.ID,
		Name:               i.Name,
		ActualAmount:       int(i.ActualAmount),
		LoanAmount:         int(i.LoanAmount),
		MonthlyPayment:     int(i.MonthlyPayment),
		Tenor:              i.Tenor,
		StartDate:          i.StartDate,
		DueDay:             i.DueDay,
		Status:             model.InstallmentStatus(i.Status),
		Notes:              i.Notes,
		CreatedAt:          i.CreatedAt,
		InterestAmount:     int(i.InterestAmount()),
		InterestPercentage: i.InterestPercentage(),
		PaidCount:          i.PaidCount(),
		RemainingPayments:  i.RemainingPayments(),
		RemainingAmount:    int(i.RemainingAmount()),
	}
	if len(i.Payments) > 0 {
		payments := make([]*model.InstallmentPayment, len(i.Payments))
		for j, p := range i.Payments {
			payments[j] = installmentPaymentToModel(&p)
		}
		inst.Payments = payments
	}
	return inst
}

func installmentPaymentToModel(p *models.InstallmentPayment) *model.InstallmentPayment {
	return &model.InstallmentPayment{
		ID:            p.ID,
		PaymentNumber: p.PaymentNumber,
		Amount:        int(p.Amount),
		PaidAt:        p.PaidAt,
		CreatedAt:     p.CreatedAt,
	}
}

func debtToModel(d *models.Debt) *model.Debt {
	debt := &model.Debt{
		ID:              d.ID,
		PersonName:      d.PersonName,
		ActualAmount:    int(d.ActualAmount),
		PaymentType:     model.DebtPaymentType(d.PaymentType),
		Status:          model.DebtStatus(d.Status),
		Notes:           d.Notes,
		DueDate:         d.DueDate,
		CreatedAt:       d.CreatedAt,
		TotalToPay:      int(d.TotalToPay()),
		PaidAmount:      int(d.PaidAmount()),
		RemainingAmount: int(d.RemainingAmount()),
	}
	if d.LoanAmount != nil {
		v := int(*d.LoanAmount)
		debt.LoanAmount = &v
	}
	if d.MonthlyPayment != nil {
		v := int(*d.MonthlyPayment)
		debt.MonthlyPayment = &v
	}
	debt.Tenor = d.Tenor
	if interestAmt := d.InterestAmount(); interestAmt != nil && *interestAmt > 0 {
		v := int(*interestAmt)
		debt.InterestAmount = &v
	}
	if interestPct := d.InterestPercentage(); interestPct != nil && *interestPct > 0 {
		debt.InterestPercentage = interestPct
	}
	if len(d.Payments) > 0 {
		payments := make([]*model.DebtPayment, len(d.Payments))
		for j, p := range d.Payments {
			payments[j] = debtPaymentToModel(&p)
		}
		debt.Payments = payments
	}
	return debt
}

func debtPaymentToModel(p *models.DebtPayment) *model.DebtPayment {
	return &model.DebtPayment{
		ID:            p.ID,
		PaymentNumber: p.PaymentNumber,
		Amount:        int(p.Amount),
		PaidAt:        p.PaidAt,
		CreatedAt:     p.CreatedAt,
	}
}

func incomeCategoryToModel(c *models.IncomeCategory) *model.IncomeCategory {
	return &model.IncomeCategory{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
	}
}

func incomeToModel(i *models.Income) *model.Income {
	inc := &model.Income{
		ID:          i.ID,
		SourceName:  i.SourceName,
		Amount:      int(i.Amount),
		IncomeType:  model.IncomeType(i.IncomeType),
		IncomeDate:  i.IncomeDate,
		IsRecurring: i.IsRecurring,
		Notes:       i.Notes,
		CreatedAt:   i.CreatedAt,
	}
	if i.Category != nil {
		inc.Category = incomeCategoryToModel(i.Category)
	}
	return inc
}

func recurringIncomeToModel(r *models.RecurringIncome) *model.RecurringIncome {
	ri := &model.RecurringIncome{
		ID:           r.ID,
		SourceName:   r.SourceName,
		Amount:       int(r.Amount),
		IncomeType:   model.IncomeType(r.IncomeType),
		RecurringDay: r.RecurringDay,
		IsActive:     r.IsActive,
		Notes:        r.Notes,
		CreatedAt:    r.CreatedAt,
	}
	if r.Category != nil {
		ri.Category = incomeCategoryToModel(r.Category)
	}
	return ri
}

func balanceSummaryToModel(b *services.BalanceSummary) *model.BalanceSummary {
	return &model.BalanceSummary{
		TotalIncome:             int(b.TotalIncome),
		TotalExpense:            int(b.TotalExpense),
		TotalInstallmentPayment: int(b.TotalInstallmentPayment),
		TotalDebtPayment:        int(b.TotalDebtPayment),
		NetBalance:              int(b.NetBalance),
		Status:                  model.BalanceStatus(b.Status),
	}
}

func balanceReportToModel(r *services.BalanceReport) *model.BalanceReport {
	report := &model.BalanceReport{
		PeriodLabel: r.PeriodLabel,
		StartDate:   r.StartDate,
		EndDate:     r.EndDate,
		NetBalance:  int(r.NetBalance),
		Status:      model.BalanceStatus(r.Status),
		Income: &model.IncomeBreakdown{
			Total: int(r.Income.Total),
			Count: r.Income.Count,
		},
		Expense: &model.ExpenseBreakdown{
			Total: int(r.Expense.Total),
			Count: r.Expense.Count,
		},
		Installment: &model.BalanceBreakdown{
			Total: int(r.Installment.Total),
			Count: r.Installment.Count,
		},
		Debt: &model.BalanceBreakdown{
			Total: int(r.Debt.Total),
			Count: r.Debt.Count,
		},
	}

	// Income by category
	if len(r.Income.ByCategory) > 0 {
		cats := make([]*model.IncomeCategorySummary, len(r.Income.ByCategory))
		for i, c := range r.Income.ByCategory {
			cats[i] = &model.IncomeCategorySummary{
				Category:    incomeCategoryToModel(&c.Category),
				TotalAmount: int(c.TotalAmount),
				IncomeCount: c.IncomeCount,
			}
		}
		report.Income.ByCategory = cats
	}

	// Income by type
	if len(r.Income.ByType) > 0 {
		types := make([]*model.IncomeTypeSummary, len(r.Income.ByType))
		for i, t := range r.Income.ByType {
			types[i] = &model.IncomeTypeSummary{
				IncomeType:  model.IncomeType(t.IncomeType),
				TotalAmount: int(t.TotalAmount),
				IncomeCount: t.IncomeCount,
			}
		}
		report.Income.ByType = types
	}

	// Expense by category
	if len(r.Expense.ByCategory) > 0 {
		cats := make([]*model.CategorySummary, len(r.Expense.ByCategory))
		for i, c := range r.Expense.ByCategory {
			cats[i] = &model.CategorySummary{
				Category:     categoryToModel(&c.Category),
				TotalAmount:  int(c.TotalAmount),
				ExpenseCount: c.ExpenseCount,
			}
		}
		report.Expense.ByCategory = cats
	}

	return report
}

func accountToModel(a *models.Account) *model.Account {
	acc := &model.Account{
		ID:             a.ID.String(),
		Name:           a.Name,
		AccountType:    model.AccountType(a.AccountType),
		CurrentBalance: int(a.CurrentBalance),
		IsDefault:      a.IsDefault,
		CreatedAt:      a.CreatedAt,
	}
	if a.ReferenceID != nil {
		refID := a.ReferenceID.String()
		acc.ReferenceID = &refID
	}
	if a.ReferenceType != nil {
		acc.ReferenceType = a.ReferenceType
	}
	return acc
}

func transactionToModel(t *models.Transaction) *model.Transaction {
	tx := &model.Transaction{
		ID:              t.ID.String(),
		TransactionDate: t.TransactionDate,
		Description:     t.Description,
		CreatedAt:       t.CreatedAt,
	}
	if t.ReferenceID != nil {
		refID := t.ReferenceID.String()
		tx.ReferenceID = &refID
	}
	if t.ReferenceType != nil {
		tx.ReferenceType = t.ReferenceType
	}
	if len(t.Entries) > 0 {
		entries := make([]*model.TransactionEntry, len(t.Entries))
		for i, e := range t.Entries {
			entries[i] = transactionEntryToModel(&e)
		}
		tx.Entries = entries
	}
	return tx
}

func transactionEntryToModel(e *models.TransactionEntry) *model.TransactionEntry {
	entry := &model.TransactionEntry{
		ID:     e.ID.String(),
		Debit:  int(e.Debit),
		Credit: int(e.Credit),
	}
	if e.Account != nil {
		entry.Account = accountToModel(e.Account)
	}
	return entry
}

func dashboardToModel(d *services.Dashboard) *model.Dashboard {
	dash := &model.Dashboard{
		TotalActiveDebt:                   int(d.TotalActiveDebt),
		TotalActiveInstallment:            int(d.TotalActiveInstallment),
		TotalExpenseThisMonth:             int(d.TotalExpenseThisMonth),
		TotalIncomeThisMonth:              int(d.TotalIncomeThisMonth),
		TotalSavingsContributionThisMonth: int(d.TotalSavingsContributionThisMonth),
		BalanceSummary:                    balanceSummaryToModel(&d.BalanceSummary),
	}
	if len(d.ActiveSavingsGoals) > 0 {
		goals := make([]*model.SavingsGoal, len(d.ActiveSavingsGoals))
		for i, g := range d.ActiveSavingsGoals {
			goals[i] = savingsGoalToModel(&g)
		}
		dash.ActiveSavingsGoals = goals
	}
	if len(d.ExpensesByCategory) > 0 {
		cats := make([]*model.CategorySummary, len(d.ExpensesByCategory))
		for i, cs := range d.ExpensesByCategory {
			cats[i] = &model.CategorySummary{
				Category:     categoryToModel(&cs.Category),
				TotalAmount:  int(cs.TotalAmount),
				ExpenseCount: cs.ExpenseCount,
			}
		}
		dash.ExpensesByCategory = cats
	}
	if len(d.RecentExpenses) > 0 {
		exps := make([]*model.Expense, len(d.RecentExpenses))
		for i, exp := range d.RecentExpenses {
			exps[i] = expenseToModel(&exp)
		}
		dash.RecentExpenses = exps
	}
	return dash
}

func upcomingPaymentsReportToModel(report *services.UpcomingPaymentsReport) *model.UpcomingPaymentsReport {
	installments := make([]*model.UpcomingInstallmentPayment, len(report.Installments))
	for i, inst := range report.Installments {
		installments[i] = &model.UpcomingInstallmentPayment{
			InstallmentID:     inst.InstallmentID.String(),
			Name:              inst.Name,
			MonthlyPayment:    int(inst.MonthlyPayment),
			DueDay:            inst.DueDay,
			DueDate:           inst.DueDate.Format("2006-01-02"),
			RemainingAmount:   int(inst.RemainingAmount),
			RemainingPayments: inst.RemainingPayments,
		}
	}

	debts := make([]*model.UpcomingDebtPayment, len(report.Debts))
	for i, debt := range report.Debts {
		debts[i] = &model.UpcomingDebtPayment{
			DebtID:          debt.DebtID.String(),
			PersonName:      debt.PersonName,
			MonthlyPayment:  int(debt.MonthlyPayment),
			DueDate:         debt.DueDate.Format("2006-01-02"),
			RemainingAmount: int(debt.RemainingAmount),
			PaymentType:     debt.PaymentType,
		}
	}

	return &model.UpcomingPaymentsReport{
		Installments:     installments,
		Debts:            debts,
		TotalInstallment: int(report.TotalInstallment),
		TotalDebt:        int(report.TotalDebt),
		TotalPayments:    int(report.TotalPayments),
	}
}

func actualPaymentsReportToModel(report *services.ActualPaymentsReport) *model.ActualPaymentsReport {
	installments := make([]*model.ActualInstallmentPayment, len(report.Installments))
	for i, inst := range report.Installments {
		installments[i] = &model.ActualInstallmentPayment{
			InstallmentID:   inst.InstallmentID.String(),
			Name:            inst.Name,
			Amount:          float64(inst.Amount),
			TransactionDate: inst.TransactionDate,
			Description:     inst.Description,
		}
	}

	debts := make([]*model.ActualDebtPayment, len(report.Debts))
	for i, debt := range report.Debts {
		debts[i] = &model.ActualDebtPayment{
			DebtID:          debt.DebtID.String(),
			PersonName:      debt.PersonName,
			Amount:          float64(debt.Amount),
			TransactionDate: debt.TransactionDate,
			Description:     debt.Description,
		}
	}

	return &model.ActualPaymentsReport{
		Installments:     installments,
		Debts:            debts,
		TotalInstallment: float64(report.TotalInstallment),
		TotalDebt:        float64(report.TotalDebt),
		TotalPayments:    float64(report.TotalPayments),
	}
}

func calculateExpenseSummary(expenses []models.Expense) *model.ExpenseSummary {
	var total int64
	count := len(expenses)
	categoryMap := make(map[string]*model.ExpenseByCategoryGroup)

	for _, exp := range expenses {
		total += exp.Total()

		// Group by category
		if exp.Category != nil {
			catID := exp.Category.ID.String()
			if group, exists := categoryMap[catID]; exists {
				group.TotalAmount += int(exp.Total())
				group.Count++
			} else {
				categoryMap[catID] = &model.ExpenseByCategoryGroup{
					Category:    categoryToModel(exp.Category),
					TotalAmount: int(exp.Total()),
					Count:       1,
				}
			}
		}
	}

	// Convert map to slice
	byCategory := make([]*model.ExpenseByCategoryGroup, 0, len(categoryMap))
	for _, group := range categoryMap {
		byCategory = append(byCategory, group)
	}

	return &model.ExpenseSummary{
		Total:      int(total),
		Count:      count,
		ByCategory: byCategory,
	}
}

func calculateIncomeSummary(incomes []models.Income) *model.IncomeSummary {
	var total int64
	count := len(incomes)
	categoryMap := make(map[string]*model.IncomeByCategoryGroup)
	typeMap := make(map[model.IncomeType]*model.IncomeByTypeGroup)

	for _, inc := range incomes {
		total += inc.Amount

		// Group by category
		if inc.Category != nil {
			catID := inc.Category.ID.String()
			if group, exists := categoryMap[catID]; exists {
				group.TotalAmount += int(inc.Amount)
				group.Count++
			} else {
				categoryMap[catID] = &model.IncomeByCategoryGroup{
					Category:    incomeCategoryToModel(inc.Category),
					TotalAmount: int(inc.Amount),
					Count:       1,
				}
			}
		}

		// Group by type
		incomeType := model.IncomeType(inc.IncomeType)
		if group, exists := typeMap[incomeType]; exists {
			group.TotalAmount += int(inc.Amount)
			group.Count++
		} else {
			typeMap[incomeType] = &model.IncomeByTypeGroup{
				IncomeType:  incomeType,
				TotalAmount: int(inc.Amount),
				Count:       1,
			}
		}
	}

	// Convert maps to slices
	byCategory := make([]*model.IncomeByCategoryGroup, 0, len(categoryMap))
	for _, group := range categoryMap {
		byCategory = append(byCategory, group)
	}

	byType := make([]*model.IncomeByTypeGroup, 0, len(typeMap))
	for _, group := range typeMap {
		byType = append(byType, group)
	}

	return &model.IncomeSummary{
		Total:      int(total),
		Count:      count,
		ByCategory: byCategory,
		ByType:     byType,
	}
}

func savingsGoalToModel(g *models.SavingsGoal) *model.SavingsGoal {
	goal := &model.SavingsGoal{
		ID:              g.ID,
		Name:            g.Name,
		TargetAmount:    int(g.TargetAmount),
		CurrentAmount:   int(g.CurrentAmount),
		TargetDate:      g.TargetDate,
		Icon:            g.Icon,
		Status:          model.SavingsGoalStatus(g.Status),
		Notes:           g.Notes,
		Progress:        g.Progress(),
		RemainingAmount: int(g.RemainingAmount()),
		MonthlyTarget:   int(g.MonthlyTarget()),
		CreatedAt:       g.CreatedAt,
	}
	if len(g.Contributions) > 0 {
		contributions := make([]*model.SavingsContribution, len(g.Contributions))
		for i, c := range g.Contributions {
			contributions[i] = savingsContributionToModel(&c)
		}
		goal.Contributions = contributions
	}
	return goal
}

func savingsContributionToModel(c *models.SavingsContribution) *model.SavingsContribution {
	contribution := &model.SavingsContribution{
		ID:               c.ID,
		Amount:           int(c.Amount),
		ContributionDate: c.ContributionDate,
		Notes:            c.Notes,
		CreatedAt:        c.CreatedAt,
	}
	if c.SavingsGoal != nil {
		contribution.SavingsGoal = savingsGoalToModel(c.SavingsGoal)
	}
	return contribution
}

func incomeBreakdownToIncomeSummary(b *services.IncomeBreakdown) *model.IncomeSummary {
	byCategory := make([]*model.IncomeByCategoryGroup, len(b.ByCategory))
	for i, cs := range b.ByCategory {
		byCategory[i] = &model.IncomeByCategoryGroup{
			Category:    incomeCategoryToModel(&cs.Category),
			TotalAmount: int(cs.TotalAmount),
			Count:       cs.IncomeCount,
		}
	}
	byType := make([]*model.IncomeByTypeGroup, len(b.ByType))
	for i, ts := range b.ByType {
		byType[i] = &model.IncomeByTypeGroup{
			IncomeType:  model.IncomeType(ts.IncomeType),
			TotalAmount: int(ts.TotalAmount),
			Count:       ts.IncomeCount,
		}
	}
	return &model.IncomeSummary{
		Total:      int(b.Total),
		Count:      b.Count,
		ByCategory: byCategory,
		ByType:     byType,
	}
}

func expenseBreakdownToExpenseSummary(b *services.ExpenseBreakdown) *model.ExpenseSummary {
	byCategory := make([]*model.ExpenseByCategoryGroup, len(b.ByCategory))
	for i, cs := range b.ByCategory {
		byCategory[i] = &model.ExpenseByCategoryGroup{
			Category:    categoryToModel(&cs.Category),
			TotalAmount: int(cs.TotalAmount),
			Count:       cs.ExpenseCount,
		}
	}
	return &model.ExpenseSummary{
		Total:      int(b.Total),
		Count:      b.Count,
		ByCategory: byCategory,
	}
}

func historySummaryToModel(r *services.HistorySummaryResult) *model.HistorySummary {
	return &model.HistorySummary{
		AvailableMonths:          r.AvailableMonths,
		SelectedMonth:            r.SelectedMonth,
		IncomeSummary:            incomeBreakdownToIncomeSummary(&r.IncomeSummary),
		ExpenseSummary:           expenseBreakdownToExpenseSummary(&r.ExpenseSummary),
		Payments:                 actualPaymentsReportToModel(r.Payments),
		TotalSavingsContribution: int(r.TotalSavingsContribution),
	}
}

func forecastSummaryToModel(r *services.ForecastSummaryResult) *model.ForecastSummary {
	return &model.ForecastSummary{
		AvailableMonths:          r.AvailableMonths,
		SelectedMonth:            r.SelectedMonth,
		IncomeSummary:            incomeBreakdownToIncomeSummary(&r.IncomeSummary),
		ExpenseSummary:           expenseBreakdownToExpenseSummary(&r.ExpenseSummary),
		Payments:                 upcomingPaymentsReportToModel(r.Payments),
		TotalSavingsContribution: int(r.TotalSavingsContribution),
	}
}
