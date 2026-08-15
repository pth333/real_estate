package model

import "time"

// ViewHistory lưu lịch sử xem chi tiết bất động sản
type ViewHistory struct {
	ID              uint64    `gorm:"primaryKey"`
	UserID          *uint64   `gorm:"column:user_id;index"`                     // NULL nếu là guest
	SessionID       *string   `gorm:"column:session_id;type:varchar(64);index"` // UUID của guest ở FE
	RealEstateID    uint64    `gorm:"column:real_estate_id;index;not null"`
	DurationSeconds int       `gorm:"column:duration_seconds;default:0;not null"`
	ViewedAt        time.Time `gorm:"column:viewed_at;default:CURRENT_TIMESTAMP"`
}

func (ViewHistory) TableName() string {
	return "view_history"
}
