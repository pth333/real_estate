package model

import "time"

type User struct {
	ID            uint64    `gorm:"primaryKey"`
	Name          string    `gorm:"column:name;uniqueIndex"`
	Email         string    `gorm:"column:email;uniqueIndex"`
	Password      string    `gorm:"column:password"`
	Phone         string    `gorm:"column:phone;uniqueIndex"`
	PhoneVerified int       `gorm:"column:phone_verified;default:0"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`

	// Vai trò KHÔNG còn lưu ở đây — xem bảng nối user_roles (một user có nhiều role).
	// Phần quyền chi tiết nằm ở role_permissions → permissions.

	// ── Thông tin môi giới (dùng khi user có role BROKER) ──
	BankAccount string `gorm:"column:bank_account;size:50"`
	BankName    string `gorm:"column:bank_name;size:100"`
	// Tiền bảo lãnh — bị trừ khi môi giới bùng (NO_SHOW_BROKER)
	GuaranteeDeposit float64 `gorm:"column:guarantee_deposit;type:decimal(15,2);default:0"`
	RatingAvg        float64 `gorm:"column:rating_avg;type:decimal(3,2);default:5"`
	TotalReviews     int     `gorm:"column:total_reviews;default:0"`
	IsActive         int     `gorm:"column:is_active;default:1"`
}
