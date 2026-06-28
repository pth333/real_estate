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
import { onMounted, onUnmounted } from 'vue'
import DefaultLayout from './layouts/DefaultLayout.vue'
import NotificationToast from '@/components/notification/NotificationToast.vue'
import { useNotificationStore } from '@/stores/notification'

const notifStore = useNotificationStore()

onMounted(() => {
  // Kết nối SSE để nhận notification real-time
  // TODO: thay user_id động sau khi có auth
  notifStore.connectSSE(1)
  notifStore.fetchList(1)
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
