package services

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type Config struct {
	DB                *gorm.DB
	Repos             *repository.Repositories
	Redis             *redis.Client
	JWTSecret         string
	ResendAPIKey      string
	FrontendURL       string
	EmailTemplatesDir string
}

type Services struct {
	Auth            *AuthService
	User            *UserService
	Category        *CategoryService
	Expense         *ExpenseService
	ExpenseTemplate *ExpenseTemplateService
	Installment     *InstallmentService
	Debt            *DebtService
	Dashboard       *DashboardService
	Email           *EmailService
	Notification    *NotificationService
	IncomeCategory  *IncomeCategoryService
	Income          *IncomeService
	RecurringIncome *RecurringIncomeService
	Balance         *BalanceService
	Account         *AccountService
	Ledger          *LedgerService
}

func NewServices(cfg Config) *Services {
	emailService := NewEmailService(cfg.ResendAPIKey, cfg.EmailTemplatesDir)
	accountService := NewAccountService(cfg.Repos.Account)
	ledgerService := NewLedgerService(cfg.DB, cfg.Repos.Account, cfg.Repos.Transaction, cfg.Repos.TransactionEntry)

	return &Services{
		Auth:            NewAuthService(cfg.Repos.User, cfg.Repos.PasswordResetToken, cfg.Repos.TwoFACode, emailService, cfg.JWTSecret, cfg.FrontendURL, accountService),
		User:            NewUserService(cfg.Repos.User),
		Category:        NewCategoryService(cfg.Repos.Category, accountService),
		Expense:         NewExpenseService(cfg.Repos.Expense, cfg.Repos.Category, cfg.Repos.Account, ledgerService),
		ExpenseTemplate: NewExpenseTemplateService(cfg.Repos.ExpenseTemplate, cfg.Repos.Expense, cfg.Repos.Category),
		Installment:     NewInstallmentService(cfg.Repos.Installment, cfg.Repos.InstallmentPayment, cfg.Repos.Account, accountService, ledgerService),
		Debt:            NewDebtService(cfg.Repos.Debt, cfg.Repos.DebtPayment, cfg.Repos.Account, accountService, ledgerService),
		Dashboard:       NewDashboardService(cfg.Repos, cfg.Redis, ledgerService),
		Email:           emailService,
		Notification:    NewNotificationService(cfg.Repos, emailService),
		IncomeCategory:  NewIncomeCategoryService(cfg.Repos.IncomeCategory, accountService),
		Income:          NewIncomeService(cfg.Repos.Income, cfg.Repos.IncomeCategory, cfg.Repos.Account, ledgerService),
		RecurringIncome: NewRecurringIncomeService(cfg.Repos.RecurringIncome, cfg.Repos.Income, cfg.Repos.IncomeCategory),
		Balance:         NewBalanceService(cfg.Repos, ledgerService),
		Account:         accountService,
		Ledger:          ledgerService,
	}
}
