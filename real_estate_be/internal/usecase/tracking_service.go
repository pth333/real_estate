package usecase

import (
	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
)

type TrackingService struct {
	searchHistoryRepo repo.ISearchHistoryRepository
	viewHistoryRepo   repo.IViewHistoryRepository
}

type ITrackingService interface {
	RecordSearch(req dto.TrackingSearchRequest) error
	RecordView(req dto.TrackingViewRequest) error
	MergeSession(sessionID string, userID uint64) error
}

func NewTrackingService(
	searchHistoryRepo repo.ISearchHistoryRepository,
	viewHistoryRepo repo.IViewHistoryRepository,
) ITrackingService {
	return &TrackingService{
		searchHistoryRepo: searchHistoryRepo,
		viewHistoryRepo:   viewHistoryRepo,
	}
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

// RecordView lưu lịch sử xem chi tiết BĐS kèm rule lọc nhiễu
func (s *TrackingService) RecordView(req dto.TrackingViewRequest) error {
	duration := req.DurationSeconds

	// Rule lọc nhiễu:
	// 1. duration < 5s -> bỏ qua (vào nhầm)
	if duration < 5 {
		return nil
	}

	// 2. duration > 3600s -> cap lại 3600 (để tab rồi đi chỗ khác)
	if duration > 3600 {
		duration = 3600
	}

	var uID *uint64
	if req.UserID > 0 {
		uID = &req.UserID
	}

	var sID *string
	if req.SessionID != "" {
		sID = &req.SessionID
	}

	view := &model.ViewHistory{
		UserID:          uID,
		SessionID:       sID,
		RealEstateID:    req.RealEstateID,
		DurationSeconds: duration,
	}

	return s.viewHistoryRepo.Create(view)
}

// MergeSession thực hiện sáp nhập dữ liệu từ guest session sang user đăng nhập
func (s *TrackingService) MergeSession(sessionID string, userID uint64) error {
	if sessionID == "" || userID == 0 {
		return nil
	}
	return s.viewHistoryRepo.MergeSession(sessionID, userID)
}
