package repo

import (
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type ManagerPostRepository interface {
	GetManagerPostsList(userID uint64, search string, offset, limit int) ([]model.RealEstate, int64, error)
	GetByID(id uint64) (*model.RealEstate, error)
	DeleteManagerPost(id uint64, userID uint64) error
	UnlinkImages(realEstateID uint64) error
}

type managerPostRepo struct {
	db *gorm.DB
}

func NewManagerPostRepository(db *gorm.DB) ManagerPostRepository {
	return &managerPostRepo{db: db}
}

func (r *managerPostRepo) GetManagerPostsList(userID uint64, search string, offset, limit int) ([]model.RealEstate, int64, error) {
	var (
		items []model.RealEstate
		total int64
	)

	query := r.db.Model(&model.RealEstate{}).Where("user_id = ?", userID)

	// if status != "" && status != "all" {
	// 	query = query.Where("balcony_direction = ?", status) // Cột balcony_direction đang được dùng tạm làm status bài đăng
	// }

	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	// Đếm tổng số bài viết của manager
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Lấy danh sách kèm theo ảnh
	err := query.Preload("Images").
		Preload("Category").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&items).Error

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *managerPostRepo) GetByID(id uint64) (*model.RealEstate, error) {
	var item model.RealEstate
	err := r.db.Preload("Images").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *managerPostRepo) DeleteManagerPost(id uint64, userID uint64) error {
	// Chỉ cho phép xóa nếu bài viết đó do chính manager (userID) sở hữu
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.RealEstate{}).Error
}

func (r *managerPostRepo) UnlinkImages(realEstateID uint64) error {
	// Gỡ liên kết ảnh cũ của bài viết này
	return r.db.Model(&model.Image{}).
		Where("real_estate_id = ?", realEstateID).
		Update("real_estate_id", nil).Error
}
