package controller

import (
	"real_estate_be/internal/dto"
	"real_estate_be/internal/global"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"
	"real_estate_be/pkg/jwt"

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

// setRefreshCookie — ghi cookie refresh_token dùng chung cho login/refresh/logout.
// httpOnly để JS không đọc được; Secure bật/tắt theo config `server.cookie_secure`
// (chỉ bật khi chạy HTTPS, bật nhầm lúc chạy HTTP thì browser sẽ không lưu cookie).
func setRefreshCookie(c *fiber.Ctx, value string, maxAge int) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    value,
		HTTPOnly: true,
		Secure:   global.Config.Server.CookieSecure,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   maxAge,
	})
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

	accessToken, refreshToken, user, err := h.service.Login(req)
	if err != nil {
		return response.Unauthorized(c, err.Error(), err)
	}

	// Set refresh token vào http-only cookie (7 ngày)
	setRefreshCookie(c, refreshToken, 7*24*3600)

	return response.OK(c, fiber.Map{
		"token": accessToken,
		"user":  user,
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

	// Set refresh token mới vào cookie (xoay vòng token)
	setRefreshCookie(c, newRefresh, 7*24*3600)

	return response.OK(c, fiber.Map{
		"token": newAccess,
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	// Xoá cookie refresh token (MaxAge = -1)
	setRefreshCookie(c, "", -1)

	return response.OK(c, fiber.Map{
		"message": "Logged out",
	})
}

// GetUserCurrentInfo — trả thông tin user đang đăng nhập (id, tên, email, roles[], permissions[]).
// FE gọi lúc khởi động và lưu vào state global để ẩn/hiện UI + chặn ở tầng route.
func (h *UserHandler) GetUserCurrentInfo(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint64)
	if !ok || userID == 0 {
		return response.Unauthorized(c, "Unauthorized", nil)
	}

	user, err := h.service.GetUserCurrentInfo(userID)
	if err != nil {
		return response.Unauthorized(c, "Không đọc được thông tin tài khoản", err.Error())
	}

	return response.OK(c, user)
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

	// Trích xuất email từ access token (nếu người dùng đang đăng nhập)
	currentUserEmail := jwt.ExtractEmailFromHeader(c.Get("Authorization"))

	if err := h.service.VerifyOTP(req, currentUserEmail); err != nil {
		return response.Unauthorized(c, "Xác thực OTP thất bại", err.Error())
	}

	return response.OK(c, fiber.Map{
		"message": "Xác thực số điện thoại thành công",
	})
}
