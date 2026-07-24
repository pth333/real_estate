<template>

  <div
    class="absolute right-0 top-full z-50 mt-2 w-80 rounded-lg border border-gray-200 bg-white shadow-lg"
  >

    <div class="flex items-center justify-between border-b px-4 py-3">

      <h3 class="text-sm font-semibold">Thông báo</h3>
       <button class="text-xs text-gray-400 hover:text-gray-600" @click="$emit('close')"> ✕ </button
      >
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
        :class="{ 'bg-blue-50': !notif.is_read }"
        @click="handleClick(notif)"
      >
         <!-- Dot cho chưa đọc --> <span
          v-if="!notif.is_read"
          class="mt-1.5 size-2 shrink-0 rounded-full bg-blue-500"
        />
        <div class="min-w-0 flex-1">

          <p class="truncate text-sm font-medium">{{ notif.title }}</p>

          <p class="mt-0.5 truncate text-xs text-gray-500"> {{ notif.message }} </p>

          <p class="mt-1 text-[10px] text-gray-400"> {{ fromNow(notif.created_at) }} </p>

        </div>

      </li>

    </ul>

  </div>

</template>

<script setup lang="ts">
import type { NotificationItem } from '~/types/real_estate'
import { useNotificationStore } from '~/stores/notification'
import { useAuthStore } from '~/stores/auth'

defineEmits<{ close: [] }>()

const store = useNotificationStore()
const authStore = useAuthStore()
const { fromNow } = useRelativeTime()

onMounted(() => {
  store.fetchList(authStore.user?.id ?? 1)
})

function handleClick(notif: NotificationItem) {
  if (!notif.is_read) {
    store.markAsRead(notif.id)
  }
}

</script>

