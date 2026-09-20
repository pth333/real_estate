package dto

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SendOTPRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" validate:"required"`
	OTP   string `json:"otp" validate:"required"`
}

type UserResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	// Một user có thể giữ NHIỀU role (bảng nối user_roles)
	Roles []string `json:"roles"`
	// Quyền chi tiết gom từ các role — FE dùng để ẩn/hiện chức năng
	Permissions []string `json:"permissions"`
}
