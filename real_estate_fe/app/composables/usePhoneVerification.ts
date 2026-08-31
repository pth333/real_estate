// Quản lý trạng thái xác thực số điện thoại (dùng cho đăng tin)
import { useState } from '#app'

const KEY = 'phone_verified'
const PHONE_KEY = 'verified_phone'
const PHONES_LIST_KEY = 'verified_phones_list'

export const usePhoneVerification = () => {
  // Sử dụng useState để đảm bảo an toàn SSR (tránh rò rỉ dữ liệu giữa các user)
  // Khởi tạo trực tiếp từ localStorage ngay khi chạy ở Client-side
  const phoneVerified = useState('phone_verified', () => {
    if (import.meta.client) {
      return localStorage.getItem(KEY) === 'true'
    }
    return false
  })

  const verifiedPhone = useState('verified_phone', () => {
    if (import.meta.client) {
      return localStorage.getItem(PHONE_KEY) || ''
    }
    return ''
  })

  const verifiedPhones = useState<string[]>('verified_phones_list', () => {
    if (import.meta.client) {
      try {
        const listStr = localStorage.getItem(PHONES_LIST_KEY)
        if (listStr) {
          return JSON.parse(listStr)
        }
        const currentPhone = localStorage.getItem(PHONE_KEY)
        if (currentPhone) {
          localStorage.setItem(PHONES_LIST_KEY, JSON.stringify([currentPhone]))
          return [currentPhone]
        }
      } catch (e) {
        console.error('Lỗi parse verified_phones_list:', e)
      }
    }
    return []
  })

  const showOTPModal = useState('show_otp_modal', () => false)

  function setPhoneVerified(phone: string) {
    if (import.meta.client) {
      localStorage.setItem(KEY, 'true')
      localStorage.setItem(PHONE_KEY, phone)

      // Thêm vào danh sách các số điện thoại đã xác thực nếu chưa có
      if (!verifiedPhones.value.includes(phone)) {
        verifiedPhones.value.push(phone)
        localStorage.setItem(PHONES_LIST_KEY, JSON.stringify(verifiedPhones.value))
      }
    }
    phoneVerified.value = true
    verifiedPhone.value = phone
  }

  function clearPhoneVerified() {
    if (import.meta.client) {
      localStorage.removeItem(KEY)
      localStorage.removeItem(PHONE_KEY)
      localStorage.removeItem(PHONES_LIST_KEY)
    }
    phoneVerified.value = false
    verifiedPhone.value = ''
    verifiedPhones.value = []
  }

  return {
    phoneVerified,
    verifiedPhone,
    verifiedPhones,
    showOTPModal,
    setPhoneVerified,
    clearPhoneVerified,
  }
}
