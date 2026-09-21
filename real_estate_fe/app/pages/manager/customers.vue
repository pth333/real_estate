<template>
  <div class="flex h-full min-h-0 flex-col gap-4 rounded-xl border border-gray-200 bg-white p-5 lg:p-6">
    <!-- Tiêu đề trang đã có ở khu vực quản lý nên không lặp lại ở đây -->
    <div class="flex flex-1 items-center justify-center py-10">
      <n-empty description="Tính năng quản lý khách hàng đang được xây dựng" />
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  alias: [
    "/nguoi-ban/quan-ly-khach-hang"
  ],
  requiresAuth: true,
});
</script>
