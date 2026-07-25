package controller

import (
	"real_estate_be/internal/dto"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service usecase.AuthServiceInterface
}

func NewUserHandler(service usecase.AuthServiceInterface) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := h.service.Register(req); err != nil {
		return response.BadRequest(c, "Register failed", err.Error())
	}

	return response.Created(c, "Dang ky thanh cong", nil)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	accessToken, refreshToken, err := h.service.Login(req)
	if err != nil {
		return response.Unauthorized(c, "Login failed", err.Error())
	}

	// Set refresh token vào http-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false, // true nếu dùng HTTPS
		SameSite: "Lax",
		Path:     "/api/2026/auth",
		MaxAge:   7 * 24 * 3600, // 7 ngày
	})

	return response.OK(c, fiber.Map{
		"token": accessToken,
	})
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	// Lấy refresh token từ cookie
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return response.Unauthorized(c, "Missing refresh token", nil)
	}

	newAccess, newRefresh, err := h.service.RefreshToken(refreshToken)
	if err != nil {
		return response.Unauthorized(c, "Refresh token failed", err.Error())
	}

	// Set refresh token mới vào cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefresh,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/api/2026/auth",
		MaxAge:   7 * 24 * 3600,
	})

	return response.OK(c, fiber.Map{
		"token": newAccess,
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	// Xoá cookie refresh token
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/api/2026/auth",
		MaxAge:   -1,
	})

	return response.OK(c, fiber.Map{
		"message": "Logged out",
	})
}

func (h *UserHandler) SendOTP(c *fiber.Ctx) error {
	var req dto.SendOTPRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if req.Phone == "" {
		return response.BadRequest(c, "Số điện thoại không được để trống", nil)
	}

	if err := h.service.SendOTP(req); err != nil {
		return response.InternalServerError(c, "Gửi OTP thất bại", err.Error())
	}

	return response.OK(c, fiber.Map{
		"message": "Mã OTP đã được gửi",
	})
}

func (h *UserHandler) VerifyOTP(c *fiber.Ctx) error {
	var req dto.VerifyOTPRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := h.service.VerifyOTP(req); err != nil {
		return response.Unauthorized(c, "Xác thực OTP thất bại", err.Error())
	}

	return response.OK(c, fiber.Map{
		"message": "Xác thực số điện thoại thành công",
	})
}
