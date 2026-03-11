package dto

type RealEstateSearchRequest struct {
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	District string `json:"district,omitempty"`
}
