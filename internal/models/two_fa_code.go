package models

import (
	"time"

	"github.com/google/uuid"
)

type TwoFACode struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Code      string    `gorm:"type:varchar(6);not null" json:"code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`

	User User `gorm:"foreignKey:UserID"`
}

func (TwoFACode) TableName() string {
	return "two_fa_codes"
}
