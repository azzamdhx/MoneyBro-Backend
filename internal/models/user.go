package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email             string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash      string     `gorm:"type:varchar(255);not null" json:"-"`
	Name              string     `gorm:"type:varchar(100);not null" json:"name"`
	TwoFAEnabled      bool       `gorm:"default:false" json:"two_fa_enabled"`
	NotifyInstallment bool       `gorm:"default:true" json:"notify_installment"`
	NotifyDebt        bool       `gorm:"default:true" json:"notify_debt"`
	NotifyDaysBefore  int        `gorm:"default:3" json:"notify_days_before"`
	CreatedAt         time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}
