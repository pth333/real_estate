package dto

type CategoryResponse struct {
	ID       int64              `json:"id"`
	Name     string             `json:"name"`
	Slug     string             `json:"slug"`
	Type     string             `json:"type"`
	Children []CategoryResponse `json:"children,omitempty"`
}

