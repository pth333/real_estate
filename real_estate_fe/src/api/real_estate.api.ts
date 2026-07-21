import api from './service.api'
import type { RealEstateSearchRequest } from '@/types/real_estate'

export const realEstateApi = {
  getList: async (payload: RealEstateSearchRequest) => {
    return (
      (await api.request({ url: '/real-estate/list', data: payload, method: 'post' })) ||
      {}
    )
  },
  getSummary: async (from?: string, to?: string) => {
    return (await api.request({ url: '/dashboard/summary', params: { from, to } })) || {}
  },
}
