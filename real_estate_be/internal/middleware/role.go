package middleware

import (
	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/response"

	"github.com/gofiber/fiber/v2"
)

// RequireRole chặn các API theo vai trò người dùng (CUSTOMER / BROKER / ADMIN).
// Phải chạy SAU AuthMiddleware để đọc được email từ token.
// Ngoài ra nạp sẵn user_id + role vào Locals cho controller dùng lại.
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email, ok := c.Locals("email").(string)
		if !ok || email == "" {
			return response.Unauthorized(c, "Unauthorized", nil)
		}

		var user model.User
		if err := global.DB.Where("email = ?", email).First(&user).Error; err != nil {
			return response.Unauthorized(c, "Không tìm thấy thông tin tài khoản", err.Error())
		}

		if !hasRole(roles, user.Role) {
			return response.Forbidden(c, "Bạn không có quyền thực hiện thao tác này", nil)
		}

		c.Locals("user_id", user.ID)
		c.Locals("role", user.Role)
		return c.Next()
	}
}

func hasRole(roles []string, role string) bool {
	for _, item := range roles {
		if item == role {
			return true
		}
	}
	return false
}
