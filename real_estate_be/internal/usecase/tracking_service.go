package usecase

import (
	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
)

type TrackingService struct {
	searchHistoryRepo repo.ISearchHistoryRepository
}

type ITrackingService interface {
	RecordSearch(req dto.TrackingSearchRequest) error
}

func NewTrackingService(searchHistoryRepo repo.ISearchHistoryRepository) ITrackingService {
	return &TrackingService{searchHistoryRepo: searchHistoryRepo}
}

// RecordSearch lưu lịch sử tìm kiếm vào bảng search_history
func (s *TrackingService) RecordSearch(req dto.TrackingSearchRequest) error {

	var priceMin, priceMax *float64

	if req.Filters.PriceRange != nil {
		priceMin = req.Filters.PriceRange.MinPrice
		priceMax = req.Filters.PriceRange.MaxPrice
	}

	var query, province, ward *string
	if req.Query != "" {
		query = &req.Query
	}
	if req.Filters.Location != nil && req.Filters.Location.City != nil {
		province = req.Filters.Location.City
	}
	if req.Filters.Location != nil && req.Filters.Location.Ward != nil {
		ward = req.Filters.Location.Ward
	}

	history := &model.SearchHistory{
		UserID:   req.UserID,
		Query:    query,
		Province: province,
		Ward:     ward,
		PriceMin: priceMin,
		PriceMax: priceMax,
	}

	return s.searchHistoryRepo.Create(history)
}
