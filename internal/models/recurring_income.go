package models

import (
	"time"

	"github.com/google/uuid"
)

type RecurringIncome struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	CategoryID   uuid.UUID  `gorm:"type:uuid;not null" json:"category_id"`
	SourceName   string     `gorm:"type:varchar(255);not null" json:"source_name"`
	Amount       int64      `gorm:"not null" json:"amount"`
	IncomeType   IncomeType `gorm:"type:varchar(50);not null" json:"income_type"`
	RecurringDay int        `gorm:"not null" json:"recurring_day"`
	IsActive     bool       `gorm:"not null;default:true" json:"is_active"`
	Notes        *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time  `gorm:"default:now()" json:"created_at"`

	User     *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category *IncomeCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (RecurringIncome) TableName() string {
	return "recurring_incomes"
}
