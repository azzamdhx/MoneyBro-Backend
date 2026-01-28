package services

import (
	"time"

	"github.com/google/uuid"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
	"github.com/azzamdhx/moneybro/backend/internal/utils"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) UpdateProfile(userID uuid.UUID, name, password *string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if name != nil {
		if err := utils.ValidateName(*name); err != nil {
			return nil, err
		}
		user.Name = *name
	}

	if password != nil {
		if err := utils.ValidatePassword(*password); err != nil {
			return nil, err
		}
		hashedPassword, err := utils.HashPassword(*password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hashedPassword
	}

	now := time.Now()
	user.UpdatedAt = &now

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
