package crawler

// import thêm:
import (
	"context"
	"log"
	"sync"
	"time"

	"real_estate_be/internal/crawler/mapper"
	provider "real_estate_be/internal/crawler/provider"
	"real_estate_be/internal/repo"

	"gorm.io/gorm"
)

type pageResult struct {
	page int
	data []provider.BatDongSanRaw
	err  error
}

func Run(db *gorm.DB) {
	repo := repo.NewRealEstateRepository(db)

	crawler := provider.NewBatDongSanCrawler()
	defer crawler.Close()

	const workers = 3
	const batchSize = 12 // mỗi vòng crawl 12 page song song

	startPage := 0

	for {
		ctx, cancel := context.WithCancel(context.Background())

		jobs := make(chan int, batchSize)
		results := make(chan pageResult, batchSize)

		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for p := range jobs {
					if ctx.Err() != nil {
						return
					}
					raw, err := crawler.CrawlPage(p)
					select {
					case results <- pageResult{page: p, data: raw, err: err}:
					case <-ctx.Done():
						return
					}
				}
			}()
		}

		for p := startPage; p < startPage+batchSize; p++ {
			jobs <- p
		}
		close(jobs)

		go func() {
			wg.Wait()
			close(results)
		}()

		buffer := map[int]pageResult{}
		nextPage := startPage
		stopAll := false

		for r := range results {
			buffer[r.page] = r

			for {
				cur, ok := buffer[nextPage]
				if !ok {
					break
				}
				delete(buffer, nextPage)

				log.Println("➡ Crawling page:", cur.page)

				if cur.err != nil {
					log.Println("❌ Crawl error:", cur.err)
					stopAll = true
					cancel()
					break
				}
				if len(cur.data) == 0 {
					log.Println("🛑 Empty page → STOP")
					stopAll = true
					cancel()
					break
				}

				for _, raw := range cur.data {
					if !isToday(raw.PostedDate) {
						log.Println("🛑 Reached older listing → STOP ALL")
						stopAll = true
						cancel()
						break
					}

					item := mapper.MapBatDongSan(raw)
					if err := repo.Create(&item); err != nil {
						log.Println("Insert error:", err)
					}
				}

				if stopAll {
					break
				}
				nextPage++
			}

			if stopAll {
				break
			}
		}

		cancel()
		if stopAll {
			break
		}
		startPage += batchSize
		time.Sleep(2 * time.Second)
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
