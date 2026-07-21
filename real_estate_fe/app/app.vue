<template>
  <div id="app" class="min-h-screen bg-limestone text-graphite">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <NotificationToast />
    <ToastDisplay />
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore } from "~/app/stores/notification";

const notifStore = useNotificationStore();

onMounted(() => {
  // Notification sẽ kết nối SSE sau khi auth hoàn thiện
  // TODO: thêm kết nối SSE realtime
});

onUnmounted(() => {
  notifStore.disconnectSSE();
});
</script>

<style>
body {
  margin: 0;
  font-family: "Inter", system-ui, sans-serif;
}
</style>
