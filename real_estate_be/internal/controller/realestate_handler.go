package controller

import (
	"sort"
	"strconv"
	"strings"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"
	"real_estate_be/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type RealEstateHandler struct {
	service usecase.IRealEstateService
	repo    repo.RealEstateRepository
}

func NewRealEstateHandler(
	service usecase.IRealEstateService,
	repo repo.RealEstateRepository,
) *RealEstateHandler {
	return &RealEstateHandler{
		service: service,
		repo:    repo,
	}
}

func (h *RealEstateHandler) List(c *fiber.Ctx) error {
	var req dto.RealEstateSearchRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	userID := h.getUserIDFromHeader(c)
	sessionID := c.Query("session_id")
	if sessionID == "" {
		sessionID = c.Get("X-Session-ID")
	}

	data, total, err := h.service.ListRealEstate(req, userID, sessionID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"total": total,
		"data":  data,
	})
}

func (h *RealEstateHandler) ListBySEOURL(c *fiber.Ctx) error {
	req := dto.RealEstateSearchRequest{}

	catParam := c.Params("category") // "nha-dat-ban-ha-noi"

	if catParam != "" {
		// Lấy tất cả categories từ DB (nên cache)
		categories, err := h.repo.GetCategory()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		// Sort giảm dần theo độ dài slug → tránh prefix conflict
		sort.Slice(categories, func(i, j int) bool {
			return len(categories[i].Slug) > len(categories[j].Slug)
		})

		for _, category := range categories {
			if !strings.HasPrefix(catParam, category.Slug) {
				continue
			}

			// if catParam == category.Slug {
			req.Slug = category.Slug
			// break
			// }

			// Tách location slug (phần sau category.Slug + "-")
			locationSlug := strings.TrimPrefix(catParam, category.Slug)
			locationSlug = strings.TrimPrefix(locationSlug, "-")

			if locationSlug != "" {

				province, err := h.repo.GetProvinceBySlug(locationSlug)

				if err != nil {
					// Slug không phải city → bỏ qua, không báo lỗi
					break
				}
				req.Filter.City = province
			}
			break
		}
	}

	segmentRaw := c.Params("*") // "gia-1-den-3-ty/dien-tich-50-100"
	if segmentRaw != "" {
		segments := strings.Split(segmentRaw, "/")
		for _, seg := range segments {
			filterRange, err := h.repo.GetFilterRangeBySlug(seg)
			if err != nil {
				continue // Segment không match → bỏ qua
			}
			switch filterRange.Type {
			case "price":
				if filterRange.MinVal != nil {
					req.Filter.MinPrice = *filterRange.MinVal
				}
				if filterRange.MaxVal != nil {
					req.Filter.MaxPrice = *filterRange.MaxVal
				}
			case "area":
				if filterRange.MinVal != nil {
					req.Filter.MinAcreage = *filterRange.MinVal
				}
				if filterRange.MaxVal != nil {
					req.Filter.MaxAcreage = *filterRange.MaxVal
				}
			}
		}
	}

	// ── Bộ lọc nâng cao (query string) ──
	applyAdvancedFilter(c, &req)

	req.Search = c.Query("search", "")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(c.Query("size", "12"))
	if size < 1 {
		size = 12
	}
	req.Page = page
	req.Size = size

	userID := h.getUserIDFromHeader(c)
	sessionID := c.Query("session_id")
	if sessionID == "" {
		sessionID = c.Get("X-Session-ID")
	}

	data, total, err := h.service.ListRealEstateByCategory(req, userID, sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"total": total,
		"data":  data,
	})
}

// applyAdvancedFilter đọc các query string của toàn bộ filter (phẳng):
func applyAdvancedFilter(c *fiber.Ctx, req *dto.RealEstateSearchRequest) {
	q := c.Request().URI().QueryArgs()

	parsePositiveInt := func(name string) *int {
		v := q.Peek(name)
		if len(v) == 0 {
			return nil
		}
		n, err := strconv.Atoi(string(v))
		if err != nil || n < 1 {
			return nil
		}
		return &n
	}
	req.Filter.Bedrooms = parsePositiveInt("bedrooms")
	req.Filter.Bathrooms = parsePositiveInt("bathrooms")

	req.Filter.HouseDirection = c.Query("house_direction")
	req.Filter.BalconyDirection = c.Query("balcony_direction")
	req.Filter.LegalDocs = c.Query("legal_docs")
	req.Filter.Interior = c.Query("interior")
}

// Detail trả về chi tiết 1 tin đăng theo ID.
func (h *RealEstateHandler) Detail(c *fiber.Ctx) error {
	raw := c.Params("id")
	id, err := strconv.ParseUint(raw, 10, 64)

	if err != nil || id == 0 {
		return response.BadRequest(c, "Mã tin đăng không hợp lệ", nil)
	}

	item, err := h.service.GetByID(id)
	if err != nil {
		return response.InternalServerError(c, "Lấy thông tin chi tiết thất bại", err.Error())
	}
	if item == nil {
		return response.Error(c, fiber.StatusNotFound, "Không tìm thấy tin đăng", nil)
	}

	return response.OK(c, item)
}

func (h *RealEstateHandler) ListTopCity(c *fiber.Ctx) error {
	// Mặc định lấy 5 thành phố có nhiều BĐS nhất
	limit := 5

	data, err := h.service.GetTopCity(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Category segment chung cho toàn khối: category đầu tiên trong menu
	categorySlug, err := h.service.GetFirstCategorySlug()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result := make([]dto.TopCityResponse, 0, len(data))
	for _, city := range data {
		result = append(result, dto.TopCityResponse{
			Name:         city.City,
			CategorySlug: categorySlug,
			CitySlug:     h.service.ToSlug(city.City),
			Count:        city.Count,
			Image:        city.Image,
		})
	}

	return response.OK(c, result)
}

func (h *RealEstateHandler) ListCity(c *fiber.Ctx) error {
	provinces, err := h.service.GetListCity()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	options := make([]dto.ProvinceResponse, len(provinces))
	for i, province := range provinces {
		options[i] = dto.ProvinceResponse{
			Name: province.Name,
			Code: province.Code,
		}
	}

	return c.JSON(fiber.Map{
		"data": options,
	})
}

func (h *RealEstateHandler) ListProject(c *fiber.Ctx) error {
	provinceCode := c.Query("province")
	wardCode := c.Query("ward")

	projects, err := h.service.GetListProject(provinceCode, wardCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	options := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		options[i] = dto.ProjectResponse{
			ID:   project.ID,
			Name: project.Name,
		}
	}

	return c.JSON(fiber.Map{
		"data": options,
	})
}

func (h *RealEstateHandler) ListWard(c *fiber.Ctx) error {
	provinceCode := c.Query("code")
	if provinceCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing district code",
		})
	}

	wards, err := h.service.GetListWard(provinceCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	options := make([]dto.ProvinceResponse, len(wards))
	for i, ward := range wards {
		options[i] = dto.ProvinceResponse{
			Name: ward.Name,
			Code: ward.Code,
		}
	}

	return c.JSON(fiber.Map{
		"data": options,
	})
}

func (h *RealEstateHandler) ListRealEstateTypes(c *fiber.Ctx) error {
	types, err := h.service.GetListRealEstateTypes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	options := make([]dto.CategoryResponse, len(types))
	for i, t := range types {
		options[i] = dto.CategoryResponse{
			ID:   t.ID,
			Name: t.Name,
		}
	}

	return response.OK(c, options)
}

// CreateRealEstate — tạo tin đăng từ payload FE (đã qua AuthMiddleware → lấy email)
func (h *RealEstateHandler) CreatePost(c *fiber.Ctx) error {
	var req dto.CreateRealEstateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Lấy email từ token đã lưu trong AuthMiddleware
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		return response.Unauthorized(c, "Unauthorized", nil)
	}

	// Tìm user theo email để lấy user_id
	user, err := h.service.GetUserByEmail(email)
	if err != nil {
		return response.Unauthorized(c, "User not found", err.Error())
	}

	id, err := h.service.CreateRealEstate(req, user.ID)
	if err != nil {
		return response.InternalServerError(c, "Create real estate failed", err.Error())
	}

	return response.Created(c, "Tạo tin đăng thành công", fiber.Map{
		"id": id,
	})
}

// UpdatePost - Cập nhật tin đăng theo ID và phân quyền sở hữu
func (h *RealEstateHandler) UpdatePost(c *fiber.Ctx) error {
	postID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || postID == 0 {
		return response.BadRequest(c, "ID bài viết không hợp lệ", err.Error())
	}

	var req dto.CreateRealEstateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Lấy email từ token đã lưu trong AuthMiddleware
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		return response.Unauthorized(c, "Unauthorized", nil)
	}

	// Tìm user theo email để lấy user_id
	user, err := h.service.GetUserByEmail(email)
	if err != nil {
		return response.Unauthorized(c, "User not found", err.Error())
	}

	err = h.service.UpdateRealEstate(postID, req, user.ID)
	if err != nil {
		return response.InternalServerError(c, "Cập nhật bài viết thất bại", err.Error())
	}

	return response.OK(c, "Cập nhật bài viết thành công")
}

func (h *RealEstateHandler) ListProjectsByProjectCategory(c *fiber.Ctx) error {
	categorySlug := c.Params("category_slug")
	if categorySlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing category slug",
		})
	}

	category, err := h.service.GetCategoryBySlug(categorySlug)
	if err != nil || category == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Không tìm thấy danh mục",
		})
	}

	projects, err := h.service.GetProjectsByCategoryID(category.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	options := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		options[i] = dto.ProjectResponse{
			ID:          project.ID,
			Name:        project.Name,
			Slug:        project.Slug,
			Status:      project.Status,
			FullAddress: project.FullAddress,
			TotalAreaHA: project.TotalAreaHA,
			TotalUnits:  project.TotalUnits,
			PriceMin:    project.PriceMin,
			PriceMax:    project.PriceMax,
			ViewCount:   uint32(project.ViewCount),
		}
	}

	return c.JSON(fiber.Map{
		"data": options,
	})
}

func (h *RealEstateHandler) IncrementProjectView(c *fiber.Ctx) error {
	raw := c.Params("id")
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid project id",
		})
	}

	err = h.service.IncrementProjectView(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Incremented project view count successfully",
	})
}

func (h *RealEstateHandler) ListFeaturedProjects(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	projects, err := h.service.GetFeaturedProjects(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	options := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		options[i] = dto.ProjectResponse{
			ID:          project.ID,
			Name:        project.Name,
			Slug:        project.Slug,
			Status:      project.Status,
			FullAddress: project.FullAddress,
			TotalAreaHA: project.TotalAreaHA,
			TotalUnits:  project.TotalUnits,
			PriceMin:    project.PriceMin,
			PriceMax:    project.PriceMax,
			ViewCount:   uint32(project.ViewCount),
		}
	}

	return c.JSON(fiber.Map{
		"data": options,
	})
}

func (h *RealEstateHandler) GetProjectDetail(c *fiber.Ctx) error {
	raw := c.Params("id")
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid project id",
		})
	}

	project, err := h.service.GetProjectByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Project not found",
		})
	}

	res := dto.ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Slug:        project.Slug,
		Status:      project.Status,
		FullAddress: project.FullAddress,
		TotalAreaHA: project.TotalAreaHA,
		TotalUnits:  project.TotalUnits,
		PriceMin:    project.PriceMin,
		PriceMax:    project.PriceMax,
		ViewCount:   uint32(project.ViewCount),
	}

	return c.JSON(fiber.Map{
		"data": res,
	})
}

// helper trích xuất UserID từ Authorization Header nếu có cho API recommend
func (h *RealEstateHandler) getUserIDFromHeader(c *fiber.Ctx) uint64 {
	return jwt.ExtractUserIDFromHeader(c.Get("Authorization"))
}

// GetRecommendations lấy danh sách gợi ý BĐS (Public)
func (h *RealEstateHandler) GetRecommendations(c *fiber.Ctx) error {
	sessionID := c.Query("session_id")
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Tự động nhận diện user_id từ Token nếu có đăng nhập
	userID := h.getUserIDFromHeader(c)

	props, err := h.service.GetRecommendations(userID, sessionID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch recommendations: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": props,
	})
}
