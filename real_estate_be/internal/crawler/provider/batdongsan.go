package crawler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type BatDongSanRaw struct {
	Title      string
	Price      string
	Address    string
	Acreage    string
	SourceURL  string
	PostedDate string
}

func buildBatDongSanURL(page int) string {
	base := "https://batdongsan.com.vn/nha-dat-ban-ha-noi"
	if page == 0 {
		return base
	}
	return fmt.Sprintf("%s/p%d", base, page)
}

func CrawlBatDongSanHanoi(page int) ([]BatDongSanRaw, error) {
	url := buildBatDongSanURL(page)
	// Chrome options (đỡ bị block nhẹ)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64)"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	var titles, prices, addresses, areas, sourceURLs, postedDates []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.js__card-title`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),

		chromedp.Evaluate(`Array.from(document.querySelectorAll('.js__card-title')).map(e => e.innerText.trim())`, &titles),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('.re__card-config-price')).map(e => e.innerText.trim())`, &prices),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('.re__card-location')).map(e => e.innerText.trim())`, &addresses),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('.re__card-config-area')).map(e => e.innerText.trim())`, &areas),
		chromedp.Evaluate(`
							Array.from(document.querySelectorAll("a.js__product-link-for-product-id"))
							.map(a => "https://batdongsan.com.vn" + a.getAttribute("href"))
							`, &sourceURLs),
		chromedp.Evaluate(`
							Array.from(
								document.querySelectorAll('.re__card-published-info-published-at')
							).map(e => e.getAttribute("aria-label"))
						`, &postedDates),
	)

	if err != nil {
		log.Println("❌ Crawl error:", err)
		return nil, err
	}

	minLen := min(len(titles), len(prices), len(addresses), len(areas), len(sourceURLs),
		len(postedDates))
	results := make([]BatDongSanRaw, 0, minLen)

	for i := 0; i < minLen; i++ {
		results = append(results, BatDongSanRaw{
			Title:      cleanText(titles[i]),
			Price:      cleanText(prices[i]),
			Address:    cleanText(addresses[i]),
			Acreage:    cleanText(areas[i]),
			SourceURL:  strings.TrimSpace(sourceURLs[i]),
			PostedDate: cleanText(postedDates[i]),
		})
	}

	return results, nil
}

func min(nums ...int) int {
	m := nums[0]
	for _, v := range nums {
		if v < m {
			m = v
		}
	}
	return m
}
func cleanText(s string) string {
	// bỏ xuống dòng, tab
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")

	// gom nhiều space thành 1
	s = strings.Join(strings.Fields(s), " ")

	return strings.TrimSpace(s)
}
