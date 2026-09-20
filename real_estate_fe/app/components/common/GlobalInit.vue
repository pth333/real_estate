<template>
  <slot />
</template>

<script setup lang="ts">
import { useMessage } from "naive-ui"
import { useAuthStore } from "~/stores/auth"

// Global để sử dụng ở bất kỳ component, plugin, interceptor nào (Chỉ chạy ở Client-side)
if (import.meta.client) {
  window.message = useMessage()
}

// Nạp thông tin user hiện tại (roles[] + permissions[]) vào state global ngay khi
// app khởi động, để middleware route và UI luôn có quyền MỚI NHẤT — admin vừa đổi
// role thì chỉ cần tải lại trang, không phải đăng nhập lại.
// Dùng ensureAccessLoaded (có chống gọi trùng) thay vì fetchCurrentUser trực tiếp,
// vì middleware route cũng đã yêu cầu cùng dữ liệu này lúc tải trang.
const authStore = useAuthStore()

onMounted(() => {
  authStore.ensureAccessLoaded()
})
</script>
