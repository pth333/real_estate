package dto

type ImageResponse struct {
	ID           uint64 `json:"id"`
	Key          string `json:"key"`
	Filename     string `json:"file_name"`
	FileType     string `json:"file_type"`
	FileSize     int64  `json:"file_size"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
}
