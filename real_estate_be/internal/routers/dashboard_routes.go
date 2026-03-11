package routers

import (
	"real_estate_be/internal/delivery/https"

	"github.com/gofiber/fiber/v2"
)

func InitDashboardRoutes(Router fiber.Router, handler *https.DashboardHandler) {
	dashRouter := Router.Group("/dashboard")
	{
		dashRouter.Get("/summary", handler.Summary)
		dashRouter.Post("/list-real-estate", handler.ListRealEstate)
	}
}
