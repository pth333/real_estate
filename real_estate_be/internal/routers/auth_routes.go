package routers

import (
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

func InitAuthRoutes(Router fiber.Router) {
	// Auth
	authController, err := wire.InitializeAuthHandler()

	if err != nil {
		panic(err)
	}

	authRouter := Router.Group("/auth")
	{
		authRouter.Post("/register", authController.Register)
		authRouter.Post("/login", authController.Login)
	}
}
