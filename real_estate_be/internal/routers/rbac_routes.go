package routers

import (
	"real_estate_be/internal/middleware"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

// InitRbacRoutes — API quản trị role/permission và gán role cho user.
//
// Cùng lưu ý như deposit_routes.go: middleware gắn TRỰC TIẾP trên từng route,
// không tạo group có middleware, để tránh Fiber v2 append handler vào group cha.
func InitRbacRoutes(Router fiber.Router) {
	rbacHandler, err := wire.InitializeAdminRbacHandler()
	if err != nil {
		panic(err)
	}

	adminGroup := Router.Group("/admin")

	adminGroup.Get("/roles",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionAdminRoleList),
		rbacHandler.ListRoles)
	adminGroup.Get("/permissions",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionAdminRoleList),
		rbacHandler.ListPermissions)
	adminGroup.Get("/users",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionAdminUserList),
		rbacHandler.ListUsers)
	adminGroup.Put("/users/:id/roles",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionAdminUserAssignRole),
		rbacHandler.SetUserRoles)
}
