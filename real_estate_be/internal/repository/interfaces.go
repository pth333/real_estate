package repository

import (
	"real_estate_be/internal/delivery/https/dto"
	model "real_estate_be/internal/models"
)

type RealEstateRepository interface {
	Create(item *model.RealEstate) error
	CreateBatch(items []*model.RealEstate) error
}

type DashboardRepository interface {
	Summary(startDate, endDate string) (model.DashboardSummary, error)
	ListRealEstate(
		req dto.RealEstateSearchRequest,
		offset int,
		limit int,
	) ([]model.RealEstate, int64, error)
	// Latest(limit int) ([]model.RealEstate, error)
}

type UserRepository interface {
	Register(item *model.User) error
	FindByEmail(email string) (*model.User, error)
	// FindById(id uint) (*model.User, error)
}
