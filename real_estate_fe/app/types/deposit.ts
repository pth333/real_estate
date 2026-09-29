/**
 * Kiểu dữ liệu luồng đặt lịch escrow (khớp DTO backend internal/dto/deposit.go).
 */

export type DepositStatus =
  | 'AWAITING_PAYMENT'
  | 'PENDING'
  | 'BROKER_REJECTED'
  | 'BROKER_CONFIRMED'
  | 'CHECKED_IN'
  | 'PENDING_PURCHASE_APPROVAL'
  | 'VISITED_BOUGHT'
  | 'VISITED_NOT_BUY'
  | 'NO_SHOW_CUSTOMER'
  | 'NO_SHOW_BROKER'
  | 'DISPUTE'
  | 'REFUNDED'
  | 'COMPLETED'
  | 'CANCELLED'

export type PaymentMethod = 'VNPAY' | 'MOMO' | 'ZALOPAY'

export type DisputeResolution = 'REFUND_CUSTOMER' | 'TRANSFER_BROKER' | 'SPLIT'

/** Loại tài liệu khách phải xuất trình khi khai đã mua nhà */
export type PurchaseProof = 'PURCHASE_CONTRACT' | 'DEPOSIT_SLIP' | 'PAYMENT_SLIP'

export const PURCHASE_PROOF_OPTIONS: { label: string; value: PurchaseProof; hint: string }[] = [
  { label: 'Hợp đồng mua bán', value: 'PURCHASE_CONTRACT', hint: 'Hợp đồng mua bán đã ký giữa 2 bên' },
  { label: 'Phiếu đặt cọc mua nhà', value: 'DEPOSIT_SLIP', hint: 'Phiếu đặt cọc mua nhà có chữ ký' },
  { label: 'Biên nhận chuyển tiền', value: 'PAYMENT_SLIP', hint: 'Biên nhận/ủy nhiệm chi đã chuyển tiền mua' },
]

export const PURCHASE_PROOF_LABEL: Record<string, string> = {
  PURCHASE_CONTRACT: 'Hợp đồng mua bán',
  DEPOSIT_SLIP: 'Phiếu đặt cọc mua nhà',
  PAYMENT_SLIP: 'Biên nhận chuyển tiền',
}

/** Vai trò của người đang xem đơn đặt lịch */
export type DepositActorRole = 'CUSTOMER' | 'BROKER' | 'ADMIN'

/** Vị trí GPS của người bấm check-in — bằng chứng bổ trợ cho OTP */
export interface CheckinLocation {
  latitude: number
  longitude: number
  /** Sai số GPS (mét) — 0 nghĩa là trình duyệt không báo */
  accuracy: number
}

export interface DepositDispute {
  id: number
  deposit_id: number
  raised_by: 'CUSTOMER' | 'BROKER' | 'SYSTEM'
  reason: string
  evidence_urls: string[]
  status: 'OPEN' | 'REVIEWING' | 'RESOLVED'
  resolution: string
  resolved_by: number | null
  resolved_at: string
  evidence_deadline: string
  created_at: string
  deposit?: Deposit
}

export interface DepositRating {
  id: number
  deposit_id: number
  broker_id: number
  customer_id: number
  rating: number
  comment: string
  created_at: string
}

export interface Deposit {
  id: number
  status: DepositStatus
  customer_id: number
  customer_name: string
  customer_phone: string
  customer_email: string
  broker_id: number
  broker_name: string
  broker_phone: string
  broker_email: string

  real_estate_id: number
  real_estate_title: string
  real_estate_address: string
  real_estate_slug: string
  real_estate_thumbnail: string
  project_id: number | null

  /** Số tiền khách đã trả cho buổi xem — chính là phí môi giới (không còn tiền cọc) */
  amount: number
  broker_fee: number
  refund_amount: number | null
  penalty_amount: number | null

  viewing_date: string
  viewing_start: string
  viewing_end: string

  payment_method: PaymentMethod
  payment_ref: string
  paid_at: string

  broker_checkin: boolean
  customer_checkin: boolean
  /** Thời điểm 2 bên xác nhận có mặt (bằng chứng vị trí) */
  broker_checkin_at: string
  customer_checkin_at: string
  broker_checkin_accuracy: number | null
  customer_checkin_accuracy: number | null
  /** Khoảng cách giữa 2 bên lúc check-in (mét) — null khi thiếu vị trí 1 bên */
  checkin_distance_meters: number | null
  /** Thời điểm khách đặt cọc mua BĐS — có mốc này (trong thời gian giữ phí) thì phí được hoàn */
  purchase_deposit_at: string
  /** true khi 2 bên ở gần nhau ⇒ buổi xem đã diễn ra */
  checkin_matched: boolean
  broker_report: string
  customer_report: string
  /** Bằng chứng ảnh kèm báo cáo mua/không mua của từng bên */
  broker_report_evidence: string[]
  customer_report_evidence: string[]
  /** Loại tài liệu khách xuất trình khi khai đã mua nhà */
  customer_purchase_proof: string
  broker_confirmed_at: string
  reject_reason: string
  report_deadline: string

  can_confirm: boolean
  can_checkin: boolean
  can_report: boolean
  /** Admin còn phải duyệt tài liệu mua nhà của đơn này */
  can_approve_purchase: boolean
  has_dispute: boolean
  has_rating: boolean

  /** Kết quả admin duyệt tài liệu mua nhà */
  purchase_decision_by: number | null
  purchase_decision_at: string
  purchase_decision_note: string

  dispute?: DepositDispute
  rating?: DepositRating

  created_at: string
  updated_at: string
}

export interface CreateDepositPayload {
  real_estate_id: number
  viewing_date: string
  viewing_start: string
  viewing_end: string
  /** Họ tên người liên hệ cho buổi xem (khách xác nhận lại, có thể khác hồ sơ) */
  contact_name: string
  /** SĐT người liên hệ cho buổi xem */
  contact_phone: string
  payment_method: PaymentMethod
}

/** Phí môi giới hệ thống áp dụng theo giá BĐS (khách không sửa được) */
export interface BookingOptions {
  real_estate_id: number
  real_estate_title: string
  price_vnd: number
  broker_fee: number
  /** Nhãn phân khúc giá đang áp dụng, VD "Từ 3 đến 5 tỷ" */
  policy_label: string
  /** true khi BĐS chưa có giá nên phải dùng mức mặc định */
  is_fallback: boolean
}

export interface CreateDepositResult {
  deposit: Deposit
  payment_url: string
  gateway: string
  is_mock: boolean
}

export interface CheckinOtpResult {
  otp: string
  expires_at: string
  /** Cảnh báo khi chưa ghi nhận được vị trí môi giới tại BĐS (rỗng = đã có bằng chứng) */
  location_warning: string
}

export interface PaymentCallbackResult {
  deposit_id: number
  status: DepositStatus
  success: boolean
  message: string
}

export interface EscrowSummary {
  holding_amount: number
  refunded_amount: number
  transferred_amount: number
  penalty_amount: number
  status_counts: Record<string, number>
  open_disputes: number
}

// Nhãn + màu tag Naive UI cho từng trạng thái đặt lịch
export const DEPOSIT_STATUS_META: Record<
  DepositStatus,
  { label: string; type: 'default' | 'info' | 'success' | 'warning' | 'error' }
> = {
  AWAITING_PAYMENT: { label: 'Chờ thanh toán', type: 'warning' },
  PENDING: { label: 'Chờ môi giới xác nhận', type: 'info' },
  BROKER_REJECTED: { label: 'Môi giới từ chối', type: 'error' },
  BROKER_CONFIRMED: { label: 'Đã xác nhận lịch', type: 'info' },
  CHECKED_IN: { label: 'Đã check-in', type: 'info' },
  PENDING_PURCHASE_APPROVAL: { label: 'Chờ duyệt tài liệu mua', type: 'warning' },
  VISITED_BOUGHT: { label: 'Đã mua nhà', type: 'success' },
  VISITED_NOT_BUY: { label: 'Đã xem - không mua', type: 'success' },
  NO_SHOW_CUSTOMER: { label: 'Khách không đến', type: 'error' },
  NO_SHOW_BROKER: { label: 'Môi giới không đến', type: 'error' },
  DISPUTE: { label: 'Đang tranh chấp', type: 'warning' },
  REFUNDED: { label: 'Đã hoàn tiền', type: 'success' },
  COMPLETED: { label: 'Đã tất toán', type: 'success' },
  CANCELLED: { label: 'Đã huỷ', type: 'default' },
}

export const PAYMENT_METHOD_LABEL: Record<PaymentMethod, string> = {
  VNPAY: 'VNPay',
  MOMO: 'Momo',
  ZALOPAY: 'ZaloPay',
}
