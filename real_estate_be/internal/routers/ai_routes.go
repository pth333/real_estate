package routers

import (
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

func InitAIRoutes(Router fiber.Router) {
	aiHandler, err := wire.InitializeAIHandler()
	if err != nil {
		panic(err)
	}

	aiRouter := Router.Group("/ai")
	{
		aiRouter.Post("/generate-content", aiHandler.GenerateContent)
	}
}
