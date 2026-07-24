/**
 * $fetch wrapper có interceptor:
 * - Tự động gắn Authorization header từ auth store
 * - Tự động refresh token khi 401
 * - Retry request sau khi refresh thành công
 * - Queue các request bị 401 trong lúc đang refresh
 */

interface RequestConfig {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE"
  body?: any
  params?: Record<string, any>
  headers?: Record<string, string>
  timeout?: number
}

type QueueItem = {
  resolve: (value: unknown) => void
  reject: (err: unknown) => void
  retry: () => Promise<unknown>
}

// ✅ FIX 1: Dùng useState thay vì module-level variable → SSR-safe, scoped per request
const useRefreshState = () => useState<boolean>("api:isRefreshing", () => false)
const useFailedQueue = () => useState<QueueItem[]>("api:failedQueue", () => [])

function processQueue(error: unknown) {
  const queue = useFailedQueue()
  const items = [...queue.value]
  queue.value = []

  items.forEach(({ resolve, reject, retry }) => {
    if (error) reject(error)
    else retry().then(resolve).catch(reject)
  })
}

async function refreshAuthToken(): Promise<void> {
  const config = useRuntimeConfig()
  const res = await $fetch<{ success: boolean; data?: { token?: string } }>(
    `${config.public.apiBaseUrl}/auth/refresh`,
    { method: "POST", credentials: "include" },
  )

  if (!res.success || !res.data?.token) {
    throw new Error("Refresh token failed")
  }

  const authStore = useAuthStore()
  authStore.token = res.data.token
}

function doFetch<T>(
  url: string,
  config: RequestConfig,
  headers: Record<string, string>,
): Promise<T> {
  return $fetch<T>(url, {
    method: config.method || "GET",
    body: config.body,
    params: config.params,
    headers,
    credentials: "include",
    timeout: config.timeout ?? 15_000,
  })
}

function buildHeaders(config: RequestConfig): Record<string, string> {
  const authStore = useAuthStore()
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...config.headers,
  }
  if (authStore.token) {
    headers["Authorization"] = `Bearer ${authStore.token}`
  }
  return headers
}

function getFullUrl(url: string): string {
  if (url.startsWith("http")) return url
  const config = useRuntimeConfig()
  return `${config.public.apiBaseUrl}${url}`
}

export const api = {
  async request<T = unknown>(url: string, config: RequestConfig = {}): Promise<T> {
    const fullUrl = getFullUrl(url)
    const headers = buildHeaders(config)
    const isRefreshing = useRefreshState()
    const failedQueue = useFailedQueue()

    try {
      return await doFetch<T>(fullUrl, config, headers)
    } catch (err: any) {
      // Không xử lý 401 cho auth endpoints → tránh vòng lặp vô tận
      if (err?.response?.status !== 401 || url.includes("/auth/")) {
        throw err
      }

      // ✅ FIX 3: Queue đúng cách — resolve/reject sau khi retry() thật sự chạy xong
      if (isRefreshing.value) {
        return new Promise<T>((resolve, reject) => {
          failedQueue.value.push({
            resolve: resolve as (value: unknown) => void,
            reject,
            // retry dùng lại token mới (đã được set trong authStore bởi refreshAuthToken)
            retry: () => {
              const freshHeaders = buildHeaders(config)
              return doFetch<T>(fullUrl, config, freshHeaders)
            },
          })
        })
      }

      // Bắt đầu refresh
      isRefreshing.value = true
      try {
        await refreshAuthToken()
        processQueue(null) // ✅ Cho queue chạy lại với token mới
        // Retry chính request này
        const freshHeaders = buildHeaders(config)
        return await doFetch<T>(fullUrl, config, freshHeaders)
      } catch (refreshErr) {
        processQueue(refreshErr)
        const authStore = useAuthStore()
        authStore.clearSession()
        navigateTo("/login")
        throw refreshErr
      } finally {
        isRefreshing.value = false
      }
    }
  },

  get<T = unknown>(url: string, config?: Omit<RequestConfig, "method" | "body">) {
    return this.request<T>(url, { method: "GET", ...config })
  },

  post<T = unknown>(url: string, body?: unknown, config?: Omit<RequestConfig, "method" | "body">) {
    return this.request<T>(url, { method: "POST", body, ...config })
  },

  put<T = unknown>(url: string, body?: unknown, config?: Omit<RequestConfig, "method" | "body">) {
    return this.request<T>(url, { method: "PUT", body, ...config })
  },

  patch<T = unknown>(url: string, body?: unknown, config?: Omit<RequestConfig, "method" | "body">) {
    return this.request<T>(url, { method: "PATCH", body, ...config })
  },

  delete<T = unknown>(url: string, config?: Omit<RequestConfig, "method" | "body">) {
    return this.request<T>(url, { method: "DELETE", ...config })
  },
}

export default defineNuxtPlugin(() => {
  return {
    provide: { api },
  }
})