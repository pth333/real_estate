package repo

import (
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type viewHistoryRepo struct {
	db *gorm.DB
}

type IViewHistoryRepository interface {
	Create(item *model.ViewHistory) error
	MergeSession(sessionID string, userID uint64) error
}

func NewViewHistoryRepository(db *gorm.DB) IViewHistoryRepository {
	return &viewHistoryRepo{db: db}
}

func (r *viewHistoryRepo) Create(item *model.ViewHistory) error {
	return r.db.Create(item).Error
}

func (r *viewHistoryRepo) MergeSession(sessionID string, userID uint64) error {
	// UPDATE view_history SET user_id = ? WHERE session_id = ? AND user_id IS NULL
	// Sáp nhập lịch sử xem của guest session sang tài khoản thật
	return r.db.Model(&model.ViewHistory{}).
		Where("session_id = ? AND user_id IS NULL", sessionID).
		Update("user_id", userID).Error
}
