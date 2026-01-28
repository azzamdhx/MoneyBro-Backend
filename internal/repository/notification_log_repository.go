package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
)

type notificationLogRepository struct {
	db *gorm.DB
}

func NewNotificationLogRepository(db *gorm.DB) NotificationLogRepository {
	return &notificationLogRepository{db: db}
}

func (r *notificationLogRepository) Create(log *models.NotificationLog) error {
	return r.db.Create(log).Error
}

func (r *notificationLogRepository) ExistsForToday(userID, referenceID uuid.UUID, notificationType models.NotificationType) (bool, error) {
	today := time.Now().Format("2006-01-02")
	var count int64
	err := r.db.Model(&models.NotificationLog{}).
		Where("user_id = ? AND reference_id = ? AND type = ? AND DATE(sent_at) = ?", userID, referenceID, notificationType, today).
		Count(&count).Error
	return count > 0, err
}
