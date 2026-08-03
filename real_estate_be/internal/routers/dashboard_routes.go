package routers

import (
	"real_estate_be/internal/middleware"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

func InitDashboardRoutes(Router fiber.Router) {
	// Dashboard
	dashboardHandler, err := wire.InitializeDashboardHandler()
	if err != nil {
		panic(err)
	}
	dashRouter := Router.Group("/dashboard", middleware.AuthMiddleware)
	{
		dashRouter.Get("/summary", dashboardHandler.Summary)
	}
}
