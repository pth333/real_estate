<template>
  <div class="flex flex-col min-h-screen gap-3">
    <AppHeader />

    <main class="flex-1 p-3">
      <slot />
    </main>

    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/common/AppHeader.vue'
import AppFooter from '@/components/common/AppFooter.vue'
</script>
