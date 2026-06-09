// real_estate_be/internal/crawler/provider/batdongsan.go
package crawler

import (
	"context"
	"fmt"
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

type BatDongSanCrawler struct {
	allocCtx    context.Context
	allocCancel context.CancelFunc

	browserCtx    context.Context
	browserCancel context.CancelFunc
}

func NewBatDongSanCrawler() *BatDongSanCrawler {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64)"),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, browserCancel := chromedp.NewContext(allocCtx)

	// warm up browser
	_ = chromedp.Run(browserCtx)

	return &BatDongSanCrawler{
		allocCtx:      allocCtx,
		allocCancel:   allocCancel,
		browserCtx:    browserCtx,
		browserCancel: browserCancel,
	}
}

func (c *BatDongSanCrawler) Close() {
	if c.browserCancel != nil {
		c.browserCancel()
	}
	if c.allocCancel != nil {
		c.allocCancel()
	}
}

func buildBatDongSanURL(page int) string {
	base := "https://batdongsan.com.vn/nha-dat-ban-ha-noi"
	if page == 0 {
		return base
	}
	return fmt.Sprintf("%s/p%d", base, page)
}

// func CrawlBatDongSanHanoi(page int) ([]BatDongSanRaw, error) {
// 	c := NewBatDongSanCrawler()
// 	defer c.Close()
// 	return c.CrawlPage(page)
// }

func (c *BatDongSanCrawler) CrawlPage(page int) ([]BatDongSanRaw, error) {
	url := buildBatDongSanURL(page)

	// Create a dedicated tab context per page so concurrent workers do not
	// interfere with each other by sharing the same tab/session state.
	tabCtx, tabCancel := chromedp.NewContext(c.browserCtx)
	defer tabCancel()

	ctx, cancel := context.WithTimeout(tabCtx, 120*time.Second)
	defer cancel()

	var titles, prices, addresses, areas, sourceURLs, postedDates []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.js__card-title`, chromedp.ByQuery),
		chromedp.Sleep(1200*time.Millisecond),

		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('.js__card-title'))
				.filter(e => !e.closest('#listing-verified-similar, #listingVerifiedSimilar'))
				.map(e => e.innerText.trim())
		`, &titles),

		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('.re__card-config-price'))
				.filter(e => !e.closest('#listing-verified-similar, #listingVerifiedSimilar'))
				.map(e => e.innerText.trim())
		`, &prices),

		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('.re__card-location'))
				.filter(e => !e.closest('#listing-verified-similar, #listingVerifiedSimilar'))
				.map(e => {
					const clone = e.cloneNode(true);
					clone.querySelectorAll('.re__card-config-dot').forEach(dot => dot.remove());
					return clone.textContent.replace(/\s+/g, ' ').trim();
				})
		`, &addresses),

		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('.re__card-config-area'))
				.filter(e => !e.closest('#listing-verified-similar, #listingVerifiedSimilar'))
				.map(e => e.innerText.trim())
		`, &areas),

		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('a.js__product-link-for-product-id'))
				.filter(a => !a.closest('#listing-verified-similar, #listingVerifiedSimilar'))
				.map(a => "https://batdongsan.com.vn" + a.getAttribute("href"))
		`, &sourceURLs),

		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('.re__card-published-info-published-at'))
				.filter(e => !e.closest('#listing-verified-similar, #listingVerifiedSimilar'))
				.map(e => e.getAttribute("aria-label"))
		`, &postedDates),
	)
	if err != nil {
		return nil, err
	}

	minLen := min(len(titles), len(prices), len(addresses), len(areas), len(sourceURLs), len(postedDates))
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
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimSpace(s)
}
