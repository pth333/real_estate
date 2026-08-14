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
	GetCategoryIdBySlug(slug string) (int64, error)
	GetCategoryBySlug(slug string) (*model.Category, error)
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) GetCategoryIdBySlug(slug string) (int64, error) {
	var id int64
	if err := r.db.
		Model(&model.Category{}).
		Select("id").
		Where("slug = ?", slug).
		First(&id).
		Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *CategoryRepository) GetCategoryBySlug(slug string) (*model.Category, error) {
	var category model.Category
	if err := r.db.Where("slug = ?", slug).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
