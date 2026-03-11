package initialize

import (
	"real_estate_be/internal/delivery/https"
	"real_estate_be/internal/global"
	"real_estate_be/internal/repository/mysql"
	"real_estate_be/internal/routers"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func InitRouter() *fiber.App {
	app := fiber.New()

	MainGroup := app.Group("/v1/2026")
	{
		// Dashboard
		dashboardRepo := mysql.NewDashboardRepository(global.DB)
		dashboardService := usecase.NewDashboardService(dashboardRepo)
		dashboardHandler := https.NewDashboardHandler(dashboardService)

		routers.InitDashboardRoutes(MainGroup, dashboardHandler)
	}

	return app
}
