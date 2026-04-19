package routers

import (
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

func InitDashboardRoutes(Router fiber.Router) {
	// Dashboard
	dashboardHandler, err := wire.InitializeDashboardHandler()
	if err != nil {
		panic(err)
	}
	dashRouter := Router.Group("/dashboard")
	{
		dashRouter.Get("/summary", dashboardHandler.Summary)
		dashRouter.Post("/list-real-estate", dashboardHandler.ListRealEstate)
	}
}
