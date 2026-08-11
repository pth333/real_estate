<template>
  <div class="flex flex-col gap-4">
    <div v-for="i in count" :key="i" class="animate-pulse bg-white shadow-sm">
      <!-- Grid ảnh skeleton -->
      <div class="grid h-60 grid-cols-[2fr_1fr] gap-0.5 bg-gray-100">
        <div class="col-span-2 row-span-2 bg-gray-200" />
        <div class="grid grid-rows-2 gap-0.5">
          <div class="bg-gray-200" />
          <div class="bg-gray-200" />
        </div>
      </div>
      <!-- Content skeleton -->
      <div class="space-y-3 p-4">
        <div class="h-4 w-3/4 bg-gray-200" />
        <div class="h-4 w-1/2 bg-gray-200" />
        <div class="flex gap-3">
          <div class="h-6 w-20 bg-gray-200" />
          <div class="h-6 w-16 bg-gray-200" />
        </div>
        <div class="flex items-center gap-2 border-t border-gray-100 pt-3">
          <div class="h-10 w-10 bg-gray-200" />
          <div class="space-y-1.5">
            <div class="h-3 w-24 bg-gray-200" />
            <div class="h-3 w-16 bg-gray-200" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  count?: number
}>()
</script>
