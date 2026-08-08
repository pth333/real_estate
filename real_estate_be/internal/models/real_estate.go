package model

import "time"

type RealEstate struct {
	ID uint64 `gorm:"primaryKey"`

	// User đăng tin (nullable cho tin từ crawler)
	UserID *uint64 `gorm:"column:user_id;index"`

	Title string `gorm:"column:title"`
	// Slug riêng cho từng listing: "{title-slug}-rs{id}" (URL trang chi tiết SEO).
	Slug     string  `gorm:"column:slug;uniqueIndex;size:255"`
	PriceVND float64 `gorm:"column:price_vnd"`

	Address  string `gorm:"column:address"`
	District string `gorm:"column:district;index"`
	City     string `gorm:"column:city;index"`

	Acreage    float64 `gorm:"column:acreage"`
	PricePerM2 float64 `gorm:"column:price_per_m2"`

	CategoryID *int64 `gorm:"column:category_id"`

	Description string `gorm:"column:description;type:text"`
	Bedrooms    *int   `gorm:"column:bedrooms"`
	Bathrooms   *int   `gorm:"column:bathrooms"`
	// Tiện ích lưu dạng JSON string: `["camera","bao_ve","pccc"]`
	Amenities string `gorm:"column:amenities;type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationship (không tạo FK constraint trong DB)
	User     *User     `gorm:"foreignKey:UserID;references:ID"`
	Category *Category `gorm:"foreignKey:CategoryID;references:ID"`
	Images   []Image   `gorm:"foreignKey:RealEstateID"`
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

// CityStat — thống kê 1 thành phố có nhiều BĐS
type CityStat struct {
	City  string `json:"city" gorm:"column:city"`
	Count int64  `json:"count" gorm:"column:total"`
	// Ảnh đại diện thành phố lấy từ bảng provinces (theo name)
	Image string `json:"image" gorm:"column:image"`
}

type Province struct {
	Code  string `gorm:"column:code;size:20" json:"code"`
	Name  string `gorm:"column:name;size:100" json:"name"`
	Image string `gorm:"column:image;type:text" json:"image"`
	Slug  string `gorm:"column:slug" json:"slug"`
}

type Ward struct {
	Code         string `gorm:"column:code;size:20" json:"code"`
	Name         string `gorm:"column:name;size:100" json:"name"`
	ProvinceCode string `gorm:"column:province_code;size:20" json:"province_code"`
}
