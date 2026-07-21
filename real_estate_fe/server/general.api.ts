import type { RealEstateSearchRequest } from '~/types/real_estate'
import useApi from "~/app/composables/useApi";

const api = useApi();

export const general = {
  GetAllCategory: async () => {
    return (await api.request({ url: '/category' })) || {}
  },
  GetListByCategory: async (payload: RealEstateSearchRequest, slug: string, page: number) => {
    return (
      (await api.request({
        url: `/real-estate/${slug}/${page}`,
        data: payload,
        method: 'post',
      })) || {}
    )
  },
}
