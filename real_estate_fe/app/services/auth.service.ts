/**
 * AuthService — toàn bộ API xác thực người dùng (đăng nhập/đăng ký/refresh/logout/OTP).
 * Dùng qua composable useAuthService().
 */
import { BaseService, defineService } from '~/services/api'
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  SendOtpResult,
  VerifyOtpResult,
} from '~/types/auth'

export class AuthService extends BaseService {
  /**
   * Đăng nhập — trả nguyên envelope (không bóc `.data`) vì store cần đọc
   * cả `success` và `message` để quyết định báo lỗi.
   */
  login(payload: LoginRequest): Promise<AuthResponse> {
    return this.post<AuthResponse>('/auth/login', payload)
  }

  /** Đăng ký tài khoản — cũng trả nguyên envelope để store báo lỗi */
  register(payload: RegisterRequest): Promise<AuthResponse> {
    return this.post<AuthResponse>('/auth/register', payload)
  }

  /** Làm mới access token — trả nguyên envelope để store kiểm tra `success` */
  refresh(): Promise<AuthResponse> {
    return this.post<AuthResponse>('/auth/refresh')
  }

  /** Đăng xuất — không cần đọc body, chỉ cần gọi thành công */
  async logout(): Promise<void> {
    await this.post('/auth/logout')
  }

  /** Gửi OTP xác thực số điện thoại trước khi đăng tin */
  sendOtp(phone: string): Promise<SendOtpResult> {
    return this.post<SendOtpResult>('/auth/send-otp', { phone })
  }

  /** Xác thực mã OTP của số điện thoại */
  verifyOtp(phone: string, otp: string): Promise<VerifyOtpResult> {
    return this.post<VerifyOtpResult>('/auth/verify-otp', { phone, otp })
  }
}

export const useAuthService = defineService(AuthService)
