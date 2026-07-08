package repo

import (
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

type ICategoryRepository interface {
	GetAll() ([]model.Category, error)
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
