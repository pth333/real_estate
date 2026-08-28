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

	// ── Ảnh dự án (image_projects) ──
	CreateProjectImage(image *model.ImageProject) error
	// LinkToProject gán project_id cho danh sách ảnh dự án đã upload
	LinkToProject(imageIDs []uint64, projectID uint64) error
	GetImagesByProjectID(projectID uint64) ([]dto.ImageResponse, error)
	// GetProjectCoverImages trả ảnh đại diện (ảnh đầu tiên) cho từng project
	GetProjectCoverImages(projectIDs []uint64) (map[uint64]string, error)
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

// ── Ảnh dự án (image_projects) ──

func (r *imageRepository) CreateProjectImage(image *model.ImageProject) error {
	return r.db.Create(image).Error
}

func (r *imageRepository) LinkToProject(imageIDs []uint64, projectID uint64) error {
	if len(imageIDs) == 0 {
		return nil
	}
	return r.db.Model(&model.ImageProject{}).
		Where("id IN ?", imageIDs).
		Update("project_id", projectID).
		Error
}

func (r *imageRepository) GetImagesByProjectID(projectID uint64) ([]dto.ImageResponse, error) {
	var items []dto.ImageResponse
	err := r.db.Model(&model.ImageProject{}).
		Where("project_id = ?", projectID).
		Order("id").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetProjectCoverImages lấy ảnh đại diện (ảnh đầu tiên theo id) cho danh sách
// dự án — 1 query duy nhất để tránh N+1.
func (r *imageRepository) GetProjectCoverImages(projectIDs []uint64) (map[uint64]string, error) {
	result := make(map[uint64]string)
	if len(projectIDs) == 0 {
		return result, nil
	}

	var items []struct {
		ProjectID uint64
		URL       string
	}
	err := r.db.Model(&model.ImageProject{}).
		Select("project_id, url").
		Where("project_id IN ?", projectIDs).
		Order("id").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}

	// Ảnh có id nhỏ nhất = ảnh đại diện (chỉ lấy 1 lần mỗi project)
	for _, it := range items {
		if _, ok := result[it.ProjectID]; !ok {
			result[it.ProjectID] = it.URL
		}
	}
	return result, nil
}
