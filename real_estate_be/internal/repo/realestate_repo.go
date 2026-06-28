package repo

import (
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type realEstateRepo struct {
	db *gorm.DB
}

type RealEstateRepository interface {
	Create(item *model.RealEstate) error
	CreateBatch(items []*model.RealEstate) error
}

func NewRealEstateRepository(db *gorm.DB) RealEstateRepository {
	return &realEstateRepo{db: db}
}

func (r *realEstateRepo) Create(item *model.RealEstate) error {
	return r.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "source_url"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"title",
				"price_vnd",
				"acreage",
				"price_per_m2",
				"address",
				"district",
				"city",
				"type_of_real_estate",
				"crawled_at",
				"updated_at",
			}),
		}).
		Create(item).
		Error
}

func (r *realEstateRepo) CreateBatch(items []*model.RealEstate) error {
	return r.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "source_url"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"title",
				"price_vnd",
				"acreage",
				"price_per_m2",
				"address",
				"district",
				"city",
				"type_of_real_estate",
				"crawled_at",
				"updated_at",
			}),
		}).
		CreateInBatches(items, 100).
		Error
}
