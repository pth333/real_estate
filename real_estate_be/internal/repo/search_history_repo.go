package repo

import (
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type searchHistoryRepo struct {
	db *gorm.DB
}

type ISearchHistoryRepository interface {
	Create(item *model.SearchHistory) error
	MergeSession(sessionID string, userID uint64) error
}

func NewSearchHistoryRepository(db *gorm.DB) ISearchHistoryRepository {
	return &searchHistoryRepo{db: db}
}

func (r *searchHistoryRepo) Create(item *model.SearchHistory) error {
	return r.db.Create(item).Error
}

func (r *searchHistoryRepo) MergeSession(sessionID string, userID uint64) error {
	return r.db.Model(&model.SearchHistory{}).
		Where("session_id = ? AND user_id IS NULL", sessionID).
		Update("user_id", userID).Error
}
