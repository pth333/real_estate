package dto

type PresignRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type PresignResponse struct {
	UploadURL string `json:"upload_url"`
	Key       string `json:"key"`
	ExpiresAt string `json:"expires_at"`
}

type ConfirmUploadRequest struct {
	Key string `json:"key"`
	// Kind xác định loại record lưu: "" (mặc định) = ảnh tin đăng (bảng images),
	// "project" = ảnh dự án (bảng image_projects).
	Kind string `json:"kind"`
}

type ConfirmUploadResponse struct {
	ImageID     uint64 `json:"image_id"`
	PublicURL   string `json:"public_url"`
	Key         string `json:"key"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
}
