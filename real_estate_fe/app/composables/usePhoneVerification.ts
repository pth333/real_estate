// Quản lý trạng thái xác thực số điện thoại (dùng cho đăng tin)
const KEY = 'phone_verified'

export const usePhoneVerification = () => {
  const phoneVerified = ref(localStorage.getItem(KEY) === 'true')

  function setPhoneVerified() {
    localStorage.setItem(KEY, 'true')
    phoneVerified.value = true
  }

  function clearPhoneVerified() {
    localStorage.removeItem(KEY)
    phoneVerified.value = false
  }

  return {
    phoneVerified,
    setPhoneVerified,
    clearPhoneVerified,
  }
}
