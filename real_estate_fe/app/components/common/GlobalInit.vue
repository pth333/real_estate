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

// Nạp phiên đăng nhập (user + roles + permissions) ngay khi app khởi động.
// init() tự lo thứ tự: refresh token TRƯỚC (nếu cần) rồi mới gọi user-current-info,
// và chống gọi trùng nên middleware route gọi lại cũng không phát sinh request thứ 2.
const authStore = useAuthStore()

onMounted(() => {
  authStore.init()
})
</script>
