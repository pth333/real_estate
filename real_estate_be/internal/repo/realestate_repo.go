package repo

import (
	"errors"
	"strconv"

	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type realEstateRepo struct {
	db *gorm.DB
}

type RealEstateRepository interface {
	Create(item *model.RealEstate) error
	// Save cập nhật 1 bản ghi (dùng để gắn slug sau khi có ID).
	Save(item *model.RealEstate) error
	CreateBatch(items []*model.RealEstate) error
	GetList(req dto.RealEstateSearchRequest, offset, limit int) ([]dto.RealEstateResponse, int64, error)
	GetListByCategory(offset int, req dto.RealEstateSearchRequest, limit int) ([]dto.RealEstateResponse, int64, error)
	GetListCity() ([]model.Province, error)
	GetListProject(provinceCode, wardCode string) ([]model.RealEstateProject, error)
	GetProjectsByCategoryID(categoryID int64) ([]model.RealEstateProject, error)
	GetListWard(provinceCode string) ([]model.Ward, error)
	GetListRealEstateTypes() ([]model.Category, error)
	// Lấy tên tỉnh/thành từ code (VD "79" → "Hồ Chí Minh")
	// GetProvinceNameByCode(name string) (string, error)
	// Lấy tên phường/xã từ code (VD "76049" → "Phường 12")
	// GetWardNameByCode(anm string) (string, error)
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
	IncrementProjectView(id uint64) error
	GetFeaturedProjects(limit int) ([]model.RealEstateProject, error)
	GetProjectByID(id uint64) (*model.RealEstateProject, error)
	// Gợi ý
	GetTrending(limit int) ([]dto.RealEstateResponse, error)
	GetByIDs(ids []uint64) ([]dto.RealEstateResponse, error)
	GetRecommendationsBasic(userID uint64, sessionID string, limit int) ([]dto.RealEstateResponse, error)
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

	// Điều kiện lọc chung (used cho COUNT và query)
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

	// Lọc tìm kiếm theo từ khóa
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		where += " AND (re.title LIKE ? OR re.city LIKE ? OR re.district LIKE ? OR re.address LIKE ?)"
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// ── Bộ lọc nâng cao ──

	if req.Filter.Bedrooms != nil && *req.Filter.Bedrooms > 0 {
		where += " AND re.bedrooms = ?"
		args = append(args, *req.Filter.Bedrooms)
	}
	if req.Filter.Bathrooms != nil && *req.Filter.Bathrooms > 0 {
		where += " AND re.bathrooms = ?"
		args = append(args, *req.Filter.Bathrooms)
	}
	if req.Filter.HouseDirection != "" {
		where += " AND re.house_direction = ?"
		args = append(args, req.Filter.HouseDirection)
	}
	if req.Filter.BalconyDirection != "" {
		where += " AND re.balcony_direction = ?"
		args = append(args, req.Filter.BalconyDirection)
	}
	if req.Filter.LegalDocs != "" {
		where += " AND re.legal_docs = ?"
		args = append(args, req.Filter.LegalDocs)
	}
	if req.Filter.Interior != "" {
		where += " AND re.interior = ?"
		args = append(args, req.Filter.Interior)
	}

	// Đếm tổng số bản ghi (không cần join)
	if err := r.db.Raw(
		"SELECT COUNT(*) FROM real_estates re "+where,
		args...,
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// deferred join: subquery lấy id trang hiện tại, rồi LEFT JOIN + gom ảnh
	listArgs := append(append([]interface{}{}, args...), limit, offset)
	err := r.db.Raw(
		"SELECT re.id, re.title, re.slug, re.price_vnd, re.address, re.district, re.city, "+
			"re.acreage, re.price_per_m2, re.bedrooms, re.bathrooms, re.description, re.created_at, "+
			"re.house_direction, re.balcony_direction, re.floors, re.legal_docs, re.interior, "+
			"re.price_electricity, re.price_water, re.price_internet, re.amenities, re.latitude, re.longitude, "+
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
	).Scan(&items).Error
	if err != nil {
		return nil, 0, err
	}

	// for i := range items {
	// 	toResponse(&items[i])
	// }
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
	if req.Filter.MinAcreage != 0 {
		where += " AND re.acreage >= ?"
		args = append(args, req.Filter.MinAcreage)
	}
	if req.Filter.MaxAcreage != 0 {
		where += " AND re.acreage <= ?"
		args = append(args, req.Filter.MaxAcreage)
	}
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		where += " AND (re.title LIKE ? OR re.city LIKE ? OR re.district LIKE ? OR re.address LIKE ?)"
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// ── Bộ lọc nâng cao ──

	if req.Filter.Bedrooms != nil && *req.Filter.Bedrooms > 0 {
		where += " AND re.bedrooms = ?"
		args = append(args, *req.Filter.Bedrooms)
	}
	if req.Filter.Bathrooms != nil && *req.Filter.Bathrooms > 0 {
		where += " AND re.bathrooms = ?"
		args = append(args, *req.Filter.Bathrooms)
	}
	if req.Filter.HouseDirection != "" {
		where += " AND re.house_direction = ?"
		args = append(args, req.Filter.HouseDirection)
	}
	if req.Filter.BalconyDirection != "" {
		where += " AND re.balcony_direction = ?"
		args = append(args, req.Filter.BalconyDirection)
	}
	if req.Filter.LegalDocs != "" {
		where += " AND re.legal_docs = ?"
		args = append(args, req.Filter.LegalDocs)
	}
	if req.Filter.Interior != "" {
		where += " AND re.interior = ?"
		args = append(args, req.Filter.Interior)
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
			"re.house_direction, re.balcony_direction, re.floors, re.legal_docs, re.interior, "+
			"re.price_electricity, re.price_water, re.price_internet, re.amenities, re.latitude, re.longitude, "+
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
	).Debug().Scan(&items).Error
	if err != nil {
		return nil, 0, err
	}

	// for i := range items {
	// 	toResponse(&items[i])
	// }
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

func (r *realEstateRepo) GetListProject(provinceCode, wardCode string) ([]model.RealEstateProject, error) {
	var projects []model.RealEstateProject
	query := r.db.Select("id, name").Model(&model.RealEstateProject{})

	// Lọc theo mã tỉnh/thành (VD "79" → province_id = 79)
	if provinceID, err := strconv.ParseUint(provinceCode, 10, 64); err == nil && provinceCode != "" {
		query = query.Where("province_id = ?", provinceID)
	}

	// Lọc theo mã phường/xã (VD "27184" → ward_id = 27184)
	if wardID, err := strconv.ParseUint(wardCode, 10, 64); err == nil && wardCode != "" {
		query = query.Where("ward_id = ?", wardID)
	}

	if err := query.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *realEstateRepo) GetProjectsByCategoryID(categoryID int64) ([]model.RealEstateProject, error) {
	var projects []model.RealEstateProject
	// Lấy trực tiếp dự án có category_id khớp với danh mục dự án được chọn trên menu
	err := r.db.Where("category_id = ?", categoryID).Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
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
	err := r.db.
		Table("real_estates re").
		Select(`re.*,
        CONCAT(COALESCE(c.id, ''), '-', COALESCE(c.name, '')) AS real_estate_type,
        COALESCE(u.name, '') AS agent_name,
        COALESCE(u.phone, '') AS agent_phone,
        COALESCE(u.email, '') AS agent_email`).
		Joins("LEFT JOIN categories c ON c.id = re.category_id").
		Joins("LEFT JOIN users u ON u.id = re.user_id").
		Where("re.id = ?", id).
		Scan(&item).Error
	if err != nil {
		return nil, err
	}

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

	if err := r.db.Select("slug").Find(&categories).Error; err != nil {
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

func (r *realEstateRepo) IncrementProjectView(id uint64) error {
	return r.db.Model(&model.RealEstateProject{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *realEstateRepo) GetFeaturedProjects(limit int) ([]model.RealEstateProject, error) {
	var projects []model.RealEstateProject
	err := r.db.Model(&model.RealEstateProject{}).
		Order("view_count DESC, id DESC").
		Limit(limit).
		Find(&projects).Error
	return projects, err
}

func (r *realEstateRepo) GetProjectByID(id uint64) (*model.RealEstateProject, error) {
	var project model.RealEstateProject
	err := r.db.Model(&model.RealEstateProject{}).Where("id = ?", id).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetTrending lấy danh sách BĐS nổi bật (nhiều lượt xem nhất hoặc mới nhất)
func (r *realEstateRepo) GetTrending(limit int) ([]dto.RealEstateResponse, error) {
	var items []dto.RealEstateResponse

	// Query lấy top real_estate_id được xem nhiều nhất từ view_history
	var trendingIDs []uint64
	r.db.Model(&model.ViewHistory{}).
		Select("real_estate_id, SUM(duration_seconds) as total_duration").
		Group("real_estate_id").
		Order("total_duration DESC, real_estate_id DESC").
		Limit(limit).
		Pluck("real_estate_id", &trendingIDs)

	if len(trendingIDs) > 0 {
		items, err := r.GetByIDs(trendingIDs)
		if err == nil && len(items) > 0 {
			return items, nil
		}
	}

	// Fallback nếu chưa có lượt xem nào: Lấy danh sách tin mới đăng nhất
	err := r.db.Raw(
		"SELECT re.id, re.title, re.slug, re.price_vnd, re.address, re.district, re.city, "+
			"re.acreage, re.price_per_m2, re.bedrooms, re.bathrooms, re.description, re.created_at, "+
			"re.house_direction, re.balcony_direction, re.floors, re.legal_docs, re.interior, "+
			"re.price_electricity, re.price_water, re.price_internet, re.amenities, re.latitude, re.longitude, "+
			"COALESCE(GROUP_CONCAT(DISTINCT img.url ORDER BY img.id SEPARATOR '|'), '') AS image_urls, "+
			"COALESCE(u.name, '') AS agent_name, "+
			"COALESCE(u.phone, '') AS agent_phone "+
			"FROM real_estates re "+
			"LEFT JOIN images img ON img.real_estate_id = re.id "+
			"LEFT JOIN users u ON u.id = re.user_id "+
			"GROUP BY re.id "+
			"ORDER BY re.created_at DESC, re.id DESC "+
			"LIMIT ?",
		limit,
	).Scan(&items).Error

	if err != nil {
		return nil, err
	}

	// for i := range items {
	// 	toResponse(&items[i])
	// }
	return items, nil
}

// GetByIDs lấy thông tin chi tiết của danh sách ID BĐS (dành cho Cache Hit)
func (r *realEstateRepo) GetByIDs(ids []uint64) ([]dto.RealEstateResponse, error) {
	if len(ids) == 0 {
		return []dto.RealEstateResponse{}, nil
	}

	var items []dto.RealEstateResponse
	err := r.db.Raw(
		"SELECT re.id, re.title, re.slug, re.price_vnd, re.address, re.district, re.city, "+
			"re.acreage, re.price_per_m2, re.bedrooms, re.bathrooms, re.description, re.created_at, "+
			"re.house_direction, re.balcony_direction, re.floors, re.legal_docs, re.interior, "+
			"re.price_electricity, re.price_water, re.price_internet, re.amenities, re.latitude, re.longitude, "+
			"COALESCE(GROUP_CONCAT(DISTINCT img.url ORDER BY img.id SEPARATOR '|'), '') AS image_urls, "+
			"COALESCE(u.name, '') AS agent_name, "+
			"COALESCE(u.phone, '') AS agent_phone "+
			"FROM real_estates re "+
			"LEFT JOIN images img ON img.real_estate_id = re.id "+
			"LEFT JOIN users u ON u.id = re.user_id "+
			"WHERE re.id IN (?) "+
			"GROUP BY re.id",
		ids,
	).Scan(&items).Error

	if err != nil {
		return nil, err
	}

	// Sắp xếp lại theo đúng thứ tự của mảng ids truyền vào (giữ tính chất ranking)
	idMap := make(map[uint64]dto.RealEstateResponse)
	// for _, item := range items {
	// 	toResponse(&item)
	// 	idMap[item.ID] = item
	// }

	sortedItems := make([]dto.RealEstateResponse, 0, len(items))
	for _, id := range ids {
		if item, exists := idMap[id]; exists {
			sortedItems = append(sortedItems, item)
		}
	}

	return sortedItems, nil
}

// GetRecommendationsBasic gợi ý cơ bản dựa trên lịch sử xem gần nhất (khi chưa có Python ML)
func (r *realEstateRepo) GetRecommendationsBasic(userID uint64, sessionID string, limit int) ([]dto.RealEstateResponse, error) {
	var items []dto.RealEstateResponse

	// Lấy tối đa 3 tin đã xem gần nhất từ view_history để làm gợi ý đầu vào
	var watchedIDs []uint64
	query := r.db.Model(&model.ViewHistory{}).
		Order("viewed_at DESC, id DESC").
		Limit(3)

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	} else if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	} else {
		// Trả về Trending nếu không có bất kỳ thông tin nào
		return r.GetTrending(limit)
	}

	query.Pluck("real_estate_id", &watchedIDs)

	if len(watchedIDs) == 0 {
		// Chưa xem tin nào -> Cold Start -> Trả về Trending
		return r.GetTrending(limit)
	}

	// Lấy chi tiết các tin đã xem để phân tích
	var watchedProps []model.RealEstate
	if err := r.db.Where("id IN (?)", watchedIDs).Find(&watchedProps).Error; err != nil || len(watchedProps) == 0 {
		return r.GetTrending(limit)
	}

	// Phân tích các tiêu chí: Lấy danh sách các quận (district) và khoảng giá trung bình
	var districts []string
	var totalPrice float64
	for _, prop := range watchedProps {
		if prop.District != "" {
			districts = append(districts, prop.District)
		}
		totalPrice += prop.PriceVND
	}

	avgPrice := totalPrice / float64(len(watchedProps))
	minPrice := avgPrice * 0.7
	maxPrice := avgPrice * 1.3

	// Query tìm các tin BĐS phù hợp (cùng quận hoặc tầm giá tương đương), loại trừ các tin đã xem
	err := r.db.Raw(
		"SELECT re.id, re.title, re.slug, re.price_vnd, re.address, re.district, re.city, "+
			"re.acreage, re.price_per_m2, re.bedrooms, re.bathrooms, re.description, re.created_at, "+
			"re.house_direction, re.balcony_direction, re.floors, re.legal_docs, re.interior, "+
			"re.price_electricity, re.price_water, re.price_internet, re.amenities, re.latitude, re.longitude, "+
			"COALESCE(GROUP_CONCAT(DISTINCT img.url ORDER BY img.id SEPARATOR '|'), '') AS image_urls, "+
			"COALESCE(u.name, '') AS agent_name, "+
			"COALESCE(u.phone, '') AS agent_phone "+
			"FROM real_estates re "+
			"LEFT JOIN images img ON img.real_estate_id = re.id "+
			"LEFT JOIN users u ON u.id = re.user_id "+
			"WHERE re.id NOT IN (?) AND (re.district IN (?) OR (re.price_vnd BETWEEN ? AND ?)) "+
			"GROUP BY re.id "+
			"ORDER BY re.created_at DESC, re.id DESC "+
			"LIMIT ?",
		watchedIDs,
		districts,
		minPrice,
		maxPrice,
		limit,
	).Scan(&items).Error

	if err != nil || len(items) == 0 {
		return r.GetTrending(limit)
	}

	// for i := range items {
	// 	toResponse(&items[i])
	// }
	return items, nil
}
