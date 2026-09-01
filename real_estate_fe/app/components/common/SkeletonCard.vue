<template>
  <!-- 1. Layout mặc định (type = 'list'): Thẻ danh sách tìm kiếm BĐS -->
  <div v-if="type === 'list'" class="flex flex-col gap-4">
    <div v-for="i in count" :key="i" class="animate-pulse bg-white shadow-sm border border-gray-100 rounded-lg overflow-hidden">
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
        <div class="h-4 w-3/4 bg-gray-200 rounded" />
        <div class="h-4 w-1/2 bg-gray-200 rounded" />
        <div class="flex gap-3">
          <div class="h-6 w-20 bg-gray-200 rounded" />
          <div class="h-6 w-16 bg-gray-200 rounded" />
        </div>
        <div class="flex items-center gap-2 border-t border-gray-100 pt-3">
          <div class="h-10 w-10 bg-gray-200 rounded-full" />
          <div class="space-y-1.5">
            <div class="h-3 w-24 bg-gray-200 rounded" />
            <div class="h-3 w-16 bg-gray-200 rounded" />
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- 2. Layout Card đơn giản (type = 'card'): Thẻ BĐS dành cho bạn trong grid 4 cột -->
  <template v-else-if="type === 'card'">
    <div v-for="i in count" :key="i" class="bg-white border border-gray-100 rounded-lg shadow-sm overflow-hidden animate-pulse flex flex-col">
      <!-- Ảnh bìa aspect-4/3 -->
      <div class="relative aspect-4/3 bg-gray-200"></div>
      <!-- Nội dung -->
      <div class="p-3 flex flex-col gap-2.5 flex-grow">
        <div class="h-4 bg-gray-200 rounded w-5/6"></div>
        <div class="h-4 bg-gray-200 rounded w-2/3"></div>
        <div class="flex items-center gap-2 mt-1">
          <div class="h-3.5 bg-gray-200 rounded w-20"></div>
          <div class="h-3.5 bg-gray-200 rounded w-12"></div>
        </div>
        <div class="h-3 bg-gray-100 rounded w-1/2"></div>
        <div class="flex items-center justify-between mt-2 pt-1.5 border-t border-gray-50">
          <div class="h-3 bg-gray-100 rounded w-16"></div>
          <div class="h-8 w-8 rounded-full bg-gray-200"></div>
        </div>
      </div>
    </div>
  </template>

  <!-- 3. Layout Dự án nổi bật (type = 'project'): Thẻ lưới cho dự án -->
  <template v-else-if="type === 'project'">
    <div v-for="i in count" :key="i" class="border border-gray-100 rounded-lg overflow-hidden animate-pulse flex flex-col bg-white shadow-sm">
      <div class="aspect-[4/3] bg-gray-200 w-full rounded-t-lg"></div>
      <div class="p-4 flex flex-col gap-2 flex-grow">
        <div class="h-3.5 bg-gray-200 rounded w-1/3 mb-1"></div>
        <div class="h-4 bg-gray-200 rounded w-5/6"></div>
        <div class="h-3 bg-gray-200 rounded w-1/2"></div>
        <div class="h-3 bg-gray-100 rounded w-2/3"></div>
      </div>
    </div>
  </template>

  <!-- 4. Layout BĐS theo địa điểm (type = 'area'): 1 thẻ to trái + lưới 2x2 phải -->
  <div v-else-if="type === 'area'" class="flex gap-3" style="height: 360px;">
    <!-- Cột trái: Skeleton Featured -->
    <div class="relative overflow-hidden rounded-lg animate-pulse bg-gray-200 shrink-0" style="flex: 0 0 45%;">
    </div>
    <!-- Cột phải: Skeleton 2x2 grid -->
    <div class="flex-1 grid grid-cols-2 grid-rows-2 gap-3">
        <div v-for="i in 4" :key="i" class="relative overflow-hidden rounded-lg animate-pulse bg-gray-200">
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  count?: number
  type?: 'list' | 'card' | 'project' | 'area'
}>(), {
  count: 1,
  type: 'list'
})
</script>
