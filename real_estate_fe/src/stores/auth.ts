import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { authApi } from '@/api/auth.api'
import type { LoginRequest, RegisterRequest } from '@/types/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const email = ref<string | null>(localStorage.getItem('user_email'))
  const name = ref<string | null>(localStorage.getItem('user_name'))

  const isAuthenticated = computed(() => !!token.value)
  const userName = computed(() => name.value ?? '')
  const userEmail = computed(() => email.value ?? '')

  /** Lưu token & user vào localStorage */
  function setSession(tok: string, userEmail?: string, userName?: string) {
    token.value = tok
    email.value = userEmail ?? null
    name.value = userName ?? null
    localStorage.setItem('token', tok)
    if (userEmail) localStorage.setItem('user_email', userEmail)
    if (userName) localStorage.setItem('user_name', userName)
  }

  /* Xoá session */
  function clearSession() {
    token.value = null
    email.value = null
    name.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user_email')
    localStorage.removeItem('user_name')
  }

  /** Login */
  async function login(payload: LoginRequest) {
    const res = await authApi.login(payload)

    if (!res.success || !res.data?.token) {
      throw new Error(res.message || 'Đăng nhập thất bại')
    }

    setSession(res.data.token, payload.email)

    return res
  }

  /** Đăng ký */
  async function register(payload: RegisterRequest) {
    const res = await authApi.register(payload)

    if (!res.success) {
      throw new Error(res.message || 'Đăng ký thất bại')
    }

    return res
  }

  /** Refresh token — gọi backend lấy access token mới từ http-only cookie */
  async function refreshToken() {
    try {
      const res = await authApi.refreshToken()
      if (res.success && res.data?.token) {
        token.value = res.data.token
        localStorage.setItem('token', res.data.token)
        return true
      }
    } catch {
      // refresh thất bại → clear session
      clearSession()
    }
    return false
  }

  /** Logout */
  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // ignore
    }
    clearSession()
  }

  return {
    token,
    email,
    name,
    isAuthenticated,
    userName,
    userEmail,
    login,
    register,
    refreshToken,
    logout,
  }
})
