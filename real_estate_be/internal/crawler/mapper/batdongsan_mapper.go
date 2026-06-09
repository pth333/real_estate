package mapper

import (
	crawler "real_estate_be/internal/crawler/provider"
	model "real_estate_be/internal/models"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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

func parsePostedDate(raw string) *time.Time {
	t, err := time.Parse("02/01/2006", raw)
	if err != nil {
		return nil
	}
	return &t
}

func MapBatDongSan(raw crawler.BatDongSanRaw) model.RealEstateModel {

	price := parsePrice(raw.Price)
	area := parseArea(raw.Acreage)
	district, city := parseAddress(raw.Address)

	var pricePerM2 float64
	if price > 0 && area > 0 {
		pricePerM2 = price / area
	}

	return model.RealEstateModel{
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
		PublishedAt:      parsePostedDate(raw.PostedDate),
		CrawledAt:        time.Now(),
	}
}
