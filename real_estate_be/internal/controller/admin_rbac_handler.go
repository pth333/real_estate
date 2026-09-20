package controller

import (
	"strconv"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

// AdminRbacHandler — quản trị role/permission và gán role cho user.
type AdminRbacHandler struct {
	service usecase.IRbacService
}

func NewAdminRbacHandler(service usecase.IRbacService) *AdminRbacHandler {
	return &AdminRbacHandler{service: service}
}

// ListRoles — toàn bộ role kèm permission được gán.
func (h *AdminRbacHandler) ListRoles(c *fiber.Ctx) error {
	roles, err := h.service.ListRoles()
	if err != nil {
		return response.InternalServerError(c, "Lấy danh sách role thất bại", err.Error())
	}
	return response.OK(c, roles)
}

// ListPermissions — toàn bộ permission theo module.
func (h *AdminRbacHandler) ListPermissions(c *fiber.Ctx) error {
	permissions, err := h.service.ListPermissions()
	if err != nil {
		return response.InternalServerError(c, "Lấy danh sách quyền thất bại", err.Error())
	}
	return response.OK(c, permissions)
}

// ListUsers — danh sách người dùng kèm role đang giữ.
func (h *AdminRbacHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	users, total, err := h.service.ListUsers(c.Query("search", ""), page, size)
	if err != nil {
		return response.InternalServerError(c, "Lấy danh sách người dùng thất bại", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "", users, fiber.Map{"total": total, "page": page, "size": size})
}

// SetUserRoles — gán lại danh sách role cho 1 user.
func (h *AdminRbacHandler) SetUserRoles(c *fiber.Ctx) error {
	userID, err := parseIDParam(c)
	if err != nil || userID == 0 {
		return response.BadRequest(c, "ID người dùng không hợp lệ", nil)
	}

	var req dto.SetUserRolesRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Dữ liệu gửi lên không hợp lệ", err.Error())
	}

	result, err := h.service.SetUserRoles(userID, req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}
