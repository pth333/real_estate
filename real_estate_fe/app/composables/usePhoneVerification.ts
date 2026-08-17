// Quản lý trạng thái xác thực số điện thoại (dùng cho đăng tin)
import { onMounted } from 'vue'
import { useState } from '#app'

const KEY = 'phone_verified'
const PHONE_KEY = 'verified_phone'
const PHONES_LIST_KEY = 'verified_phones_list'

export const usePhoneVerification = () => {
  // Sử dụng useState để đảm bảo an toàn SSR (tránh rò rỉ dữ liệu giữa các user)
  const phoneVerified = useState('phone_verified', () => false)
  const verifiedPhone = useState('verified_phone', () => '')
  const verifiedPhones = useState<string[]>('verified_phones_list', () => [])
  const showOTPModal = useState('show_otp_modal', () => false)

  // Khởi tạo trạng thái từ localStorage (chỉ chạy ở client-side khi mounted)
  onMounted(() => {
    if (import.meta.client) {
      phoneVerified.value = localStorage.getItem(KEY) === 'true'
      verifiedPhone.value = localStorage.getItem(PHONE_KEY) || ''
      try {
        const listStr = localStorage.getItem(PHONES_LIST_KEY)
        if (listStr) {
          verifiedPhones.value = JSON.parse(listStr)
        } else if (verifiedPhone.value) {
          verifiedPhones.value = [verifiedPhone.value]
          localStorage.setItem(PHONES_LIST_KEY, JSON.stringify(verifiedPhones.value))
        }
      } catch (e) {
        console.error('Lỗi parse verified_phones_list:', e)
        verifiedPhones.value = verifiedPhone.value ? [verifiedPhone.value] : []
      }
    }
  })

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
