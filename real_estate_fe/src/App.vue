<template>
  <div id="app" class="min-h-screen bg-gray-50 text-gray-800">
    <DefaultLayout>
      <router-view />
    </DefaultLayout>
    <!-- Toast notifications layer -->
    <NotificationToast />
  </div>
</template>

<script setup lang="ts">
import DefaultLayout from './layouts/DefaultLayout.vue'
import NotificationToast from '@/components/notification/NotificationToast.vue'
import { useNotificationStore } from '@/stores/notification'
import { useAuthStore } from './stores/auth.ts'

const notifStore = useNotificationStore()
const authStore = useAuthStore()

onMounted(() => {
  // Kết nối SSE để nhận notification real-time
  // TODO: thay user_id động sau khi có auth
  notifStore.connectSSE(authStore.userId)
  notifStore.fetchList(authStore.userId)
})

onUnmounted(() => {
  notifStore.disconnectSSE()
})
</script>

<style>
/* CSS global nhẹ */
body {
  margin: 0;
  font-family: 'Inter', system-ui, sans-serif;
}
</style>
