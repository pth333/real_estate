package routers

import (
	"real_estate_be/internal/middleware"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

func InitTrackingRoutes(Router fiber.Router) {
	trackingHandler, err := wire.InitializeTrackingHandler()
	if err != nil {
		panic(err)
	}

	trackingRouter := Router.Group("/tracking")
	{
		// 1. Lưu thời gian xem BĐS (Public - Khách & Thành viên đều có thể gọi)
		trackingRouter.Post("/view", trackingHandler.RecordView)

		// 2. Lưu lịch sử tìm kiếm (Public)
		trackingRouter.Post("/search", trackingHandler.RecordSearch)

		// 3. Sáp nhập session (Yêu cầu đăng nhập để lấy UserID thật)
		trackingRouter.Post("/merge", middleware.AuthMiddleware, trackingHandler.MergeSession)
	}
}
