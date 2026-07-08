package usecase

import (
	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
)

type DashboardService struct {
	repo repo.IDashboardRepository
}

type IDashboardService interface {
	// FilterByPrice(minPrice, maxPrice *float64) []model.RealEstate
	Summary(from, to string) (model.DashboardSummary, error)
	ListRealEstate(req dto.RealEstateSearchRequest) ([]model.RealEstate, int64, error)
}

func NewDashboardService(repo repo.IDashboardRepository) IDashboardService {
	return &DashboardService{
		repo: repo,
	}
}

func (s *DashboardService) Summary(from, to string) (model.DashboardSummary, error) {
	start := from + " 00:00:00"
	end := to + " 23:59:59"

	return s.repo.GetSummary(start, end)
}

func (s *DashboardService) ListRealEstate(
	req dto.RealEstateSearchRequest,
) ([]model.RealEstate, int64, error) {

	offset := (req.Page - 1) * req.Size

	return s.repo.GetListRealEstate(req, offset, req.Size)
}

// func (s *DashboardService) FilterByPrice(minPrice, maxPrice *float64) []model.RealEstate {
// 	return s.repo.FilterByPrice(minPrice, maxPrice)
// }
