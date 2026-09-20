package repo

import (
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type INotificationLogRepository interface {
	Create(logItem *model.NotificationLog) error
	ListByDeposit(depositID uint64) ([]model.NotificationLog, error)
}

type notificationLogRepo struct {
	db *gorm.DB
}

func NewNotificationLogRepository(db *gorm.DB) INotificationLogRepository {
	return &notificationLogRepo{db: db}
}

func (r *notificationLogRepo) Create(logItem *model.NotificationLog) error {
	return r.db.Create(logItem).Error
}

func (r *notificationLogRepo) ListByDeposit(depositID uint64) ([]model.NotificationLog, error) {
	var items []model.NotificationLog
	err := r.db.Where("deposit_id = ?", depositID).Order("created_at ASC").Find(&items).Error
	return items, err
}
