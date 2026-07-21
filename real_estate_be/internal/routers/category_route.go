package routers

import (
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

func InitCategoryRoutes(Router fiber.Router) {
	// Category
	categoryController, err := wire.InitializeCategoryHandler()

	if err != nil {
		panic(err)
	}

	categoryRouter := Router.Group("/category")
	{
		categoryRouter.Get("/", categoryController.GetAllCategories)

	}
}
