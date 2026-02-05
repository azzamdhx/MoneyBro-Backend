package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type ExpenseTemplateGroupService struct {
	groupRepo      repository.ExpenseTemplateGroupRepository
	expenseService *ExpenseService
	categoryRepo   repository.CategoryRepository
}

func NewExpenseTemplateGroupService(
	groupRepo repository.ExpenseTemplateGroupRepository,
	expenseService *ExpenseService,
	categoryRepo repository.CategoryRepository,
) *ExpenseTemplateGroupService {
	return &ExpenseTemplateGroupService{
		groupRepo:      groupRepo,
		expenseService: expenseService,
		categoryRepo:   categoryRepo,
	}
}

type CreateExpenseTemplateGroupInput struct {
	Name         string
	RecurringDay *int
	Notes        *string
	Items        []CreateExpenseTemplateItemInput
}

type CreateExpenseTemplateItemInput struct {
	CategoryID uuid.UUID
	ItemName   string
	UnitPrice  int64
	Quantity   int
}

func (s *ExpenseTemplateGroupService) Create(userID uuid.UUID, input CreateExpenseTemplateGroupInput) (*models.ExpenseTemplateGroup, error) {
	if input.Name == "" {
		return nil, errors.New("group name is required")
	}
	if input.RecurringDay != nil && (*input.RecurringDay < 1 || *input.RecurringDay > 31) {
		return nil, errors.New("recurring day must be between 1 and 31")
	}
	if len(input.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	group := &models.ExpenseTemplateGroup{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         input.Name,
		RecurringDay: input.RecurringDay,
		Notes:        input.Notes,
	}

	if err := s.groupRepo.Create(group); err != nil {
		return nil, err
	}

	// Add items
	for _, itemInput := range input.Items {
		if itemInput.ItemName == "" {
			continue
		}
		if itemInput.UnitPrice <= 0 {
			continue
		}
		quantity := itemInput.Quantity
		if quantity <= 0 {
			quantity = 1
		}

		item := &models.ExpenseTemplateItem{
			ID:         uuid.New(),
			GroupID:    group.ID,
			CategoryID: itemInput.CategoryID,
			ItemName:   itemInput.ItemName,
			UnitPrice:  itemInput.UnitPrice,
			Quantity:   quantity,
		}

		if err := s.groupRepo.AddItem(item); err != nil {
			return nil, err
		}
	}

	return s.groupRepo.GetByID(group.ID)
}

func (s *ExpenseTemplateGroupService) GetByID(id uuid.UUID) (*models.ExpenseTemplateGroup, error) {
	return s.groupRepo.GetByID(id)
}

func (s *ExpenseTemplateGroupService) GetByUserID(userID uuid.UUID) ([]models.ExpenseTemplateGroup, error) {
	return s.groupRepo.GetByUserID(userID)
}

type UpdateExpenseTemplateGroupInput struct {
	Name         *string
	RecurringDay *int
	Notes        *string
}

func (s *ExpenseTemplateGroupService) Update(id uuid.UUID, input UpdateExpenseTemplateGroupInput) (*models.ExpenseTemplateGroup, error) {
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
	if input.Notes != nil {
		group.Notes = input.Notes
	}

	if err := s.groupRepo.Update(group); err != nil {
		return nil, err
	}

	return s.groupRepo.GetByID(group.ID)
}

func (s *ExpenseTemplateGroupService) Delete(id uuid.UUID) error {
	return s.groupRepo.Delete(id)
}

// Item operations
func (s *ExpenseTemplateGroupService) AddItem(groupID uuid.UUID, input CreateExpenseTemplateItemInput) (*models.ExpenseTemplateGroup, error) {
	if input.ItemName == "" {
		return nil, errors.New("item name is required")
	}
	if input.UnitPrice <= 0 {
		return nil, errors.New("unit price must be positive")
	}
	quantity := input.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	item := &models.ExpenseTemplateItem{
		ID:         uuid.New(),
		GroupID:    groupID,
		CategoryID: input.CategoryID,
		ItemName:   input.ItemName,
		UnitPrice:  input.UnitPrice,
		Quantity:   quantity,
	}

	if err := s.groupRepo.AddItem(item); err != nil {
		return nil, err
	}

	return s.groupRepo.GetByID(groupID)
}

type UpdateExpenseTemplateItemInput struct {
	CategoryID *uuid.UUID
	ItemName   *string
	UnitPrice  *int64
	Quantity   *int
}

func (s *ExpenseTemplateGroupService) UpdateItem(itemID uuid.UUID, input UpdateExpenseTemplateItemInput) (*models.ExpenseTemplateItem, error) {
	item, err := s.groupRepo.GetItemByID(itemID)
	if err != nil {
		return nil, err
	}

	if input.CategoryID != nil {
		item.CategoryID = *input.CategoryID
	}
	if input.ItemName != nil {
		item.ItemName = *input.ItemName
	}
	if input.UnitPrice != nil {
		item.UnitPrice = *input.UnitPrice
	}
	if input.Quantity != nil {
		item.Quantity = *input.Quantity
	}

	if err := s.groupRepo.UpdateItem(item); err != nil {
		return nil, err
	}

	return s.groupRepo.GetItemByID(itemID)
}

func (s *ExpenseTemplateGroupService) DeleteItem(itemID uuid.UUID) error {
	return s.groupRepo.DeleteItem(itemID)
}

// CreateExpensesFromGroup creates expenses from all items in a template group
func (s *ExpenseTemplateGroupService) CreateExpensesFromGroup(userID uuid.UUID, groupID uuid.UUID, expenseDate *time.Time) ([]models.Expense, error) {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return nil, err
	}

	if group.UserID != userID {
		return nil, errors.New("template group not found")
	}

	var expenses []models.Expense
	for _, item := range group.Items {
		input := CreateExpenseInput{
			CategoryID:  item.CategoryID,
			ItemName:    item.ItemName,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
			Notes:       group.Notes,
			ExpenseDate: expenseDate,
		}

		expense, err := s.expenseService.Create(userID, input)
		if err != nil {
			return expenses, err
		}
		expenses = append(expenses, *expense)
	}

	return expenses, nil
}
