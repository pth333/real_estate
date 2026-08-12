package routers

import (
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

func InitNotificationRoutes(Router fiber.Router) {
	notificationHandler, err := wire.InitializeNotificationHandler()
	if err != nil {
		panic(err)
	}

	r := Router.Group("/notifications")
	r.Get("/stream", notificationHandler.Stream)
	r.Get("/", notificationHandler.GetNotifications)
}
