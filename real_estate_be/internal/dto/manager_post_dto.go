package dto

type ManagerPostListResponse struct {
	ID        uint64  `json:"id"`
	Title     string  `json:"title"`
	Slug      string  `json:"slug"`
	Thumbnail string  `json:"thumbnail"`
	Type      string  `json:"type"`
	Price     float64 `json:"price"`
	Unit      string  `json:"unit"`
	Area      float64 `json:"area"`
	CreatedAt string  `json:"created_at"`
}
