package routers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	// each entity exports a constructor + Register method
	NewAuthRoutes(app, db).Register()
	NewDashboardRoutes(app, db).Register()
}
