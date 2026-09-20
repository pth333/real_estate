/**
 * TrackingService — API theo dõi hành vi (session merging sau khi đăng nhập).
 * Dùng qua composable useTrackingService().
 */
import { BaseService, defineService } from '~/services/api'

export class TrackingService extends BaseService {
  /** Gộp session ẩn danh vào tài khoản vừa đăng nhập */
  async mergeSession(sessionId: string): Promise<void> {
    await this.post('/tracking/merge', { session_id: sessionId })
  }
}

export const useTrackingService = defineService(TrackingService)
