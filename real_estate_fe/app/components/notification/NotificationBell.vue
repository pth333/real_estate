<template>
  <div class="relative">
    <button class="relative p-2 text-gray-600 transition hover:bg-gray-100 hover:text-gray-800" @click="togglePanel">
      <IconBell class="size-6" />
      <span v-if="store.unreadCount > 0"
        class="absolute -right-1 -top-1 flex size-5 items-center justify-center bg-red-500 text-[11px] font-bold text-white">
        {{ store.unreadCount > 99 ? '99+' : store.unreadCount }}
      </span>
    </button>
    <NotificationPanel v-if="showPanel" @close="showPanel = false" />
  </div>
</template>

<script setup lang="ts">
const store = useNotificationStore()
const showPanel = ref(false)

function togglePanel() {
  showPanel.value = !showPanel.value
}
</script>
