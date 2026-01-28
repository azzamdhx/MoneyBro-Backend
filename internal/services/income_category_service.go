package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
)

type IncomeCategoryService struct {
	incomeCategoryRepo repository.IncomeCategoryRepository
}

func NewIncomeCategoryService(incomeCategoryRepo repository.IncomeCategoryRepository) *IncomeCategoryService {
	return &IncomeCategoryService{incomeCategoryRepo: incomeCategoryRepo}
}

func (s *IncomeCategoryService) Create(userID uuid.UUID, name string) (*models.IncomeCategory, error) {
	if name == "" {
		return nil, errors.New("income category name is required")
	}

	category := &models.IncomeCategory{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
	}

	if err := s.incomeCategoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *IncomeCategoryService) GetByID(id uuid.UUID) (*models.IncomeCategory, error) {
	return s.incomeCategoryRepo.GetByID(id)
}

func (s *IncomeCategoryService) GetByUserID(userID uuid.UUID) ([]models.IncomeCategory, error) {
	return s.incomeCategoryRepo.GetByUserID(userID)
}

func (s *IncomeCategoryService) Update(id uuid.UUID, name string) (*models.IncomeCategory, error) {
	category, err := s.incomeCategoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name == "" {
		return nil, errors.New("income category name is required")
	}

	category.Name = name

	if err := s.incomeCategoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *IncomeCategoryService) Delete(id uuid.UUID) error {
	return s.incomeCategoryRepo.Delete(id)
}
