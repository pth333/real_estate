package routers

import (
	"real_estate_be/internal/middleware"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

func InitAIRoutes(Router fiber.Router) {
	aiHandler, err := wire.InitializeAIHandler()
	if err != nil {
		panic(err)
	}

	aiRouter := Router.Group("/ai", middleware.AuthMiddleware)
	{
		aiRouter.Post("/generate-content", aiHandler.GenerateContent)
	}
}
