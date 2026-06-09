package main

import (
	"real_estate_be/internal/crawler"
	"real_estate_be/internal/global"
	"real_estate_be/internal/initialize"
)

func main() {
	initialize.RunCrawler()
	crawler.Run(global.DB)
}
