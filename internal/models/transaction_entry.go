package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionEntry struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TransactionID uuid.UUID `gorm:"type:uuid;not null" json:"transaction_id"`
	AccountID     uuid.UUID `gorm:"type:uuid;not null" json:"account_id"`
	Debit         int64     `gorm:"not null;default:0" json:"debit"`
	Credit        int64     `gorm:"not null;default:0" json:"credit"`
	CreatedAt     time.Time `gorm:"default:now()" json:"created_at"`

	Transaction *Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
	Account     *Account     `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

func (TransactionEntry) TableName() string {
	return "transaction_entries"
}

func (e *TransactionEntry) Amount() int64 {
	if e.Debit > 0 {
		return e.Debit
	}
	return e.Credit
}

func (e *TransactionEntry) IsDebit() bool {
	return e.Debit > 0
}

func (e *TransactionEntry) IsCredit() bool {
	return e.Credit > 0
}
