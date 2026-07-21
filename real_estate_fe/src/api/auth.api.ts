import api from './service.api'
import type { LoginRequest, RegisterRequest, AuthResponse } from '@/types/auth'

export const authApi = {
  login: async (payload: LoginRequest) => {
    const res = await api.request<AuthResponse>({ url: '/auth/login', data: payload, method: 'post' })
    return res.data
  },
  register: async (payload: RegisterRequest) => {
    const res = await api.request<AuthResponse>({ url: '/auth/register', data: payload, method: 'post' })
    return res.data
  },
  refreshToken: async () => {
    const res = await api.request<AuthResponse>({ url: '/auth/refresh', data: {}, method: 'post' })
    return res.data
  },
  logout: async () => {
    const res = await api.request<AuthResponse>({ url: '/auth/logout', data: {}, method: 'post' })
    return res.data
  },
}
