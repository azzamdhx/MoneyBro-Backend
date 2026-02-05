package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type RecurringIncomeService struct {
	recurringIncomeRepo repository.RecurringIncomeRepository
	incomeService       *IncomeService
	incomeCategoryRepo  repository.IncomeCategoryRepository
}

func NewRecurringIncomeService(
	recurringIncomeRepo repository.RecurringIncomeRepository,
	incomeService *IncomeService,
	incomeCategoryRepo repository.IncomeCategoryRepository,
) *RecurringIncomeService {
	return &RecurringIncomeService{
		recurringIncomeRepo: recurringIncomeRepo,
		incomeService:       incomeService,
		incomeCategoryRepo:  incomeCategoryRepo,
	}
}

type CreateRecurringIncomeInput struct {
	CategoryID   uuid.UUID
	SourceName   string
	Amount       int64
	IncomeType   models.IncomeType
	RecurringDay int
	Notes        *string
}

type UpdateRecurringIncomeInput struct {
	CategoryID   *uuid.UUID
	SourceName   *string
	Amount       *int64
	IncomeType   *models.IncomeType
	RecurringDay *int
	IsActive     *bool
	Notes        *string
}

func (s *RecurringIncomeService) Create(userID uuid.UUID, input CreateRecurringIncomeInput) (*models.RecurringIncome, error) {
	if input.SourceName == "" {
		return nil, errors.New("source name is required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if input.RecurringDay < 1 || input.RecurringDay > 31 {
		return nil, errors.New("recurring day must be between 1 and 31")
	}

	_, err := s.incomeCategoryRepo.GetByID(input.CategoryID)
	if err != nil {
		return nil, errors.New("invalid income category")
	}

	recurringIncome := &models.RecurringIncome{
		ID:           uuid.New(),
		UserID:       userID,
		CategoryID:   input.CategoryID,
		SourceName:   input.SourceName,
		Amount:       input.Amount,
		IncomeType:   input.IncomeType,
		RecurringDay: input.RecurringDay,
		IsActive:     true,
		Notes:        input.Notes,
	}

	if err := s.recurringIncomeRepo.Create(recurringIncome); err != nil {
		return nil, err
	}

	return s.recurringIncomeRepo.GetByID(recurringIncome.ID)
}

func (s *RecurringIncomeService) GetByID(id uuid.UUID) (*models.RecurringIncome, error) {
	return s.recurringIncomeRepo.GetByID(id)
}

func (s *RecurringIncomeService) GetByUserID(userID uuid.UUID, isActive *bool) ([]models.RecurringIncome, error) {
	return s.recurringIncomeRepo.GetByUserID(userID, isActive)
}

func (s *RecurringIncomeService) Update(id uuid.UUID, input UpdateRecurringIncomeInput) (*models.RecurringIncome, error) {
	recurringIncome, err := s.recurringIncomeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.CategoryID != nil {
		_, err := s.incomeCategoryRepo.GetByID(*input.CategoryID)
		if err != nil {
			return nil, errors.New("invalid income category")
		}
		recurringIncome.CategoryID = *input.CategoryID
	}

	if input.SourceName != nil {
		if *input.SourceName == "" {
			return nil, errors.New("source name is required")
		}
		recurringIncome.SourceName = *input.SourceName
	}

	if input.Amount != nil {
		if *input.Amount <= 0 {
			return nil, errors.New("amount must be greater than 0")
		}
		recurringIncome.Amount = *input.Amount
	}

	if input.IncomeType != nil {
		recurringIncome.IncomeType = *input.IncomeType
	}

	if input.RecurringDay != nil {
		if *input.RecurringDay < 1 || *input.RecurringDay > 31 {
			return nil, errors.New("recurring day must be between 1 and 31")
		}
		recurringIncome.RecurringDay = *input.RecurringDay
	}

	if input.IsActive != nil {
		recurringIncome.IsActive = *input.IsActive
	}

	if input.Notes != nil {
		recurringIncome.Notes = input.Notes
	}

	if err := s.recurringIncomeRepo.Update(recurringIncome); err != nil {
		return nil, err
	}

	return s.recurringIncomeRepo.GetByID(id)
}

func (s *RecurringIncomeService) Delete(id uuid.UUID) error {
	return s.recurringIncomeRepo.Delete(id)
}

func (s *RecurringIncomeService) CreateIncomeFromRecurring(userID uuid.UUID, recurringID uuid.UUID, incomeDate *time.Time) (*models.Income, error) {
	recurring, err := s.recurringIncomeRepo.GetByID(recurringID)
	if err != nil {
		return nil, err
	}

	if recurring.UserID != userID {
		return nil, errors.New("recurring income not found")
	}

	// Use IncomeService.Create() to ensure ledger entry is created
	input := CreateIncomeInput{
		CategoryID:  recurring.CategoryID,
		SourceName:  recurring.SourceName,
		Amount:      recurring.Amount,
		IncomeType:  recurring.IncomeType,
		IncomeDate:  incomeDate,
		IsRecurring: true,
		Notes:       recurring.Notes,
	}

	return s.incomeService.Create(userID, input)
}
