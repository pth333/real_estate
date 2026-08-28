package model

import "time"

// Favorite — bất động sản yêu thích của người dùng.
type Favorite struct {
	ID           uint64    `gorm:"primaryKey"`
	UserID       uint64    `gorm:"column:user_id;index"`
	RealEstateID uint64    `gorm:"column:real_estate_id;index"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}
