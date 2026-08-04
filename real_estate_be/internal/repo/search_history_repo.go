package repo

import (
	"real_estate_be/internal/models"

	"gorm.io/gorm"
)

type searchHistoryRepo struct {
	db *gorm.DB
}

type ISearchHistoryRepository interface {
	Create(item *model.SearchHistory) error
}

func NewSearchHistoryRepository(db *gorm.DB) ISearchHistoryRepository {
	return &searchHistoryRepo{db: db}
}

func (r *searchHistoryRepo) Create(item *model.SearchHistory) error {
	return r.db.Create(item).Error
}
