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
  }
  error?: string
}

export interface UserInfo {
  id: number
  email: string
  name: string
}
