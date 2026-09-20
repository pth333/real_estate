<template>
  <div id="app" class="min-h-screen bg-limestone text-graphite">
    <NuxtErrorBoundary>
      <NuxtLayout>
        <NConfigProvider :theme-overrides="themeOverrides">
          <NNotificationProvider>
            <NMessageProvider>
              <NDialogProvider>
                <GlobalInit>
                  <NuxtPage />
                </GlobalInit>
              </NDialogProvider>
            </NMessageProvider>
          </NNotificationProvider>
        </NConfigProvider>
      </NuxtLayout>
    </NuxtErrorBoundary>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore } from "~/stores/notification";
import { useAuthStore } from "~/stores/auth";
import type { GlobalThemeOverrides } from "naive-ui";

const notifStore = useNotificationStore();
const authStore = useAuthStore();

const themeOverrides: GlobalThemeOverrides = {
  common: {
    borderRadius: "6px",
    borderRadiusSmall: "4px",
    // Đồng bộ màu thương hiệu với Tailwind emerald-600 dùng trong toàn site.
    // Trước đây Naive UI dùng màu primary mặc định (#18a058) nên nút Naive lệch
    // tông so với các khối emerald của Tailwind.
    primaryColor: "#059669",
    primaryColorHover: "#10b981",
    primaryColorPressed: "#047857",
    primaryColorSuppl: "#10b981",
  },
};

// Thông báo là dữ liệu cá nhân → chỉ kết nối SSE khi đã đăng nhập.
// Khách vãng lai vẫn xem web bình thường, không bị gọi API thông báo (tránh 401 lặp lại).
watch(
  () => authStore.isAuthenticated,
  (authenticated) => {
    if (authenticated) notifStore.connectSSE();
    else notifStore.disconnectSSE();
  },
  { immediate: true },
);

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
