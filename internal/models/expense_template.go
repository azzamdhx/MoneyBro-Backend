package models

import (
	"time"

	"github.com/google/uuid"
)

// ExpenseTemplateGroup represents a group/collection of expense template items
type ExpenseTemplateGroup struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	RecurringDay *int      `gorm:"type:int" json:"recurring_day,omitempty"`
	Notes        *string   `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time `gorm:"default:now()" json:"created_at"`

	User  *User                 `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items []ExpenseTemplateItem `gorm:"foreignKey:GroupID" json:"items,omitempty"`
}

func (ExpenseTemplateGroup) TableName() string {
	return "expense_template_groups"
}

func (g *ExpenseTemplateGroup) Total() int64 {
	var total int64
	for _, item := range g.Items {
		total += item.Total()
	}
	return total
}

// ExpenseTemplateItem represents a single item within a template group
type ExpenseTemplateItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	GroupID    uuid.UUID `gorm:"type:uuid;not null" json:"group_id"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	ItemName   string    `gorm:"type:varchar(255);not null" json:"item_name"`
	UnitPrice  int64     `gorm:"not null" json:"unit_price"`
	Quantity   int       `gorm:"not null;default:1" json:"quantity"`
	CreatedAt  time.Time `gorm:"default:now()" json:"created_at"`

	Group    *ExpenseTemplateGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Category *Category             `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (ExpenseTemplateItem) TableName() string {
	return "expense_template_items"
}

func (i *ExpenseTemplateItem) Total() int64 {
	return i.UnitPrice * int64(i.Quantity)
}
