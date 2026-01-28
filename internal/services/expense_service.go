package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type ExpenseService struct {
	expenseRepo  repository.ExpenseRepository
	categoryRepo repository.CategoryRepository
}

func NewExpenseService(expenseRepo repository.ExpenseRepository, categoryRepo repository.CategoryRepository) *ExpenseService {
	return &ExpenseService{
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
	}
}

type CreateExpenseInput struct {
	CategoryID  uuid.UUID
	ItemName    string
	UnitPrice   int64
	Quantity    int
	Notes       *string
	ExpenseDate *time.Time
}

func (s *ExpenseService) Create(userID uuid.UUID, input CreateExpenseInput) (*models.Expense, error) {
	if input.ItemName == "" {
		return nil, errors.New("item name is required")
	}
	if input.UnitPrice <= 0 {
		return nil, errors.New("unit price must be positive")
	}
	if input.Quantity <= 0 {
		input.Quantity = 1
	}

	expense := &models.Expense{
		ID:          uuid.New(),
		UserID:      userID,
		CategoryID:  input.CategoryID,
		ItemName:    input.ItemName,
		UnitPrice:   input.UnitPrice,
		Quantity:    input.Quantity,
		Notes:       input.Notes,
		ExpenseDate: input.ExpenseDate,
	}

	if err := s.expenseRepo.Create(expense); err != nil {
		return nil, err
	}

	return s.expenseRepo.GetByID(expense.ID)
}

func (s *ExpenseService) GetByID(id uuid.UUID) (*models.Expense, error) {
	return s.expenseRepo.GetByID(id)
}

func (s *ExpenseService) GetByUserID(userID uuid.UUID, filter *repository.ExpenseFilter) ([]models.Expense, error) {
	return s.expenseRepo.GetByUserID(userID, filter)
}

type UpdateExpenseInput struct {
	CategoryID  *uuid.UUID
	ItemName    *string
	UnitPrice   *int64
	Quantity    *int
	Notes       *string
	ExpenseDate *time.Time
}

func (s *ExpenseService) Update(id uuid.UUID, input UpdateExpenseInput) (*models.Expense, error) {
	expense, err := s.expenseRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.CategoryID != nil {
		expense.CategoryID = *input.CategoryID
	}
	if input.ItemName != nil {
		expense.ItemName = *input.ItemName
	}
	if input.UnitPrice != nil {
		expense.UnitPrice = *input.UnitPrice
	}
	if input.Quantity != nil {
		expense.Quantity = *input.Quantity
	}
	if input.Notes != nil {
		expense.Notes = input.Notes
	}
	if input.ExpenseDate != nil {
		expense.ExpenseDate = input.ExpenseDate
	}

	if err := s.expenseRepo.Update(expense); err != nil {
		return nil, err
	}

	return s.expenseRepo.GetByID(expense.ID)
}

func (s *ExpenseService) Delete(id uuid.UUID) error {
	return s.expenseRepo.Delete(id)
}
