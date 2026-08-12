package model

import (
	"time"
)

type Notification struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	ListingID uint64    `gorm:"column:listing_id;uniqueIndex" json:"listing_id"` // Link tới BĐS, unique để tránh trùng
	Type      string    `gorm:"column:type" json:"type"`                         // ví dụ: "new_listing"
	Payload   string    `gorm:"column:payload;type:json" json:"payload"`         // JSON chứa title, slug, price, acreage, address
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
