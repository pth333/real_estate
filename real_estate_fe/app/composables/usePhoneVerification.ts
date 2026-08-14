// Quản lý trạng thái xác thực số điện thoại (dùng cho đăng tin)
import { ref } from 'vue'

const KEY = 'phone_verified'
const PHONE_KEY = 'verified_phone'
const PHONES_LIST_KEY = 'verified_phones_list'

const phoneVerified = ref(false)
const verifiedPhone = ref('')
const verifiedPhones = ref<string[]>([])
const showOTPModal = ref(false)

// Khởi tạo trạng thái từ localStorage (chỉ chạy ở client-side)
if (typeof window !== 'undefined') {
  phoneVerified.value = localStorage.getItem(KEY) === 'true'
  verifiedPhone.value = localStorage.getItem(PHONE_KEY) || ''
  try {
    const listStr = localStorage.getItem(PHONES_LIST_KEY)
    if (listStr) {
      verifiedPhones.value = JSON.parse(listStr)
    } else if (verifiedPhone.value) {
      // Nếu có sdt cũ nhưng chưa có danh sách thì tạo danh sách chứa sdt cũ
      verifiedPhones.value = [verifiedPhone.value]
      localStorage.setItem(PHONES_LIST_KEY, JSON.stringify(verifiedPhones.value))
    }
  } catch (e) {
    console.error('Lỗi parse verified_phones_list:', e)
    verifiedPhones.value = verifiedPhone.value ? [verifiedPhone.value] : []
  }
}

export const usePhoneVerification = () => {
  function setPhoneVerified(phone: string) {
    if (typeof window !== 'undefined') {
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
    if (typeof window !== 'undefined') {
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
