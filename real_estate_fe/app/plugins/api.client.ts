/**
 * $fetch wrapper có interceptor:
 * - Tự động gắn Authorization header từ auth store
 * - Tự động refresh token khi 401
 * - Queue các request bị 401 trong lúc đang refresh
 * - Hiển thị lỗi global qua naive-ui
 */

interface RequestConfig {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  body?: any;
  params?: Record<string, any>;
  headers?: Record<string, string>;
  timeout?: number;
  /** Nếu true → không hiển thị toast lỗi */
  silent?: boolean;
}

type QueueItem = {
  resolve: (value: unknown) => void;
  reject: (err: unknown) => void;
  retry: () => Promise<unknown>;
};

// ── Refresh token ────────────────────────────────────────
const useRefreshState = () =>
  useState<boolean>("api:isRefreshing", () => false);
const useFailedQueue = () => useState<QueueItem[]>("api:failedQueue", () => []);

function processQueue(error: unknown) {
  const queue = useFailedQueue();
  const items = [...queue.value];
  queue.value = [];
  items.forEach(({ resolve, reject, retry }) => {
    if (error) reject(error);
    else retry().then(resolve).catch(reject);
  });
}

function buildHeaders(
  config: RequestConfig,
  token?: string | null,
): Record<string, string> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...config.headers,
  };
  if (token) headers["Authorization"] = `Bearer ${token}`;
  return headers;
}

function getFullUrl(url: string): string {
  if (url.startsWith("http")) return url;
  const config = useRuntimeConfig();
  return `${config.public.apiBaseUrl}${url}`;
}

// ── API instance ─────────────────────────────────────────
export const api = {
  async request<T = unknown>(
    url: string,
    config: RequestConfig = {},
  ): Promise<T> {
    const authStore = useAuthStore();
    const fullUrl = getFullUrl(url);
    const token = authStore.token ?? undefined;
    const headers = buildHeaders(config, token);
    const isRefreshing = useRefreshState();
    const failedQueue = useFailedQueue();

    try {
      const result = await $fetch<T>(fullUrl, {
        method: config.method || "GET",
        body: config.body,
        params: config.params,
        headers,
        credentials: "include",
        timeout: config.timeout ?? 15_000,
      });

      // Kiểm tra business error: API trả về { status: false, message: "..." }
      if (result && typeof result === "object" && "status" in (result as any)) {
        const r = result as Record<string, any>;
        if (r.status === false) {
          const msg = r.message || r.error || "Yêu cầu thất bại";
          if (!config.silent) window.message?.warning(msg);
          const bizErr: any = new Error(msg);
          bizErr.__business = true;
          bizErr.data = result;
          throw bizErr;
        }
      }

      return result;
    } catch (err: any) {
      // Lỗi business đã xử lý ở trên → chỉ throw
      if (err?.__business) throw err;

      // 401 → refresh token
      if (err?.response?.status === 401 && !url.includes("/auth/")) {
        if (isRefreshing.value) {
          return new Promise<T>((resolve, reject) => {
            failedQueue.value.push({
              resolve: resolve as (value: unknown) => void,
              reject,
              retry: () => {
                return this.request<T>(url, config);
              },
            });
          });
        }

        isRefreshing.value = true;
        try {
          // Gọi refresh token qua store
          await authStore.refreshToken();

          // Retry tất cả request trong queue
          processQueue(null);

          // Retry request hiện tại với token mới
          const freshHeaders = buildHeaders(config, authStore.token);
          return await $fetch<T>(fullUrl, {
            method: config.method || "GET",
            body: config.body,
            params: config.params,
            headers: freshHeaders,
            credentials: "include",
            timeout: config.timeout ?? 15_000,
          });
        } catch (refreshErr) {
          processQueue(refreshErr);
          window.message?.warning(
            "Phiên đăng nhập hết hạn, vui lòng đăng nhập lại",
          );
          navigateTo("/login");
          throw refreshErr;
        } finally {
          isRefreshing.value = false;
        }
      }

      // HTTP / network error
      if (!config.silent) {
        const msg =
          err?.data?.message ||
          err?.data?.error ||
          err?.message ||
          "Có lỗi xảy ra";
        const status = err?.response?.status;
        if (status >= 500) window.message?.error(msg);
        else if (status >= 400) window.message?.warning(msg);
        else window.message?.error(msg);
      }

      throw err;
    }
  },

  get<T = unknown>(
    url: string,
    config?: Omit<RequestConfig, "method" | "body">,
  ) {
    return this.request<T>(url, { method: "GET", ...config });
  },

  post<T = unknown>(
    url: string,
    body?: unknown,
    config?: Omit<RequestConfig, "method" | "body">,
  ) {
    return this.request<T>(url, { method: "POST", body, ...config });
  },

  put<T = unknown>(
    url: string,
    body?: unknown,
    config?: Omit<RequestConfig, "method" | "body">,
  ) {
    return this.request<T>(url, { method: "PUT", body, ...config });
  },

  patch<T = unknown>(
    url: string,
    body?: unknown,
    config?: Omit<RequestConfig, "method" | "body">,
  ) {
    return this.request<T>(url, { method: "PATCH", body, ...config });
  },

  delete<T = unknown>(
    url: string,
    config?: Omit<RequestConfig, "method" | "body">,
  ) {
    return this.request<T>(url, { method: "DELETE", ...config });
  },
};

export default defineNuxtPlugin(() => {
  return { provide: { api } };
});
