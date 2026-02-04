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
		TwoFAEnabled:      u.TwoFAEnabled,
		NotifyInstallment: u.NotifyInstallment,
		NotifyDebt:        u.NotifyDebt,
		NotifyDaysBefore:  u.NotifyDaysBefore,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
	}
}

func categoryToModel(c *models.Category) *model.Category {
	return &model.Category{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
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

func expenseTemplateToModel(t *models.ExpenseTemplate) *model.ExpenseTemplate {
	tmpl := &model.ExpenseTemplate{
		ID:           t.ID,
		ItemName:     t.ItemName,
		UnitPrice:    int(t.UnitPrice),
		Quantity:     t.Quantity,
		Total:        int(t.Total()),
		RecurringDay: t.RecurringDay,
		Notes:        t.Notes,
		CreatedAt:    t.CreatedAt,
	}
	if t.Category != nil {
		tmpl.Category = categoryToModel(t.Category)
	}
	return tmpl
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
		TotalActiveDebt:        int(d.TotalActiveDebt),
		TotalActiveInstallment: int(d.TotalActiveInstallment),
		TotalExpenseThisMonth:  int(d.TotalExpenseThisMonth),
		TotalIncomeThisMonth:   int(d.TotalIncomeThisMonth),
		BalanceSummary:         balanceSummaryToModel(&d.BalanceSummary),
	}
	if len(d.UpcomingInstallments) > 0 {
		insts := make([]*model.Installment, len(d.UpcomingInstallments))
		for i, inst := range d.UpcomingInstallments {
			insts[i] = installmentToModel(&inst)
		}
		dash.UpcomingInstallments = insts
	}
	if len(d.UpcomingDebts) > 0 {
		debts := make([]*model.Debt, len(d.UpcomingDebts))
		for i, debt := range d.UpcomingDebts {
			debts[i] = debtToModel(&debt)
		}
		dash.UpcomingDebts = debts
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
