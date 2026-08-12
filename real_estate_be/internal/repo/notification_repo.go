package repo

import (
	model "real_estate_be/internal/models"
	"gorm.io/gorm"
)

type INotificationRepository interface {
	Create(notification *model.Notification) error
	GetLatest(limit int) ([]model.Notification, error)
}

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) INotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *model.Notification) error {
	// Sử dụng Clause ON CONFLICT (MySQL: ON DUPLICATE KEY UPDATE) để tránh lỗi khi trùng listing_id
	return r.db.Save(notification).Error
}

func (r *NotificationRepository) GetLatest(limit int) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.Order("created_at DESC").Limit(limit).Find(&notifications).Error
	return notifications, err
}
