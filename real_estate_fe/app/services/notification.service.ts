/**
 * NotificationService — API thông báo (kết nối SSE nằm ở store, không thuộc service).
 * Dùng qua composable useNotificationService().
 */
import { BaseService, defineService } from '~/services/api'
import type { NotificationItem } from '~/types/real_estate'

export class NotificationService extends BaseService {
  /** Danh sách thông báo đã lưu (envelope → bóc `.data`) */
  getNotifications(): Promise<NotificationItem[]> {
    return this.getData<NotificationItem[]>('/notifications')
  }
}

export const useNotificationService = defineService(NotificationService)
