package routers

import (
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
		trackingRouter.Post("/search", trackingHandler.RecordSearch)
	}
}
