package model

import "time"

type RealEstate struct {
	ID uint64 `gorm:"primaryKey"`

	Title    string  `gorm:"column:title"`
	PriceVND float64 `gorm:"column:price_vnd"`

	Address  string `gorm:"column:address"`
	District string `gorm:"column:district;index"`
	City     string `gorm:"column:city;index"`

	Acreage    float64 `gorm:"column:acreage"`
	PricePerM2 float64 `gorm:"column:price_per_m2"`

	TypeOfRealEstate string `gorm:"column:type_of_real_estate"`

	Source    string `gorm:"column:source;index"`
	SourceURL string `gorm:"column:source_url;uniqueIndex"`

	CrawledAt time.Time `gorm:"column:crawled_at;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type DashboardSummary struct {
	TotalPosts int64   `json:"total_posts" gorm:"column:total_posts"`
	AvgPriceM2 float64 `json:"avg_price_m2" gorm:"column:avg_price_m2"`
	MaxPriceM2 float64 `json:"max_price_m2" gorm:"column:max_price_m2"`
	MinPriceM2 float64 `json:"min_price_m2" gorm:"column:min_price_m2"`
}

type DistrictStat struct {
	District   string  `json:"district"`
	TotalPosts int64   `json:"total_posts"`
	AvgPriceM2 float64 `json:"avg_price_m2"`
}
