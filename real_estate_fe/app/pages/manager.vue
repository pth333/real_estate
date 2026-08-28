<template>
  <div class="h-screen bg-gray-50 flex flex-col font-sans overflow-hidden">
    <!-- Header của hệ thống quản lý (Cố định) -->
    <div class="h-[65px] bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between flex-shrink-0">
      <div class="flex items-center gap-3">
        <n-avatar round size="medium" class="bg-red-50 text-red-500">
          <template #icon>
            <n-icon>
              <IconBuilding />
            </n-icon>
          </template>
        </n-avatar>
        <div class="flex flex-col leading-tight">
          <span class="text-lg font-bold text-gray-900">Hệ thống quản trị BĐS</span>
          <span class="text-xs text-gray-400">Quản lý tin đăng &amp; khách hàng</span>
        </div>
      </div>

      <!-- Nút thoát về Trang chủ -->
      <n-button text class="text-emerald-600 hover:text-emerald-700" @click="goToHome">
        Quay lại trang chủ
        <template #icon>
          <n-icon>
            <IconArrowRight />
          </n-icon>
        </template>
      </n-button>
    </div>

    <!-- Layout chính dạng Sidebar + Content -->
    <div class="flex flex-1 overflow-hidden">
      <!-- Sidebar điều hướng (Cố định, tự động đồng bộ activeKey theo route) -->
      <ManagerSidebar :active-key="activeKey" class="flex-shrink-0 h-full" />

      <!-- Nội dung động của các trang con -->
      <div class="flex-1 p-6 overflow-y-auto flex flex-col h-full">
        <div class="w-full flex-1 flex flex-col gap-4 min-h-0">
          <div class="flex flex-col gap-1 flex-shrink-0">
            <!-- Breadcrumb kiểu dự án -->
            <nav class="text-xs text-emerald-600">
              <span>Quản lý</span>
              <span class="mx-1 text-gray-400">/</span>
              <span class="text-gray-500">{{ currentTitle }}</span>
            </nav>
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
const activeKey = computed<"projects" | "posts" | "customers" | "favorites">(() => {
  if (route.path.includes("quan-ly-du-an") || route.path.includes("projects")) {
    return "projects";
  }
  if (route.path.includes("tao-du-an")) {
    return "projects";
  }
  if (route.path.includes("quan-ly-khach-hang") || route.path.includes("customers")) {
    return "customers";
  }
  if (route.path.includes("quan-ly-yeu-thich") || route.path.includes("favorites")) {
    return "favorites";
  }
  return "posts";
});

const currentTitle = computed(() => {
  if (activeKey.value === "projects") return "Danh sách dự án";
  if (activeKey.value === "posts") return "Danh sách bài đăng của bạn";
  if (activeKey.value === "customers") return "Danh sách khách hàng đăng ký";
  return "Danh mục bất động sản yêu thích";
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
