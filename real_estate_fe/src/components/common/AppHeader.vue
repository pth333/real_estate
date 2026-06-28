<template>
  <header class="bg-white shadow-sm border-b border-gray-200">
    <div class="container mx-auto flex justify-between items-center p-4">
      <h1 class="text-xl font-bold text-blue-600">RealEstate</h1>

      <nav class="flex items-center gap-6 text-gray-700">
        <NotificationBell />
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import NotificationBell from '@/components/notification/NotificationBell.vue'
</script>
