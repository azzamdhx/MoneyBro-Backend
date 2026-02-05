package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type ExpenseTemplateService struct {
	templateRepo   repository.ExpenseTemplateRepository
	expenseService *ExpenseService
	categoryRepo   repository.CategoryRepository
}

func NewExpenseTemplateService(
	templateRepo repository.ExpenseTemplateRepository,
	expenseService *ExpenseService,
	categoryRepo repository.CategoryRepository,
) *ExpenseTemplateService {
	return &ExpenseTemplateService{
		templateRepo:   templateRepo,
		expenseService: expenseService,
		categoryRepo:   categoryRepo,
	}
}

type CreateExpenseTemplateInput struct {
	CategoryID   uuid.UUID
	ItemName     string
	UnitPrice    int64
	Quantity     int
	RecurringDay *int
	Notes        *string
}

func (s *ExpenseTemplateService) Create(userID uuid.UUID, input CreateExpenseTemplateInput) (*models.ExpenseTemplate, error) {
	if input.ItemName == "" {
		return nil, errors.New("item name is required")
	}
	if input.UnitPrice <= 0 {
		return nil, errors.New("unit price must be positive")
	}
	if input.Quantity <= 0 {
		input.Quantity = 1
	}
	if input.RecurringDay != nil && (*input.RecurringDay < 1 || *input.RecurringDay > 31) {
		return nil, errors.New("recurring day must be between 1 and 31")
	}

	template := &models.ExpenseTemplate{
		ID:           uuid.New(),
		UserID:       userID,
		CategoryID:   input.CategoryID,
		ItemName:     input.ItemName,
		UnitPrice:    input.UnitPrice,
		Quantity:     input.Quantity,
		RecurringDay: input.RecurringDay,
		Notes:        input.Notes,
	}

	if err := s.templateRepo.Create(template); err != nil {
		return nil, err
	}

	return s.templateRepo.GetByID(template.ID)
}

func (s *ExpenseTemplateService) GetByID(id uuid.UUID) (*models.ExpenseTemplate, error) {
	return s.templateRepo.GetByID(id)
}

func (s *ExpenseTemplateService) GetByUserID(userID uuid.UUID) ([]models.ExpenseTemplate, error) {
	return s.templateRepo.GetByUserID(userID)
}

func (s *ExpenseTemplateService) Update(id uuid.UUID, input CreateExpenseTemplateInput) (*models.ExpenseTemplate, error) {
	template, err := s.templateRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	template.CategoryID = input.CategoryID
	template.ItemName = input.ItemName
	template.UnitPrice = input.UnitPrice
	template.Quantity = input.Quantity
	template.RecurringDay = input.RecurringDay
	template.Notes = input.Notes

	if err := s.templateRepo.Update(template); err != nil {
		return nil, err
	}

	return s.templateRepo.GetByID(template.ID)
}

func (s *ExpenseTemplateService) Delete(id uuid.UUID) error {
	return s.templateRepo.Delete(id)
}

func (s *ExpenseTemplateService) CreateExpenseFromTemplate(userID uuid.UUID, templateID uuid.UUID, expenseDate *time.Time) (*models.Expense, error) {
	template, err := s.templateRepo.GetByID(templateID)
	if err != nil {
		return nil, err
	}

	if template.UserID != userID {
		return nil, errors.New("template not found")
	}

	// Use ExpenseService.Create() to ensure ledger entry is created
	input := CreateExpenseInput{
		CategoryID:  template.CategoryID,
		ItemName:    template.ItemName,
		UnitPrice:   template.UnitPrice,
		Quantity:    template.Quantity,
		Notes:       template.Notes,
		ExpenseDate: expenseDate,
	}

	return s.expenseService.Create(userID, input)
}
