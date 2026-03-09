package model

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey"`
	Name      string    `gorm:"column:name;uniqueIndex"`
	Email     string    `gorm:"column:email;uniqueIndex"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
