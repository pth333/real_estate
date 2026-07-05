import api from './service.api'
import type { LoginRequest, RegisterRequest, AuthResponse } from '@/types/auth'

export const authApi = {
  login(payload: LoginRequest): Promise<AuthResponse> {
    return api.post('/auth/login', payload).then((res) => res.data)
  },

  register(payload: RegisterRequest): Promise<AuthResponse> {
    return api.post('/auth/register', payload).then((res) => res.data)
  },

  refreshToken(): Promise<AuthResponse> {
    return api.post('/auth/refresh', {}).then((res) => res.data)
  },

  logout(): Promise<AuthResponse> {
    return api.post('/auth/logout', {}).then((res) => res.data)
  },
}
