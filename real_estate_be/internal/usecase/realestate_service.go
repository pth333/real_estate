package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
	"real_estate_be/pkg/kafka"
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
	producer     *kafka.Producer
}

type IRealEstateService interface {
	ListRealEstate(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error)
	ListRealEstateByCategory(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error)
	GetByID(id uint64) (*dto.RealEstateResponse, error)
	GetListCity() ([]model.Province, error)
	GetListProject(provinceCode, wardCode string) ([]model.RealEstateProject, error)
	GetCategoryBySlug(slug string) (*model.Category, error)
	GetProjectsByCategoryID(categoryID int64) ([]model.RealEstateProject, error)
	GetListWard(provinceCode string) ([]model.Ward, error)
	GetListRealEstateTypes() ([]model.Category, error)
	GetUserByEmail(email string) (*model.User, error)
	CreateRealEstate(req dto.CreateRealEstateRequest, userID uint64) (uint64, error)
	GetTopCity(limit int) ([]model.CityStat, error)
	GetFirstCategorySlug() (string, error)
	ToSlug(city string) string
	GenerateListingSlug(title string, id uint64) string
	IncrementProjectView(id uint64) error
	GetFeaturedProjects(limit int) ([]model.RealEstateProject, error)
	GetProjectByID(id uint64) (*model.RealEstateProject, error)
	GetRecommendations(userID uint64, sessionID string, limit int) ([]dto.RealEstateResponse, error)
}

func NewRealEstateService(repo repo.RealEstateRepository, categoryRepo repo.ICategoryRepository, imageRepo repo.ImageRepository, userRepo repo.IUserRepository, producer *kafka.Producer) IRealEstateService {
	return &RealEstateService{repo: repo, categoryRepo: categoryRepo, imageRepo: imageRepo, userRepo: userRepo, producer: producer}
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
	if req.Slug != "" {
		if id, err := s.categoryRepo.GetCategoryIdBySlug(req.Slug); err == nil && id > 0 {
			return s.repo.GetListByCategory(offset, req, limit)
		}
	}
	return s.repo.GetList(req, offset, limit)
}

func (s *RealEstateService) GetListCity() ([]model.Province, error) {
	return s.repo.GetListCity()
}

func (s *RealEstateService) GetListProject(provinceCode, wardCode string) ([]model.RealEstateProject, error) {
	return s.repo.GetListProject(provinceCode, wardCode)
}

func (s *RealEstateService) GetCategoryBySlug(slug string) (*model.Category, error) {
	return s.categoryRepo.GetCategoryBySlug(slug)
}

func (s *RealEstateService) GetProjectsByCategoryID(categoryID int64) ([]model.RealEstateProject, error) {
	return s.repo.GetProjectsByCategoryID(categoryID)
}

func (s *RealEstateService) GetListWard(provinceCode string) ([]model.Ward, error) {
	return s.repo.GetListWard(provinceCode)
}
func (s *RealEstateService) GetListRealEstateTypes() ([]model.Category, error) {
	return s.repo.GetListRealEstateTypes()
}

func (s *RealEstateService) GetTopCity(limit int) ([]model.CityStat, error) {
	cities, err := s.repo.GetTopCityByCount(limit)
	if err != nil {
		return nil, err
	}
	return cities, nil
}

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

func (s *RealEstateService) CreateRealEstate(req dto.CreateRealEstateRequest, userID uint64) (uint64, error) {
	var categoryID *int64
	fmt.Println("type", req.RealEstateType)
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

	amenitiesJSON, err := json.Marshal(req.Amenities)
	if err != nil {
		return 0, fmt.Errorf("lỗi mã hoá tiện ích: %w", err)
	}

	address := strings.TrimSpace(strings.Join([]string{
		req.DetailAddress, wardName, cityName,
	}, " "))

	estate := &model.RealEstate{
		ProjectID:        req.ProjectID,
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

	if estate.Slug == "" {
		estate.Slug = s.GenerateListingSlug(estate.Title, estate.ID)
		if err := s.repo.Save(estate); err != nil {
			return 0, err
		}
	}

	if err := s.imageRepo.LinkToRealEstate(req.ImageIDs, estate.ID); err != nil {
		return 0, err
	}

	// Gửi Kafka event chuẩn cho Notify
	if s.producer != nil {
		topic := global.Config.Kafka.Topics.RealEstateNotified
		if topic == "" {
			topic = "real_estate.notified.v1" // Fallback topic name chuẩn
		}

		event := kafka.NewRealEstateNewListingEvent(*estate)
		key := strconv.FormatUint(estate.ID, 10) // Key là ID BĐS

		if err := s.producer.Publish(context.Background(), topic, key, event); err != nil {
			log.Printf("⚠️ [Kafka] publish notify error: %v", err)
		} else {
			log.Printf("✅ [Kafka] published new listing notify event for ID: %d to topic: %s", estate.ID, topic)
		}
	}

	return estate.ID, nil
}

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

func (s *RealEstateService) GenerateListingSlug(title string, id uint64) string {
	return fmt.Sprintf("%s-rs%d", s.ToSlug(title), id)
}

func (s *RealEstateService) IncrementProjectView(id uint64) error {
	return s.repo.IncrementProjectView(id)
}

func (s *RealEstateService) GetFeaturedProjects(limit int) ([]model.RealEstateProject, error) {
	return s.repo.GetFeaturedProjects(limit)
}

func (s *RealEstateService) GetProjectByID(id uint64) (*model.RealEstateProject, error) {
	return s.repo.GetProjectByID(id)
}

// GetRecommendations lấy danh sách gợi ý BĐS hỗ trợ phân tầng Cache Redis và fallback Trending
func (s *RealEstateService) GetRecommendations(userID uint64, sessionID string, limit int) ([]dto.RealEstateResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	ctx := context.Background()
	var cacheKey string

	// 1. Định nghĩa Cache Key dựa trên User/Session
	if userID > 0 {
		cacheKey = fmt.Sprintf("rec:user:%d", userID)
	} else if sessionID != "" {
		cacheKey = fmt.Sprintf("rec:session:%s", sessionID)
	}

	// 2. Thử lấy danh sách ID từ Redis Cache
	if cacheKey != "" && global.RedisClient != nil {
		cachedData, err := global.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil && cachedData != "" {
			var cachedIDs []uint64
			if err := json.Unmarshal([]byte(cachedData), &cachedIDs); err == nil && len(cachedIDs) > 0 {
				// Query DB lấy chi tiết tin từ danh sách ID đã cache (đảm bảo tính chất ranking)
				props, err := s.repo.GetByIDs(cachedIDs)
				if err == nil && len(props) > 0 {
					return props, nil
				}
			}
		}
	}

	// 3. Cache Miss hoặc lỗi Redis -> Gọi repo tính toán gợi ý cơ bản / fallback Trending
	props, err := s.repo.GetRecommendationsBasic(userID, sessionID, limit)
	if err != nil {
		return nil, err
	}

	// 4. Nếu lấy được kết quả và có cacheKey -> Cache danh sách ID vào Redis
	if len(props) > 0 && cacheKey != "" && global.RedisClient != nil {
		ids := make([]uint64, len(props))
		for i, prop := range props {
			ids[i] = prop.ID
		}

		encoded, err := json.Marshal(ids)
		if err == nil {
			ttl := 30 * time.Minute
			if userID > 0 {
				ttl = time.Hour // User đã đăng nhập giữ lâu hơn
			}
			global.RedisClient.Set(ctx, cacheKey, string(encoded), ttl)
		}
	}

	return props, nil
}
