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

			// Bất động sản yêu thích (cần đăng nhập)
			authGroup.Post("/favorite/:id", realEstateHandler.ToggleFavorite)
			authGroup.Get("/favorites", realEstateHandler.ListFavorites)
		}

		// Route khong can auth
		realEstateRouter.Get("/list/top-city", realEstateHandler.ListTopCity)

		realEstateRouter.Get("/list/city", realEstateHandler.ListCity)

		realEstateRouter.Get("/list/ward", realEstateHandler.ListWard)

		realEstateRouter.Get("/list/project", realEstateHandler.ListProject)

		realEstateRouter.Get("/project/featured", realEstateHandler.ListFeaturedProjects)

		realEstateRouter.Get("/list/types", realEstateHandler.ListRealEstateTypes)

		//dự án
		realEstateRouter.Get("/project-category/:category_slug", realEstateHandler.ListProjectsByProjectCategory)
		realEstateRouter.Post("/project/view/:id", realEstateHandler.IncrementProjectView)
		realEstateRouter.Get("/project/detail/:id", realEstateHandler.GetProjectDetail)

		// Gợi ý BĐS (Public)
		realEstateRouter.Get("/recommend", realEstateHandler.GetRecommendations)

		// Detail đặt trước wildcard SEO URL
		realEstateRouter.Get("/detail/:id", realEstateHandler.Detail)

		//real estate SEO URL
		realEstateRouter.Get("/:category", realEstateHandler.ListBySEOURL)
		realEstateRouter.Get("/:category/*", realEstateHandler.ListBySEOURL)

	}
}
