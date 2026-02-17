package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type savingsContributionRepository struct {
	db *gorm.DB
}

func NewSavingsContributionRepository(db *gorm.DB) SavingsContributionRepository {
	return &savingsContributionRepository{db: db}
}

func (r *savingsContributionRepository) Create(contribution *models.SavingsContribution) error {
	return r.db.Create(contribution).Error
}

func (r *savingsContributionRepository) GetByID(id uuid.UUID) (*models.SavingsContribution, error) {
	var contribution models.SavingsContribution
	err := r.db.Preload("SavingsGoal").First(&contribution, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contribution, nil
}

func (r *savingsContributionRepository) GetBySavingsGoalID(goalID uuid.UUID) ([]models.SavingsContribution, error) {
	var contributions []models.SavingsContribution
	err := r.db.Where("savings_goal_id = ?", goalID).Order("contribution_date DESC").Find(&contributions).Error
	return contributions, err
}

func (r *savingsContributionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.SavingsContribution{}, "id = ?", id).Error
}
