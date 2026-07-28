package repo

import (
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type ImageRepository interface {
	Create(image *model.Image) error
	FindByKey(key string) (*model.Image, error)
}

type imageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{db: db}
}

func (r *imageRepository) Create(image *model.Image) error {
	return r.db.Create(image).Error
}

func (r *imageRepository) FindByKey(key string) (*model.Image, error) {
	var image model.Image
	err := r.db.Where("`key` = ?", key).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}
