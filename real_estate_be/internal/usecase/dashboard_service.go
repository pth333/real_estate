package usecase

import (
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
)

type DashboardService struct {
	repo repo.IDashboardRepository
}

type IDashboardService interface {
	Summary(from, to string) (model.DashboardSummary, error)
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
