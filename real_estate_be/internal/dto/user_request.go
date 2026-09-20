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
	// Vai trò dùng để FE điều hướng: CUSTOMER / BROKER / ADMIN
	Role string `json:"role"`
}
