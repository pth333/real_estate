package routers

import "github.com/gofiber/fiber/v2"

type DashboardRoutes struct {
}

func (r *DashboardRoutes) InitDashboardRoutes(Router fiber.Router) {
	dashRouter := Router.Group("/dashboard")
	{
		dashRouter.Get("/summary")
		dashRouter.Post("/list-real-estate")
	}
}
