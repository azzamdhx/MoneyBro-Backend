package graph

import (
	"github.com/azzamdhx/moneybro/backend/internal/graph/model"
	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/services"
)

func userToModel(u *models.User) *model.User {
	return &model.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
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

func dashboardToModel(d *services.Dashboard) *model.Dashboard {
	dash := &model.Dashboard{
		TotalActiveDebt:        int(d.TotalActiveDebt),
		TotalActiveInstallment: int(d.TotalActiveInstallment),
		TotalExpenseThisMonth:  int(d.TotalExpenseThisMonth),
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
