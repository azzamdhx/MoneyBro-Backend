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
	PasswordResetToken PasswordResetTokenRepository
	TwoFACode          TwoFACodeRepository
	Account            AccountRepository
	Transaction        TransactionRepository
	TransactionEntry   TransactionEntryRepository
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
		PasswordResetToken: NewPasswordResetTokenRepository(db),
		TwoFACode:          NewTwoFACodeRepository(db),
		Account:            NewAccountRepository(db),
		Transaction:        NewTransactionRepository(db),
		TransactionEntry:   NewTransactionEntryRepository(db),
	}
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAllWithNotificationsEnabled() ([]models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	DeleteAllUserData(userID uuid.UUID) error
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
	GetByID(id uuid.UUID) (*models.InstallmentPayment, error)
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
	GetByID(id uuid.UUID) (*models.DebtPayment, error)
	GetByDebtID(debtID uuid.UUID) ([]models.DebtPayment, error)
	GetLastPaymentNumber(debtID uuid.UUID) (int, error)
}

type NotificationLogRepository interface {
	Create(log *models.NotificationLog) error
	GetByUserID(userID uuid.UUID) ([]models.NotificationLog, error)
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

type PasswordResetTokenRepository interface {
	Create(token *models.PasswordResetToken) error
	GetByToken(token string) (*models.PasswordResetToken, error)
	GetValidByToken(token string) (*models.PasswordResetToken, error)
	MarkAsUsed(id uuid.UUID) error
	DeleteExpiredTokens() error
	DeleteByUserID(userID uuid.UUID) error
}

type TwoFACodeRepository interface {
	Create(code *models.TwoFACode) error
	GetValidByUserIDAndCode(userID uuid.UUID, code string) (*models.TwoFACode, error)
	MarkAsUsed(id uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error
	DeleteExpiredCodes() error
}

type AccountRepository interface {
	Create(account *models.Account) error
	GetByID(id uuid.UUID) (*models.Account, error)
	GetByUserID(userID uuid.UUID) ([]models.Account, error)
	GetByUserIDAndType(userID uuid.UUID, accountType models.AccountType) ([]models.Account, error)
	GetByUserIDAndTypeAndReferenceType(userID uuid.UUID, accountType models.AccountType, referenceType string) ([]models.Account, error)
	GetDefaultByUserID(userID uuid.UUID) (*models.Account, error)
	GetByReference(referenceID uuid.UUID, referenceType string) (*models.Account, error)
	Update(account *models.Account) error
	UpdateBalance(id uuid.UUID, balance int64) error
	AddToBalance(id uuid.UUID, amount int64) error
	Delete(id uuid.UUID) error
	DeleteByReference(referenceID uuid.UUID, referenceType string) error
}

type TransactionRepository interface {
	Create(tx *models.Transaction) error
	GetByID(id uuid.UUID) (*models.Transaction, error)
	GetByUserID(userID uuid.UUID) ([]models.Transaction, error)
	GetByUserIDAndDateRange(userID uuid.UUID, startDate, endDate string) ([]models.Transaction, error)
	GetByUserIDAndDateRangeAndReferenceType(userID uuid.UUID, startDate, endDate, referenceType string) ([]models.Transaction, error)
	GetByReference(referenceID uuid.UUID, referenceType string) (*models.Transaction, error)
	Delete(id uuid.UUID) error
	DeleteByReference(referenceID uuid.UUID, referenceType string) error
}

type TransactionEntryRepository interface {
	Create(entry *models.TransactionEntry) error
	CreateBatch(entries []models.TransactionEntry) error
	GetByTransactionID(transactionID uuid.UUID) ([]models.TransactionEntry, error)
	GetByAccountID(accountID uuid.UUID) ([]models.TransactionEntry, error)
	GetByAccountIDAndDateRange(accountID uuid.UUID, startDate, endDate string) ([]models.TransactionEntry, error)
	GetByUserIDAndDateRange(userID uuid.UUID, startDate, endDate string) ([]models.TransactionEntry, error)
	DeleteByTransactionID(transactionID uuid.UUID) error
	SumByAccountID(accountID uuid.UUID) (debit int64, credit int64, err error)
	SumByAccountIDAndDateRange(accountID uuid.UUID, startDate, endDate string) (debit int64, credit int64, err error)
}
