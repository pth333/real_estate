package repo

import (
	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type ImageRepository interface {
	Create(image *model.Image) error
	FindByKey(key string) (*model.Image, error)
	// LinkToRealEstate gán real_estate_id cho danh sách ảnh đã upload
	LinkToRealEstate(imageIDs []uint64, realEstateID uint64) error
	GetImagesByRealEstateID(id uint64) ([]dto.ImageResponse, error)
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

func (r *imageRepository) LinkToRealEstate(imageIDs []uint64, realEstateID uint64) error {
	if len(imageIDs) == 0 {
		return nil
	}
	return r.db.Model(&model.Image{}).
		Where("id IN ?", imageIDs).
		Update("real_estate_id", realEstateID).
		Error
}

func (r *imageRepository) FindByKey(key string) (*model.Image, error) {
	var image model.Image
	err := r.db.Where("`key` = ?", key).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *imageRepository) GetImagesByRealEstateID(id uint64) ([]dto.ImageResponse, error) {
	var items []dto.ImageResponse
	err := r.db.Model(&model.Image{}).
		Where("real_estate_id = ?", id).
		Order("id").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
