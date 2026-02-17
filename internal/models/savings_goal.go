package models

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type SavingsGoalStatus string

const (
	SavingsGoalStatusActive    SavingsGoalStatus = "ACTIVE"
	SavingsGoalStatusCompleted SavingsGoalStatus = "COMPLETED"
	SavingsGoalStatusCancelled SavingsGoalStatus = "CANCELLED"
)

type SavingsGoal struct {
	ID            uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID        uuid.UUID         `gorm:"type:uuid;not null" json:"user_id"`
	Name          string            `gorm:"type:varchar(255);not null" json:"name"`
	TargetAmount  int64             `gorm:"not null" json:"target_amount"`
	CurrentAmount int64             `gorm:"not null;default:0" json:"current_amount"`
	TargetDate    time.Time         `gorm:"type:date;not null" json:"target_date"`
	Icon          *string           `gorm:"type:varchar(50)" json:"icon,omitempty"`
	Status        SavingsGoalStatus `gorm:"type:varchar(20);not null;default:'ACTIVE'" json:"status"`
	Notes         *string           `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt     time.Time         `gorm:"default:now()" json:"created_at"`

	User          *User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Contributions []SavingsContribution `gorm:"foreignKey:SavingsGoalID" json:"contributions,omitempty"`
}

func (SavingsGoal) TableName() string {
	return "savings_goals"
}

func (s *SavingsGoal) Progress() float64 {
	if s.TargetAmount == 0 {
		return 0
	}
	return float64(s.CurrentAmount) / float64(s.TargetAmount) * 100
}

func (s *SavingsGoal) RemainingAmount() int64 {
	remaining := s.TargetAmount - s.CurrentAmount
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (s *SavingsGoal) MonthlyTarget() int64 {
	now := time.Now()
	if s.TargetDate.Before(now) {
		return s.RemainingAmount()
	}

	months := monthsBetween(now, s.TargetDate)
	if months <= 0 {
		return s.RemainingAmount()
	}

	return int64(math.Ceil(float64(s.RemainingAmount()) / float64(months)))
}

func monthsBetween(start, end time.Time) int {
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	total := years*12 + months

	if end.Day() > start.Day() {
		total++
	}

	if total < 0 {
		return 0
	}
	return total
}
