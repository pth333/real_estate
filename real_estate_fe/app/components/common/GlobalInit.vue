<template>
  <slot />
</template>

<script setup lang="ts">
import { useMessage } from "naive-ui"

// Global để sử dụng ở bất kỳ component, plugin, interceptor nào (Chỉ chạy ở Client-side)
if (import.meta.client) {
  window.message = useMessage()
}
</script>
