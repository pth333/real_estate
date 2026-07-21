import axios from 'axios'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  const apiBaseUrl = config.public.apiBaseUrl

  const api = axios.create({
    baseURL: apiBaseUrl,
    headers: {
      'Content-Type': 'application/json',
    },
    withCredentials: true,
    timeout: 5000,
  })

  /** Flag ngăn vòng lặp refresh vô hạn */
  let isRefreshing = false
  /** Hàng đợi các request đang chờ refresh xong */
  let failedQueue: Array<{
    resolve: (token: string) => void
    reject: (err: unknown) => void
  }> = []

  function processQueue(error: unknown, token: string | null) {
    failedQueue.forEach((prom) => {
      if (error) {
        prom.reject(error)
      } else {
        prom.resolve(token!)
      }
    })
    failedQueue = []
  }

  api.interceptors.request.use(
    (config) => {
      // Chỉ chạy phía client
      if (import.meta.client) {
        const authStore = useAuthStore()
        if (authStore.token) {
          config.headers.Authorization = `Bearer ${authStore.token}`
        }
      }
      return config
    },
    (error) => Promise.reject(error),
  )

  api.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config

      if (error.response?.status !== 401 || originalRequest._retry) {
        return Promise.reject(error)
      }

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({
            resolve: (token: string) => {
              originalRequest.headers.Authorization = `Bearer ${token}`
              resolve(api(originalRequest))
            },
            reject,
          })
        })
      }

      isRefreshing = true
      originalRequest._retry = true

      try {
        const res = await axios.post(
          `${apiBaseUrl}/auth/refresh`,
          {},
          { withCredentials: true },
        )

        const newToken: string = res.data.data?.token
        if (!newToken) {
          throw new Error('No token in refresh response')
        }

        // Cập nhật token trong store
        const authStore = useAuthStore()
        authStore.token = newToken

        processQueue(null, newToken)

        originalRequest.headers.Authorization = `Bearer ${newToken}`
        return api(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError, null)
        const authStore = useAuthStore()
        authStore.clearSession()
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    },
  )

  return {
    provide: {
      api,
    },
  }
})
