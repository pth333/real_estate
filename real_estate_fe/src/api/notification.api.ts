import api from './service.api'

export const notificationApi = {
  getList: async (userID: number, page = 1, limit = 20) => {
    return (
      (await api.request({
        url: '/notifications',
        params: { user_id: userID, page, limit },
      })) || {}
    )
  },
  markAsRead: async (id: number) => {
    return (await api.request({ url: `/notifications/${id}/read`, method: 'patch' })) || {}
  },
}
