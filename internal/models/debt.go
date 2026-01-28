package models

import (
	"time"

	"github.com/google/uuid"
)

type DebtStatus string
type DebtPaymentType string

const (
	DebtStatusActive    DebtStatus = "ACTIVE"
	DebtStatusCompleted DebtStatus = "COMPLETED"

	DebtPaymentTypeOneTime     DebtPaymentType = "ONE_TIME"
	DebtPaymentTypeInstallment DebtPaymentType = "INSTALLMENT"
)

type Debt struct {
	ID             uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	PersonName     string          `gorm:"type:varchar(255);not null" json:"person_name"`
	ActualAmount   int64           `gorm:"not null" json:"actual_amount"`
	LoanAmount     *int64          `json:"loan_amount,omitempty"`
	PaymentType    DebtPaymentType `gorm:"type:varchar(20);not null" json:"payment_type"`
	MonthlyPayment *int64          `json:"monthly_payment,omitempty"`
	Tenor          *int            `json:"tenor,omitempty"`
	DueDate        *time.Time      `gorm:"type:date" json:"due_date,omitempty"`
	Status         DebtStatus      `gorm:"type:varchar(20);not null;default:'ACTIVE'" json:"status"`
	Notes          *string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt      time.Time       `gorm:"default:now()" json:"created_at"`

	User     *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Payments []DebtPayment `gorm:"foreignKey:DebtID" json:"payments,omitempty"`
}

func (Debt) TableName() string {
	return "debts"
}

func (d *Debt) InterestAmount() *int64 {
	if d.LoanAmount == nil {
		return nil
	}
	amount := *d.LoanAmount - d.ActualAmount
	return &amount
}

func (d *Debt) InterestPercentage() *float64 {
	if d.LoanAmount == nil || d.ActualAmount == 0 {
		return nil
	}
	interest := *d.LoanAmount - d.ActualAmount
	percentage := float64(interest) / float64(d.ActualAmount) * 100
	return &percentage
}

func (d *Debt) TotalToPay() int64 {
	if d.LoanAmount != nil {
		return *d.LoanAmount
	}
	return d.ActualAmount
}

func (d *Debt) PaidAmount() int64 {
	var total int64
	for _, p := range d.Payments {
		total += p.Amount
	}
	return total
}

func (d *Debt) RemainingAmount() int64 {
	return d.TotalToPay() - d.PaidAmount()
}
