package crawler

import (
	"context"
	"log"
	"sync"
	"time"

	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
	kafkapkg "real_estate_be/pkg/kafka"

	"gorm.io/gorm"
)

// pageResult chứa kết quả crawl của 1 trang.
type pageResult struct {
	page int
	data []model.RealEstate
	err  error
}

const (
	workers   = 3
	batchSize = 12
)

// Run thực hiện crawl theo batch, mỗi batch gồm N trang chạy song song.
// Kết quả được xử lý theo thứ tự trang. Dừng khi hết tin hôm nay.
func Run(ctx context.Context, crawl ICrawler, db *gorm.DB, producer *kafkapkg.Producer) {
	re := repo.NewRealEstateRepository(db)

	for page := 0; ; page += batchSize {
		if ctx.Err() != nil {
			log.Println("⏹ Crawler stopped by context")
			return
		}

		// Crawl batch N trang song song
		pages := crawlBatch(ctx, crawl, page, batchSize)

		// Xử lý kết quả theo thứ tự trang (0, 1, 2...)
		more := processPages(ctx, pages, re, producer)
		if !more {
			break
		}

		time.Sleep(2 * time.Second)
	}

	log.Println("✅ Crawl today listings DONE")
}

// crawlBatch crawl các trang từ start đến start+count, trả về map[page]result.
func crawlBatch(ctx context.Context, crawl ICrawler, start, count int) map[int]pageResult {
	ctxInner, cancel := context.WithCancel(ctx)
	defer cancel()

	// Worker pool
	jobs := make(chan int, count)
	results := make(chan pageResult, count)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				if ctxInner.Err() != nil {
					return
				}
				data, err := crawl.CrawlPage(p)
				select {
				case results <- pageResult{page: p, data: data, err: err}:
				case <-ctxInner.Done():
					return
				}
			}
		}()
	}

	// Gửi job
	for p := start; p < start+count; p++ {
		jobs <- p
	}
	close(jobs)

	// Chờ worker xong rồi đóng results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Gom kết quả vào map
	out := make(map[int]pageResult, count)
	for r := range results {
		out[r.page] = r
	}
	return out
}

// processPages xử lý kết quả theo thứ tự trang.
// Trả về false nếu cần dừng (lỗi, hết dữ liệu, gặp tin cũ).
func processPages(ctx context.Context, pages map[int]pageResult, re repo.RealEstateRepository, producer *kafkapkg.Producer) bool {
	for p := 0; ; p++ {
		res, ok := pages[p]
		if !ok {
			// Đã xử lý hết pages có kết quả
			break
		}

		log.Println("➡ Processing page:", p)

		if res.err != nil {
			log.Println("❌ Crawl error:", res.err)
			return false
		}
		if len(res.data) == 0 {
			log.Println("🛑 Empty page → STOP")
			return false
		}

		if !processItems(ctx, res.data, re, producer) {
			return false
		}
	}
	return true
}

// processItems lưu từng item vào DB và publish Kafka.
// Trả về false nếu gặp item cũ (cần dừng).
func processItems(ctx context.Context, items []model.RealEstate, re repo.RealEstateRepository, producer *kafkapkg.Producer) bool {
	for i := range items {
		if !isTodayFromPosted(&items[i]) {
			log.Println("🛑 Reached older listing → STOP ALL")
			return false
		}

		// Lưu DB (UPSERT)
		if err := re.Create(&items[i]); err != nil {
			log.Println("Insert error:", err)
			continue
		}

		// Publish Kafka
		if producer != nil {
			event := kafkapkg.NewRealEstateCrawledEvent(items[i])
			if err := producer.Publish(ctx, items[i].SourceURL, event); err != nil {
				log.Printf("⚠️ Kafka publish error: %v", err)
			}
		}
	}
	return true
}

// isTodayFromPosted kiểm tra bài đăng có phải hôm nay không.
// Tạm thời luôn true vì model chưa lưu postedDate.
func isTodayFromPosted(_ *model.RealEstate) bool {
	return true
}

// isToday kiểm tra chuỗi ngày có phải hôm nay không.
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

// Scheduler chạy crawl định kỳ.
type Scheduler struct {
	Interval time.Duration
	Crawler  ICrawler
	DB       *gorm.DB
	Producer *kafkapkg.Producer
}

func NewScheduler(interval time.Duration, crawl ICrawler, db *gorm.DB, producer *kafkapkg.Producer) *Scheduler {
	return &Scheduler{
		Interval: interval,
		Crawler:  crawl,
		DB:       db,
		Producer: producer,
	}
}

// Start bắt đầu scheduler, block cho đến khi ctx kết thúc.
func (s *Scheduler) Start(ctx context.Context) {
	log.Printf("⏰ [Crawler Scheduler] starting, interval=%v", s.Interval)

	// Chạy ngay lần đầu
	Run(ctx, s.Crawler, s.DB, s.Producer)

	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("⏹ Crawler scheduler stopped")
			return
		case <-ticker.C:
			log.Println("⏰ [Crawler Scheduler] wake up, starting crawl...")
			Run(ctx, s.Crawler, s.DB, s.Producer)
		}
	}
}
