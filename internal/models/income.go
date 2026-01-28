package models

import (
	"time"

	"github.com/google/uuid"
)

type IncomeType string

const (
	IncomeTypeSalary     IncomeType = "SALARY"
	IncomeTypeFreelance  IncomeType = "FREELANCE"
	IncomeTypeInvestment IncomeType = "INVESTMENT"
	IncomeTypeGift       IncomeType = "GIFT"
	IncomeTypeBonus      IncomeType = "BONUS"
	IncomeTypeRefund     IncomeType = "REFUND"
	IncomeTypeBusiness   IncomeType = "BUSINESS"
	IncomeTypeOther      IncomeType = "OTHER"
)

type Income struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	CategoryID  uuid.UUID  `gorm:"type:uuid;not null" json:"category_id"`
	SourceName  string     `gorm:"type:varchar(255);not null" json:"source_name"`
	Amount      int64      `gorm:"not null" json:"amount"`
	IncomeType  IncomeType `gorm:"type:varchar(50);not null" json:"income_type"`
	IncomeDate  time.Time  `gorm:"type:date;not null" json:"income_date"`
	IsRecurring bool       `gorm:"not null;default:false" json:"is_recurring"`
	Notes       *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt   time.Time  `gorm:"default:now()" json:"created_at"`

	User     *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category *IncomeCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Income) TableName() string {
	return "incomes"
}
