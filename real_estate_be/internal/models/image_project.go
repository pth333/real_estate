package model

import "time"

// ImageProject — ảnh của dự án (bảng image_projects).
// Cấu trúc giống Image nhưng liên kết với dự án qua project_id.
type ImageProject struct {
	ID       uint64 `gorm:"primaryKey"`
	Key      string `gorm:"column:key;uniqueIndex;size:255"`
	Filename string `gorm:"column:filename;size:255"`
	FileType string `gorm:"column:file_type;size:100"`
	FileSize int64  `gorm:"column:file_size"`
	URL      string `gorm:"column:url;size:500"`
	ThumbnailURL string `gorm:"column:thumbnail_url;size:500"`

	// Liên kết ảnh với dự án (nullable — chưa thuộc dự án nào khi mới upload)
	ProjectID *uint64 `gorm:"column:project_id;index"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
