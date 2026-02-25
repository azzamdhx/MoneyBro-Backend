package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type RecurringIncomeGroupService struct {
	groupRepo          repository.RecurringIncomeGroupRepository
	incomeService      *IncomeService
	incomeCategoryRepo repository.IncomeCategoryRepository
}

func NewRecurringIncomeGroupService(
	groupRepo repository.RecurringIncomeGroupRepository,
	incomeService *IncomeService,
	incomeCategoryRepo repository.IncomeCategoryRepository,
) *RecurringIncomeGroupService {
	return &RecurringIncomeGroupService{
		groupRepo:          groupRepo,
		incomeService:      incomeService,
		incomeCategoryRepo: incomeCategoryRepo,
	}
}

type CreateRecurringIncomeGroupInput struct {
	Name         string
	RecurringDay *int
	IsActive     *bool
	Notes        *string
	Items        []CreateRecurringIncomeItemInput
}

type CreateRecurringIncomeItemInput struct {
	CategoryID uuid.UUID
	SourceName string
	Amount     int64
}

func (s *RecurringIncomeGroupService) Create(userID uuid.UUID, input CreateRecurringIncomeGroupInput) (*models.RecurringIncomeGroup, error) {
	if input.Name == "" {
		return nil, errors.New("group name is required")
	}
	if input.RecurringDay != nil && (*input.RecurringDay < 1 || *input.RecurringDay > 31) {
		return nil, errors.New("recurring day must be between 1 and 31")
	}
	if len(input.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	group := &models.RecurringIncomeGroup{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         input.Name,
		RecurringDay: input.RecurringDay,
		IsActive:     isActive,
		Notes:        input.Notes,
	}

	if err := s.groupRepo.Create(group); err != nil {
		return nil, err
	}

	// Add items
	for _, itemInput := range input.Items {
		if itemInput.SourceName == "" {
			continue
		}
		if itemInput.Amount <= 0 {
			continue
		}

		item := &models.RecurringIncomeItem{
			ID:         uuid.New(),
			GroupID:    group.ID,
			CategoryID: itemInput.CategoryID,
			SourceName: itemInput.SourceName,
			Amount:     itemInput.Amount,
		}

		if err := s.groupRepo.AddItem(item); err != nil {
			return nil, err
		}
	}

	return s.groupRepo.GetByID(group.ID)
}

func (s *RecurringIncomeGroupService) GetByID(id uuid.UUID) (*models.RecurringIncomeGroup, error) {
	return s.groupRepo.GetByID(id)
}

func (s *RecurringIncomeGroupService) GetByUserID(userID uuid.UUID, isActive *bool) ([]models.RecurringIncomeGroup, error) {
	return s.groupRepo.GetByUserID(userID, isActive)
}

type UpdateRecurringIncomeGroupInput struct {
	Name         *string
	RecurringDay *int
	IsActive     *bool
	Notes        *string
}

func (s *RecurringIncomeGroupService) Update(id uuid.UUID, input UpdateRecurringIncomeGroupInput) (*models.RecurringIncomeGroup, error) {
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		group.Name = *input.Name
	}
	if input.RecurringDay != nil {
		group.RecurringDay = input.RecurringDay
	}
	if input.IsActive != nil {
		group.IsActive = *input.IsActive
	}
	if input.Notes != nil {
		group.Notes = input.Notes
	}

	if err := s.groupRepo.Update(group); err != nil {
		return nil, err
	}

	return s.groupRepo.GetByID(group.ID)
}

func (s *RecurringIncomeGroupService) Delete(id uuid.UUID) error {
	return s.groupRepo.Delete(id)
}

// Item operations
func (s *RecurringIncomeGroupService) AddItem(groupID uuid.UUID, input CreateRecurringIncomeItemInput) (*models.RecurringIncomeGroup, error) {
	if input.SourceName == "" {
		return nil, errors.New("source name is required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	item := &models.RecurringIncomeItem{
		ID:         uuid.New(),
		GroupID:    groupID,
		CategoryID: input.CategoryID,
		SourceName: input.SourceName,
		Amount:     input.Amount,
	}

	if err := s.groupRepo.AddItem(item); err != nil {
		return nil, err
	}

	return s.groupRepo.GetByID(groupID)
}

type UpdateRecurringIncomeItemInput struct {
	CategoryID *uuid.UUID
	SourceName *string
	Amount     *int64
}

func (s *RecurringIncomeGroupService) UpdateItem(itemID uuid.UUID, input UpdateRecurringIncomeItemInput) (*models.RecurringIncomeItem, error) {
	item, err := s.groupRepo.GetItemByID(itemID)
	if err != nil {
		return nil, err
	}

	if input.CategoryID != nil {
		item.CategoryID = *input.CategoryID
	}
	if input.SourceName != nil {
		item.SourceName = *input.SourceName
	}
	if input.Amount != nil {
		item.Amount = *input.Amount
	}

	if err := s.groupRepo.UpdateItem(item); err != nil {
		return nil, err
	}

	return s.groupRepo.GetItemByID(itemID)
}

func (s *RecurringIncomeGroupService) DeleteItem(itemID uuid.UUID) error {
	return s.groupRepo.DeleteItem(itemID)
}

// CreateIncomesFromGroup creates incomes from all items in a recurring income group
func (s *RecurringIncomeGroupService) CreateIncomesFromGroup(userID uuid.UUID, groupID uuid.UUID, incomeDate *time.Time) ([]models.Income, error) {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return nil, err
	}

	if group.UserID != userID {
		return nil, errors.New("recurring income group not found")
	}

	var incomes []models.Income
	for _, item := range group.Items {
		input := CreateIncomeInput{
			CategoryID:  item.CategoryID,
			SourceName:  item.SourceName,
			Amount:      item.Amount,
			IncomeDate:  incomeDate,
			IsRecurring: true,
			Notes:       group.Notes,
		}

		income, err := s.incomeService.Create(userID, input)
		if err != nil {
			return incomes, err
		}
		incomes = append(incomes, *income)
	}

	return incomes, nil
}
