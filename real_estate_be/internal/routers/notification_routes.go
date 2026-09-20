package routers

import (
	"real_estate_be/internal/middleware"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

// InitNotificationRoutes — thông báo realtime (SSE) + danh sách thông báo.
// Tất cả đều CẦN ĐĂNG NHẬP vì là dữ liệu cá nhân.
// Trước đây nhóm này không khai middleware mà "ăn" AuthMiddleware thừa hưởng từ
// group cha — sửa lại cho tường minh (xem ghi chú ở real_estate_routes.go).
func InitNotificationRoutes(Router fiber.Router) {
	notificationHandler, err := wire.InitializeNotificationHandler()
	if err != nil {
		panic(err)
	}

	r := Router.Group("/notifications")
	r.Get("/stream", middleware.AuthMiddleware, notificationHandler.Stream)
	r.Get("/", middleware.AuthMiddleware, notificationHandler.GetNotifications)
}
