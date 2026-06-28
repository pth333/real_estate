import api from './service.api'
import type { NotificationItem } from '@/types/real_estate'

export default {
  getList(userID: number, page = 1, limit = 20) {
    return api.get<{
      data: NotificationItem[]
      total: number
      page: number
    }>('/notifications', {
      params: { user_id: userID, page, limit },
    })
  },

  markAsRead(id: number) {
    return api.patch(`/notifications/${id}/read`)
  },
}
