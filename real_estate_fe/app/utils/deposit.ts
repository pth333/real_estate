/**
 * Helper dùng chung cho UI luồng đặt lịch.
 */
import type { CheckinLocation, Deposit } from '~/types/deposit'

/**
 * Khung giờ xem nhà cố định 1 tiếng.
 * Khách chỉ chọn giờ bắt đầu, giờ kết thúc luôn = giờ bắt đầu + 1 tiếng.
 */
export class ViewingSlot {
  start: string
  end: string

  constructor(start: string, end: string) {
    this.start = start
    this.end = end
  }

  /** "09:00" + "10:00" → "09:00 - 10:00" */
  get label(): string {
    return `${this.start} - ${this.end}`
  }
}

/** Giờ làm việc cho phép xem nhà + độ dài mỗi khung (khung cố định 1 tiếng) */
export const VIEWING_SLOT_FIRST_HOUR = 8
export const VIEWING_SLOT_LAST_HOUR = 18
export const VIEWING_SLOT_DURATION_HOURS = 1

/**
 * Danh sách khung giờ xem nhà của một ngày.
 * - Mỗi khung cố định 1 tiếng: khách chỉ chọn giờ bắt đầu, giờ kết thúc tự suy ra.
 * - Nếu là ngày hôm nay: chỉ giữ khung bắt đầu SAU thời điểm hiện tại.
 */
export function buildViewingSlots(dateStr: string | null, now: Date = new Date()): ViewingSlot[] {
  const slots: ViewingSlot[] = []
  const isToday = dateStr === formatIsoDate(now)
  const lastStartHour = VIEWING_SLOT_LAST_HOUR - VIEWING_SLOT_DURATION_HOURS

  // Duyệt từng mốc giờ chẵn trong giờ làm việc, mỗi mốc là một khung dài 1 tiếng
  for (let hour = VIEWING_SLOT_FIRST_HOUR; hour <= lastStartHour; hour += VIEWING_SLOT_DURATION_HOURS) {
    // Khung đã bắt đầu (hoặc đang diễn ra) trong ngày hôm nay thì không cho chọn nữa
    if (isToday && hour <= now.getHours()) continue
    slots.push(new ViewingSlot(formatHour(hour), formatHour(hour + VIEWING_SLOT_DURATION_HOURS)))
  }

  return slots
}

/** Date → "2026-03-05" theo giờ local (khớp định dạng value-format của n-date-picker) */
function formatIsoDate(date: Date): string {
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${date.getFullYear()}-${month}-${day}`
}

/** 9 → "09:00" */
function formatHour(hour: number): string {
  return `${String(hour).padStart(2, '0')}:00`
}

/** Chuẩn hoá SĐT khách nhập: bỏ khoảng trắng/dấu chấm/gạch nối, +84 → 0 */
export function normalizeVietnamesePhone(phone: string): string {
  const digits = phone.trim().replace(/[\s.\-()]/g, '')
  return digits.startsWith('+84') ? `0${digits.slice(3)}` : digits
}

/** SĐT di động VN hợp lệ: 10 số, bắt đầu bằng 0 */
export function isVietnamesePhone(phone: string): boolean {
  return /^0\d{9}$/.test(normalizeVietnamesePhone(phone))
}

/**
 * Lấy vị trí GPS hiện tại của trình duyệt để làm bằng chứng check-in.
 * Trả về null khi khách từ chối quyền định vị hoặc máy không lấy được vị trí —
 * trường hợp này KHÔNG chặn check-in bằng OTP, chỉ là thiếu bằng chứng bổ trợ.
 */
export function getCurrentLocation(): Promise<CheckinLocation | null> {
  if (typeof navigator === 'undefined' || !navigator.geolocation) {
    return Promise.resolve(null)
  }

  return new Promise((resolve) => {
    navigator.geolocation.getCurrentPosition(
      (position) =>
        resolve({
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
          accuracy: position.coords.accuracy,
        }),
      () => resolve(null),
      { enableHighAccuracy: true, timeout: 10000, maximumAge: 30000 },
    )
  })
}

/** "85 m" / "1.2 km" — hiển thị khoảng cách 2 bên lúc check-in */
export function formatDistance(meters: number | null | undefined): string {
  if (meters === null || meters === undefined) return '—'
  if (meters < 1000) return `${Math.round(meters)} m`
  return `${(meters / 1000).toFixed(1)} km`
}

export interface ReportOption {
  label: string
  value: string
  description: string
}

/** Định dạng tiền chính xác theo VNĐ: 5000000 → "5.000.000 đ" */
export function formatVnd(amount: number | null | undefined): string {
  if (amount === null || amount === undefined) return '—'
  return `${amount.toLocaleString('vi-VN')} đ`
}

/** "2026-03-05" + "09:00" → "09:00 05/03/2026" */
export function formatSlot(deposit: Deposit): string {
  return `${deposit.viewing_start} - ${deposit.viewing_end} ${formatViewingDate(deposit.viewing_date)}`
}

/** "2026-03-05" → "05/03/2026" */
export function formatViewingDate(dateStr: string): string {
  if (!dateStr) return '—'
  const [year, month, day] = dateStr.split('-')
  if (!year || !month || !day) return dateStr
  return `${day}/${month}/${year}`
}

/**
 * Tuỳ chọn báo cáo ĐIỂM DANH khi 2 bên không check-in được bằng OTP/vị trí.
 * Phần mua/không mua đã bỏ: phí hoàn hay không do khách ĐẶT CỌC MUA bất động sản quyết định.
 */
export function reportOptions(deposit: Deposit, isBroker: boolean): ReportOption[] {
  if (deposit.status !== 'BROKER_CONFIRMED') return []
  return [
    {
      label: isBroker ? 'Tôi có mặt dẫn khách' : 'Tôi có mặt tại buổi xem',
      value: 'ATTENDED',
      description: 'Xác nhận bạn đã đến đúng hẹn',
    },
    {
      label: isBroker ? 'Tôi không đến được' : 'Tôi không đến',
      value: 'NO_SHOW',
      description: 'Xác nhận bạn đã không có mặt',
    },
  ]
}

/** Nhãn hiển thị của giá trị báo cáo */
export function reportLabel(value: string, checkedIn: boolean): string {
  if (!value) return '—'
  if (checkedIn) {
    if (value === 'BOUGHT') return 'Đã mua nhà'
    if (value === 'NOT_BUY') return 'Không mua'
  } else {
    if (value === 'ATTENDED') return 'Có mặt'
    if (value === 'NO_SHOW') return 'Không đến'
  }
  return value
}
