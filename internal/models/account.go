package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeAsset     AccountType = "ASSET"
	AccountTypeLiability AccountType = "LIABILITY"
	AccountTypeIncome    AccountType = "INCOME"
	AccountTypeExpense   AccountType = "EXPENSE"
)

type Account struct {
	ID             uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	Name           string      `gorm:"type:varchar(100);not null" json:"name"`
	AccountType    AccountType `gorm:"type:varchar(20);not null" json:"account_type"`
	CurrentBalance int64       `gorm:"not null;default:0" json:"current_balance"`
	IsDefault      bool        `gorm:"not null;default:false" json:"is_default"`
	ReferenceID    *uuid.UUID  `gorm:"type:uuid" json:"reference_id,omitempty"`
	ReferenceType  *string     `gorm:"type:varchar(50)" json:"reference_type,omitempty"`
	CreatedAt      time.Time   `gorm:"default:now()" json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Account) TableName() string {
	return "accounts"
}
