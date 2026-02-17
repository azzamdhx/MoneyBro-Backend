package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeInstallmentReminder NotificationType = "INSTALLMENT_REMINDER"
	NotificationTypeDebtReminder        NotificationType = "DEBT_REMINDER"
	NotificationTypeSavingsGoalReminder NotificationType = "SAVINGS_GOAL_REMINDER"
)

type NotificationLog struct {
	ID           uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	Type         NotificationType `gorm:"type:varchar(50);not null" json:"type"`
	ReferenceID  uuid.UUID        `gorm:"type:uuid;not null" json:"reference_id"`
	SentAt       time.Time        `gorm:"not null" json:"sent_at"`
	EmailSubject *string          `gorm:"type:varchar(255)" json:"email_subject,omitempty"`
	CreatedAt    time.Time        `gorm:"default:now()" json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (NotificationLog) TableName() string {
	return "notification_logs"
}
