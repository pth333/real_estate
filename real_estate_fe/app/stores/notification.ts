import { defineStore } from "pinia";
import type { NotificationItem, NotificationSSEPayload } from "~/types/real_estate";
import { useNotificationService } from "~/services/notification.service";
import { NotificationStream } from "~/services/notification-stream";

export const useNotificationStore = defineStore("notification", () => {
  const items = ref<NotificationItem[]>([]);
  const unreadCount = ref(0);
  const loading = ref(false);
  const connected = ref(false);

  // Client SSE dùng fetch (EventSource không gửi được header Authorization → luôn 401)
  let stream: NotificationStream | null = null;

  async function fetchList() {
    loading.value = true;
    try {
      const notificationService = useNotificationService();
      items.value = await notificationService.getNotifications();

      if (import.meta.client) {
        // Load trạng thái đọc từ localStorage để tính unread
        const lastRead = localStorage.getItem("last_notif_read_at") || "0";
        const readIdsStr = localStorage.getItem("read_notification_ids") || "[]";
        const readIds = JSON.parse(readIdsStr) as number[];

        unreadCount.value = items.value.filter(n => {
          const isReadById = readIds.includes(Number(n.id));
          const isReadByTime = new Date(n.created_at).getTime() <= parseInt(lastRead);
          return !isReadById && !isReadByTime;
        }).length;
      } else {
        unreadCount.value = 0;
      }
    } catch (e) {
      console.error("Lỗi tải notifications:", e);
    } finally {
      loading.value = false;
    }
  }

  /** Xử lý 1 thông báo đẩy từ SSE */
  function pushRealtimeNotification(payload: NotificationSSEPayload) {
    items.value.unshift({
      id: Date.now(), // Fake ID cho client
      type: "new_listing",
      payload: payload,
      created_at: new Date().toISOString(),
    });

    window.message?.success(
      `BĐS mới: ${payload.title} - ${(payload.price / 1_000_000_000).toFixed(1)} tỷ`,
      { duration: 5000, closable: true },
    );

    unreadCount.value++;
  }

  /**
   * Kết nối nhận thông báo realtime. Chỉ chạy khi đã đăng nhập
   * (thông báo là dữ liệu cá nhân, khách vãng lai không gọi).
   */
  function connectSSE() {
    if (import.meta.server || stream) return;

    const authStore = useAuthStore();
    const config = useRuntimeConfig();

    stream = new NotificationStream({
      url: `${config.public.apiBaseUrl}/notifications/stream`,
      getToken: () => authStore.token ?? null,
      refreshToken: () => authStore.refreshToken(),
      handlers: {
        onMessage: (payload) => pushRealtimeNotification(payload as NotificationSSEPayload),
        onOpen: () => {
          connected.value = true;
        },
        onClose: () => {
          connected.value = false;
        },
      },
    });

    stream.start();
  }

  function disconnectSSE() {
    stream?.stop();
    stream = null;
    connected.value = false;
  }

  function markAllAsRead() {
    const lastNotif = items.value[0];
    if (lastNotif) {
      localStorage.setItem("last_notif_read_at", new Date(lastNotif.created_at).getTime().toString());
      // Xoá bớt danh sách ID đã đọc riêng lẻ vì mốc thời gian đã bao quát tất cả
      localStorage.removeItem("read_notification_ids");
      unreadCount.value = 0;
    }
  }

  function markAsRead(notifId: number) {
    const readIdsStr = localStorage.getItem("read_notification_ids") || "[]";
    let readIds: number[] = [];
    try {
      readIds = JSON.parse(readIdsStr);
    } catch (e) {
      readIds = [];
    }

    if (!readIds.includes(Number(notifId))) {
      readIds.push(Number(notifId));
      localStorage.setItem("read_notification_ids", JSON.stringify(readIds));

      // Giảm unreadCount nếu nó > 0
      if (unreadCount.value > 0) {
        unreadCount.value--;
      }
    }
  }

  return {
    items, unreadCount, loading, connected,
    fetchList, markAllAsRead, markAsRead, connectSSE, disconnectSSE,
  };
});
