<template>
  <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
    <!-- Loading skeleton -->
    <template v-if="loading">
      <div v-for="i in 4" :key="i" class="rounded-lg bg-white p-4 shadow-sm">
        <div class="h-4 w-24 animate-pulse rounded bg-gray-200"></div>
        <div class="mt-2 h-8 w-28 animate-pulse rounded bg-gray-200"></div>
      </div>
    </template>

    <!-- Error state -->
    <template v-else-if="error">
      <div class="col-span-full rounded-lg bg-red-50 p-4 text-center text-sm text-red-600">
        {{ error }}
      </div>
    </template>

    <!-- Data -->
    <template v-else>
      <!-- Tổng số bài đăng -->
      <div class="rounded-lg bg-white p-4 shadow-sm">
        <p class="text-sm text-gray-500">Tổng số BĐS</p>
        <p class="mt-1 text-2xl font-bold text-blue-600">
          {{ store.summary?.total_posts ?? '—' }}
        </p>
      </div>

      <!-- Giá trung bình / m² -->
      <div class="rounded-lg bg-white p-4 shadow-sm">
        <p class="text-sm text-gray-500">Giá trung bình / m²</p>
        <p class="mt-1 text-2xl font-bold text-green-600">
          {{ formatCardPrice(store.summary?.avg_price_m2) }}
        </p>
      </div>

      <!-- Giá cao nhất / m² -->
      <div class="rounded-lg bg-white p-4 shadow-sm">
        <p class="text-sm text-gray-500">Giá cao nhất / m²</p>
        <p class="mt-1 text-2xl font-bold text-orange-600">
          {{ formatCardPrice(store.summary?.max_price_m2) }}
        </p>
      </div>

      <!-- Giá thấp nhất / m² -->
      <div class="rounded-lg bg-white p-4 shadow-sm">
        <p class="text-sm text-gray-500">Giá thấp nhất / m²</p>
        <p class="mt-1 text-2xl font-bold text-red-600">
          {{ formatCardPrice(store.summary?.min_price_m2) }}
        </p>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useRealEstateStore } from '~/stores/real_estate'

const props = withDefaults(defineProps<{
  loading?: boolean
  error?: string | null
}>(), {
  loading: false,
  error: null,
})

const store = useRealEstateStore()
const { formatPriceCompact } = usePriceFormatter()

function formatCardPrice(price?: number): string {
  if (price == null) return '—'
  return formatPriceCompact(price)
}
</script>
