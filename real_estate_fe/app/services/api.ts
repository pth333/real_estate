/**
 * Tầng nền cho mọi service gọi API.
 *
 * Mục đích: gom toàn bộ việc gọi HTTP về các class service (xem deposit.service.ts),
 * component/store chỉ gọi hàm — không tự dựng URL, không tự bóc envelope.
 *
 * Cách dùng:
 *   export class FooService extends BaseService {
 *     getFoo(id: number) { return this.getData<Foo>(`/foo/${id}`) }
 *   }
 *   export const useFooService = defineService(FooService)
 */
import type { api } from '~/plugins/api.client'

/** Kiểu của $api (instance trong plugins/api.client.ts) */
export type ApiClient = typeof api

/** Query params — nhận mọi giá trị, rỗng thì bỏ qua */
export type QueryParams = Record<string, string | number | boolean | null | undefined>

/** Envelope chuẩn backend trả về (internal/response/response.go) */
export interface ApiEnvelope<T> {
  success: boolean
  message?: string
  data: T
  meta?: ApiMeta
  error?: unknown
}

export interface ApiMeta {
  total?: number
  page?: number
  size?: number
}

/** Kết quả danh sách đã chuẩn hoá: items + tổng số bản ghi */
export interface ListResult<T> {
  items: T[]
  total: number
}

/**
 * BaseService giữ $api và các helper dùng chung.
 * - getData/postData/putData: bóc sẵn `.data` trong envelope
 * - get/post/put/remove: trả nguyên response, dùng cho endpoint không theo envelope
 */
export abstract class BaseService {
  constructor(protected readonly api: ApiClient) {}

  protected get<T>(url: string, params?: QueryParams, silent = false): Promise<T> {
    return this.api.get<T>(url, { ...(params ? { params } : {}), ...(silent ? { silent } : {}) })
  }

  protected post<T>(url: string, body?: unknown, params?: QueryParams): Promise<T> {
    return this.api.post<T>(url, body, params ? { params } : undefined)
  }

  protected put<T>(url: string, body?: unknown): Promise<T> {
    return this.api.put<T>(url, body)
  }

  protected async remove(url: string): Promise<void> {
    await this.api.delete(url)
  }

  /** GET và trả về `data` trong envelope. silent = true để không hiện toast lỗi. */
  protected async getData<T>(url: string, params?: QueryParams, silent = false): Promise<T> {
    const res = await this.get<ApiEnvelope<T>>(url, params, silent)
    return res.data
  }

  /** POST và trả về `data` trong envelope */
  protected async postData<T>(url: string, body?: unknown, params?: QueryParams): Promise<T> {
    const res = await this.post<ApiEnvelope<T>>(url, body, params)
    return res.data
  }

  /** PUT và trả về `data` trong envelope */
  protected async putData<T>(url: string, body?: unknown): Promise<T> {
    const res = await this.put<ApiEnvelope<T>>(url, body)
    return res.data
  }

  /** GET danh sách: trả về `data` (mảng) + `meta.total` */
  protected async getList<T>(url: string, params?: QueryParams): Promise<ListResult<T>> {
    const res = await this.get<ApiEnvelope<T[]>>(url, params)
    return { items: res.data ?? [], total: res.meta?.total ?? 0 }
  }
}

/**
 * defineService — tạo composable lấy $api từ Nuxt app rồi khởi tạo service.
 * Dùng trong component/store: `const fooService = useFooService()`.
 */
export function defineService<T extends BaseService>(
  ServiceClass: new (api: ApiClient) => T,
): () => T {
  return () => {
    const { $api } = useNuxtApp()
    return new ServiceClass($api)
  }
}
