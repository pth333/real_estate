package usecase

import (
	"errors"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/repo"
)

type IManagerPostUseCase interface {
	GetManagerPostsList(userID uint64, search string, page, size int) ([]dto.ManagerPostListResponse, int64, error)
	DeleteManagerPost(userID uint64, postID uint64) error
}

type managerPostUseCase struct {
	managerRepo repo.ManagerPostRepository
}

func NewManagerPostUseCase(managerRepo repo.ManagerPostRepository) IManagerPostUseCase {
	return &managerPostUseCase{managerRepo: managerRepo}
}

func (u *managerPostUseCase) GetManagerPostsList(userID uint64, search string, page, size int) ([]dto.ManagerPostListResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	items, total, err := u.managerRepo.GetManagerPostsList(userID, search, offset, size)
	if err != nil {
		return nil, 0, err
	}

	dtoItems := make([]dto.ManagerPostListResponse, 0, len(items))
	for _, item := range items {
		thumbnail := ""
		if len(item.Images) > 0 {
			thumbnail = item.Images[0].URL // Sử dụng trường URL ảnh
		}

		// Ánh xạ tên loại BĐS
		realEstateType := "Bất động sản"
		if item.Category != nil {
			realEstateType = item.Category.Name
		}

		// Lấy trạng thái lưu tạm ở cột BalconyDirection
		// statusVal := "pending"
		// if item.BalconyDirection != "" {
		// 	statusVal = item.BalconyDirection
		// }

		dtoItems = append(dtoItems, dto.ManagerPostListResponse{
			ID:        item.ID,
			Title:     item.Title,
			Slug:      item.Slug,
			Thumbnail: thumbnail,
			Type:      realEstateType,
			Price:     item.PriceVND,
			Unit:      "vnd",
			Area:      item.Acreage,
			CreatedAt: item.CreatedAt.Format("02-01-2006"),
		})
	}

	return dtoItems, total, err
}

func (u *managerPostUseCase) DeleteManagerPost(userID uint64, postID uint64) error {
	post, err := u.managerRepo.GetByID(postID)
	if err != nil {
		return err
	}

	if post.UserID == nil || *post.UserID != userID {
		return errors.New("bạn không có quyền xóa bài viết này")
	}

	// Gỡ liên kết ảnh trước khi xóa bài viết
	_ = u.managerRepo.UnlinkImages(postID)

	return u.managerRepo.DeleteManagerPost(postID, userID)
}
