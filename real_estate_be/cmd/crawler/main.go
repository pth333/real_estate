package main

import (
	"real_estate_be/database"
	"real_estate_be/internal/config"
	"real_estate_be/internal/crawler"
)

func main() {

	// Load .env
	cfg := config.LoadConfig()

	// Connect MySQL
	db := database.Connect(cfg)

	// Migrate
	// db.AutoMigrate(&model.RealEstate{})
	// Run crawler
	crawler.Run(db)

}
