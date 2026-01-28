package services

import (
	"github.com/redis/go-redis/v9"

	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type Config struct {
	Repos        *repository.Repositories
	Redis        *redis.Client
	JWTSecret    string
	ResendAPIKey string
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
}

func NewServices(cfg Config) *Services {
	emailService := NewEmailService(cfg.ResendAPIKey)

	return &Services{
		Auth:            NewAuthService(cfg.Repos.User, cfg.JWTSecret),
		User:            NewUserService(cfg.Repos.User),
		Category:        NewCategoryService(cfg.Repos.Category),
		Expense:         NewExpenseService(cfg.Repos.Expense, cfg.Repos.Category),
		ExpenseTemplate: NewExpenseTemplateService(cfg.Repos.ExpenseTemplate, cfg.Repos.Expense, cfg.Repos.Category),
		Installment:     NewInstallmentService(cfg.Repos.Installment, cfg.Repos.InstallmentPayment),
		Debt:            NewDebtService(cfg.Repos.Debt, cfg.Repos.DebtPayment),
		Dashboard:       NewDashboardService(cfg.Repos, cfg.Redis),
		Email:           emailService,
		Notification:    NewNotificationService(cfg.Repos, emailService),
	}
}
