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
	GetListCity() ([]model.Province, error)
	GetListWard(provinceCode string) ([]model.Ward, error)
	GetListRealEstateTypes() ([]model.Category, error)
	GetUserByEmail(email string) (*model.User, error)
	CreateRealEstate(req dto.CreateRealEstateRequest, userID uint64) (uint64, error)
	GetTopCity(limit int) ([]model.CityStat, error)
	ToSlug(city string) string
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

// ListRealEstateByCategory lấy category từ slug, nếu slug là category hợp lệ
// thì query BĐS theo category. Nếu slug không khớp category nào (VD slug SEO mỹ
// thuật) thì bỏ qua lọc category, fallback query theo payload filter.
func (s *RealEstateService) ListRealEstateByCategory(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error) {
	limit := req.Size
	if limit < 1 {
		limit = 10
	}
	offset := (req.Page - 1) * limit

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
		UserID:      &userID,
		Title:       req.Title,
		PriceVND:    priceVND,
		Address:     address,
		District:    wardName,
		City:        cityName,
		Acreage:     req.Area,
		PricePerM2:  pricePerM2,
		CategoryID:  categoryID,
		Description: req.Description,
		Bedrooms:    req.BedroomCount,
		Bathrooms:   req.BathroomCount,
		Amenities:   string(amenitiesJSON),
	}

	if err := s.repo.Create(estate); err != nil {
		return 0, err
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
			if s, ok := accents[r]; ok {
				b.WriteString(s)
			}
		}
	}
	return strings.Trim(b.String(), "-")
}
