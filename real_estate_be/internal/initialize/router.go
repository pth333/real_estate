package initialize

import (
	"real_estate_be/internal/routers"

	"github.com/gofiber/fiber/v2"
)

func InitRouter() *fiber.App {
	app := fiber.New()

	MainGroup := app.Group("/v1/2026")
	{
		// Dashboard
		routers.InitDashboardRoutes(MainGroup)
		//Auth
		routers.InitAuthRoutes(MainGroup)
	}

	return app
}
