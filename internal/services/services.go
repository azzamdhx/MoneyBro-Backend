package services

import (
	"github.com/redis/go-redis/v9"

	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type Config struct {
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
}

func NewServices(cfg Config) *Services {
	emailService := NewEmailService(cfg.ResendAPIKey, cfg.EmailTemplatesDir)

	return &Services{
		Auth:            NewAuthService(cfg.Repos.User, cfg.Repos.PasswordResetToken, emailService, cfg.JWTSecret, cfg.FrontendURL),
		User:            NewUserService(cfg.Repos.User),
		Category:        NewCategoryService(cfg.Repos.Category),
		Expense:         NewExpenseService(cfg.Repos.Expense, cfg.Repos.Category),
		ExpenseTemplate: NewExpenseTemplateService(cfg.Repos.ExpenseTemplate, cfg.Repos.Expense, cfg.Repos.Category),
		Installment:     NewInstallmentService(cfg.Repos.Installment, cfg.Repos.InstallmentPayment),
		Debt:            NewDebtService(cfg.Repos.Debt, cfg.Repos.DebtPayment),
		Dashboard:       NewDashboardService(cfg.Repos, cfg.Redis),
		Email:           emailService,
		Notification:    NewNotificationService(cfg.Repos, emailService),
		IncomeCategory:  NewIncomeCategoryService(cfg.Repos.IncomeCategory),
		Income:          NewIncomeService(cfg.Repos.Income, cfg.Repos.IncomeCategory),
		RecurringIncome: NewRecurringIncomeService(cfg.Repos.RecurringIncome, cfg.Repos.Income, cfg.Repos.IncomeCategory),
		Balance:         NewBalanceService(cfg.Repos),
	}
}
