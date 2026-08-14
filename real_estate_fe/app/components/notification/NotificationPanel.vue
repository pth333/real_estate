<template>
  <div class="w-85 border border-gray-200 bg-white shadow-xl rounded-lg overflow-hidden flex flex-col">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-gray-100 px-4 py-3 bg-gray-50">
      <h3 class="text-sm font-semibold text-gray-800 flex items-center gap-1.5">
        Thông báo
        <n-badge v-if="store.unreadCount > 0" dot type="error" />
      </h3>
      <div class="flex items-center gap-2">
        <n-button
          v-if="store.unreadCount > 0"
          type="primary"
          text
          size="tiny"
          class="text-xs text-emerald-600 hover:text-emerald-700 font-medium"
          @click="store.markAllAsRead"
        >
          Đọc tất cả
        </n-button>
        <n-button
          circle
          quaternary
          size="small"
          class="text-gray-400 hover:text-gray-600"
          @click="$emit('close')"
        >
          <template #icon>
            <span class="text-xs">✕</span>
          </template>
        </n-button>
      </div>
    </div>

    <!-- Body content -->
    <div class="flex-1 min-h-[150px] flex flex-col justify-center">
      <!-- Loading state -->
      <div v-if="store.loading" class="flex flex-col items-center justify-center py-8 gap-2">
        <n-spin size="medium" stroke="var(--n-color-target)" />
        <span class="text-xs text-gray-400">Đang tải thông báo...</span>
      </div>

      <!-- Empty state -->
      <div v-else-if="store.items.length === 0" class="py-10">
        <n-empty description="Chưa có thông báo nào" />
      </div>

      <!-- Notification list -->
      <n-scrollbar v-else style="max-height: 360px">
        <n-list hoverable clickable class="divide-y divide-gray-50">
          <n-list-item
            v-for="notif in store.items"
            :key="notif.id"
            class="transition-colors duration-200 hover:bg-gray-50 cursor-pointer"
            @click="handleClick(notif)"
          >
            <n-thing>
              <!-- Title -->
              <template #title>
                <div class="flex items-start justify-between gap-2">
                  <span class="text-sm font-semibold text-gray-800 line-clamp-2 leading-snug">
                    {{ notif.payload?.title || "Bất động sản mới" }}
                  </span>
                  <!-- Badge for unread notification. Let's say if we can track specific read status, else we assume all currently loaded are highlighted depending on timestamp -->
                  <span v-if="isUnread(notif)" class="flex h-2 w-2 translate-y-1.5 rounded-full bg-emerald-500 shrink-0" />
                </div>
              </template>

              <!-- Description / Address -->
              <template #description>
                <span class="text-xs text-gray-500 line-clamp-1 mt-0.5">
                  {{ notif.payload?.address }}
                </span>
              </template>

              <!-- Custom price & time info at the bottom -->
              <div class="mt-2 flex items-center justify-between">
                <span class="text-xs font-bold text-red-500">
                  {{ formatPrice(notif.payload?.price || 0) }}
                </span>
                <span class="text-[10px] text-gray-400">
                  {{ fromNow(notif.created_at) }}
                </span>
              </div>
            </n-thing>
          </n-list-item>
        </n-list>
      </n-scrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore } from "~/stores/notification";
import type { NotificationItem } from "~/types/real_estate";
import { formatPrice, fromNow } from "~/utils/format";

defineEmits<{ close: [] }>();

const store = useNotificationStore();

onMounted(() => {
  store.fetchList();
});

function isUnread(notif: NotificationItem): boolean {
  if (typeof window !== "undefined") {
    const lastRead = localStorage.getItem("last_notif_read_at") || "0";
    const readIdsStr = localStorage.getItem("read_notification_ids") || "[]";
    let readIds: number[] = [];
    try {
      readIds = JSON.parse(readIdsStr);
    } catch (e) {
      readIds = [];
    }

    const isReadById = readIds.includes(Number(notif.id));
    const isReadByTime = new Date(notif.created_at).getTime() <= parseInt(lastRead);
    return !isReadById && !isReadByTime;
  }
  return false;
}

function handleClick(notif: NotificationItem) {
  // Đánh dấu đã đọc trên client tab hiện tại ngay lập tức
  store.markAsRead(notif.id);

  if (notif.payload?.slug) {
    // Nếu là URL internal hoặc external
    window.open(notif.payload.slug, "_blank");
  }
}
</script>
