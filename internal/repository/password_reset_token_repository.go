package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *passwordResetTokenRepository) GetByToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Preload("User").First(&resetToken, "token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

func (r *passwordResetTokenRepository) GetValidByToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Preload("User").
		Where("token = ? AND expires_at > ? AND used_at IS NULL", token, time.Now()).
		First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

func (r *passwordResetTokenRepository) MarkAsUsed(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used_at", &now).Error
}

func (r *passwordResetTokenRepository) DeleteExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.PasswordResetToken{}).Error
}

func (r *passwordResetTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.PasswordResetToken{}).Error
}
