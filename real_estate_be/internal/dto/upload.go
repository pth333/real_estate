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
}

type ConfirmUploadResponse struct {
	ImageID   uint64 `json:"image_id"`
	PublicURL string `json:"public_url"`
	Key       string `json:"key"`
}
