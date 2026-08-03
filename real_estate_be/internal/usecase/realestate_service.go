package usecase

import (
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
}

func NewRealEstateService(repo repo.RealEstateRepository, categoryRepo repo.ICategoryRepository, imageRepo repo.ImageRepository, userRepo repo.IUserRepository) IRealEstateService {
	return &RealEstateService{repo: repo, categoryRepo: categoryRepo, imageRepo: imageRepo, userRepo: userRepo}
}

// mapToResponse chuyển model.RealEstate → dto.RealEstateResponse
func mapSliceToResponse(data []model.RealEstate) []dto.RealEstateResponse {
	result := make([]dto.RealEstateResponse, len(data))
	for i, m := range data {
		result[i] = dto.RealEstateResponse{
			ID:               m.ID,
			Title:            m.Title,
			PriceVND:         m.PriceVND,
			Address:          m.Address,
			District:         m.District,
			City:             m.City,
			Acreage:          m.Acreage,
			PricePerM2:       m.PricePerM2,
			TypeOfRealEstate: m.TypeOfRealEstate,
		}
	}
	return result
}

func (s *RealEstateService) ListRealEstate(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error) {
	limit := req.Size
	if limit < 1 {
		limit = 10
	}
	offset := (req.Page - 1) * limit
	data, total, err := s.repo.GetList(req, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return mapSliceToResponse(data), total, nil
}

// ListRealEstateByCategory lấy category từ slug, sau đó query BĐS theo category_id
func (s *RealEstateService) ListRealEstateByCategory(req dto.RealEstateSearchRequest) ([]dto.RealEstateResponse, int64, error) {
	limit := req.Size
	if limit < 1 {
		limit = 10
	}
	offset := (req.Page - 1) * limit
	data, total, err := s.repo.GetListByCategory(offset, req, limit)
	if err != nil {
		return nil, 0, err
	}
	return mapSliceToResponse(data), total, nil
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

func (s *RealEstateService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepo.FindByEmail(email)
}

// CreateRealEstate tạo tin đăng: lưu RealEstate rồi gắn các ảnh đã upload vào tin
func (s *RealEstateService) CreateRealEstate(req dto.CreateRealEstateRequest, userID uint64) (uint64, error) {
	// Map category theo real_estate_type (FE gửi thẳng id, dạng chuỗi số)
	var categoryID *uint64
	if req.RealEstateType != "" {
		if id, err := strconv.ParseUint(req.RealEstateType, 10, 64); err == nil {
			categoryID = &id
		}
	}

	// Tính giá theo đơn vị
	var priceVND float64
	switch req.Unit {
	case "usd":
		priceVND = req.Price * USDToVND
	case "eur":
		priceVND = req.Price * EURToVND
	default: // vnd
		priceVND = req.Price
	}

	// Địa chỉ ghép từ province + ward + detail
	address := strings.TrimSpace(strings.Join([]string{
		req.DetailAddress, req.Ward, req.Province,
	}, " "))

	var pricePerM2 float64
	if req.Area > 0 {
		pricePerM2 = priceVND / req.Area
	}

	estate := &model.RealEstate{
		UserID:           &userID,
		Title:            req.Title,
		PriceVND:         priceVND,
		Address:          address,
		District:         req.Ward,
		City:             req.Province,
		Acreage:          req.Area,
		PricePerM2:       pricePerM2,
		CategoryID:       categoryID,
		TypeOfRealEstate: req.RealEstateType,
		Description:      req.Description,
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
