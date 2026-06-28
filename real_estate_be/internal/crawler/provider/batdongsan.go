// Package provider chứa các implementation cụ thể của ICrawler.
package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	model "real_estate_be/internal/models"

	"github.com/chromedp/chromedp"
)

// BatDongSanCrawler crawl dữ liệu BĐS từ batdongsan.com.vn.
type BatDongSanCrawler struct {
	allocCtx    context.Context
	allocCancel context.CancelFunc

	browserCtx    context.Context
	browserCancel context.CancelFunc
}

// BatDongSanRaw là dữ liệu thô từ DOM.
type BatDongSanRaw struct {
	Title      string
	Price      string
	Address    string
	Acreage    string
	SourceURL  string
	PostedDate string
}

func NewBatDongSanCrawler() *BatDongSanCrawler {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64)"),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, browserCancel := chromedp.NewContext(allocCtx)

	// Warm up browser
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

// CrawlPage implements ICrawler, trả về []model.RealEstateModel.
func (c *BatDongSanCrawler) CrawlPage(page int) ([]model.RealEstate, error) {
	raw, err := c.crawlRaw(page)
	if err != nil {
		return nil, err
	}

	items := make([]model.RealEstate, 0, len(raw))
	for _, r := range raw {
		items = append(items, mapBatDongSan(r))
	}
	return items, nil
}

// crawlRaw crawl dữ liệu thô từ DOM.
func (c *BatDongSanCrawler) crawlRaw(page int) ([]BatDongSanRaw, error) {
	url := buildBatDongSanURL(page)

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

// ── Mapper functions ──

func mapBatDongSan(raw BatDongSanRaw) model.RealEstate {
	price := parsePrice(raw.Price)
	area := parseArea(raw.Acreage)
	district, city := parseAddress(raw.Address)

	var pricePerM2 float64
	if price > 0 && area > 0 {
		pricePerM2 = price / area
	}

	return model.RealEstate{
		Title:            raw.Title,
		PriceVND:         price,
		Address:          raw.Address,
		District:         district,
		City:             city,
		Acreage:          area,
		PricePerM2:       pricePerM2,
		TypeOfRealEstate: "",
		Source:           "batdongsan.com.vn",
		SourceURL:        raw.SourceURL,
		CrawledAt:        time.Now(),
	}
}

func parsePrice(price string) float64 {
	price = strings.ToLower(strings.ReplaceAll(price, ",", "."))

	re := regexp.MustCompile(`([\d.]+)\s*(tỷ|triệu)?`)
	m := re.FindStringSubmatch(price)
	if len(m) < 2 {
		return 0
	}

	val, _ := strconv.ParseFloat(m[1], 64)

	switch m[2] {
	case "tỷ":
		return val * 1_000_000_000
	case "triệu":
		return val * 1_000_000
	default:
		return val
	}
}

func parseArea(area string) float64 {
	area = strings.ToLower(strings.ReplaceAll(area, ",", "."))

	re := regexp.MustCompile(`([\d.]+)`)
	m := re.FindStringSubmatch(area)
	if len(m) < 2 {
		return 0
	}

	val, _ := strconv.ParseFloat(m[1], 64)
	return val
}

func parseAddress(addr string) (district, city string) {
	parts := strings.Split(addr, ",")
	if len(parts) >= 2 {
		district = strings.TrimSpace(parts[0])
		city = strings.TrimSpace(parts[len(parts)-1])
	}
	return
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
