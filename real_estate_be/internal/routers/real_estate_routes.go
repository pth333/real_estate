package routers

import (
	"real_estate_be/internal/middleware"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

// InitRealEstateRoutes — API bất động sản.

func InitRealEstateRoutes(Router fiber.Router) {
	realEstateHandler, err := wire.InitializeRealEstateHandler()
	if err != nil {
		panic(err)
	}

	realEstateRouter := Router.Group("/real-estate")

	// ── Cần đăng nhập: dữ liệu cá nhân ──
	realEstateRouter.Post("/favorite/:id", middleware.AuthMiddleware, realEstateHandler.ToggleFavorite)
	realEstateRouter.Get("/favorites", middleware.AuthMiddleware, realEstateHandler.ListFavorites)

	// ── Public: danh sách & tìm kiếm ──
	realEstateRouter.Post("/list", realEstateHandler.List)

	// ── Public: dữ liệu danh mục (menu lọc, tỉnh/phường, loại BĐS) ──
	realEstateRouter.Get("/list/top-city", realEstateHandler.ListTopCity)
	realEstateRouter.Get("/list/city", realEstateHandler.ListCity)
	realEstateRouter.Get("/list/ward", realEstateHandler.ListWard)
	realEstateRouter.Get("/list/project", realEstateHandler.ListProject)
	realEstateRouter.Get("/list/types", realEstateHandler.ListRealEstateTypes)

	// ── Public: dự án ──
	realEstateRouter.Get("/project/featured", realEstateHandler.ListFeaturedProjects)
	realEstateRouter.Get("/project-category/:category_slug", realEstateHandler.ListProjectsByProjectCategory)
	realEstateRouter.Post("/project/view/:id", realEstateHandler.IncrementProjectView)
	realEstateRouter.Get("/project/detail/:id", realEstateHandler.GetProjectDetail)
	realEstateRouter.Get("/project/:id/listings", realEstateHandler.GetRealEstateListingsByProjectID)

	// ── Public: gợi ý BĐS ──
	realEstateRouter.Get("/recommend", realEstateHandler.GetRecommendations)

	// ── Public: chi tiết (đặt trước wildcard SEO URL) ──
	realEstateRouter.Get("/detail/:id", realEstateHandler.Detail)

	// ── Public: SEO URL theo danh mục ──
	realEstateRouter.Get("/:category", realEstateHandler.ListBySEOURL)
	realEstateRouter.Get("/:category/*", realEstateHandler.ListBySEOURL)
}
