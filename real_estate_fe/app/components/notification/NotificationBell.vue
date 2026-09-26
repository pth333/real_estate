<template>
  <!-- Không set width cứng để panel tự co theo màn hình (xem NotificationPanel) -->
  <n-popover v-model:show="showPopover" trigger="click" placement="bottom-end" raw :show-arrow="false">
    <template #trigger>
      <button class="relative flex items-center justify-center rounded-full p-2 text-gray-600 transition hover:bg-gray-100 hover:text-gray-800">
        <n-badge :value="store.unreadCount" :max="9" :show="store.unreadCount > 0">
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

onMounted(() => {
  store.fetchList()
})
</script>
