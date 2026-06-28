package repo

import (
	"real_estate_be/internal/controller/dto"
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type dashboardRepo struct {
	db *gorm.DB
}

type IDashboardRepository interface {
	GetSummary(startDate, endDate string) (model.DashboardSummary, error)
	GetListRealEstate(
		req dto.RealEstateSearchRequest,
		offset int,
		limit int,
	) ([]model.RealEstate, int64, error)
}

func NewDashboardRepository(db *gorm.DB) IDashboardRepository {
	return &dashboardRepo{db: db}
}

func (r *dashboardRepo) GetSummary(startDate, endDate string) (model.DashboardSummary, error) {
	var result model.DashboardSummary

	err := r.db.
		Table("real_estates").
		Select(`
			COUNT(*) AS total_posts,
			AVG(price_vnd) AS avg_price_m2,
			MAX(price_vnd) AS max_price_m2,
			MIN(price_vnd) AS min_price_m2
		`).
		Where("crawled_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&result).
		Error

	return result, err
}

func (r *dashboardRepo) GetListRealEstate(
	req dto.RealEstateSearchRequest,
	offset int,
	limit int,
) ([]model.RealEstate, int64, error) {
	var (
		items []model.RealEstate
		total int64
	)

	db := r.db.Model(&model.RealEstate{})

	if req.Filter.District != "" {
		db = db.Where("district = ?", req.Filter.District)
	}

	if req.Filter.MinPrice != 0 {
		db = db.Where("price_vnd >= ?", req.Filter.MinPrice)
	}

	if req.Filter.MaxPrice != 0 {
		db = db.Where("price_vnd <= ?", req.Filter.MaxPrice)
	}

	if err := r.db.Model(&model.RealEstate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&items).
		Error

	return items, total, err
}
