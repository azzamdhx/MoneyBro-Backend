package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TransactionDate time.Time `gorm:"type:date;not null" json:"transaction_date"`
	Description     string    `gorm:"type:varchar(255);not null" json:"description"`
	ReferenceID     *uuid.UUID `gorm:"type:uuid" json:"reference_id,omitempty"`
	ReferenceType   *string    `gorm:"type:varchar(50)" json:"reference_type,omitempty"`
	CreatedAt       time.Time  `gorm:"default:now()" json:"created_at"`

	User    *User               `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Entries []TransactionEntry  `gorm:"foreignKey:TransactionID" json:"entries,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func (t *Transaction) TotalDebit() int64 {
	var total int64
	for _, entry := range t.Entries {
		total += entry.Debit
	}
	return total
}

func (t *Transaction) TotalCredit() int64 {
	var total int64
	for _, entry := range t.Entries {
		total += entry.Credit
	}
	return total
}

func (t *Transaction) IsBalanced() bool {
	return t.TotalDebit() == t.TotalCredit()
}
