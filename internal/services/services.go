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
	Auth                 *AuthService
	User                 *UserService
	Category             *CategoryService
	Expense              *ExpenseService
	ExpenseTemplateGroup *ExpenseTemplateGroupService
	Installment          *InstallmentService
	Debt                 *DebtService
	Dashboard            *DashboardService
	Email                *EmailService
	Notification         *NotificationService
	IncomeCategory       *IncomeCategoryService
	Income               *IncomeService
	RecurringIncome      *RecurringIncomeService
	Balance              *BalanceService
	Account              *AccountService
	Ledger               *LedgerService
	UpcomingPayments     *UpcomingPaymentsService
	ActualPayments       *ActualPaymentsService
	SavingsGoal          *SavingsGoalService
	MonthlySummary       *MonthlySummaryService
}

func NewServices(cfg Config) *Services {
	emailService := NewEmailService(cfg.ResendAPIKey, cfg.EmailTemplatesDir)
	accountService := NewAccountService(cfg.Repos.Account)
	ledgerService := NewLedgerService(cfg.DB, cfg.Repos.Account, cfg.Repos.Transaction, cfg.Repos.TransactionEntry)

	// Create services that will be dependencies for others
	incomeService := NewIncomeService(cfg.Repos.Income, cfg.Repos.IncomeCategory, cfg.Repos.Account, ledgerService)
	expenseService := NewExpenseService(cfg.Repos.Expense, cfg.Repos.Category, cfg.Repos.Account, ledgerService)

	return &Services{
		Auth:                 NewAuthService(cfg.Repos.User, cfg.Repos.PasswordResetToken, cfg.Repos.TwoFACode, cfg.Repos.RefreshToken, emailService, cfg.JWTSecret, cfg.FrontendURL, accountService),
		User:                 NewUserService(cfg.Repos.User),
		Category:             NewCategoryService(cfg.Repos.Category, accountService),
		Expense:              expenseService,
		ExpenseTemplateGroup: NewExpenseTemplateGroupService(cfg.Repos.ExpenseTemplateGroup, expenseService, cfg.Repos.Category),
		Installment:          NewInstallmentService(cfg.Repos.Installment, cfg.Repos.InstallmentPayment, cfg.Repos.Account, accountService, ledgerService),
		Debt:                 NewDebtService(cfg.Repos.Debt, cfg.Repos.DebtPayment, cfg.Repos.Account, accountService, ledgerService),
		Dashboard:            NewDashboardService(cfg.Repos, cfg.Redis, ledgerService),
		Email:                emailService,
		Notification:         NewNotificationService(cfg.Repos, emailService),
		IncomeCategory:       NewIncomeCategoryService(cfg.Repos.IncomeCategory, accountService),
		Income:               incomeService,
		RecurringIncome:      NewRecurringIncomeService(cfg.Repos.RecurringIncome, incomeService, cfg.Repos.IncomeCategory),
		Balance:              NewBalanceService(cfg.Repos, ledgerService),
		Account:              accountService,
		Ledger:               ledgerService,
		UpcomingPayments:     NewUpcomingPaymentsService(cfg.Repos),
		ActualPayments:       NewActualPaymentsService(cfg.Repos),
		SavingsGoal:          NewSavingsGoalService(cfg.Repos.SavingsGoal, cfg.Repos.SavingsContribution, cfg.Repos.Account, accountService, ledgerService),
		MonthlySummary:       NewMonthlySummaryService(cfg.Repos, NewUpcomingPaymentsService(cfg.Repos), NewActualPaymentsService(cfg.Repos)),
	}
}
