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

type IManagerPostUseCase interface {
	GetManagerPostsList(userID uint64, search string, page, size int) ([]dto.ManagerPostListResponse, int64, error)
	CreateRealEstate(req dto.CreateRealEstateRequest, userID uint64) (uint64, error)
	UpdateRealEstate(id uint64, req dto.CreateRealEstateRequest) error
	DeleteManagerPost(postID uint64) error
	GenerateListingSlug(title string, id uint64) string
	CreateProject(req dto.CreateProjectRequest) (uint64, error)
	UpdateProject(id uint64, req dto.CreateProjectRequest) error
	GetProjectDetail(id uint64) (*dto.ManagerProjectDetailResponse, error)
	ListProjects(search string, page, size int) ([]dto.ManagerProjectResponse, int64, error)
}

type managerPostUseCase struct {
	managerRepo    repo.ManagerPostRepository
	realEstateRepo repo.RealEstateRepository
	imageRepo      repo.ImageRepository
	producer       *kafka.Producer
}

func NewManagerPostUseCase(managerRepo repo.ManagerPostRepository, realEstateRepo repo.RealEstateRepository, imageRepo repo.ImageRepository, producer *kafka.Producer) IManagerPostUseCase {
	return &managerPostUseCase{
		managerRepo:    managerRepo,
		realEstateRepo: realEstateRepo,
		imageRepo:      imageRepo,
		producer:       producer,
	}
}

func (u *managerPostUseCase) GenerateListingSlug(title string, id uint64) string {
	slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	return fmt.Sprintf("%s-rs%d", slug, id)
}
func (u *managerPostUseCase) GetManagerPostsList(userID uint64, search string, page, size int) ([]dto.ManagerPostListResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	items, total, err := u.managerRepo.GetManagerPostsList(userID, search, offset, size)
	if err != nil {
		return nil, 0, err
	}

	dtoItems := make([]dto.ManagerPostListResponse, 0, len(items))
	for _, item := range items {
		thumbnail := ""
		if len(item.Images) > 0 {
			thumbnail = item.Images[0].URL // Sử dụng trường URL ảnh
		}

		// Ánh xạ tên loại BĐS
		realEstateType := "Bất động sản"
		if item.Category != nil {
			realEstateType = item.Category.Name
		}

		// Lấy trạng thái lưu tạm ở cột BalconyDirection
		// statusVal := "pending"
		// if item.BalconyDirection != "" {
		// 	statusVal = item.BalconyDirection
		// }

		dtoItems = append(dtoItems, dto.ManagerPostListResponse{
			ID:        item.ID,
			Title:     item.Title,
			Slug:      item.Slug,
			Thumbnail: thumbnail,
			Type:      realEstateType,
			Price:     item.PriceVND,
			Unit:      "vnd",
			Area:      item.Acreage,
			CreatedAt: item.CreatedAt.Format("02-01-2006"),
		})
	}

	return dtoItems, total, err
}

func (s *managerPostUseCase) CreateRealEstate(req dto.CreateRealEstateRequest, userID uint64) (uint64, error) {
	var categoryID *int64
	fmt.Println("type", req.RealEstateType)
	if req.RealEstateType != "" {
		if id, err := strconv.ParseInt(req.RealEstateType, 10, 64); err == nil {
			categoryID = &id
		}
	}

	pricePerM2 := req.PricePerM2
	switch req.Unit {
	case "usd":
		pricePerM2 = req.PricePerM2 * USDToVND
	case "eur":
		pricePerM2 = req.PricePerM2 * EURToVND
	}
	priceVND := pricePerM2 * req.Area

	amenitiesJSON, err := json.Marshal(req.Amenities)
	if err != nil {
		return 0, fmt.Errorf("lỗi mã hoá tiện ích: %w", err)
	}

	cityName := req.Province
	if name, err := s.realEstateRepo.GetProvinceNameByCode(req.Province); err == nil && name != "" {
		cityName = name
	}
	wardName := req.Ward
	if name, err := s.realEstateRepo.GetWardNameByCode(req.Ward); err == nil && name != "" {
		wardName = name
	}

	// address := strings.TrimSpace(strings.Join([]string{
	// 	req.DetailAddress, wardName, cityName,
	// }, " "))

	estate := &model.RealEstate{
		ProjectID:        req.ProjectID,
		UserID:           &userID,
		Title:            req.Title,
		PriceVND:         priceVND,
		Address:          req.DetailAddress,
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
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
	}

	if err := s.realEstateRepo.Create(estate); err != nil {
		return 0, err
	}

	if estate.Slug == "" {
		estate.Slug = s.GenerateListingSlug(estate.Title, estate.ID)
		if err := s.realEstateRepo.Save(estate); err != nil {
			return 0, err
		}
	}

	imageIDs := make([]uint64, 0, len(req.Images))

	for _, image := range req.Images {
		imageIDs = append(imageIDs, image.ID)
	}

	if err := s.imageRepo.LinkToRealEstate(imageIDs, estate.ID); err != nil {
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

func (s *managerPostUseCase) UpdateRealEstate(id uint64, req dto.CreateRealEstateRequest) error {
	var rawEstate model.RealEstate
	if err := global.DB.First(&rawEstate, id).Error; err != nil {
		return fmt.Errorf("không tìm thấy BDS cần cập nhật")
	}
	var categoryID *int64
	if req.RealEstateType != "" {
		if catID, err := strconv.ParseInt(req.RealEstateType, 10, 64); err == nil {
			categoryID = &catID
		}
	}

	pricePerM2 := req.PricePerM2
	switch req.Unit {
	case "usd":
		pricePerM2 = req.PricePerM2 * USDToVND
	case "eur":
		pricePerM2 = req.PricePerM2 * EURToVND
	}
	priceVND := pricePerM2 * req.Area

	amenitiesJSON, err := json.Marshal(req.Amenities)
	if err != nil {
		return fmt.Errorf("lỗi mã hoá tiện ích: %w", err)
	}

	// Đổi code tỉnh/phường → tên trước khi lưu (FE gửi code, cột city/district cần tên).
	// Không tìm thấy code (giá trị đã là tên, VD khi edit tin cũ) → giữ nguyên.
	cityName := req.Province
	if name, err := s.realEstateRepo.GetProvinceNameByCode(req.Province); err == nil && name != "" {
		cityName = name
	}
	wardName := req.Ward
	if name, err := s.realEstateRepo.GetWardNameByCode(req.Ward); err == nil && name != "" {
		wardName = name
	}

	// address := strings.TrimSpace(strings.Join([]string{
	// 	req.DetailAddress, wardName, cityName,
	// }, " "))

	rawEstate.ProjectID = req.ProjectID
	rawEstate.Title = req.Title
	rawEstate.PriceVND = priceVND
	rawEstate.Address = req.DetailAddress
	rawEstate.District = wardName
	rawEstate.City = cityName
	rawEstate.Acreage = req.Area
	rawEstate.PricePerM2 = pricePerM2
	rawEstate.CategoryID = categoryID
	rawEstate.Description = req.Description
	rawEstate.Bedrooms = req.BedroomCount
	rawEstate.Bathrooms = req.BathroomCount
	rawEstate.Amenities = string(amenitiesJSON)
	rawEstate.HouseDirection = req.HouseDirection
	rawEstate.BalconyDirection = req.BalconyDirection
	rawEstate.Floors = req.FloorCount
	rawEstate.LegalDocs = req.LegalDocs
	rawEstate.Interior = req.Interior
	rawEstate.PriceElectricity = req.PriceElectricity
	rawEstate.PriceWater = req.PriceWater
	rawEstate.PriceInternet = req.PriceInternet
	rawEstate.Latitude = req.Latitude
	rawEstate.Longitude = req.Longitude

	if err := s.realEstateRepo.Save(&rawEstate); err != nil {
		return err
	}
	jsonData, err := json.MarshalIndent(req.Images, " ", " ")

	fmt.Println("Image: ", string(jsonData))

	imageIDs := make([]uint64, 0, len(req.Images))
	for _, image := range req.Images {
		imageIDs = append(imageIDs, image.ID)
	}

	// Gỡ liên kết toàn bộ ảnh cũ ra trước khi lưu mới
	global.DB.Model(&model.Image{}).Where("real_estate_id = ?", id).Update("real_estate_id", nil)
	if err := s.imageRepo.LinkToRealEstate(imageIDs, id); err != nil {
		return err
	}

	return nil
}

func (u *managerPostUseCase) DeleteManagerPost(postID uint64) error {
	_, err := u.managerRepo.GetByID(postID)
	if err != nil {
		return err
	}

	// Gỡ liên kết ảnh trước khi xóa bài viết
	_ = u.managerRepo.UnlinkImages(postID)

	return u.managerRepo.DeleteManagerPost(postID)
}

// CreateProject tạo dự án mới. province/ward là MÃ tỉnh/phường (VD "000331"),
// lưu dạng string để giữ số 0 đầu và map đúng với bảng provinces/wards.
func (u *managerPostUseCase) CreateProject(req dto.CreateProjectRequest) (uint64, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return 0, fmt.Errorf("tên dự án không được để trống")
	}
	if req.ProvinceCode == "" {
		return 0, fmt.Errorf("vui lòng chọn tỉnh/thành phố")
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	project := &model.RealEstateProject{
		Name:            name,
		AlternativeName: req.AlternativeName,
		Slug:            slugify(name),
		Status:          status,
		FullAddress:     req.FullAddress,
		ProvinceCode:    req.ProvinceCode,
		WardCode:        req.WardCode,
		CategoryID:      req.CategoryID,
		TotalAreaHA:     req.TotalAreaHA,
		TotalUnits:      req.TotalUnits,
		PriceMin:        req.PriceMin,
		PriceMax:        req.PriceMax,
	}

	// Parse ngày bắt đầu / bàn giao (format "2006-01-02")
	if req.ConstructionStartDate != nil && *req.ConstructionStartDate != "" {
		if t, err := time.Parse("2006-01-02", *req.ConstructionStartDate); err == nil {
			project.ConstructionStartDate = &t
		}
	}
	if req.HandoverDate != nil && *req.HandoverDate != "" {
		if t, err := time.Parse("2006-01-02", *req.HandoverDate); err == nil {
			project.HandoverDate = &t
		}
	}

	if err := u.realEstateRepo.CreateProject(project); err != nil {
		return 0, err
	}

	// Liên kết ảnh dự án đã upload với dự án vừa tạo
	if err := u.imageRepo.LinkToProject(req.ImageIDs, project.ID); err != nil {
		return 0, fmt.Errorf("liên kết ảnh dự án thất bại: %w", err)
	}

	return project.ID, nil
}

// ListProjects trả về danh sách dự án (phân trang, lọc theo tên/địa chỉ).
func (u *managerPostUseCase) ListProjects(search string, page, size int) ([]dto.ManagerProjectResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	projects, total, err := u.realEstateRepo.ListProjects(size, offset, search)
	if err != nil {
		return nil, 0, err
	}

	items := make([]dto.ManagerProjectResponse, 0, len(projects))
	for _, p := range projects {
		items = append(items, dto.ManagerProjectResponse{
			ID:              p.ID,
			Name:            p.Name,
			Slug:            p.Slug,
			AlternativeName: p.AlternativeName,
			Status:          p.Status,
			FullAddress:     p.FullAddress,
			TotalAreaHA:     p.TotalAreaHA,
			TotalUnits:      p.TotalUnits,
			PriceMin:        p.PriceMin,
			PriceMax:        p.PriceMax,
			CreatedAt:       p.CreatedAt.Format("02-01-2006"),
		})
	}

	// Gắn ảnh đại diện (ảnh đầu tiên) cho từng dự án
	if len(items) > 0 {
		ids := make([]uint64, len(items))
		for i := range items {
			ids[i] = items[i].ID
		}
		if covers, err := u.imageRepo.GetProjectCoverImages(ids); err == nil {
			for i := range items {
				if url, ok := covers[items[i].ID]; ok {
					items[i].Thumbnail = url
				}
			}
		}
	}

	return items, total, nil
}

// GetProjectDetail trả về chi tiết 1 dự án để điền form chỉnh sửa.
func (u *managerPostUseCase) GetProjectDetail(id uint64) (*dto.ManagerProjectDetailResponse, error) {
	project, err := u.realEstateRepo.GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	resp := &dto.ManagerProjectDetailResponse{
		ID:              project.ID,
		Name:            project.Name,
		Slug:            project.Slug,
		AlternativeName: project.AlternativeName,
		Status:          project.Status,
		FullAddress:     project.FullAddress,
		ProvinceCode:    project.ProvinceCode,
		WardCode:        project.WardCode,
		CategoryID:      project.CategoryID,
		TotalAreaHA:     project.TotalAreaHA,
		TotalUnits:      project.TotalUnits,
		PriceMin:        project.PriceMin,
		PriceMax:        project.PriceMax,
	}
	if project.ConstructionStartDate != nil {
		resp.ConstructionStartDate = project.ConstructionStartDate.Format("2006-01-02")
	}
	if project.HandoverDate != nil {
		resp.HandoverDate = project.HandoverDate.Format("2006-01-02")
	}
	// Kèm danh sách ảnh dự án (để hiển thị khi chỉnh sửa)
	resp.Images, _ = u.imageRepo.GetImagesByProjectID(id)
	return resp, nil
}

// UpdateProject cập nhật thông tin dự án.
func (u *managerPostUseCase) UpdateProject(id uint64, req dto.CreateProjectRequest) error {
	project, err := u.realEstateRepo.GetProjectByID(id)
	if err != nil {
		return fmt.Errorf("không tìm thấy dự án")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return fmt.Errorf("tên dự án không được để trống")
	}

	project.Name = name
	project.AlternativeName = req.AlternativeName
	project.Status = req.Status
	if project.Status == "" {
		project.Status = "active"
	}
	project.FullAddress = req.FullAddress

	// Cập nhật vị trí (mã tỉnh/phường, giữ số 0 đầu)
	project.ProvinceCode = req.ProvinceCode
	project.WardCode = req.WardCode
	project.CategoryID = req.CategoryID

	project.TotalAreaHA = req.TotalAreaHA
	project.TotalUnits = req.TotalUnits
	project.PriceMin = req.PriceMin
	project.PriceMax = req.PriceMax

	// Cập nhật mốc thời gian (xóa nếu không gửi)
	project.ConstructionStartDate = nil
	project.HandoverDate = nil
	if req.ConstructionStartDate != nil && *req.ConstructionStartDate != "" {
		if t, err := time.Parse("2006-01-02", *req.ConstructionStartDate); err == nil {
			project.ConstructionStartDate = &t
		}
	}
	if req.HandoverDate != nil && *req.HandoverDate != "" {
		if t, err := time.Parse("2006-01-02", *req.HandoverDate); err == nil {
			project.HandoverDate = &t
		}
	}

	if err := u.realEstateRepo.UpdateProject(project); err != nil {
		return err
	}

	// Liên kết ảnh dự án đã upload với dự án
	return u.imageRepo.LinkToProject(req.ImageIDs, id)
}
