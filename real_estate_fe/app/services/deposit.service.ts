/**
 * DepositService — gom toàn bộ API của luồng đặt cọc escrow vào 1 class.
 * Dùng qua composable useDepositService() để lấy sẵn $api từ Nuxt app.
 */
import type { api } from '~/plugins/api.client'
import type {
  BookingOptions,
  CheckinOtpResult,
  CreateDepositPayload,
  CreateDepositResult,
  Deposit,
  DepositDispute,
  DepositListResult,
  DisputeListResult,
  DisputeResolution,
  EscrowSummary,
  PaymentCallbackResult,
} from '~/types/deposit'

type ApiClient = typeof api

interface ApiEnvelope<T> {
  success: boolean
  message?: string
  data: T
  meta?: { total: number; page: number; size: number }
}

export class DepositService {
  constructor(private readonly api: ApiClient) {}

  // ── Khách hàng ────────────────────────────────────────

  /** Mức cọc + phí môi giới hệ thống đề xuất theo giá BĐS */
  async getBookingOptions(realEstateId: number): Promise<BookingOptions> {
    const res = await this.api.get<ApiEnvelope<BookingOptions>>('/deposits/booking-options', {
      params: { real_estate_id: realEstateId },
    })
    return res.data
  }

  /** Tạo đơn đặt cọc, trả về URL thanh toán của cổng */
  async createDeposit(payload: CreateDepositPayload): Promise<CreateDepositResult> {
    const res = await this.api.post<ApiEnvelope<CreateDepositResult>>('/deposits', payload)
    return res.data
  }

  /** Danh sách đơn đặt cọc của khách đang đăng nhập */
  async getMyDeposits(params: { status?: string; page?: number; size?: number }): Promise<DepositListResult> {
    const res = await this.api.get<ApiEnvelope<Deposit[]>>('/deposits/my', {
      params: { ...params },
    })
    return { items: res.data ?? [], total: res.meta?.total ?? 0 }
  }

  /** Khách nhập OTP do môi giới hiển thị để check-in */
  async checkin(depositId: number, otp: string): Promise<Deposit> {
    const res = await this.api.post<ApiEnvelope<Deposit>>(`/deposits/${depositId}/checkin`, { otp })
    return res.data
  }

  /** Khách đánh giá môi giới sau buổi xem */
  async rateBroker(depositId: number, rating: number, comment: string): Promise<void> {
    await this.api.post(`/deposits/${depositId}/rating`, { rating, comment })
  }

  // ── Môi giới ──────────────────────────────────────────

  async getBrokerDeposits(params: { status?: string; page?: number; size?: number }): Promise<DepositListResult> {
    const res = await this.api.get<ApiEnvelope<Deposit[]>>('/broker/deposits', {
      params: { ...params },
    })
    return { items: res.data ?? [], total: res.meta?.total ?? 0 }
  }

  async confirmDeposit(depositId: number): Promise<Deposit> {
    const res = await this.api.post<ApiEnvelope<Deposit>>(`/broker/deposits/${depositId}/confirm`)
    return res.data
  }

  async rejectDeposit(depositId: number, reason: string): Promise<Deposit> {
    const res = await this.api.post<ApiEnvelope<Deposit>>(`/broker/deposits/${depositId}/reject`, { reason })
    return res.data
  }

  /** Môi giới sinh OTP check-in tại chỗ (hiệu lực 10 phút) */
  async generateOtp(depositId: number): Promise<CheckinOtpResult> {
    const res = await this.api.post<ApiEnvelope<CheckinOtpResult>>(`/broker/deposits/${depositId}/otp`)
    return res.data
  }

  // ── Dùng chung 2 bên ──────────────────────────────────

  async getDeposit(depositId: number): Promise<Deposit> {
    const res = await this.api.get<ApiEnvelope<Deposit>>(`/deposits/${depositId}`)
    return res.data
  }

  /**
   * Báo cáo kết quả buổi xem (mỗi bên chỉ gửi được 1 lần).
   * - Đã check-in: BOUGHT / NOT_BUY — bắt buộc kèm ảnh bằng chứng;
   *   khách khai BOUGHT còn phải kèm purchaseProof (loại tài liệu mua bán).
   * - Chưa check-in: ATTENDED / NO_SHOW (báo về chính mình)
   */
  async submitReport(
    depositId: number,
    report: string,
    evidenceUrls: string[] = [],
    purchaseProof = '',
  ): Promise<Deposit> {
    const res = await this.api.post<ApiEnvelope<Deposit>>(`/deposits/${depositId}/report`, {
      report,
      evidence_urls: evidenceUrls,
      purchase_proof: purchaseProof,
    })
    return res.data
  }

  async openDispute(depositId: number, reason: string, evidenceUrls: string[]): Promise<DepositDispute> {
    const res = await this.api.post<ApiEnvelope<DepositDispute>>(`/deposits/${depositId}/dispute`, {
      reason,
      evidence_urls: evidenceUrls,
    })
    return res.data
  }

  async addEvidence(disputeId: number, evidenceUrls: string[]): Promise<DepositDispute> {
    const res = await this.api.post<ApiEnvelope<DepositDispute>>(`/disputes/${disputeId}/evidence`, {
      evidence_urls: evidenceUrls,
    })
    return res.data
  }

  // ── Cổng thanh toán ───────────────────────────────────

  /** Xác nhận thanh toán từ return URL / cổng mock */
  async confirmPayment(params: Record<string, string>): Promise<PaymentCallbackResult> {
    const res = await this.api.get<ApiEnvelope<PaymentCallbackResult>>('/deposits/payment/callback', {
      params,
      silent: true,
    })
    return res.data
  }

  // ── Admin ─────────────────────────────────────────────

  async getAllDeposits(params: { status?: string; page?: number; size?: number }): Promise<DepositListResult> {
    const res = await this.api.get<ApiEnvelope<Deposit[]>>('/admin/deposits', {
      params: { ...params },
    })
    return { items: res.data ?? [], total: res.meta?.total ?? 0 }
  }

  async getEscrowSummary(): Promise<EscrowSummary> {
    const res = await this.api.get<ApiEnvelope<EscrowSummary>>('/admin/escrow')
    return res.data
  }

  /**
   * Admin duyệt/từ chối tài liệu mua nhà.
   * Duyệt → hoàn 100% tiền cọc và trừ 1 căn tồn kho của dự án.
   * Từ chối → đơn chuyển sang tranh chấp để xử lý tiếp (tiền vẫn freeze).
   */
  async decidePurchase(depositId: number, approved: boolean, note: string): Promise<void> {
    await this.api.post(`/admin/deposits/${depositId}/purchase-decision`, { approved, note })
  }

  async getDisputes(params: { status?: string; page?: number; size?: number }): Promise<DisputeListResult> {
    const res = await this.api.get<ApiEnvelope<DepositDispute[]>>('/admin/disputes', {
      params: { ...params },
    })
    return { items: res.data ?? [], total: res.meta?.total ?? 0 }
  }

  async getDispute(disputeId: number): Promise<DepositDispute> {
    const res = await this.api.get<ApiEnvelope<DepositDispute>>(`/admin/disputes/${disputeId}`)
    return res.data
  }

  /** Admin ra quyết định xử lý tranh chấp và release tiền */
  async resolveDispute(
    disputeId: number,
    payload: { resolution: DisputeResolution; note: string; split_customer_percent: number },
  ): Promise<DepositDispute> {
    const res = await this.api.post<ApiEnvelope<DepositDispute>>(`/admin/disputes/${disputeId}/resolve`, payload)
    return res.data
  }
}

export function useDepositService(): DepositService {
  const { $api } = useNuxtApp()
  return new DepositService($api)
}
