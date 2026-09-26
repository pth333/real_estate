<template>
  <!-- Nền + màu chữ gốc dùng token Tailwind thật.
       Trước đây là bg-limestone/text-graphite nhưng các token này CHƯA TỪNG được
       định nghĩa trong main.css nên không render ra style nào. -->
  <div id="app" class="min-h-screen bg-white text-gray-900">
    <NuxtErrorBoundary>
      <NuxtLayout>
        <NConfigProvider :theme-overrides="themeOverrides" :breakpoints="BREAKPOINTS">
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

/**
 * Breakpoint cho n-grid (responsive="screen").
 * Giữ nguyên bộ key gốc của Naive UI và bổ sung md/lg trùng khớp 1:1 với Tailwind
 * (md = 768px tablet, lg = 1024px desktop) để grid và class Tailwind luôn đổi cùng lúc.
 */
const BREAKPOINTS: Record<string, number> = {
  xs: 0,
  s: 640,
  md: 768,
  m: 1024,
  lg: 1024,
  l: 1280,
  xl: 1536,
  "2xl": 1920,
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
