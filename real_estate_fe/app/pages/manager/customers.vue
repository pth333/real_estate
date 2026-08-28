<template>
  <div class="bg-white rounded-lg border border-gray-200 p-6 flex flex-col gap-4 h-full min-h-0">
    <div>
      <p class="text-sm text-gray-500 mb-1">Quản lý khách hàng</p>
      <h2 class="text-lg font-bold text-gray-800">Danh sách khách hàng đăng ký</h2>
    </div>

    <div class="flex-1 flex items-center justify-center py-10">
      <n-empty description="Tính năng quản lý khách hàng đang được xây dựng" />
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  alias: [
    "/nguoi-ban/quan-ly-khach-hang"
  ],
});
</script>
