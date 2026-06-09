package crawler

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"real_estate_be/internal/crawler/mapper"
	provider "real_estate_be/internal/crawler/provider"
	"real_estate_be/internal/repo"
	"real_estate_be/pkg/kafka"

	"gorm.io/gorm"
)

type pageResult struct {
	page int
	data []provider.BatDongSanRaw
	err  error
}

func Run(db *gorm.DB) {
	repoInstance := repo.NewRealEstateRepository(db)

	producer, err := kafka.NewProducer()
	if err != nil {
		log.Println("Kafka producer error:", err)
		return
	}

	crawler := provider.NewBatDongSanCrawler()
	defer crawler.Close()
	defer producer.Close()

	const (
		workerCount = 3
		batchSize   = 12
	)

	// Graceful shutdown via OS signal
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	currentPage := 0
	for {
		shouldStop := crawlBatch(
			ctx,
			crawler,
			&repoInstance,
			producer,
			currentPage,
			workerCount,
			batchSize,
		)

		if shouldStop {
			break
		}

		currentPage += batchSize

		// Check nếu có signal trong lúc sleep
		select {
		case <-ctx.Done():
			log.Println("⏹ Received shutdown signal")
			return
		case <-time.After(2 * time.Second):
		}
	}

	log.Println("✅ Crawl today listings DONE")
}

func crawlBatch(
	ctx context.Context,
	crawler *provider.BatDongSanCrawler,
	repository *repo.RealEstateRepository,
	producer *kafka.Producer,
	startPage, workerCount, batchSize int,
) bool {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pageJobs := make(chan int, batchSize)
	pageResults := make(chan pageResult, batchSize)
	startCrawlWorkers(ctx, crawler, pageJobs, pageResults, workerCount)

	for page := startPage; page < startPage+batchSize; page++ {
		select {
		case pageJobs <- page:
		case <-ctx.Done():
			return true
		}
	}
	close(pageJobs)

	bufferByPage := make(map[int]pageResult, batchSize)
	expectedPage := startPage

	for result := range pageResults {
		bufferByPage[result.page] = result

		for {
			pageResult, exists := bufferByPage[expectedPage]
			if !exists {
				break
			}
			delete(bufferByPage, expectedPage)

			if stop := handlePageResult(ctx, repository, producer, pageResult); stop {
				cancel()
				return true
			}

			expectedPage++
		}
	}

	return false
}

func startCrawlWorkers(
	ctx context.Context,
	crawler *provider.BatDongSanCrawler,
	pageJobs <-chan int,
	pageResults chan<- pageResult,
	workerCount int,
) {
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range pageJobs {
				if ctx.Err() != nil {
					return
				}

				raw, err := crawler.CrawlPage(page)
				select {
				case pageResults <- pageResult{page: page, data: raw, err: err}:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(pageResults)
	}()
}

func handlePageResult(
	ctx context.Context,
	repository *repo.RealEstateRepository,
	producer *kafka.Producer,
	result pageResult,
) bool {
	if result.err != nil {
		log.Println("❌ Crawl error:", result.err)
		return true
	}

	if len(result.data) == 0 {
		log.Println("🛑 Empty page → STOP")
		return true
	}

	log.Println("➡ Crawling page:", result.page)

	for _, raw := range result.data {
		if !isToday(raw.PostedDate) {
			log.Println("🛑 Reached older listing → STOP ALL")
			return true
		}

		item := mapper.MapBatDongSan(raw)
		if err := (*repository).Create(&item); err != nil {
			log.Println("Insert error:", err)
			continue
		}

		if producer != nil {
			event := kafka.NewRealEstateCrawledEvent(item)
			if err := producer.Publish(ctx, item.SourceURL, event); err != nil {
				log.Println("Publish error:", err)
			}
		}
	}

	return false
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
