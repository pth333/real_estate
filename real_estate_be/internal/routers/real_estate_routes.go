package routers

import (
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

func InitRealEstateRoutes(Router fiber.Router) {
	realEstateHandler, err := wire.InitializeRealEstateHandler()
	if err != nil {
		panic(err)
	}

	realEstateRouter := Router.Group("/real-estate")
	{
		realEstateRouter.Post("/list", realEstateHandler.List)
		realEstateRouter.Post("/:slug/:page", realEstateHandler.ListRealEsateByCategory)
	}
}
