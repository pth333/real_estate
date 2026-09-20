/**
 * DepositService — toàn bộ API của luồng đặt cọc escrow.
 * Dùng qua composable useDepositService().
 */
import { BaseService, defineService } from '~/services/api'
import type {
  BookingOptions,
  CheckinOtpResult,
  CreateDepositPayload,
  CreateDepositResult,
  Deposit,
  DepositDispute,
  DisputeResolution,
  EscrowSummary,
  PaymentCallbackResult,
} from '~/types/deposit'

export class DepositService extends BaseService {
  // ── Khách hàng ────────────────────────────────────────

  /** Mức cọc + phí môi giới hệ thống đề xuất theo giá BĐS */
  getBookingOptions(realEstateId: number): Promise<BookingOptions> {
    return this.getData<BookingOptions>('/deposits/booking-options', { real_estate_id: realEstateId })
  }

  /** Tạo đơn đặt cọc, trả về URL thanh toán của cổng */
  createDeposit(payload: CreateDepositPayload): Promise<CreateDepositResult> {
    return this.postData<CreateDepositResult>('/deposits', payload)
  }

  /** Danh sách đơn đặt cọc của khách đang đăng nhập */
  getMyDeposits(params: { status?: string; page?: number; size?: number }) {
    return this.getList<Deposit>('/deposits/my', params)
  }

  /** Khách nhập OTP do môi giới hiển thị để check-in */
  checkin(depositId: number, otp: string): Promise<Deposit> {
    return this.postData<Deposit>(`/deposits/${depositId}/checkin`, { otp })
  }

  /** Khách đánh giá môi giới sau buổi xem */
  async rateBroker(depositId: number, rating: number, comment: string): Promise<void> {
    await this.post(`/deposits/${depositId}/rating`, { rating, comment })
  }

  // ── Môi giới ──────────────────────────────────────────

  getBrokerDeposits(params: { status?: string; page?: number; size?: number }) {
    return this.getList<Deposit>('/broker/deposits', params)
  }

  confirmDeposit(depositId: number): Promise<Deposit> {
    return this.postData<Deposit>(`/broker/deposits/${depositId}/confirm`)
  }

  rejectDeposit(depositId: number, reason: string): Promise<Deposit> {
    return this.postData<Deposit>(`/broker/deposits/${depositId}/reject`, { reason })
  }

  /** Môi giới sinh OTP check-in tại chỗ (hiệu lực 10 phút) */
  generateOtp(depositId: number): Promise<CheckinOtpResult> {
    return this.postData<CheckinOtpResult>(`/broker/deposits/${depositId}/otp`)
  }

  // ── Dùng chung 2 bên ──────────────────────────────────

  getDeposit(depositId: number): Promise<Deposit> {
    return this.getData<Deposit>(`/deposits/${depositId}`)
  }

  /**
   * Báo cáo kết quả buổi xem (mỗi bên chỉ gửi được 1 lần).
   * - Đã check-in: BOUGHT / NOT_BUY — bắt buộc kèm ảnh bằng chứng;
   *   khách khai BOUGHT còn phải kèm purchaseProof (loại tài liệu mua bán).
   * - Chưa check-in: ATTENDED / NO_SHOW (báo về chính mình)
   */
  submitReport(
    depositId: number,
    report: string,
    evidenceUrls: string[] = [],
    purchaseProof = '',
  ): Promise<Deposit> {
    return this.postData<Deposit>(`/deposits/${depositId}/report`, {
      report,
      evidence_urls: evidenceUrls,
      purchase_proof: purchaseProof,
    })
  }

  openDispute(depositId: number, reason: string, evidenceUrls: string[]): Promise<DepositDispute> {
    return this.postData<DepositDispute>(`/deposits/${depositId}/dispute`, {
      reason,
      evidence_urls: evidenceUrls,
    })
  }

  addEvidence(disputeId: number, evidenceUrls: string[]): Promise<DepositDispute> {
    return this.postData<DepositDispute>(`/disputes/${disputeId}/evidence`, {
      evidence_urls: evidenceUrls,
    })
  }

  // ── Cổng thanh toán ───────────────────────────────────

  /** Xác nhận thanh toán từ return URL / cổng mock (silent: trang kết quả tự hiển thị lỗi) */
  confirmPayment(params: Record<string, string>): Promise<PaymentCallbackResult> {
    return this.getData<PaymentCallbackResult>('/deposits/payment/callback', params, true)
  }

  // ── Admin ─────────────────────────────────────────────

  getAllDeposits(params: { status?: string; page?: number; size?: number }) {
    return this.getList<Deposit>('/admin/deposits', params)
  }

  getEscrowSummary(): Promise<EscrowSummary> {
    return this.getData<EscrowSummary>('/admin/escrow')
  }

  /**
   * Admin duyệt/từ chối tài liệu mua nhà.
   * Duyệt → hoàn 100% tiền cọc và trừ 1 căn tồn kho của dự án.
   * Từ chối → đơn chuyển sang tranh chấp để xử lý tiếp (tiền vẫn freeze).
   */
  async decidePurchase(depositId: number, approved: boolean, note: string): Promise<void> {
    await this.post(`/admin/deposits/${depositId}/purchase-decision`, { approved, note })
  }

  getDisputes(params: { status?: string; page?: number; size?: number }) {
    return this.getList<DepositDispute>('/admin/disputes', params)
  }

  getDispute(disputeId: number): Promise<DepositDispute> {
    return this.getData<DepositDispute>(`/admin/disputes/${disputeId}`)
  }

  /** Admin ra quyết định xử lý tranh chấp và release tiền */
  resolveDispute(
    disputeId: number,
    payload: { resolution: DisputeResolution; note: string; split_customer_percent: number },
  ): Promise<DepositDispute> {
    return this.postData<DepositDispute>(`/admin/disputes/${disputeId}/resolve`, payload)
  }
}

export const useDepositService = defineService(DepositService)

