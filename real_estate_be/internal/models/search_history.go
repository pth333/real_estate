package model

import "time"

// SearchHistory lưu lịch sử tìm kiếm tối giản kèm BĐS kết quả
type SearchHistory struct {
	ID           uint64    `gorm:"primaryKey"`
	UserID       *uint64   `gorm:"column:user_id;index"`
	SessionID    *string   `gorm:"column:session_id;type:varchar(64);index"`
	Query        string    `gorm:"column:query;type:varchar(255);not null"`
	RealEstateID *uint64   `gorm:"column:real_estate_id;index"` // ID của BĐS xuất hiện trong kết quả tìm kiếm
	SearchedAt   time.Time `gorm:"column:searched_at;default:CURRENT_TIMESTAMP"`
}

func (SearchHistory) TableName() string {
	return "search_history"
}
