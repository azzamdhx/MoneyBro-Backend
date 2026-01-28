package models

import (
	"time"

	"github.com/google/uuid"
)

type PasswordResetToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `gorm:"default:now()" json:"created_at"`

	User User `gorm:"foreignKey:UserID"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
