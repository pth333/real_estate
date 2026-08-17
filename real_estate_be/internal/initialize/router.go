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
		AllowOrigins:     "http://localhost:5173,http://localhost:3000,http://127.0.0.1:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	MainGroup := app.Group("/api/2026")
	{
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
	}

	return app
}
