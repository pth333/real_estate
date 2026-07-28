<template>
  <slot />
</template>

<script setup lang="ts">
import { useMessage } from "naive-ui"

// Global để sử dụng ở bất kỳ component, plugin, interceptor nào
window.message = useMessage()
</script>
