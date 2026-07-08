package model

type Category struct {
	ID       int64        `gorm:"primaryKey"`
	ParentID *int64       `gorm:"column:parent_id"`
	Name     string       `gorm:"column:name;uniqueIndex"`
	Slug     string       `gorm:"column:slug;uniqueIndex"`
	Children []*Category  `gorm:"-" json:"children,omitempty"`
}
