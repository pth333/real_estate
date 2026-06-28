import api from './service.api'
import type {
  DashboardSummary,
  PaginatedResponse,
  RealEstateModel,
  RealEstateSearchRequest,
} from '@/types/real_estate'

export default {
  getList(payload: RealEstateSearchRequest) {
    return api.post<PaginatedResponse<RealEstateModel>>(
      '/dashboard/list-real-estate',
      payload,
    )
  },

  getSummary(from?: string, to?: string) {
    return api.get<DashboardSummary>('/dashboard/summary', {
      params: { from, to },
    })
  },
}
