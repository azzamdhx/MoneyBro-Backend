package models

import (
	"time"

	"github.com/google/uuid"
)

type InstallmentPayment struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	InstallmentID uuid.UUID `gorm:"type:uuid;not null" json:"installment_id"`
	PaymentNumber int       `gorm:"not null" json:"payment_number"`
	Amount        int64     `gorm:"not null" json:"amount"`
	PaidAt        time.Time `gorm:"type:date;not null" json:"paid_at"`
	CreatedAt     time.Time `gorm:"default:now()" json:"created_at"`

	Installment *Installment `gorm:"foreignKey:InstallmentID" json:"installment,omitempty"`
}

func (InstallmentPayment) TableName() string {
	return "installment_payments"
}
