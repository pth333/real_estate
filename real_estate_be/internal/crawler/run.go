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

const (
	workerCount = 3
	batchSize   = 12
	batchDelay  = 2 * time.Second
)

type pageResult struct {
	page int
	data []provider.BatDongSanRaw
	err  error
}

// Run starts the crawl for today's listings.
// It works in batches: each batch crawls N pages concurrently,
// then processes results sequentially in page order.
// Stops when encountering an error, empty page, or old listing.
func Run(db *gorm.DB) {
	repo := repo.NewRealEstateRepository(db)

	producer, err := kafka.NewProducer()
	if err != nil {
		log.Println("Kafka producer error:", err)
		return
	}

	crawler := provider.NewBatDongSanCrawler()
	defer crawler.Close()
	defer producer.Close()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	for startPage := 0; ; startPage += batchSize {
		if processBatch(ctx, crawler, &repo, producer, startPage) {
			break
		}

		select {
		case <-ctx.Done():
			log.Println("⏹ Received shutdown signal")
			return
		case <-time.After(batchDelay):
		}
	}

	log.Println("✅ Crawl today listings DONE")
}

// processBatch crawls pages [startPage, startPage+batchSize).
// Pages are crawled concurrently but results are processed
// sequentially in page order.
func processBatch(
	ctx context.Context,
	crawler *provider.BatDongSanCrawler,
	repository *repo.RealEstateRepository,
	producer *kafka.Producer,
	startPage int,
) bool {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := crawlPages(ctx, crawler, startPage, batchSize)
	return processResults(ctx, repository, producer, results, startPage)
}

// crawlPages fans out crawling to workerCount goroutines.
// Returns a channel that emits results as workers complete (may be out of order).
func crawlPages(
	ctx context.Context,
	crawler *provider.BatDongSanCrawler,
	startPage, count int,
) <-chan pageResult {
	jobs := make(chan int, count)
	results := make(chan pageResult, count)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go crawlWorker(ctx, crawler, jobs, results, &wg)
	}

	// Dispatch page jobs
	go func() {
		for page := startPage; page < startPage+count; page++ {
			select {
			case jobs <- page:
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	// Close results when all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// crawlWorker reads page numbers from jobs, crawls them, and sends results.
func crawlWorker(
	ctx context.Context,
	crawler *provider.BatDongSanCrawler,
	jobs <-chan int,
	results chan<- pageResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for page := range jobs {
		if ctx.Err() != nil {
			return
		}

		raw, err := crawler.CrawlPage(page)
		select {
		case results <- pageResult{page: page, data: raw, err: err}:
		case <-ctx.Done():
			return
		}
	}
}

// processResults reads results from the channel and processes them
// in page order. Since results may arrive out of order (workers finish
// at different times), we buffer them until the expected next page arrives.
func processResults(
	ctx context.Context,
	repository *repo.RealEstateRepository,
	producer *kafka.Producer,
	results <-chan pageResult,
	startPage int,
) bool {
	buffer := make(map[int]pageResult)
	expectedPage := startPage

	for result := range results {
		buffer[result.page] = result

		// Drain buffer in sequential order
		for {
			res, ok := buffer[expectedPage]
			if !ok {
				break
			}
			delete(buffer, expectedPage)
			if handlePageResult(ctx, repository, producer, res) {
				return true
			}
			expectedPage++
		}
	}

	return false
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
