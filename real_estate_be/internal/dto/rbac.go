package dto

// ── RBAC: role / permission ──────────────────────────────

// RoleResponse — 1 role kèm danh sách permission được gán
type RoleResponse struct {
	ID          uint64   `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsActive    bool     `json:"is_active"`
	Permissions []string `json:"permissions"`
}

type PermissionResponse struct {
	ID          uint64 `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Module      string `json:"module"`
	Description string `json:"description"`
}

// AdminUserResponse — 1 user trong màn quản trị, kèm các role đang giữ
type AdminUserResponse struct {
	ID        uint64   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	IsActive  bool     `json:"is_active"`
	Roles     []string `json:"roles"`
	CreatedAt string   `json:"created_at"`
}

// SetUserRolesRequest — gán lại toàn bộ role của user (một user có thể giữ nhiều role)
type SetUserRolesRequest struct {
	Roles []string `json:"roles"`
}
