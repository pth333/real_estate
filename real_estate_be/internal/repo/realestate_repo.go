package repo

import (
	"errors"
	"strings"

	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type realEstateRepo struct {
	db *gorm.DB
}

// toResponse chuyển dto.RealEstateResponse (scan từ deferred join, có ImageURLs) → Images slice
func toResponse(m *dto.RealEstateResponse) {
	if m.ImageURLs != "" {
		m.Images = strings.Split(m.ImageURLs, "|")
	}
}

type RealEstateRepository interface {
	Create(item *model.RealEstate) error
	// Save cập nhật 1 bản ghi (dùng để gắn slug sau khi có ID).
	Save(item *model.RealEstate) error
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
	// Lấy khoảng giá/diện tích theo slug SEO (filter_ranges). Không parse chuỗi.
	GetFilterRangeBySlug(slug string) (*model.FilterRange, error)
	// Lấy toàn bộ filter_ranges (giá + diện tích) để FE dựng menu lọc theo slug chuẩn.
	GetFilterRanges() ([]model.FilterRange, error)
	// Lấy 1 tin đăng theo ID (trang chi tiết -rs), kèm gom ảnh.
	GetByID(id uint64) (*dto.RealEstateResponse, error)
	GetCategory() ([]model.Category, error)
	GetProvinceBySlug(city string) (string, error)
}

func NewRealEstateRepository(db *gorm.DB) RealEstateRepository {
	return &realEstateRepo{db: db}
}

func (r *realEstateRepo) Create(item *model.RealEstate) error {
	return r.db.Create(item).Error
}

func (r *realEstateRepo) Save(item *model.RealEstate) error {
	return r.db.Save(item).Error
}

func (r *realEstateRepo) CreateBatch(items []*model.RealEstate) error {
	return r.db.CreateInBatches(items, 100).Error
}

func (r *realEstateRepo) GetList(req dto.RealEstateSearchRequest, offset, limit int) ([]dto.RealEstateResponse, int64, error) {
	var (
		items []dto.RealEstateResponse
		total int64
	)

	// Điều kiện lọc chung (dùng cho COUNT và query)
	where := "WHERE 1=1"
	args := []interface{}{}

	if req.Filter.District != "" {
		where += " AND re.district = ?"
		args = append(args, req.Filter.District)
	}
	if req.Filter.City != "" {
		where += " AND re.city = ?"
		args = append(args, req.Filter.City)
	}
	if req.Filter.MinPrice != 0 {
		where += " AND re.price_vnd >= ?"
		args = append(args, req.Filter.MinPrice)
	}
	if req.Filter.MaxPrice != 0 {
		where += " AND re.price_vnd <= ?"
		args = append(args, req.Filter.MaxPrice)
	}
	if req.Filter.MinAcreage != 0 {
		where += " AND re.acreage >= ?"
		args = append(args, req.Filter.MinAcreage)
	}
	if req.Filter.MaxAcreage != 0 {
		where += " AND re.acreage <= ?"
		args = append(args, req.Filter.MaxAcreage)
	}

	// Đếm tổng số bản ghi (không cần join)
	if err := r.db.Raw(
		"SELECT COUNT(*) FROM real_estates re "+where,
		args...,
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// deferred join: subquery lấy id trang hiện tại, rồi LEFT JOIN users + gom ảnh
	listArgs := append(append([]interface{}{}, args...), limit, offset)
	err := r.db.Raw(
		"SELECT re.id, re.title, re.slug, re.price_vnd, re.address, re.district, re.city, "+
			"re.acreage, re.price_per_m2, re.bedrooms, re.bathrooms, re.description, re.created_at, "+
			"COALESCE(GROUP_CONCAT(DISTINCT img.url ORDER BY img.id SEPARATOR '|'), '') AS image_urls, "+
			"COALESCE(u.name, '') AS agent_name, "+
			"COALESCE(u.phone, '') AS agent_phone "+
			"FROM real_estates re "+
			"LEFT JOIN images img ON img.real_estate_id = re.id "+
			"LEFT JOIN users u ON u.id = re.user_id "+
			"JOIN (SELECT re.id FROM real_estates re "+where+
			" ORDER BY re.created_at DESC, re.id DESC LIMIT ? OFFSET ?) tmp ON tmp.id = re.id "+
			"GROUP BY re.id "+
			"ORDER BY re.created_at DESC, re.id DESC",
		listArgs...,
	).Debug().Scan(&items).Error
	if err != nil {
		return nil, 0, err
	}

	for i := range items {
		toResponse(&items[i])
	}
	return items, total, nil
}

func (r *realEstateRepo) GetListByCategory(offset int, req dto.RealEstateSearchRequest, limit int) ([]dto.RealEstateResponse, int64, error) {
	var (
		items []dto.RealEstateResponse
		total int64
	)

	// điều kiện lọc chung (dùng cho COUNT và query con)
	where := "WHERE c.slug = ?"
	args := []interface{}{req.Slug}

	if req.Filter.District != "" {
		where += " AND re.district = ?"
		args = append(args, req.Filter.District)
	}
	if req.Filter.City != "" {
		where += " AND re.city = ?"
		args = append(args, req.Filter.City)
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

	// deferred join: subquery lấy id trang hiện tại, rồi LEFT JOIN users + gom ảnh
	listArgs := append(append([]interface{}{}, args...), limit, offset)
	err := r.db.Raw(
		"SELECT re.id, re.title, re.slug, re.price_vnd, re.address, re.district, re.city, "+
			"re.acreage, re.price_per_m2, re.bedrooms, re.bathrooms, re.description, re.created_at, "+
			"COALESCE(GROUP_CONCAT(DISTINCT img.url ORDER BY img.id SEPARATOR '|'), '') AS image_urls, "+
			"COALESCE(u.name, '') AS agent_name, "+
			"COALESCE(u.phone, '') AS agent_phone "+
			"FROM real_estates re "+
			"LEFT JOIN images img ON img.real_estate_id = re.id "+
			"LEFT JOIN users u ON u.id = re.user_id "+
			"JOIN (SELECT re.id FROM real_estates re JOIN categories c ON c.id = re.category_id "+
			where+" ORDER BY re.created_at DESC, re.id DESC LIMIT ? OFFSET ?) tmp ON tmp.id = re.id "+
			"GROUP BY re.id "+
			"ORDER BY re.created_at DESC, re.id DESC",
		listArgs...,
	).Scan(&items).Error
	if err != nil {
		return nil, 0, err
	}

	for i := range items {
		toResponse(&items[i])
	}
	return items, total, nil
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

// GetFilterRanges trả toàn bộ filter_ranges (type price + area), dùng để FE
// dựng menu lọc theo đúng slug chuẩn trong DB (không build slug thủ công).
func (r *realEstateRepo) GetFilterRanges() ([]model.FilterRange, error) {
	var ranges []model.FilterRange
	if err := r.db.Order("type ASC, min_val ASC").Find(&ranges).Error; err != nil {
		return nil, err
	}
	return ranges, nil
}

// GetByID lấy 1 tin đăng theo ID (trang chi tiết slug -rs{id}), dùng deferred
// join tương tự GetListByCategory để kèm agent name + gom ảnh.
func (r *realEstateRepo) GetByID(id uint64) (*dto.RealEstateResponse, error) {
	var item dto.RealEstateResponse
	err := r.db.Raw(
		"SELECT re.id, re.title, re.slug, re.price_vnd, re.address, re.district, re.city, "+
			"re.acreage, re.price_per_m2, re.bedrooms, re.bathrooms, re.description, re.created_at, "+
			"COALESCE(GROUP_CONCAT(DISTINCT img.url ORDER BY img.id SEPARATOR '|'), '') AS image_urls, "+
			"COALESCE(u.name, '') AS agent_name, "+
			"COALESCE(u.phone, '') AS agent_phone "+
			"FROM real_estates re "+
			"LEFT JOIN images img ON img.real_estate_id = re.id "+
			"LEFT JOIN users u ON u.id = re.user_id "+
			"WHERE re.id = ? "+
			"GROUP BY re.id",
		id,
	).Scan(&item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, nil
	}
	toResponse(&item)
	return &item, nil
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

func (r *realEstateRepo) GetCategory() ([]model.Category, error) {
	var categories []model.Category

	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *realEstateRepo) GetProvinceBySlug(city string) (string, error) {
	var name string

	if err := r.db.Debug().
		Model(&model.Province{}).
		Where("code_name = ?", city).
		Pluck("name", &name).Error; err != nil {
		return "", err
	}

	return name, nil
}

// GetFilterRangeBySlug tìm khoảng giá/diện tích theo slug SEO từ bảng filter_ranges.
func (r *realEstateRepo) GetFilterRangeBySlug(slug string) (*model.FilterRange, error) {
	if slug == "" {
		return nil, nil
	}
	var f model.FilterRange
	if err := r.db.Where("slug = ?", slug).First(&f).Error; err != nil {
		// Không tìm thấy → trả nil, nil (để parser bỏ qua segment)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}
