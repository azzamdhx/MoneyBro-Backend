package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type twoFACodeRepository struct {
	db *gorm.DB
}

func NewTwoFACodeRepository(db *gorm.DB) TwoFACodeRepository {
	return &twoFACodeRepository{db: db}
}

func (r *twoFACodeRepository) Create(code *models.TwoFACode) error {
	return r.db.Create(code).Error
}

func (r *twoFACodeRepository) GetValidByUserIDAndCode(userID uuid.UUID, code string) (*models.TwoFACode, error) {
	var twoFACode models.TwoFACode
	err := r.db.Preload("User").
		Where("user_id = ? AND code = ? AND expires_at > ? AND used = false", userID, code, time.Now()).
		First(&twoFACode).Error
	if err != nil {
		return nil, err
	}
	return &twoFACode, nil
}

func (r *twoFACodeRepository) MarkAsUsed(id uuid.UUID) error {
	return r.db.Model(&models.TwoFACode{}).
		Where("id = ?", id).
		Update("used", true).Error
}

func (r *twoFACodeRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.TwoFACode{}).Error
}

func (r *twoFACodeRepository) DeleteExpiredCodes() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.TwoFACode{}).Error
}
