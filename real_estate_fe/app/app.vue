<template>
  <div id="app" class="min-h-screen bg-limestone text-graphite">
    <NuxtErrorBoundary>
      <NuxtLayout>
        <NConfigProvider :theme-overrides="themeOverrides">
          <NNotificationProvider>
            <NMessageProvider>
              <GlobalInit>
                <NuxtPage />
              </GlobalInit>
            </NMessageProvider>
          </NNotificationProvider>
        </NConfigProvider>
      </NuxtLayout>
    </NuxtErrorBoundary>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore } from "~/stores/notification";
import type { GlobalThemeOverrides } from "naive-ui";

const notifStore = useNotificationStore();

const themeOverrides: GlobalThemeOverrides = {
  common: {
    borderRadius: "6px",
    borderRadiusSmall: "4px",
  },
};

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
