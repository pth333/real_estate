package crawler

import (
	model "real_estate_be/internal/models"
)

// ICrawler định nghĩa interface chung cho các provider crawl BĐS.
type ICrawler interface {
	CrawlPage(page int) ([]model.RealEstate, error)
	Close()
}

// CrawlerConfig là cấu hình dùng chung cho tất cả provider.
type CrawlerConfig struct {
	BaseURL     string
	Concurrency int
	BatchSize   int
}
