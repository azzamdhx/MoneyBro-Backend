package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type CategoryService struct {
	categoryRepo   repository.CategoryRepository
	accountService *AccountService
}

func NewCategoryService(categoryRepo repository.CategoryRepository, accountService *AccountService) *CategoryService {
	return &CategoryService{
		categoryRepo:   categoryRepo,
		accountService: accountService,
	}
}

func (s *CategoryService) Create(userID uuid.UUID, name string) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("category name is required")
	}

	category := &models.Category{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	// Create linked EXPENSE account for this category
	if _, err := s.accountService.CreateLinkedAccount(userID, name, models.AccountTypeExpense, category.ID, "category"); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) GetByID(id uuid.UUID) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *CategoryService) GetByUserID(userID uuid.UUID) ([]models.Category, error) {
	return s.categoryRepo.GetByUserID(userID)
}

func (s *CategoryService) Update(id uuid.UUID, name string) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name == "" {
		return nil, errors.New("category name is required")
	}

	category.Name = name

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) Delete(id uuid.UUID) error {
	return s.categoryRepo.Delete(id)
}
