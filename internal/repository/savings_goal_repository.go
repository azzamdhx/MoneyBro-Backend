package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type savingsGoalRepository struct {
	db *gorm.DB
}

func NewSavingsGoalRepository(db *gorm.DB) SavingsGoalRepository {
	return &savingsGoalRepository{db: db}
}

func (r *savingsGoalRepository) Create(goal *models.SavingsGoal) error {
	return r.db.Create(goal).Error
}

func (r *savingsGoalRepository) GetByID(id uuid.UUID) (*models.SavingsGoal, error) {
	var goal models.SavingsGoal
	err := r.db.Preload("Contributions").First(&goal, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (r *savingsGoalRepository) GetByUserID(userID uuid.UUID, status *models.SavingsGoalStatus) ([]models.SavingsGoal, error) {
	var goals []models.SavingsGoal
	query := r.db.Preload("Contributions").Where("user_id = ?", userID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("created_at DESC").Find(&goals).Error
	return goals, err
}

func (r *savingsGoalRepository) GetActiveByUserID(userID uuid.UUID) ([]models.SavingsGoal, error) {
	var goals []models.SavingsGoal
	err := r.db.Preload("Contributions").
		Where("user_id = ? AND status = ?", userID, models.SavingsGoalStatusActive).
		Order("target_date ASC").Find(&goals).Error
	return goals, err
}

func (r *savingsGoalRepository) GetByTargetDateRange(startDate, endDate string, status models.SavingsGoalStatus) ([]models.SavingsGoal, error) {
	var goals []models.SavingsGoal
	err := r.db.Preload("Contributions").Preload("User").
		Where("target_date BETWEEN ? AND ? AND status = ?", startDate, endDate, status).
		Find(&goals).Error
	return goals, err
}

func (r *savingsGoalRepository) Update(goal *models.SavingsGoal) error {
	return r.db.Save(goal).Error
}

func (r *savingsGoalRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.SavingsGoal{}, "id = ?", id).Error
}
