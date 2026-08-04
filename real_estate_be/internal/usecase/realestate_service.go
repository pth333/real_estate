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
}

func NewRealEstateService(repo repo.RealEstateRepository, categoryRepo repo.ICategoryRepository, imageRepo repo.ImageRepository, userRepo repo.IUserRepository) IRealEstateService {
	return &RealEstateService{repo: repo, categoryRepo: categoryRepo, imageRepo: imageRepo, userRepo: userRepo}
}

// mapToResponse chuyển model.RealEstate → dto.RealEstateResponse
func mapSliceToResponse(data []model.RealEstate) []dto.RealEstateResponse {
	result := make([]dto.RealEstateResponse, len(data))
	for i, m := range data {
		result[i] = dto.RealEstateResponse{
			ID:         m.ID,
			Title:      m.Title,
			PriceVND:   m.PriceVND,
			Address:    m.Address,
			District:   m.District,
			City:       m.City,
			Acreage:    m.Acreage,
			PricePerM2: m.PricePerM2,
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
