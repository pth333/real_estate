package mysql

import (
	"fmt"
	"real_estate_be/internal/delivery/https/dto"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repository"

	"gorm.io/gorm"
)

type dashboardRepo struct {
	db *gorm.DB
}

// đảm bảo implement interface
var _ repository.DashboardRepository = (*dashboardRepo)(nil)

func NewDashboardRepository(db *gorm.DB) repository.DashboardRepository {
	return &dashboardRepo{db: db}
}

func (r *dashboardRepo) Summary(startDate, endDate string) (model.DashboardSummary, error) {

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

func (r *dashboardRepo) ListRealEstate(
	req dto.RealEstateSearchRequest,
	offset int,
	limit int,
) ([]model.RealEstate, int64, error) {

	var (
		items []model.RealEstate
		total int64
	)

	db := r.db.Model(&model.RealEstate{})

	// // filter (nếu có)
	// if req.District != "" {
	// 	db = db.Where("district = ?", req.District)
	// }

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&items).
		Error

	fmt.Println((items))

	return items, total, err
}
