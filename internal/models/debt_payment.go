package models

import (
	"time"

	"github.com/google/uuid"
)

type DebtPayment struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	DebtID        uuid.UUID `gorm:"type:uuid;not null" json:"debt_id"`
	PaymentNumber int       `gorm:"not null" json:"payment_number"`
	Amount        int64     `gorm:"not null" json:"amount"`
	PaidAt        time.Time `gorm:"type:date;not null" json:"paid_at"`
	CreatedAt     time.Time `gorm:"default:now()" json:"created_at"`

	Debt *Debt `gorm:"foreignKey:DebtID" json:"debt,omitempty"`
}

func (DebtPayment) TableName() string {
	return "debt_payments"
}
