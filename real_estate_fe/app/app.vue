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
  // Kết nối SSE ngay khi vào app
  notifStore.connectSSE();
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
