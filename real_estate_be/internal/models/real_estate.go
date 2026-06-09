package model

import "time"

type RealEstateModel struct {
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

	PublishedAt *time.Time `gorm:"column:published_at"`

	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
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

// RealEstateEnriched lưu kết quả enrichment
type RealEstateEnriched struct {
	ID               uint64   `gorm:"primaryKey"`
	SourceURL        string   `gorm:"column:source_url;uniqueIndex;size:512"`
	TypeOfRealEstate string   `gorm:"column:type_of_real_estate;index"`
	Latitude         *float64 `gorm:"column:latitude"`
	Longitude        *float64 `gorm:"column:longitude"`
}

func (RealEstateEnriched) TableName() string {
	return "real_estate_enriched"
}
