export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
}

export interface AuthResponse {
  success: boolean
  message?: string
  data?: {
    token?: string
    user?: UserInfo
  }
  error?: string
}

/** Kết quả gửi OTP xác thực số điện thoại */
export interface SendOtpResult {
  success: boolean
  message?: string
}

/** Kết quả xác thực OTP số điện thoại */
export interface VerifyOtpResult {
  success: boolean
  message?: string
}

export interface UserInfo {
  id: number
  email: string
  name: string
  phone?: string
  /** Vai trò: CUSTOMER / BROKER / ADMIN — dùng để điều hướng giao diện */
  role?: 'CUSTOMER' | 'BROKER' | 'ADMIN'
}
