package models

import (
	"time"

	"github.com/google/uuid"
)

type SavingsContribution struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SavingsGoalID    uuid.UUID `gorm:"type:uuid;not null" json:"savings_goal_id"`
	Amount           int64     `gorm:"not null" json:"amount"`
	ContributionDate time.Time `gorm:"type:date;not null" json:"contribution_date"`
	Notes            *string   `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt        time.Time `gorm:"default:now()" json:"created_at"`

	SavingsGoal *SavingsGoal `gorm:"foreignKey:SavingsGoalID" json:"savings_goal,omitempty"`
}

func (SavingsContribution) TableName() string {
	return "savings_contributions"
}
