<template>
  <div class="h-screen bg-gray-50 flex flex-col font-sans overflow-hidden">
    <!-- Header của hệ thống quản lý (Cố định) -->
    <div class="h-[65px] bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between shadow-sm flex-shrink-0">
      <div class="flex items-center gap-3">
        <n-avatar round size="medium" class="bg-red-50 text-red-500">
          <template #icon>
            <n-icon><IconBuilding /></n-icon>
          </template>
        </n-avatar>
        <span class="text-lg font-bold text-gray-900">HỆ THỐNG QUẢN TRỊ BĐS</span>
      </div>

      <!-- Nút thoát về Trang chủ -->
      <n-button text @click="goToHome">
        Quay lại trang chủ
        <template #icon>
          <n-icon><IconArrowRight /></n-icon>
        </template>
      </n-button>
    </div>

    <!-- Layout chính dạng Sidebar + Content -->
    <div class="flex flex-1 overflow-hidden">
      <!-- Sidebar điều hướng (Cố định, tự động đồng bộ activeKey theo route) -->
      <ManagerSidebar :active-key="activeKey" class="flex-shrink-0 h-full" />

      <!-- Nội dung động của các trang con -->
      <div class="flex-1 p-6 overflow-hidden flex flex-col h-full">
        <div class="w-full flex-1 flex flex-col gap-4 min-h-0">
          <div class="flex justify-between items-center flex-shrink-0">
            <h1 class="text-xl font-bold text-gray-900">{{ currentTitle }}</h1>
          </div>

          <!-- NuxtPage render nội dung của posts.vue hoặc customers.vue -->
          <div class="flex-1 min-h-0 flex flex-col">
            <NuxtPage />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ManagerSidebar from "~/components/manager/ManagerSidebar.vue";
import IconBuilding from "~/icons/IconBuilding.vue";
import IconArrowRight from "~/icons/IconArrowRight.vue";

// Thiết lập alias cho cả folder quản lý sang /nguoi-ban và tắt layout default
definePageMeta({
  alias: "/nguoi-ban",
  layout: "empty",
});

const route = useRoute();

// Tự động suy ra activeKey dựa trên URL hiện tại
const activeKey = computed<"posts" | "customers">(() => {
  if (route.path.includes("quan-ly-khach-hang") || route.path.includes("customers")) {
    return "customers";
  }
  return "posts";
});

const currentTitle = computed(() => {
  return activeKey.value === "posts" ? "Danh sách bài đăng của bạn" : "Danh sách khách hàng đăng ký";
});

const goToHome = () => {
  navigateTo("/");
};
</script>

<style scoped>
body {
  overflow-x: hidden;
}
</style>
