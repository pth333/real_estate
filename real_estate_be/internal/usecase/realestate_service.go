package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	repo              repo.RealEstateRepository
	categoryRepo      repo.ICategoryRepository
	imageRepo         repo.ImageRepository
	userRepo          repo.IUserRepository
	searchHistoryRepo repo.ISearchHistoryRepository
	producer          *kafka.Producer
}

type IRealEstateService interface {
	ListRealEstate(req dto.RealEstateSearchRequest, userID uint64, sessionID string) ([]dto.RealEstateResponse, int64, error)
	ListRealEstateByCategory(req dto.RealEstateSearchRequest, userID uint64, sessionID string) ([]dto.RealEstateResponse, int64, error)
	GetByID(id uint64) (*dto.RealEstateResponse, error)
	GetListCity() ([]model.Province, error)
	GetListProject(provinceCode, wardCode string) ([]model.RealEstateProject, error)
	GetCategoryBySlug(slug string) (*model.Category, error)
	GetProjectsByCategoryID(categoryID int64) ([]model.RealEstateProject, error)
	GetListWard(provinceCode string) ([]model.Ward, error)
	GetListRealEstateTypes() ([]model.Category, error)
	GetUserByEmail(email string) (*model.User, error)
	GetTopCity(limit int) ([]model.CityStat, error)
	GetFirstCategorySlug() (string, error)
	ToSlug(city string) string
	IncrementProjectView(id uint64) error
	GetFeaturedProjects(limit int) ([]model.RealEstateProject, error)
	GetProjectByID(id uint64) (*model.RealEstateProject, error)
	GetRecommendations(userID uint64, sessionID string, limit int) ([]dto.RealEstateResponse, error)

	// ── Bất động sản yêu thích (favorite) ──
	ToggleFavorite(userID, realEstateID uint64) (bool, error)
	ListFavorites(userID uint64, page, size int) ([]dto.RealEstateResponse, int64, error)
}

func NewRealEstateService(
	repo repo.RealEstateRepository,
	categoryRepo repo.ICategoryRepository,
	imageRepo repo.ImageRepository,
	userRepo repo.IUserRepository,
	searchHistoryRepo repo.ISearchHistoryRepository,
	producer *kafka.Producer,
) IRealEstateService {
	return &RealEstateService{
		repo:              repo,
		categoryRepo:      categoryRepo,
		imageRepo:         imageRepo,
		userRepo:          userRepo,
		searchHistoryRepo: searchHistoryRepo,
		producer:          producer,
	}
}

func MapRealEstateResponse(data []dto.RealEstateResponse) {
	for i := range data {
		m := &data[i]

		if m.ImageURLsRaw != "" {
			m.ImageURLs = strings.Split(m.ImageURLsRaw, "|")
		} else {
			m.ImageURLs = []string{}
		}

		// AmenitiesRaw → []string
		if m.AmenitiesRaw != "" {
			var amenities []string

			if err := json.Unmarshal(
				[]byte(m.AmenitiesRaw),
				&amenities,
			); err == nil {
				m.Amenities = amenities
			} else {
				m.Amenities = []string{}
			}
		} else {
			m.Amenities = []string{}
		}
	}
}

func (s *RealEstateService) ListRealEstate(req dto.RealEstateSearchRequest, userID uint64, sessionID string) ([]dto.RealEstateResponse, int64, error) {
	limit := req.Size
	if limit < 1 {
		limit = 10
	}
	offset := (req.Page - 1) * limit
	data, total, err := s.repo.GetList(req, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Gắn cờ yêu thích cho user đã đăng nhập
	s.flagFavorites(data, userID)

	// Tự động lưu lịch sử tìm kiếm tối giản nếu tìm thấy kết quả
	if req.Search != "" && len(data) > 0 {
		s.autoRecordSearch(userID, sessionID, req.Search, data)
	}

	return data, total, nil
}

func (s *RealEstateService) GetByID(id uint64) (*dto.RealEstateResponse, error) {
	return s.repo.GetByID(id)
}

func (s *RealEstateService) ListRealEstateByCategory(req dto.RealEstateSearchRequest, userID uint64, sessionID string) ([]dto.RealEstateResponse, int64, error) {
	limit := req.Size
	if limit < 1 {
		limit = 10
	}
	offset := (req.Page - 1) * limit

	var data []dto.RealEstateResponse
	var total int64
	var err error

	if req.Slug != "" {
		if id, errCat := s.categoryRepo.GetCategoryIdBySlug(req.Slug); errCat == nil && id > 0 {
			data, total, err = s.repo.GetListByCategory(offset, req, limit)

			// jsonData, _ := json.MarshalIndent(data, " ", " ")
			// fmt.Println("JSON DATA: ", string(jsonData))
		}
	} else {
		data, total, err = s.repo.GetList(req, offset, limit)
	}

	MapRealEstateResponse(data)

	if err != nil {
		return nil, 0, err
	}

	// Tự động lưu lịch sử tìm kiếm tối giản nếu tìm thấy kết quả
	if req.Search != "" && len(data) > 0 {
		s.autoRecordSearch(userID, sessionID, req.Search, data)
	}

	return data, total, nil
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

func (s *RealEstateService) ToSlug(input string) string {
	return slugify(input)
}

// slugify bỏ dấu tiếng Việt, hạ thường và nối bằng dấu "-" (dùng chung cho cả
// slug tin đăng và slug dự án).
func slugify(input string) string {
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
		fmt.Printf("Cache data: %s, key: %s", cachedData, cacheKey)

		if err == nil && cachedData != "" {
			var cachedIDs []uint64
			if err := json.Unmarshal([]byte(cachedData), &cachedIDs); err == nil && len(cachedIDs) > 0 {
				// Query DB lấy chi tiết tin từ danh sách ID đã cache (đảm bảo tính chất ranking)
				props, err := s.repo.GetByIDs(cachedIDs)
				if err == nil && len(props) > 0 {
					s.flagFavorites(props, userID)
					return props, nil
				}
			}
		}
	}

	// 3. Cache Miss hoặc lỗi Redis -> Gọi gRPC Recommendation Service (nếu có)
	var props []dto.RealEstateResponse
	var err error
	var strategyUsed = "db_fallback"

	if global.RecommendationClient != nil {
		var ids []uint64
		ids, strategyUsed, err = global.RecommendationClient.GetRecommendations(
			ctx,
			userID,
			sessionID,
			0,   // realEstateID = 0 (không có trong ngữ cảnh GetRecommendations hiện tại)
			0.0, // lat = 0.0
			0.0, // lon = 0.0
			int32(limit),
		)
		if err == nil && len(ids) > 0 {
			// Query DB lấy chi tiết các BĐS theo thứ tự xếp hạng từ gRPC
			props, err = s.repo.GetByIDs(ids)
			if err != nil {
				log.Printf("[Recommendation] Lỗi khi lấy chi tiết BĐS từ database: %v", err)
			}
		} else if err != nil {
			log.Printf("[Recommendation] Lỗi gọi gRPC Recommendation Service: %v (đang chuyển sang fallback DB)", err)
		}
	}

	// Fallback về cơ chế DB thuần nếu không lấy được kết quả từ gRPC
	if len(props) == 0 {
		props, err = s.repo.GetRecommendationsBasic(userID, sessionID, limit)
		if err != nil {
			return nil, err
		}
		strategyUsed = "db_fallback"
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

	_ = strategyUsed // Có thể dùng để log hoặc xử lý thêm nếu cần
	s.flagFavorites(props, userID)
	return props, nil
}

// autoRecordSearch tự động lưu lịch sử tìm kiếm tối giản kèm danh sách BĐS kết quả
func (s *RealEstateService) autoRecordSearch(userID uint64, sessionID string, query string, results []dto.RealEstateResponse) {
	if query == "" || len(results) == 0 {
		return
	}

	var uID *uint64
	if userID > 0 {
		uID = &userID
	}

	var sID *string
	if sessionID != "" {
		sID = &sessionID
	}

	for _, item := range results {
		estateID := item.ID
		_ = s.searchHistoryRepo.Create(&model.SearchHistory{
			UserID:       uID,
			SessionID:    sID,
			Query:        query,
			RealEstateID: &estateID,
		})
	}
}

// ── Bất động sản yêu thích (favorite) ──

// flagFavorites gắn cờ IsFavorite cho danh sách BĐS dựa trên danh sách ID yêu thích của user.
func (s *RealEstateService) flagFavorites(items []dto.RealEstateResponse, userID uint64) {
	if userID == 0 || len(items) == 0 {
		return
	}
	ids, err := s.repo.GetFavoriteRealEstateIDs(userID)
	if err != nil {
		return
	}
	set := make(map[uint64]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	for i := range items {
		items[i].IsFavorite = set[items[i].ID]
	}
}

// ToggleFavorite thêm/bỏ 1 tin đăng vào danh mục yêu thích của user.
// Trả về trạng thái mới: true = đã yêu thích, false = đã bỏ yêu thích.
func (s *RealEstateService) ToggleFavorite(userID, realEstateID uint64) (bool, error) {
	if userID == 0 {
		return false, fmt.Errorf("người dùng chưa đăng nhập")
	}

	// Đảm bảo tin đăng tồn tại
	if _, err := s.repo.GetByID(realEstateID); err != nil {
		return false, fmt.Errorf("không tìm thấy bất động sản")
	}

	if s.repo.IsFavorite(userID, realEstateID) {
		if err := s.repo.RemoveFavorite(userID, realEstateID); err != nil {
			return false, err
		}
		return false, nil
	}

	if err := s.repo.AddFavorite(userID, realEstateID); err != nil {
		return false, err
	}
	return true, nil
}

// ListFavorites trả về danh sách BĐS yêu thích của user (phân trang).
func (s *RealEstateService) ListFavorites(userID uint64, page, size int) ([]dto.RealEstateResponse, int64, error) {
	if userID == 0 {
		return []dto.RealEstateResponse{}, 0, nil
	}
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 12
	}
	offset := (page - 1) * size
	return s.repo.ListFavoriteRealEstates(userID, size, offset)
}
