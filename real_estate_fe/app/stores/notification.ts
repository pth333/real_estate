import { defineStore } from "pinia";
import type { NotificationItem, NotificationSSEPayload } from "~/types/real_estate";

export const useNotificationStore = defineStore("notification", () => {
  const items = ref<NotificationItem[]>([]);
  const unreadCount = ref(0);
  const loading = ref(false);
  const eventSource = ref<EventSource | null>(null);
  const connected = ref(false);

  async function fetchList() {
    loading.value = true;
    try {
      const { $api } = useNuxtApp();
      const res = await $api.get<{ data: NotificationItem[] }>("/notifications");
      items.value = res.data;

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

  function connectSSE() {
    if (eventSource.value || import.meta.server) return;

    const config = useRuntimeConfig();
    const url = `${config.public.apiBaseUrl}/notifications/stream`;
    const es = new EventSource(url);

    es.onopen = () => { connected.value = true; };

    es.onmessage = (event) => {
      try {
        const payload: NotificationSSEPayload = JSON.parse(event.data);

        // Push vào list hiện tại
        items.value.unshift({
          id: Date.now(), // Fake ID cho client
          type: "new_listing",
          payload: payload,
          created_at: new Date().toISOString()
        });

        if (typeof window !== 'undefined' && (window as any).$message) {
            (window as any).$message.success(`BĐS mới: ${payload.title} - ${(payload.price / 1_000_000_000).toFixed(1)} tỷ`, {
                duration: 5000,
                closable: true
            });
        }

        unreadCount.value++;
      } catch (e) {
        console.error("SSE parse error:", e);
      }
    };

    es.onerror = () => {
        connected.value = false;
        es.close();
        eventSource.value = null;
        // Reconnect sau 5s
        setTimeout(connectSSE, 5000);
    };

    eventSource.value = es;
  }

  function disconnectSSE() {
    if (eventSource.value) {
      eventSource.value.close();
      eventSource.value = null;
      connected.value = false;
    }
  }

  return {
    items, unreadCount, loading, connected,
    fetchList, markAllAsRead, markAsRead, connectSSE, disconnectSSE,
  };
});
