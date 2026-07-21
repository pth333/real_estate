package usecase

import (
	"errors"

	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
)

type RealEstateService struct {
	repo         repo.RealEstateRepository
	categoryRepo repo.ICategoryRepository
}

type IRealEstateService interface {
	ListRealEstate(req dto.RealEstateSearchRequest) ([]model.RealEstate, int64, error)
	ListRealEstateByCategory(req dto.RealEstateSearchRequest) ([]model.RealEstate, int64, error)
}

func NewRealEstateService(repo repo.RealEstateRepository, categoryRepo repo.ICategoryRepository) IRealEstateService {
	return &RealEstateService{repo: repo, categoryRepo: categoryRepo}
}

func (s *RealEstateService) ListRealEstate(req dto.RealEstateSearchRequest) ([]model.RealEstate, int64, error) {
	offset := (req.Page - 1) * req.Size
	return s.repo.GetList(req, offset, req.Size)
}

// ListRealEstateByCategory lấy category từ slug, sau đó query BĐS theo category_id
func (s *RealEstateService) ListRealEstateByCategory(req dto.RealEstateSearchRequest) ([]model.RealEstate, int64, error) {
	// 1. Lấy category từ slug

	categoryID, err := s.categoryRepo.GetCategoryIdBySlug(req.Slug)

	if err != nil {
		return nil, 0, errors.New("Category not found")
	}
	return s.repo.GetListByCategoryID(categoryID, req, req.Size)
}
