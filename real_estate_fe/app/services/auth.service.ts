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
  UserInfo,
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

  /**
   * Làm mới access token từ cookie httpOnly `refresh_token`.
   * silent: true — đây là request nền (khởi động phiên/SSE), thất bại thì store tự xử lý,
   * KHÔNG hiện toast "Invalid or expired token" cho người dùng.
   */
  refresh(): Promise<AuthResponse> {
    return this.api.post<AuthResponse>('/auth/refresh', undefined, { silent: true })
  }

  /** Đăng xuất — không cần đọc body, chỉ cần gọi thành công (silent để không toast lỗi) */
  async logout(): Promise<void> {
    await this.api.post('/auth/logout', undefined, { silent: true })
  }

  /**
   * Thông tin user đang đăng nhập: id, tên, email, roles[], permissions[].
   * FE gọi lúc khởi động để state global luôn có quyền MỚI NHẤT
   * (admin đổi role thì không cần đăng nhập lại mới thấy đúng giao diện).
   * silent: true vì đây là request nền lúc tải trang.
   */
  getUserCurrentInfo(): Promise<UserInfo> {
    return this.getData<UserInfo>('/auth/user-current-info', undefined, true)
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
