package crawler

import (
	"log"
	"time"

	"real_estate_be/internal/crawler/mapper"
	provider "real_estate_be/internal/crawler/provider"
	"real_estate_be/internal/repository/mysql"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	repo := mysql.NewRealEstateRepository(db)

	page := 0

	for {
		log.Println("➡ Crawling page:", page)

		rawData, err := provider.CrawlBatDongSanHanoi(page)
		if err != nil {
			log.Println("❌ Crawl error:", err)
			break
		}

		if len(rawData) == 0 {
			log.Println("🛑 Empty page → STOP")
			break
		}

		stop := false

		for _, raw := range rawData {
			if !isToday(raw.PostedDate) {
				log.Println("🛑 Reached older listing → STOP ALL")
				stop = true
				break
			}

			item := mapper.MapBatDongSan(raw)
			if err := repo.Create(&item); err != nil {
				log.Println("Insert error:", err)
			}
		}

		if stop {
			break
		}

		page++
		time.Sleep(3 * time.Second) // tránh bị block
	}

	log.Println("✅ Crawl today listings DONE")
}

func isToday(posted string) bool {
	t, err := time.Parse("02/01/2006", posted)
	if err != nil {
		return false
	}

	now := time.Now()
	return t.Year() == now.Year() &&
		t.Month() == now.Month() &&
		t.Day() == now.Day()
}
