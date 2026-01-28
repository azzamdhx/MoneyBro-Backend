package models

import (
	"time"

	"github.com/google/uuid"
)

type InstallmentStatus string

const (
	InstallmentStatusActive    InstallmentStatus = "ACTIVE"
	InstallmentStatusCompleted InstallmentStatus = "COMPLETED"
)

type Installment struct {
	ID             uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID         `gorm:"type:uuid;not null" json:"user_id"`
	Name           string            `gorm:"type:varchar(255);not null" json:"name"`
	ActualAmount   int64             `gorm:"not null" json:"actual_amount"`
	LoanAmount     int64             `gorm:"not null" json:"loan_amount"`
	MonthlyPayment int64             `gorm:"not null" json:"monthly_payment"`
	Tenor          int               `gorm:"not null" json:"tenor"`
	StartDate      time.Time         `gorm:"type:date;not null" json:"start_date"`
	DueDay         int               `gorm:"not null" json:"due_day"`
	Status         InstallmentStatus `gorm:"type:varchar(20);not null;default:'ACTIVE'" json:"status"`
	Notes          *string           `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt      time.Time         `gorm:"default:now()" json:"created_at"`

	User     *User                  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Payments []InstallmentPayment `gorm:"foreignKey:InstallmentID" json:"payments,omitempty"`
}

func (Installment) TableName() string {
	return "installments"
}

func (i *Installment) InterestAmount() int64 {
	return i.LoanAmount - i.ActualAmount
}

func (i *Installment) InterestPercentage() float64 {
	if i.ActualAmount == 0 {
		return 0
	}
	return float64(i.InterestAmount()) / float64(i.ActualAmount) * 100
}

func (i *Installment) PaidCount() int {
	return len(i.Payments)
}

func (i *Installment) RemainingPayments() int {
	return i.Tenor - i.PaidCount()
}

func (i *Installment) RemainingAmount() int64 {
	return int64(i.RemainingPayments()) * i.MonthlyPayment
}
