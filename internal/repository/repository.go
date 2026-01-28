package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type Repositories struct {
	User               UserRepository
	Category           CategoryRepository
	Expense            ExpenseRepository
	ExpenseTemplate    ExpenseTemplateRepository
	Installment        InstallmentRepository
	InstallmentPayment InstallmentPaymentRepository
	Debt               DebtRepository
	DebtPayment        DebtPaymentRepository
	NotificationLog    NotificationLogRepository
	IncomeCategory     IncomeCategoryRepository
	Income             IncomeRepository
	RecurringIncome    RecurringIncomeRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:               NewUserRepository(db),
		Category:           NewCategoryRepository(db),
		Expense:            NewExpenseRepository(db),
		ExpenseTemplate:    NewExpenseTemplateRepository(db),
		Installment:        NewInstallmentRepository(db),
		InstallmentPayment: NewInstallmentPaymentRepository(db),
		Debt:               NewDebtRepository(db),
		DebtPayment:        NewDebtPaymentRepository(db),
		NotificationLog:    NewNotificationLogRepository(db),
		IncomeCategory:     NewIncomeCategoryRepository(db),
		Income:             NewIncomeRepository(db),
		RecurringIncome:    NewRecurringIncomeRepository(db),
	}
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
}

type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id uuid.UUID) (*models.Category, error)
	GetByUserID(userID uuid.UUID) ([]models.Category, error)
	Update(category *models.Category) error
	Delete(id uuid.UUID) error
}

type ExpenseRepository interface {
	Create(expense *models.Expense) error
	GetByID(id uuid.UUID) (*models.Expense, error)
	GetByUserID(userID uuid.UUID, filter *ExpenseFilter) ([]models.Expense, error)
	GetByUserIDAndDateRange(userID uuid.UUID, startDate, endDate string) ([]models.Expense, error)
	Update(expense *models.Expense) error
	Delete(id uuid.UUID) error
}

type ExpenseFilter struct {
	CategoryID *uuid.UUID
	StartDate  *string
	EndDate    *string
}

type ExpenseTemplateRepository interface {
	Create(template *models.ExpenseTemplate) error
	GetByID(id uuid.UUID) (*models.ExpenseTemplate, error)
	GetByUserID(userID uuid.UUID) ([]models.ExpenseTemplate, error)
	Update(template *models.ExpenseTemplate) error
	Delete(id uuid.UUID) error
}

type InstallmentRepository interface {
	Create(installment *models.Installment) error
	GetByID(id uuid.UUID) (*models.Installment, error)
	GetByUserID(userID uuid.UUID, status *models.InstallmentStatus) ([]models.Installment, error)
	GetByDueDay(dueDay int, status models.InstallmentStatus) ([]models.Installment, error)
	Update(installment *models.Installment) error
	Delete(id uuid.UUID) error
}

type InstallmentPaymentRepository interface {
	Create(payment *models.InstallmentPayment) error
	GetByInstallmentID(installmentID uuid.UUID) ([]models.InstallmentPayment, error)
	GetLastPaymentNumber(installmentID uuid.UUID) (int, error)
}

type DebtRepository interface {
	Create(debt *models.Debt) error
	GetByID(id uuid.UUID) (*models.Debt, error)
	GetByUserID(userID uuid.UUID, status *models.DebtStatus) ([]models.Debt, error)
	GetByDueDateRange(startDate, endDate string, status models.DebtStatus) ([]models.Debt, error)
	Update(debt *models.Debt) error
	Delete(id uuid.UUID) error
}

type DebtPaymentRepository interface {
	Create(payment *models.DebtPayment) error
	GetByDebtID(debtID uuid.UUID) ([]models.DebtPayment, error)
	GetLastPaymentNumber(debtID uuid.UUID) (int, error)
}

type NotificationLogRepository interface {
	Create(log *models.NotificationLog) error
	ExistsForToday(userID, referenceID uuid.UUID, notificationType models.NotificationType) (bool, error)
}

type IncomeCategoryRepository interface {
	Create(category *models.IncomeCategory) error
	GetByID(id uuid.UUID) (*models.IncomeCategory, error)
	GetByUserID(userID uuid.UUID) ([]models.IncomeCategory, error)
	Update(category *models.IncomeCategory) error
	Delete(id uuid.UUID) error
}

type IncomeFilter struct {
	CategoryID *uuid.UUID
	IncomeType *models.IncomeType
	StartDate  *string
	EndDate    *string
}

type IncomeRepository interface {
	Create(income *models.Income) error
	GetByID(id uuid.UUID) (*models.Income, error)
	GetByUserID(userID uuid.UUID, filter *IncomeFilter) ([]models.Income, error)
	GetByUserIDAndDateRange(userID uuid.UUID, startDate, endDate string) ([]models.Income, error)
	Update(income *models.Income) error
	Delete(id uuid.UUID) error
}

type RecurringIncomeRepository interface {
	Create(recurringIncome *models.RecurringIncome) error
	GetByID(id uuid.UUID) (*models.RecurringIncome, error)
	GetByUserID(userID uuid.UUID, isActive *bool) ([]models.RecurringIncome, error)
	GetByRecurringDay(day int, isActive bool) ([]models.RecurringIncome, error)
	Update(recurringIncome *models.RecurringIncome) error
	Delete(id uuid.UUID) error
}
