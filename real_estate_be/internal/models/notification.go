package model

import "time"

type Notification struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    uint64    `gorm:"column:user_id;index"`
	Title     string    `gorm:"column:title"`
	Message   string    `gorm:"column:message"`
	IsRead    bool      `gorm:"column:is_read;default:false"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
