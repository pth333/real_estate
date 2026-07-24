import { defineStore } from "pinia";
import type { NotificationItem, NotificationSSEPayload } from "~/types/real_estate";

export const useNotificationStore = defineStore("notification", () => {
  const items = ref<NotificationItem[]>([]);
  const total = ref(0);
  const unreadCount = ref(0);
  const loading = ref(false);
  const eventSource = ref<EventSource | null>(null);
  const connected = ref(false);

  async function fetchList(userID: number) {
    loading.value = true;
    try {
      const { $api } = useNuxtApp();
      const res = await $api.get<{ data: { data: NotificationItem[]; total: number } }>(
        "/notifications",
        { params: { user_id: userID, page: 1, limit: 20 } },
      );
      items.value = (res as any).data.data;
      total.value = res.data.total;
      unreadCount.value = res.data.data.filter((n) => !n.is_read).length;
    } catch (e) {
      console.error("Lỗi tải notifications:", e);
    } finally {
      loading.value = false;
    }
  }

  async function markAsRead(id: number) {
    try {
      const { $api } = useNuxtApp();
      await $api.patch(`/notifications/${id}/read`);
      const notif = items.value.find((n) => n.id === id);
      if (notif && !notif.is_read) {
        notif.is_read = true;
        unreadCount.value = Math.max(0, unreadCount.value - 1);
      }
    } catch (e) {
      console.error("Lỗi đánh dấu đã đọc:", e);
    }
  }

  function connectSSE(userID: number) {
    if (eventSource.value) return;

    const config = useRuntimeConfig()
    const url = `${config.public.apiBaseUrl}/notifications/stream?user_id=${userID}`;
    const es = new EventSource(url);

    es.onopen = () => { connected.value = true; };

    es.onmessage = (event) => {
      try {
        const payload: NotificationSSEPayload = JSON.parse(event.data);
        if (payload.type === "new_listing") {
          // Show toast thông qua useToast system
          const { showToast } = useToast()
          showToast('success', `BĐS mới: ${payload.title} - ${payload.price_vnd ? (payload.price_vnd / 1_000_000_000).toFixed(1) + ' tỷ' : ''}`)
          unreadCount.value++;
          fetchList(userID);
        }
      } catch { /* ignore */ }
    };

    es.onerror = () => { connected.value = false; };
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
    items, total, unreadCount, loading, connected,
    fetchList, markAsRead, connectSSE, disconnectSSE,
  };
});
