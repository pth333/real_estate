package routers

import (
	"real_estate_be/internal/middleware"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

// InitRealEstateRoutes — API bất động sản.
//
// ⚠️ Middleware gắn TRỰC TIẾP trên từng route, KHÔNG dùng group có middleware.
// Lý do (Fiber v2): `Group.Group(prefix, handlers...)` append handlers vào slice Handlers
// của group CHA, nên các route đăng ký sau trên group cha sẽ thừa hưởng middleware đó —
// đúng những route được ghi chú là "public" lại bị bắt đăng nhập.
//
// Phân loại:
//   - PUBLIC (khách vãng lai xem được): danh sách/tìm kiếm, chi tiết, dữ liệu danh mục,
//     dự án, gợi ý. Đây là phần nội dung chính của website.
//   - CẦN ĐĂNG NHẬP: yêu thích (dữ liệu cá nhân của từng user).
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
