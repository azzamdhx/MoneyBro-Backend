package models

import (
	"time"

	"github.com/google/uuid"
)

type Income struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	CategoryID  uuid.UUID  `gorm:"type:uuid;not null" json:"category_id"`
	SourceName  string     `gorm:"type:varchar(255);not null" json:"source_name"`
	Amount      int64      `gorm:"not null" json:"amount"`
	IncomeDate  time.Time  `gorm:"type:date;not null" json:"income_date"`
	IsRecurring bool       `gorm:"not null;default:false" json:"is_recurring"`
	Notes       *string    `gorm:"type:text" json:"notes,omitempty"`
	PocketID    *uuid.UUID `gorm:"type:uuid" json:"pocket_id,omitempty"`
	CreatedAt   time.Time  `gorm:"default:now()" json:"created_at"`

	User     *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category *IncomeCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Pocket   *Account        `gorm:"foreignKey:PocketID" json:"pocket,omitempty"`
}

func (Income) TableName() string {
	return "incomes"
}
