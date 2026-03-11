package usecase

import (
	"real_estate_be/internal/delivery/https/dto"
	model "real_estate_be/internal/models"
)

type DashboardService struct {
	repo mysql.IDashboardRepository
}

type IDashboardService interface {
	Summary(from, to string) (model.DashboardSummary, error)
	ListRealEstate(req dto.RealEstateSearchRequest) ([]model.RealEstate, int64, error)
}

func NewDashboardService(repo mysql.IDashboardRepository) *DashboardService {
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
