import { defineStore } from 'pinia'
import type { NotificationItem, NotificationSSEPayload } from '@/types/real_estate'
import { notificationApi } from '@/api/notification.api'

export const useNotificationStore = defineStore('notification', () => {
  const items = ref<NotificationItem[]>([])
  const total = ref(0)
  const unreadCount = ref(0)
  const loading = ref(false)
  /** Toast queue — các notification mới nhận qua SSE */
  const toasts = ref<NotificationSSEPayload[]>([])
  /** EventSource instance */
  const eventSource = ref<EventSource | null>(null)
  const connected = ref(false)

  async function fetchList(userID: number) {
    loading.value = true
    try {
      const res = await notificationApi.getList(userID)
      items.value = res.data.data
      total.value = res.data.total
      unreadCount.value = res.data.data.filter((n: any) => !n.is_read).length
    } catch (e) {
      console.error('Lỗi tải notifications:', e)
    } finally {
      loading.value = false
    }
  }

  async function markAsRead(id: number) {
    try {
      await notificationApi.markAsRead(id)
      const notif = items.value.find((n) => n.id === id)
      if (notif && !notif.is_read) {
        notif.is_read = true
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } catch (e) {
      console.error('Lỗi đánh dấu đã đọc:', e)
    }
  }

  /** Kết nối SSE stream để nhận notification real-time */
  function connectSSE(userID: number) {
    // Tránh tạo nhiều kết nối
    if (eventSource.value) return

    const url = `${
      import.meta.env.VITE_API_URL || 'http://localhost:8000'
    }/api/2026/notifications/stream?user_id=${userID}`
    const es = new EventSource(url)

    es.onopen = () => {
      connected.value = true
    }

    es.onmessage = (event) => {
      try {
        const payload: NotificationSSEPayload = JSON.parse(event.data)
        if (payload.type === 'new_listing') {
          // Thêm vào toast queue
          toasts.value.push(payload)

          // Tự động xoá toast sau 5s
          setTimeout(() => {
            const idx = toasts.value.indexOf(payload)
            if (idx !== -1) toasts.value.splice(idx, 1)
          }, 5000)

          // Tăng unread count
          unreadCount.value++

          // Refresh danh sách notification
          fetchList(userID)
        }
      } catch {
        // ignore parse error
      }
    }

    es.onerror = () => {
      connected.value = false
      // EventSource sẽ tự động reconnect
    }

    eventSource.value = es
  }

  function disconnectSSE() {
    if (eventSource.value) {
      eventSource.value.close()
      eventSource.value = null
      connected.value = false
    }
  }

  /** Xoá toast khỏi queue */
  function dismissToast(payload: NotificationSSEPayload) {
    const idx = toasts.value.indexOf(payload)
    if (idx !== -1) toasts.value.splice(idx, 1)
  }

  return {
    items,
    total,
    unreadCount,
    loading,
    toasts,
    eventSource,
    connected,
    fetchList,
    markAsRead,
    connectSSE,
    disconnectSSE,
    dismissToast,
  }
})
