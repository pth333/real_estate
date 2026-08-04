package model

import "time"

// SearchHistory lưu lịch sử tìm kiếm
type SearchHistory struct {
	ID             uint64    `gorm:"primaryKey"`
	UserID         string    `gorm:"column:user_id;type:char(10);index"`
	Query          *string   `gorm:"column:query;type:varchar(100)"`
	Province       *string   `gorm:"column:province;type:varchar(100)"`
	Ward           *string   `gorm:"column:ward;type:varchar(100)"`
	RealEstateType *int      `gorm:"column:real_estate_type"`
	PriceMin       *float64  `gorm:"column:price_min;type:bigint"`
	PriceMax       *float64  `gorm:"column:price_max;type:bigint"`
	SearchedAt     time.Time `gorm:"column:searched_at;default:CURRENT_TIMESTAMP"`
}

func (SearchHistory) TableName() string {
	return "search_history"
}
