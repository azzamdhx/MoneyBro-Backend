package models

import (
	"time"

	"github.com/google/uuid"
)

// RecurringIncomeGroup represents a group/collection of recurring income items
type RecurringIncomeGroup struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	RecurringDay *int      `gorm:"type:int" json:"recurring_day,omitempty"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	Notes        *string   `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time `gorm:"default:now()" json:"created_at"`

	User  *User                 `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items []RecurringIncomeItem `gorm:"foreignKey:GroupID" json:"items,omitempty"`
}

func (RecurringIncomeGroup) TableName() string {
	return "recurring_income_groups"
}

func (g *RecurringIncomeGroup) Total() int64 {
	var total int64
	for _, item := range g.Items {
		total += item.Amount
	}
	return total
}

// RecurringIncomeItem represents a single item within a recurring income group
type RecurringIncomeItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	GroupID    uuid.UUID `gorm:"type:uuid;not null" json:"group_id"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	SourceName string    `gorm:"type:varchar(255);not null" json:"source_name"`
	Amount     int64     `gorm:"not null" json:"amount"`
	CreatedAt  time.Time `gorm:"default:now()" json:"created_at"`

	Group    *RecurringIncomeGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Category *IncomeCategory       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (RecurringIncomeItem) TableName() string {
	return "recurring_income_items"
}
