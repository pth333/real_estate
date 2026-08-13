<template>
  <n-popover v-model:show="showPopover" trigger="click" placement="bottom-end" :width="340" raw :show-arrow="false">
    <template #trigger>
      <button class="relative p-2 text-gray-600 transition hover:bg-gray-100 hover:text-gray-800 rounded-full flex items-center justify-center">
        <n-badge :value="store.unreadCount" :max="99" :show="store.unreadCount > 0">
          <IconBell class="size-6 text-gray-600" />
        </n-badge>
      </button>
    </template>

    <!-- Custom render panel inside popover -->
    <NotificationPanel @close="showPopover = false" />
  </n-popover>
</template>

<script setup lang="ts">
import { useNotificationStore } from "~/stores/notification"

const store = useNotificationStore()
const showPopover = ref(false)
</script>
