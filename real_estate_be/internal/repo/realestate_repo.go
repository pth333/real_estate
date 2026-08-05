package repo

import (
	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type realEstateRepo struct {
	db *gorm.DB
}

// toResponse chuyển model.RealEstate (đã Preload User + Images) → dto.RealEstateResponse
func toResponse(m model.RealEstate) dto.RealEstateResponse {
	images := make([]string, 0, len(m.Images))
	for _, img := range m.Images {
		images = append(images, img.URL)
	}

	agentName, agentPhone := "", ""
	if m.User != nil {
		agentName = m.User.Name
		agentPhone = m.User.Phone
	}

	return dto.RealEstateResponse{
		ID:          m.ID,
		Title:       m.Title,
		PriceVND:    m.PriceVND,
		Address:     m.Address,
		District:    m.District,
		City:        m.City,
		Acreage:     m.Acreage,
		PricePerM2:  m.PricePerM2,
		Images:      images,
		Bedrooms:    m.Bedrooms,
		Bathrooms:   m.Bathrooms,
		Description: m.Description,
		AgentName:   agentName,
		AgentPhone:  agentPhone,
		CreatedAt:   m.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

type RealEstateRepository interface {
	Create(item *model.RealEstate) error
	CreateBatch(items []*model.RealEstate) error
	GetList(req dto.RealEstateSearchRequest, offset, limit int) ([]dto.RealEstateResponse, int64, error)
	GetListByCategory(offset int, req dto.RealEstateSearchRequest, limit int) ([]dto.RealEstateResponse, int64, error)
	GetListCity() ([]model.Province, error)
	GetListWard(provinceCode string) ([]model.Ward, error)
	GetListRealEstateTypes() ([]model.Category, error)
	// Lấy tên tỉnh/thành từ code (VD "79" → "Hồ Chí Minh")
	GetProvinceNameByCode(code string) (string, error)
	// Lấy tên phường/xã từ code (VD "76049" → "Phường 12")
	GetWardNameByCode(code string) (string, error)
	// Lấy N thành phố có nhiều BĐS nhất (theo cột city)
	GetTopCityByCount(limit int) ([]model.CityStat, error)
}

func NewRealEstateRepository(db *gorm.DB) RealEstateRepository {
	return &realEstateRepo{db: db}
}

func (r *realEstateRepo) Create(item *model.RealEstate) error {
	return r.db.Create(item).Error
}

func (r *realEstateRepo) CreateBatch(items []*model.RealEstate) error {
	return r.db.CreateInBatches(items, 100).Error
}

func (r *realEstateRepo) GetList(req dto.RealEstateSearchRequest, offset, limit int) ([]dto.RealEstateResponse, int64, error) {
	var (
		items []model.RealEstate
		total int64
	)

	// Điều kiện lọc chung (dùng cho COUNT và query)
	where := "WHERE 1=1"
	args := []interface{}{}

	if req.Filter.District != "" {
		where += " AND re.district = ?"
		args = append(args, req.Filter.District)
	}
	if req.Filter.MinPrice != 0 {
		where += " AND re.price_vnd >= ?"
		args = append(args, req.Filter.MinPrice)
	}
	if req.Filter.MaxPrice != 0 {
		where += " AND re.price_vnd <= ?"
		args = append(args, req.Filter.MaxPrice)
	}

	// Đếm tổng số bản ghi (không cần join)
	if err := r.db.Raw(
		"SELECT COUNT(*) FROM real_estates re "+where,
		args...,
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// Lấy danh sách dùng GORM + Preload (User + Images)
	db := r.db.Model(&model.RealEstate{}).
		Preload("User").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("images.id ASC")
		})

	if req.Filter.District != "" {
		db = db.Where("district = ?", req.Filter.District)
	}
	if req.Filter.MinPrice != 0 {
		db = db.Where("price_vnd >= ?", req.Filter.MinPrice)
	}
	if req.Filter.MaxPrice != 0 {
		db = db.Where("price_vnd <= ?", req.Filter.MaxPrice)
	}

	err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.RealEstateResponse, len(items))
	for i, item := range items {
		result[i] = toResponse(item)
	}
	return result, total, nil
}

func (r *realEstateRepo) GetListByCategory(offset int, req dto.RealEstateSearchRequest, limit int) ([]dto.RealEstateResponse, int64, error) {
	var (
		items []model.RealEstate
		total int64
	)

	// điều kiện lọc chung (dùng cho COUNT và query con)
	where := "WHERE c.slug = ?"
	args := []interface{}{req.Slug}

	if req.Filter.District != "" {
		where += " AND re.district = ?"
		args = append(args, req.Filter.District)
	}
	if req.Filter.MinPrice != 0 {
		where += " AND re.price_vnd >= ?"
		args = append(args, req.Filter.MinPrice)
	}
	if req.Filter.MaxPrice != 0 {
		where += " AND re.price_vnd <= ?"
		args = append(args, req.Filter.MaxPrice)
	}
	if req.Search != "" {
		where += " AND MATCH(re.title, re.address) AGAINST(? IN BOOLEAN MODE)"
		args = append(args, req.Search)
	}

	// đếm tổng số bản ghi (không cần join ảnh/user)
	if err := r.db.Raw(
		"SELECT COUNT(*) FROM real_estates re JOIN categories c ON c.id = re.category_id "+where,
		args...,
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// Lấy category id từ slug, sau đó GORM + Preload (User + Images)
	var categoryID int64
	if err := r.db.Model(&model.Category{}).Select("id").Where("slug = ?", req.Slug).First(&categoryID).Error; err != nil {
		return nil, 0, err
	}

	db := r.db.Model(&model.RealEstate{}).
		Preload("User").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("images.id ASC")
		}).
		Where("category_id = ?", categoryID)

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
		// Full-text search — khớp với COUNT query ở trên
		db = db.Where("MATCH(title, address) AGAINST(? IN BOOLEAN MODE)", req.Search)
	}

	err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.RealEstateResponse, len(items))
	for i, item := range items {
		result[i] = toResponse(item)
	}
	return result, total, nil
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

func (r *realEstateRepo) GetListRealEstateTypes() ([]model.Category, error) {
	var types []model.Category
	result := r.db.Select("id, name").Find(&types)
	if result.Error != nil {
		return nil, result.Error
	}
	return types, nil
}

func (r *realEstateRepo) GetProvinceNameByCode(code string) (string, error) {
	var name string
	if err := r.db.Model(&model.Province{}).Select("name").Where("code = ?", code).First(&name).Error; err != nil {
		return "", err
	}
	return name, nil
}

func (r *realEstateRepo) GetWardNameByCode(code string) (string, error) {
	var name string
	if err := r.db.Model(&model.Ward{}).Select("name").Where("code = ?", code).First(&name).Error; err != nil {
		return "", err
	}
	return name, nil
}

func (r *realEstateRepo) GetTopCityByCount(limit int) ([]model.CityStat, error) {
	var cities []model.CityStat
	// Subquery đếm top N thành phố trước, rồi LEFT JOIN provinces để lấy ảnh
	// → chỉ join N dòng thay vì join toàn bộ real_estates (tối ưu hơn)
	result := r.db.Table("(?) AS t", r.db.Model(&model.RealEstate{}).
		Select("city, COUNT(*) AS total").
		Where("city IS NOT NULL AND city <> ''").
		Group("city").
		Order("total DESC").
		Limit(limit)).
		Select("t.city, t.total, COALESCE(p.image, '') AS image").
		Joins("LEFT JOIN provinces p ON p.name = t.city").
		Scan(&cities)
	if result.Error != nil {
		return nil, result.Error
	}
	return cities, nil
}
