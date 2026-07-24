<template>
  <div id="app" class="min-h-screen bg-limestone text-graphite">
    <NuxtErrorBoundary>
      <NuxtLayout>
        <NuxtPage />
      </NuxtLayout>
      <ToastDisplay />

      <template #error="{ error, clearError }">
        <div class="flex min-h-screen flex-col items-center justify-center gap-4 bg-gray-50">
          <h2 class="text-xl font-semibold text-gray-800">Đã xảy ra lỗi</h2>
          <p class="text-sm text-gray-500">{{ error }}</p>
          <button class="rounded-lg bg-blue-600 px-6 py-2 text-sm font-semibold text-white hover:bg-blue-700"
            @click="clearError()">
            Thử lại
          </button>
        </div>
      </template>
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
