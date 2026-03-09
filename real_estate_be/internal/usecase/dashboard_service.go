package usecase

import (
	"real_estate_be/internal/delivery/https/dto"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repository"
)

type DashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) *DashboardService {
	return &DashboardService{
		repo: repo,
	}
}

// ===== Summary =====

func (s *DashboardService) Summary(from, to string) (model.DashboardSummary, error) {
	start := from + " 00:00:00"
	end := to + " 23:59:59"

	return s.repo.Summary(start, end)
}

func (s *DashboardService) ListRealEstate(
	req dto.RealEstateSearchRequest,
) ([]model.RealEstate, int64, error) {

	offset := (req.Page - 1) * req.Size
	return s.repo.ListRealEstate(req, offset, req.Size)
}
