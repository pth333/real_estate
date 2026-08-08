package usecase

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
)

// Tỷ giá tạm để quy đổi unit usd/eur → VND. Nên chuyển sang config sau.
const (
	USDToVND = 25000
	EURToVND = 27000
)

type RealEstateService struct {
	repo         repo.RealEstateRepository
	categoryRepo repo.ICategoryRepository
	imageRepo    repo.ImageRepository
	userRepo     repo.IUserRepository
}

type IRealEstateService interface {
	ListRealEstate(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error)
	ListRealEstateByCategory(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error)
	// GetByID lấy 1 tin đăng theo ID (trang chi tiết slug -rs{id}); trả nil nếu không tồn tại.
	GetByID(id uint64) (*dto.RealEstateResponse, error)
	GetListCity() ([]model.Province, error)
	GetListWard(provinceCode string) ([]model.Ward, error)
	GetListRealEstateTypes() ([]model.Category, error)
	GetUserByEmail(email string) (*model.User, error)
	CreateRealEstate(req dto.CreateRealEstateRequest, userID uint64) (uint64, error)
	GetTopCity(limit int) ([]model.CityStat, error)
	GetFirstCategorySlug() (string, error)
	ToSlug(city string) string

	// ApplyFilterSegment(seg string, filter *dto.Filter) error
	// GenerateListingSlug tạo slug trang chi tiết "{title-slug}-rs{id}".
	GenerateListingSlug(title string, id uint64) string
}

func NewRealEstateService(repo repo.RealEstateRepository, categoryRepo repo.ICategoryRepository, imageRepo repo.ImageRepository, userRepo repo.IUserRepository) IRealEstateService {
	return &RealEstateService{repo: repo, categoryRepo: categoryRepo, imageRepo: imageRepo, userRepo: userRepo}
}

func (s *RealEstateService) ListRealEstate(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error) {
	limit := req.Size
	if limit < 1 {
		limit = 10
	}
	offset := (req.Page - 1) * limit
	return s.repo.GetList(req, offset, limit)
}

func (s *RealEstateService) GetByID(id uint64) (*dto.RealEstateResponse, error) {
	return s.repo.GetByID(id)
}

func (s *RealEstateService) ListRealEstateByCategory(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error) {
	limit := req.Size
	if limit < 1 {
		limit = 10
	}
	offset := (req.Page - 1) * limit
	jsonData, _ := json.MarshalIndent(req, "", "  ")
	fmt.Println("Data:", string(jsonData))

	// Nếu slug khớp category trong DB → giữ query theo category
	if req.Slug != "" {
		if id, err := s.categoryRepo.GetCategoryIdBySlug(req.Slug); err == nil && id > 0 {
			return s.repo.GetListByCategory(offset, req, limit)
		}
	}

	// Không khớp category → fallback: query theo payload filter (bỏ category)
	return s.repo.GetList(req, offset, limit)
}

func (s *RealEstateService) GetListCity() ([]model.Province, error) {
	return s.repo.GetListCity()
}

func (s *RealEstateService) GetListWard(provinceCode string) ([]model.Ward, error) {
	return s.repo.GetListWard(provinceCode)
}
func (s *RealEstateService) GetListRealEstateTypes() ([]model.Category, error) {
	return s.repo.GetListRealEstateTypes()
}

// GetTopCity trả về N thành phố có nhiều BĐS nhất, kèm slug để FE điều hướng.
func (s *RealEstateService) GetTopCity(limit int) ([]model.CityStat, error) {
	cities, err := s.repo.GetTopCityByCount(limit)
	if err != nil {
		return nil, err
	}

	return cities, nil
}

// GetFirstCategorySlug lấy slug của category đầu tiên trong menu (theo name ASC),
// dùng làm segment `{category}` cho khối "Bất động sản theo địa điểm" trên trang chủ.
func (s *RealEstateService) GetFirstCategorySlug() (string, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		return "", err
	}
	for _, c := range categories {
		if c.Slug != "" {
			return c.Slug, nil
		}
	}
	return "", nil
}

func (s *RealEstateService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepo.FindByEmail(email)
}

// CreateRealEstate tạo tin đăng: lưu RealEstate rồi gắn các ảnh đã upload vào tin
func (s *RealEstateService) CreateRealEstate(req dto.CreateRealEstateRequest, userID uint64) (uint64, error) {
	// Map category theo real_estate_type (FE gửi thẳng id, dạng chuỗi số)
	var categoryID *int64
	if req.RealEstateType != "" {
		if id, err := strconv.ParseInt(req.RealEstateType, 10, 64); err == nil {
			categoryID = &id
		}
	}

	pricePerM2 := req.Price
	switch req.Unit {
	case "usd":
		pricePerM2 = req.Price * USDToVND
	case "eur":
		pricePerM2 = req.Price * EURToVND
	}
	priceVND := pricePerM2 * req.Area

	cityName, cityErr := s.repo.GetProvinceNameByCode(req.Province)
	wardName, wardErr := s.repo.GetWardNameByCode(req.Ward)

	if cityErr != nil {
		return 0, fmt.Errorf("không tìm thấy tỉnh/thành: %s", req.Province)
	}
	if wardErr != nil {
		return 0, fmt.Errorf("không tìm thấy phường/xã: %s", req.Ward)
	}

	// Tiện ích (mảng) → JSON string để lưu vào cột amenities (MySQL không có kiểu array)
	amenitiesJSON, err := json.Marshal(req.Amenities)
	if err != nil {
		return 0, fmt.Errorf("lỗi mã hoá tiện ích: %w", err)
	}

	// Địa chỉ ghép từ detail + tên phường/xã + tên tỉnh/thành
	address := strings.TrimSpace(strings.Join([]string{
		req.DetailAddress, wardName, cityName,
	}, " "))

	estate := &model.RealEstate{
		UserID:           &userID,
		Title:            req.Title,
		PriceVND:         priceVND,
		Address:          address,
		District:         wardName,
		City:             cityName,
		Acreage:          req.Area,
		PricePerM2:       pricePerM2,
		CategoryID:       categoryID,
		Description:      req.Description,
		Bedrooms:         req.BedroomCount,
		Bathrooms:        req.BathroomCount,
		Amenities:        string(amenitiesJSON),
		HouseDirection:   req.HouseDirection,
		BalconyDirection: req.BalconyDirection,
		Floors:           req.FloorCount,
		LegalDocs:        req.LegalDocs,
		Interior:         req.Interior,
		PriceElectricity: req.PriceElectricity,
		PriceWater:       req.PriceWater,
		PriceInternet:    req.PriceInternet,
	}

	if err := s.repo.Create(estate); err != nil {
		return 0, err
	}

	// Generate slug trang chi tiết sau khi có ID: "{title-slug}-rs{id}".
	if estate.Slug == "" {
		estate.Slug = s.GenerateListingSlug(estate.Title, estate.ID)
		if err := s.repo.Save(estate); err != nil {
			return 0, err
		}
	}

	// Gắn ảnh đã upload vào tin đăng
	if err := s.imageRepo.LinkToRealEstate(req.ImageIDs, estate.ID); err != nil {
		return 0, err
	}

	return estate.ID, nil
}

// ToSlug chuyển tên tiếng Việt sang slug không dấu, VD "Hồ Chí Minh" → "ho-chi-minh".
func (s *RealEstateService) ToSlug(input string) string {
	var accents = map[rune]string{
		'à': "a", 'á': "a", 'ả': "a", 'ã': "a", 'ạ': "a", 'ă': "a", 'ắ': "a", 'ằ': "a", 'ẳ': "a", 'ẵ': "a", 'ặ': "a", 'â': "a", 'ấ': "a", 'ầ': "a", 'ẩ': "a", 'ẫ': "a", 'ậ': "a",
		'è': "e", 'é': "e", 'ẻ': "e", 'ẽ': "e", 'ẹ': "e", 'ê': "e", 'ế': "e", 'ề': "e", 'ể': "e", 'ễ': "e", 'ệ': "e",
		'ì': "i", 'í': "i", 'ỉ': "i", 'ĩ': "i", 'ị': "i",
		'ò': "o", 'ó': "o", 'ỏ': "o", 'õ': "o", 'ọ': "o", 'ô': "o", 'ố': "o", 'ồ': "o", 'ổ': "o", 'ỗ': "o", 'ộ': "o", 'ơ': "o", 'ớ': "o", 'ờ': "o", 'ở': "o", 'ỡ': "o", 'ợ': "o",
		'ù': "u", 'ú': "u", 'ủ': "u", 'ũ': "u", 'ụ': "u", 'ư': "u", 'ứ': "u", 'ừ': "u", 'ử': "u", 'ữ': "u", 'ự': "u",
		'ỳ': "y", 'ý': "y", 'ỷ': "y", 'ỹ': "y", 'ỵ': "y",
		'đ': "d",
	}

	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(input)) {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '.' || r == ',' || r == '_':
			// Tránh gạch nối kép
			if b.Len() > 0 && b.String()[b.Len()-1] != '-' {
				b.WriteRune('-')
			}
		default:
			if ch, ok := accents[r]; ok {
				b.WriteString(ch)
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// Slug không khớp bất kỳ range nào → bỏ qua (không lỗi), để URL thừa
// 1 đoạn không làm hỏng trang.
// func (s *RealEstateService) ApplyFilterSegment(seg string, filter *dto.Filter) error {
// 	if seg == "" || filter == nil {
// 		return nil
// 	}
// 	f, err := s.repo.GetFilterRangeBySlug(seg)
// 	if err != nil {
// 		// Lỗi DB thực → báo lỗi; không tìm thấy → trả nil, nil.
// 		return err
// 	}
// 	if f == nil {
// 		// Slug không nằm trong filter_ranges → URL thừa → bỏ qua.
// 		return nil
// 	}
// 	switch f.Type {
// 	case "price":
// 		if f.MinVal != nil {
// 			filter.MinPrice = *f.MinVal
// 		}
// 		if f.MaxVal != nil {
// 			filter.MaxPrice = *f.MaxVal
// 		}
// 	case "area":
// 		if f.MinVal != nil {
// 			filter.MinAcreage = *f.MinVal
// 		}
// 		if f.MaxVal != nil {
// 			filter.MaxAcreage = *f.MaxVal
// 		}
// 	}
// 	return nil
// }

// GenerateListingSlug tạo slug trang chi tiết cho 1 tin đăng: "{title-slug}-rs{id}".
// VD title "Nhà phố 2 tầng Cầu Giấy" id 123 → "nha-pho-2-tang-cau-giay-rs123".
// Phần "-rs{id}" giúp backend detect (regex -rs\d+$) và đảm bảo unique theo id.
func (s *RealEstateService) GenerateListingSlug(title string, id uint64) string {
	return fmt.Sprintf("%s-rs%d", s.ToSlug(title), id)
}

func (s *RealEstateService) GetCategory() ([]model.Category, error) {
	return s.repo.GetCategory()
}
