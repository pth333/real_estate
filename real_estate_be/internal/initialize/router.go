package initialize

import (
	"real_estate_be/internal/routers"

	"github.com/gofiber/fiber/v2"
)

func InitRouter() *fiber.App {
	app := fiber.New()

	MainGroup := app.Group("/api/2026")
	{
		// Dashboard
		routers.InitDashboardRoutes(MainGroup)
		// Auth
		routers.InitAuthRoutes(MainGroup)
		// Notifications + SSE
		routers.InitNotificationRoutes(MainGroup)
	}

	return app
}
