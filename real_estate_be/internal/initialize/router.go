package initialize

import (
	"real_estate_be/internal/global"
	"real_estate_be/internal/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitRouter() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://136.85.81.118,http://soihieuland.io.vn",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	MainGroup := app.Group("/api/2026")
	{
		// ⚠️ Đặt cọc escrow đăng ký ĐẦU TIÊN.
		// Fiber v2 `Group.Group(prefix, handlers...)` append handlers vào group cha, nên các
		// nhóm bên dưới (category/upload/ai) sẽ làm MainGroup dính thêm AuthMiddleware.
		// Nhóm đặt cọc cần endpoint webhook thanh toán KHÔNG cần token, nên phải tạo group
		// lúc MainGroup còn sạch middleware.
		routers.InitDepositRoutes(MainGroup)
		// RBAC: quản trị role/permission + gán role cho user
		routers.InitRbacRoutes(MainGroup)
		// Auth
		routers.InitAuthRoutes(MainGroup)
		// Category
		routers.InitCategoryRoutes(MainGroup)
		// Real Estate
		routers.InitRealEstateRoutes(MainGroup)
		// Tracking
		routers.InitTrackingRoutes(MainGroup)
		// Notifications + SSE
		routers.InitNotificationRoutes(MainGroup)
		// Upload
		routers.InitUploadRoutes(MainGroup, global.S3Client)
		// AI
		routers.InitAIRoutes(MainGroup)
		// Manager Routes (Sử dụng Google Wire chuẩn quy hoạch)
		routers.InitManagerRoutes(MainGroup)
	}

	return app
}
