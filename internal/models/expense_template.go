package models

import (
	"time"

	"github.com/google/uuid"
)

type ExpenseTemplate struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	CategoryID   uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	ItemName     string    `gorm:"type:varchar(255);not null" json:"item_name"`
	UnitPrice    int64     `gorm:"not null" json:"unit_price"`
	Quantity     int       `gorm:"not null;default:1" json:"quantity"`
	RecurringDay *int      `gorm:"type:int" json:"recurring_day,omitempty"`
	Notes        *string   `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time `gorm:"default:now()" json:"created_at"`

	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (ExpenseTemplate) TableName() string {
	return "expense_templates"
}

func (e *ExpenseTemplate) Total() int64 {
	return e.UnitPrice * int64(e.Quantity)
}
