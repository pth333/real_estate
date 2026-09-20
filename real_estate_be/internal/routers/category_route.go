package routers

import (
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

// InitCategoryRoutes — danh mục/menu hiển thị trên header.
// PUBLIC: khách vãng lai vẫn phải thấy được menu danh mục.
// (Không truyền middleware vào Router.Group — xem ghi chú ở real_estate_routes.go)
func InitCategoryRoutes(Router fiber.Router) {
	categoryController, err := wire.InitializeCategoryHandler()
	if err != nil {
		panic(err)
	}

	categoryRouter := Router.Group("/category")
	categoryRouter.Get("/", categoryController.GetAllCategories)
}
