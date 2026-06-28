import { defineStore } from 'pinia'
import type { NotificationItem, NotificationSSEPayload } from '@/types/real_estate'
import notificationApi from '@/api/notification.api'

interface NotificationState {
  items: NotificationItem[]
  total: number
  unreadCount: number
  loading: boolean
  /** Toast queue — các notification mới nhận qua SSE */
  toasts: NotificationSSEPayload[]
  /** EventSource instance */
  eventSource: EventSource | null
  connected: boolean
}

export const useNotificationStore = defineStore('notification', {
  state: (): NotificationState => ({
    items: [],
    total: 0,
    unreadCount: 0,
    loading: false,
    toasts: [],
    eventSource: null,
    connected: false,
  }),

  actions: {
    async fetchList(userID: number) {
      this.loading = true
      try {
        const res = await notificationApi.getList(userID)
        this.items = res.data.data
        this.total = res.data.total
        this.unreadCount = res.data.data.filter((n) => !n.is_read).length
      } catch (e) {
        console.error('Lỗi tải notifications:', e)
      } finally {
        this.loading = false
      }
    },

    async markAsRead(id: number) {
      try {
        await notificationApi.markAsRead(id)
        const notif = this.items.find((n) => n.id === id)
        if (notif && !notif.is_read) {
          notif.is_read = true
          this.unreadCount = Math.max(0, this.unreadCount - 1)
        }
      } catch (e) {
        console.error('Lỗi đánh dấu đã đọc:', e)
      }
    },

    /** Kết nối SSE stream để nhận notification real-time */
    connectSSE(userID: number) {
      // Tránh tạo nhiều kết nối
      if (this.eventSource) return

      const url = `${
        import.meta.env.VITE_API_URL || 'http://localhost:8000'
      }/api/notifications/stream?user_id=${userID}`
      const es = new EventSource(url)

      es.onopen = () => {
        this.connected = true
      }

      es.onmessage = (event) => {
        try {
          const payload: NotificationSSEPayload = JSON.parse(event.data)
          if (payload.type === 'new_listing') {
            // Thêm vào toast queue
            this.toasts.push(payload)

            // Tự động xoá toast sau 5s
            setTimeout(() => {
              const idx = this.toasts.indexOf(payload)
              if (idx !== -1) this.toasts.splice(idx, 1)
            }, 5000)

            // Tăng unread count
            this.unreadCount++

            // Refresh danh sách notification
            this.fetchList(userID)
          }
        } catch {
          // ignore parse error
        }
      }

      es.onerror = () => {
        this.connected = false
        // EventSource sẽ tự động reconnect
      }

      this.eventSource = es
    },

    disconnectSSE() {
      if (this.eventSource) {
        this.eventSource.close()
        this.eventSource = null
        this.connected = false
      }
    },

    /** Xoá toast khỏi queue */
    dismissToast(payload: NotificationSSEPayload) {
      const idx = this.toasts.indexOf(payload)
      if (idx !== -1) this.toasts.splice(idx, 1)
    },
  },
})
