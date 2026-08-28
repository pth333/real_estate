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

// CreateProjectRequest — yêu cầu tạo dự án từ Manager.
// province/ward là MÃ tỉnh/phường (VD "79", "27184") để map với
// danh sách dự án khi tạo tin đăng (AddressSection lọc theo mã này).
type CreateProjectRequest struct {
	Name                   string   `json:"name"`
	AlternativeName        string   `json:"alternative_name"`
	FullAddress            string   `json:"full_address"`
	ProvinceCode           string   `json:"province"`
	WardCode               string   `json:"ward"`
	Status                 string   `json:"status"`
	TotalAreaHA            *float64 `json:"total_area_ha"`
	TotalUnits             *uint32  `json:"total_units"`
	PriceMin               *float64 `json:"price_min"`
	PriceMax               *float64 `json:"price_max"`
	ConstructionStartDate  *string  `json:"construction_start_date"` // format "2006-01-02"
	HandoverDate           *string  `json:"handover_date"`           // format "2006-01-02"
	// Danh mục dự án (category type="project")
	CategoryID *int64 `json:"category_id"`
	// Id ảnh dự án đã upload (image_projects), để liên kết với dự án
	ImageIDs []uint64 `json:"image_ids"`
}

// ManagerProjectResponse — 1 dự án trong danh sách quản lý dự án.
type ManagerProjectResponse struct {
	ID              uint64   `json:"id"`
	Name            string   `json:"name"`
	Slug            string   `json:"slug"`
	AlternativeName string   `json:"alternative_name"`
	Status          string   `json:"status"`
	FullAddress     string   `json:"full_address"`
	Thumbnail       string   `json:"thumbnail"`
	TotalAreaHA     *float64 `json:"total_area_ha"`
	TotalUnits      *uint32  `json:"total_units"`
	PriceMin        *float64 `json:"price_min"`
	PriceMax        *float64 `json:"price_max"`
	CreatedAt       string   `json:"created_at"`
}

// ManagerProjectDetailResponse — chi tiết 1 dự án (dùng cho form tạo/sửa).
// ProvinceCode/WardCode là MÃ tỉnh/phường (VD "79", "27184") để select khớp.
type ManagerProjectDetailResponse struct {
	ID                    uint64   `json:"id"`
	Name                  string   `json:"name"`
	Slug                  string   `json:"slug"`
	AlternativeName       string   `json:"alternative_name"`
	Status                string   `json:"status"`
	FullAddress           string   `json:"full_address"`
	ProvinceCode          string   `json:"province"`
	WardCode              string   `json:"ward"`
	TotalAreaHA           *float64 `json:"total_area_ha"`
	TotalUnits            *uint32  `json:"total_units"`
	PriceMin              *float64 `json:"price_min"`
	PriceMax              *float64 `json:"price_max"`
	ConstructionStartDate string   `json:"construction_start_date"` // format "2006-01-02"
	HandoverDate          string   `json:"handover_date"`           // format "2006-01-02"
	CategoryID            *int64   `json:"category_id"`
	Images                []ImageResponse `json:"images"`
}
