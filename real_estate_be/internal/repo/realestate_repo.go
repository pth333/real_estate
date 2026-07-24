package repo

import (
	"real_estate_be/internal/dto"
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
	GetList(req dto.RealEstateSearchRequest, offset, limit int) ([]model.RealEstate, int64, error)
	GetListByCategory(offset int, req dto.RealEstateSearchRequest, limit int) ([]model.RealEstate, int64, error)
	GetListCity() ([]model.Province, error)
	GetListWard(provinceCode string) ([]model.Ward, error)
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

func (r *realEstateRepo) GetList(req dto.RealEstateSearchRequest, offset, limit int) ([]model.RealEstate, int64, error) {
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

// GetListByCategory query BĐS theo category_id với filter và offset-based pagination
func (r *realEstateRepo) GetListByCategory(offset int, req dto.RealEstateSearchRequest, limit int) ([]model.RealEstate, int64, error) {
	var (
		items []model.RealEstate
		total int64
	)
	db := r.db.Model(&model.RealEstate{}).
		Joins("JOIN categories ON categories.id = real_estates.category_id").
		Where("categories.slug = ?", req.Slug)
	if req.Filter.District != "" {
		db = db.Where("district = ?", req.Filter.District)
	}
	if req.Filter.MinPrice != 0 {
		db = db.Where("price_vnd >= ?", req.Filter.MinPrice)
	}
	if req.Filter.MaxPrice != 0 {
		db = db.Where("price_vnd <= ?", req.Filter.MaxPrice)
	}

	if req.Search != "" {
		db = db.Where(
			"MATCH(title, address) AGAINST(? IN BOOLEAN MODE)",
			req.Search,
		)
	}

	if err := db.Count(&total).Error; err != nil {
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

func (r *realEstateRepo) GetListCity() ([]model.Province, error) {
	var provinces []model.Province
	result := r.db.Select("code, name").Find(&provinces)
	if result.Error != nil {
		return nil, result.Error
	}
	return provinces, nil
}

func (r *realEstateRepo) GetListWard(provinceCode string) ([]model.Ward, error) {
	var wards []model.Ward
	result := r.db.Where("province_code = ?", provinceCode).
		Select("code, name").
		Find(&wards)
	if result.Error != nil {
		return nil, result.Error
	}
	return wards, nil
}
