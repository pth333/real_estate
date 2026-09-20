/**
 * FavoriteService — API bất động sản yêu thích.
 * Dùng qua composable useFavoriteService().
 */
import { BaseService, defineService } from '~/services/api'

/** Envelope trả về của API toggle yêu thích */
export interface ToggleFavoriteResponse {
  data?: {
    is_favorite?: boolean
  }
}

export class FavoriteService extends BaseService {
  /** Toggle yêu thích 1 tin đăng, trả về trạng thái mới (không cần body) */
  async toggle(realEstateId: number): Promise<boolean> {
    const res = await this.post<ToggleFavoriteResponse>(`/real-estate/favorite/${realEstateId}`)
    return res?.data?.is_favorite ?? false
  }
}

export const useFavoriteService = defineService(FavoriteService)
