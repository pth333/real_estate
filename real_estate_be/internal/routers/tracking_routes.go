package routers

import (
	"real_estate_be/internal/middleware"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

// InitTrackingRoutes — theo dõi hành vi xem BĐS.
// Ghi nhận lượt xem là PUBLIC (khách vãng lai cũng được tính), còn merge session
// cần đăng nhập để biết user_id thật.
// (Không truyền middleware vào Router.Group — xem ghi chú ở real_estate_routes.go)
func InitTrackingRoutes(Router fiber.Router) {
	trackingHandler, err := wire.InitializeTrackingHandler()
	if err != nil {
		panic(err)
	}

	trackingRouter := Router.Group("/tracking")

	// Public: lưu thời gian xem BĐS (khách & thành viên đều gọi được)
	trackingRouter.Post("/view", trackingHandler.RecordView)

	// Cần đăng nhập: sáp nhập session ẩn danh vào tài khoản sau khi login
	trackingRouter.Post("/merge", middleware.AuthMiddleware, trackingHandler.MergeSession)
}
