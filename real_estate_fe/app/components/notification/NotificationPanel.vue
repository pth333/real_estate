<template>
  <div
    class="absolute right-0 top-full z-50 mt-2 w-80 border border-gray-200 bg-white shadow-lg"
  >
    <div class="flex items-center justify-between border-b px-4 py-3">
      <h3 class="text-sm font-semibold">Thông báo</h3>
      <div class="flex items-center gap-2">
        <button
          v-if="store.unreadCount > 0"
          class="text-[10px] text-blue-500 hover:underline"
          @click="store.markAllAsRead"
        >
          Đọc tất cả
        </button>
        <button class="text-xs text-gray-400 hover:text-gray-600" @click="$emit('close')">
          ✕
        </button>
      </div>
    </div>

    <div v-if="store.loading" class="flex items-center justify-center py-8 text-sm text-gray-400">
      Đang tải...
    </div>

    <div v-else-if="store.items.length === 0" class="py-8 text-center text-sm text-gray-400">
      Chưa có thông báo nào
    </div>

    <ul v-else class="max-h-80 divide-y overflow-y-auto">
      <li
        v-for="notif in store.items"
        :key="notif.id"
        class="flex cursor-pointer items-start gap-2 px-4 py-3 transition hover:bg-gray-50"
        @click="handleClick(notif)"
      >
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-medium">
            {{ notif.payload?.title || "Bất động sản mới" }}
          </p>
          <p class="mt-0.5 truncate text-xs text-gray-500">
            {{ notif.payload?.address }}
          </p>
          <div class="mt-1 flex items-center justify-between">
            <span class="text-[10px] font-bold text-red-500">
              {{ formatPrice(notif.payload?.price) }}
            </span>
            <span class="text-[10px] text-gray-400">
              {{ formatTime(notif.created_at) }}
            </span>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore } from "~/stores/notification";
import type { NotificationItem } from "~/types/real_estate";

defineEmits<{ close: [] }>();

const store = useNotificationStore();

onMounted(() => {
  store.fetchList();
});

function formatPrice(price: number) {
  if (!price) return "";
  if (price >= 1_000_000_000) {
    return (price / 1_000_000_000).toFixed(1) + " tỷ";
  }
  return (price / 1_000_000).toFixed(0) + " triệu";
}

function formatTime(dateStr: string) {
  const date = new Date(dateStr);
  return date.toLocaleTimeString("vi-VN", {
    hour: "2-digit",
    minute: "2-digit",
  });
}

function handleClick(notif: NotificationItem) {
  if (notif.payload?.url) {
      // Nếu là URL internal hoặc external
      window.open(notif.payload.url, '_blank');
  }
}
</script>
