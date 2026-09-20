package routers

import (
	"real_estate_be/internal/middleware"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

// InitAIRoutes — sinh nội dung tin đăng bằng AI. CẦN ĐĂNG NHẬP.
// (Không truyền middleware vào Router.Group — xem ghi chú ở real_estate_routes.go)
func InitAIRoutes(Router fiber.Router) {
	aiHandler, err := wire.InitializeAIHandler()
	if err != nil {
		panic(err)
	}

	aiRouter := Router.Group("/ai")
	aiRouter.Post("/generate-content", middleware.AuthMiddleware, aiHandler.GenerateContent)
}
