<template>
  <div id="app" class="min-h-screen bg-limestone text-graphite">
    <NuxtErrorBoundary>
      <NuxtLayout>
        <NNotificationProvider>
          <NMessageProvider>
            <GlobalInit>
              <NuxtPage />
            </GlobalInit>
          </NMessageProvider>
        </NNotificationProvider>
      </NuxtLayout>
    </NuxtErrorBoundary>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore } from "~/stores/notification";

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
