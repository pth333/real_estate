package usecase

import (
	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
)

type RealEstateService struct {
	repo         repo.RealEstateRepository
	categoryRepo repo.ICategoryRepository
}

type IRealEstateService interface {
	ListRealEstate(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error)
	ListRealEstateByCategory(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error)
	GetListCity() ([]model.Province, error)
	GetListWard(provinceCode string) ([]model.Ward, error)
	GetListRealEstateTypes() ([]model.Category, error)
}

func NewRealEstateService(repo repo.RealEstateRepository, categoryRepo repo.ICategoryRepository) IRealEstateService {
	return &RealEstateService{repo: repo, categoryRepo: categoryRepo}
}

// mapToResponse chuyển model.RealEstate → dto.RealEstateResponse
func mapSliceToResponse(data []model.RealEstate) []dto.RealEstateResponse {
	result := make([]dto.RealEstateResponse, len(data))
	for i, m := range data {
		result[i] = dto.RealEstateResponse{
			ID:               m.ID,
			Title:            m.Title,
			PriceVND:         m.PriceVND,
			Address:          m.Address,
			District:         m.District,
			City:             m.City,
			Acreage:          m.Acreage,
			PricePerM2:       m.PricePerM2,
			TypeOfRealEstate: m.TypeOfRealEstate,
		}
	}
	return result
}

func (s *RealEstateService) ListRealEstate(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error) {
	limit := req.Size
	if limit < 1 {
		limit = 10
	}
	offset := (req.Page - 1) * limit
	data, total, err := s.repo.GetList(req, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return mapSliceToResponse(data), total, nil
}

// ListRealEstateByCategory lấy category từ slug, sau đó query BĐS theo category_id
func (s *RealEstateService) ListRealEstateByCategory(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error) {
	limit := req.Size
	if limit < 1 {
		limit = 10
	}
	offset := (req.Page - 1) * limit
	data, total, err := s.repo.GetListByCategory(offset, req, limit)
	if err != nil {
		return nil, 0, err
	}
	return mapSliceToResponse(data), total, nil
}

func (s *RealEstateService) GetListCity() ([]model.Province, error) {
	return s.repo.GetListCity()
}

func (s *RealEstateService) GetListWard(provinceCode string) ([]model.Ward, error) {
	return s.repo.GetListWard(provinceCode)
}
func (s *RealEstateService) GetListRealEstateTypes() ([]model.Category, error) {
	return s.repo.GetListRealEstateTypes()
}
