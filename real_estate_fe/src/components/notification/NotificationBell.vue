<template>
  <div class="relative">
    <button
      class="relative rounded-full p-2 text-gray-600 transition hover:bg-gray-100 hover:text-gray-800"
      @click="togglePanel"
    >
      <!-- Icon chuông -->
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="size-6"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="2"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
        />
      </svg>

      <!-- Badge số chưa đọc -->
      <span
        v-if="store.unreadCount > 0"
        class="absolute -right-1 -top-1 flex size-5 items-center justify-center rounded-full bg-red-500 text-[11px] font-bold text-white"
      >
        {{ store.unreadCount > 99 ? '99+' : store.unreadCount }}
      </span>
    </button>

    <!-- Panel dropdown -->
    <NotificationPanel
      v-if="showPanel"
      @close="showPanel = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useNotificationStore } from '@/stores/notification'
import NotificationPanel from './NotificationPanel.vue'

const store = useNotificationStore()
const showPanel = ref(false)

function togglePanel() {
  showPanel.value = !showPanel.value
}
</script>
