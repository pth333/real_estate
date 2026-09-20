/**
 * RealEstateService — toàn bộ API domain bất động sản: danh sách, chi tiết,
 * gợi ý, khu vực nổi bật, yêu thích và số liệu tổng quan dashboard.
 * Dùng qua composable useRealEstateService().
 */
import { BaseService, defineService } from '~/services/api'
import type { ListResult, QueryParams } from '~/services/api'
import type {
  DashboardSummary,
  PaginatedResponse,
  RealEstateResponse,
  RealEstateSearchRequest,
  TopCityOption,
} from '~/types/real_estate'

/**
 * Params của API gợi ý: khách dùng session_id, user đã đăng nhập dùng user_id.
 * (Khai báo bằng type alias để truyền thẳng được vào query params.)
 */
export type RecommendationParams = {
  limit?: number
  session_id?: string
  user_id?: number
}

export class RealEstateService extends BaseService {
  /**
   * Tìm kiếm danh sách BĐS (body RealEstateSearchRequest).
   * Endpoint trả `{ total, data }` ở top-level, KHÔNG bọc envelope
   * → dùng post (không dùng postData).
   */
  searchList(payload: RealEstateSearchRequest): Promise<PaginatedResponse<RealEstateResponse>> {
    return this.post<PaginatedResponse<RealEstateResponse>>('/real-estate/list', payload)
  }

  /**
   * Danh sách BĐS theo SEO URL của trang danh mục (category + segment filter).
   * Endpoint trả `{ total, data }` top-level → dùng get (không dùng getData).
   */
  getListByCategoryPath(
    path: string,
    params: QueryParams,
  ): Promise<PaginatedResponse<RealEstateResponse>> {
    return this.get<PaginatedResponse<RealEstateResponse>>(`/real-estate/${path}`, params)
  }

  /** Chi tiết 1 tin đăng theo id (envelope → bóc `.data`) */
  getDetail(id: number): Promise<RealEstateResponse> {
    return this.getData<RealEstateResponse>(`/real-estate/detail/${id}`)
  }

  /** Số liệu tổng quan dashboard theo khoảng ngày (envelope → bóc `.data`) */
  getSummary(from?: string, to?: string): Promise<DashboardSummary> {
    return this.getData<DashboardSummary>('/dashboard/summary', { from, to })
  }

  /** Danh sách BĐS gợi ý (envelope → bóc `.data`) */
  getRecommendations(params: RecommendationParams = {}): Promise<RealEstateResponse[]> {
    return this.getData<RealEstateResponse[]>('/real-estate/recommend', params)
  }

  /** Top thành phố nhiều tin đăng nhất cho khối "BĐS theo địa điểm" (envelope → bóc `.data`) */
  getTopCities(): Promise<TopCityOption[]> {
    return this.getData<TopCityOption[]>('/real-estate/list/top-city')
  }

  /**
   * Danh mục BĐS yêu thích của user hiện tại.
   * Endpoint trả `{ data, total }` top-level → map sang `{ items, total }`.
   */
  async getFavorites(params: { page: number; size: number }): Promise<ListResult<RealEstateResponse>> {
    const res = await this.get<PaginatedResponse<RealEstateResponse>>('/real-estate/favorites', params)
    return { items: res.data ?? [], total: res.total ?? 0 }
  }
}

export const useRealEstateService = defineService(RealEstateService)
