package repo

import (
	"real_estate_be/internal/models"

	"gorm.io/gorm"
)

type dashboardRepo struct {
	db *gorm.DB
}

type IDashboardRepository interface {
	GetSummary(startDate, endDate string) (model.DashboardSummary, error)
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
