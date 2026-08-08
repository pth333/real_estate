package routers

import (
	"real_estate_be/internal/middleware"
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
		// Apply auth middleware cho tat ca routes trong group
		authGroup := realEstateRouter.Group("/", middleware.AuthMiddleware)
		{
			authGroup.Post("/list", realEstateHandler.List)
			authGroup.Post("/create-post", realEstateHandler.CreatePost)
		}

		// Route khong can auth
		realEstateRouter.Get("/list/top-city", realEstateHandler.ListTopCity)
		realEstateRouter.Get("/list/city", realEstateHandler.ListCity)
		realEstateRouter.Get("/list/ward", realEstateHandler.ListWard)
		realEstateRouter.Get("/list/types", realEstateHandler.ListRealEstateTypes)

		// Server-driven SEO: URL là nguồn truth, backend tự parse
		// /{category} | /{category}/{location} | /{category}/{location}/{filters}
		// Đặt sau các route static /list/... để Fiber ưu tiên match static trước.
		realEstateRouter.Get("/:category", realEstateHandler.ListBySEOURL)
		realEstateRouter.Get("/:category/*", realEstateHandler.ListBySEOURL)

	}
}
