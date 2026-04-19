package main

import (
	"real_estate_be/internal/crawler"
	"real_estate_be/internal/global"
	"real_estate_be/internal/initialize"
)

func main() {

	// Load config
	initialize.RunCrawler()

	// Connect MySQL
	// db := database.Connect(cfg)

	// Migrate
	// db.AutoMigrate(&model.RealEstate{})
	// Run crawler
	crawler.Run(global.DB)
}
