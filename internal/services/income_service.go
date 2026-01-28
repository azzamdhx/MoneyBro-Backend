package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type IncomeService struct {
	incomeRepo         repository.IncomeRepository
	incomeCategoryRepo repository.IncomeCategoryRepository
}

func NewIncomeService(incomeRepo repository.IncomeRepository, incomeCategoryRepo repository.IncomeCategoryRepository) *IncomeService {
	return &IncomeService{
		incomeRepo:         incomeRepo,
		incomeCategoryRepo: incomeCategoryRepo,
	}
}

type CreateIncomeInput struct {
	CategoryID  uuid.UUID
	SourceName  string
	Amount      int64
	IncomeType  models.IncomeType
	IncomeDate  *time.Time
	IsRecurring bool
	Notes       *string
}

type UpdateIncomeInput struct {
	CategoryID  *uuid.UUID
	SourceName  *string
	Amount      *int64
	IncomeType  *models.IncomeType
	IncomeDate  *time.Time
	IsRecurring *bool
	Notes       *string
}

func (s *IncomeService) Create(userID uuid.UUID, input CreateIncomeInput) (*models.Income, error) {
	if input.SourceName == "" {
		return nil, errors.New("source name is required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	_, err := s.incomeCategoryRepo.GetByID(input.CategoryID)
	if err != nil {
		return nil, errors.New("invalid income category")
	}

	incomeDate := time.Now()
	if input.IncomeDate != nil {
		incomeDate = *input.IncomeDate
	}

	income := &models.Income{
		ID:          uuid.New(),
		UserID:      userID,
		CategoryID:  input.CategoryID,
		SourceName:  input.SourceName,
		Amount:      input.Amount,
		IncomeType:  input.IncomeType,
		IncomeDate:  incomeDate,
		IsRecurring: input.IsRecurring,
		Notes:       input.Notes,
	}

	if err := s.incomeRepo.Create(income); err != nil {
		return nil, err
	}

	return s.incomeRepo.GetByID(income.ID)
}

func (s *IncomeService) GetByID(id uuid.UUID) (*models.Income, error) {
	return s.incomeRepo.GetByID(id)
}

func (s *IncomeService) GetByUserID(userID uuid.UUID, filter *repository.IncomeFilter) ([]models.Income, error) {
	return s.incomeRepo.GetByUserID(userID, filter)
}

func (s *IncomeService) Update(id uuid.UUID, input UpdateIncomeInput) (*models.Income, error) {
	income, err := s.incomeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.CategoryID != nil {
		_, err := s.incomeCategoryRepo.GetByID(*input.CategoryID)
		if err != nil {
			return nil, errors.New("invalid income category")
		}
		income.CategoryID = *input.CategoryID
	}

	if input.SourceName != nil {
		if *input.SourceName == "" {
			return nil, errors.New("source name is required")
		}
		income.SourceName = *input.SourceName
	}

	if input.Amount != nil {
		if *input.Amount <= 0 {
			return nil, errors.New("amount must be greater than 0")
		}
		income.Amount = *input.Amount
	}

	if input.IncomeType != nil {
		income.IncomeType = *input.IncomeType
	}

	if input.IncomeDate != nil {
		income.IncomeDate = *input.IncomeDate
	}

	if input.IsRecurring != nil {
		income.IsRecurring = *input.IsRecurring
	}

	if input.Notes != nil {
		income.Notes = input.Notes
	}

	if err := s.incomeRepo.Update(income); err != nil {
		return nil, err
	}

	return s.incomeRepo.GetByID(id)
}

func (s *IncomeService) Delete(id uuid.UUID) error {
	return s.incomeRepo.Delete(id)
}
