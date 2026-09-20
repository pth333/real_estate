package middleware

import (
	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/response"

	"github.com/gofiber/fiber/v2"
)

// RequireAuth chỉ yêu cầu ĐÃ ĐĂNG NHẬP, đồng thời nạp user_id/roles/permissions
// vào Locals để handler dùng lại (không phải query thêm).
// Dùng cho endpoint như /auth/user-current-info — cần biết user là ai nhưng không chặn theo quyền.
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if _, err := loadAccess(c); err != nil {
			return err
		}
		return c.Next()
	}
}

// RequirePermission chặn API theo PERMISSION thay vì theo tên role.
// Nhờ vậy thêm role mới chỉ cần gán permission trong DB, không phải sửa code.
// User có NHIỀU role → chỉ cần 1 role chứa permission là qua.
// Phải chạy SAU AuthMiddleware để đọc được email từ token.
func RequirePermission(permissionCodes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		access, err := loadAccess(c)
		if err != nil {
			return err
		}

		if !hasAny(access.Permissions, permissionCodes) {
			return response.Forbidden(c, "Bạn không có quyền thực hiện thao tác này", nil)
		}

		c.Locals("roles", access.Roles)
		c.Locals("permissions", access.Permissions)
		return c.Next()
	}
}

// RequireRole chặn theo ROLE CODE — dùng cho trường hợp cần đúng vai trò
// (VD chỉ ADMIN được gán role), không phải theo hành động.
func RequireRole(roleCodes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		access, err := loadAccess(c)
		if err != nil {
			return err
		}

		if !hasAny(access.Roles, roleCodes) {
			return response.Forbidden(c, "Bạn không có quyền thực hiện thao tác này", nil)
		}

		c.Locals("roles", access.Roles)
		c.Locals("permissions", access.Permissions)
		return c.Next()
	}
}

// loadAccess đọc user + role/permission và nạp user_id/roles/permissions vào Locals.
func loadAccess(c *fiber.Ctx) (access *accessInfo, err error) {
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		return nil, response.Unauthorized(c, "Unauthorized", nil)
	}

	var user model.User
	if findErr := global.DB.Where("email = ?", email).First(&user).Error; findErr != nil {
		return nil, response.Unauthorized(c, "Không tìm thấy thông tin tài khoản", findErr.Error())
	}
	if user.IsActive == 0 {
		return nil, response.Forbidden(c, "Tài khoản đã bị khoá", nil)
	}

	info, queryErr := repo.NewRbacRepository(global.DB).GetUserAccess(user.ID)
	if queryErr != nil {
		return nil, response.InternalServerError(c, "Không đọc được quyền của tài khoản", queryErr.Error())
	}

	// KHÔNG set Locals("role"): user có thể giữ nhiều role nên không tồn tại "role đại diện" đúng.
	// Nơi cần phân biệt phía của đơn (khách hay môi giới) phải suy ra từ chính bản ghi đơn.
	c.Locals("user_id", user.ID)
	return &accessInfo{UserID: user.ID, Roles: info.Roles, Permissions: info.Permissions}, nil
}

type accessInfo struct {
	UserID      uint64
	Roles       []string
	Permissions []string
}

func hasAny(granted []string, required []string) bool {
	for _, want := range required {
		for _, have := range granted {
			if want == have {
				return true
			}
		}
	}
	return false
}
