/**
 * Helper dùng chung cho UI luồng đặt cọc.
 */
import type { Deposit } from '~/types/deposit'

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
 * Tuỳ chọn báo cáo theo giai đoạn:
 * - Đã check-in → báo kết quả mua / không mua
 * - Chưa check-in → điểm danh chính mình có mặt hay không
 */
export function reportOptions(deposit: Deposit, isBroker: boolean): ReportOption[] {
  if (deposit.status === 'CHECKED_IN') {
    return [
      {
        label: isBroker ? 'Khách đã mua nhà' : 'Tôi đã mua nhà',
        value: 'BOUGHT',
        description: 'Hoàn 100% tiền cọc cho khách',
      },
      {
        label: isBroker ? 'Khách không mua' : 'Tôi không mua',
        value: 'NOT_BUY',
        description: 'Hoàn tiền cọc trừ phí môi giới',
      },
    ]
  }
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
