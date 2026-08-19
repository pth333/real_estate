package controller

import (
	"strconv"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type ManagerPostHandler struct {
	service  usecase.IManagerPostUseCase
	userRepo repo.IUserRepository // Sử dụng trực tiếp IUserRepository
}

func NewManagerPostHandler(useCase usecase.IManagerPostUseCase, userRepo repo.IUserRepository) *ManagerPostHandler {
	return &ManagerPostHandler{
		service:  useCase,
		userRepo: userRepo,
	}
}

// GetManagerPostsList - Lấy danh sách bài đăng của Manager hiện tại
func (h *ManagerPostHandler) GetManagerPostsList(c *fiber.Ctx) error {
	// Lấy email từ JWT Token lưu trong AuthMiddleware
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		return response.Unauthorized(c, "Unauthorized", nil)
	}

	// Lấy thông tin User tương ứng từ Repo
	user, err := h.userRepo.FindByEmail(email)
	if err != nil {
		return response.Unauthorized(c, "Không tìm thấy thông tin tài khoản", err.Error())
	}

	// Lấy các tham số lọc và phân trang từ query string
	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	res, total, err := h.service.GetManagerPostsList(user.ID, search, page, size)
	if err != nil {
		return response.InternalServerError(c, "Lấy danh sách bài đăng thất bại", err.Error())
	}

	return response.OK(c, fiber.Map{
		"posts": res,
		"total": total,
	})
}

// CreateRealEstate — tạo tin đăng từ payload FE (đã qua AuthMiddleware → lấy email)
func (h *ManagerPostHandler) CreatePost(c *fiber.Ctx) error {
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
	user, err := h.userRepo.FindByEmail(email)
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
func (h *ManagerPostHandler) UpdatePost(c *fiber.Ctx) error {
	realEstateID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || realEstateID == 0 {
		return response.BadRequest(c, "ID bài viết không hợp lệ", err.Error())
	}

	var req dto.CreateRealEstateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	err = h.service.UpdateRealEstate(realEstateID, req)
	if err != nil {
		return response.InternalServerError(c, "Cập nhật bài viết thất bại", err.Error())
	}

	return response.OK(c, "Cập nhật bài viết thành công")
}

// DeleteManagerPost - Xóa bài đăng của Manager hiện tại
func (h *ManagerPostHandler) DeleteManagerPost(c *fiber.Ctx) error {
	// Lấy ID bài viết cần xóa từ URL params
	postID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "ID bài viết không hợp lệ", err.Error())
	}

	err = h.service.DeleteManagerPost(postID)
	if err != nil {
		return response.InternalServerError(c, "Xóa bài viết thất bại", err.Error())
	}

	return response.OK(c, "Xóa bài viết thành công")
}
