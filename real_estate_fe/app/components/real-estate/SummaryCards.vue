<template>
  <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
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
        {{ formatPrice(store.summary?.avg_price_m2) }}
      </p>
    </div>

    <!-- Giá cao nhất / m² -->
    <div class="rounded-lg bg-white p-4 shadow-sm">
      <p class="text-sm text-gray-500">Giá cao nhất / m²</p>
      <p class="mt-1 text-2xl font-bold text-orange-600">
        {{ formatPrice(store.summary?.max_price_m2) }}
      </p>
    </div>

    <!-- Giá thấp nhất / m² -->
    <div class="rounded-lg bg-white p-4 shadow-sm">
      <p class="text-sm text-gray-500">Giá thấp nhất / m²</p>
      <p class="mt-1 text-2xl font-bold text-red-600">
        {{ formatPrice(store.summary?.min_price_m2) }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRealEstateStore } from '@/stores/real_estate'

const store = useRealEstateStore()

function formatPrice(price?: number): string {
  if (price == null) return '—'
  if (price >= 1_000_000_000) return (price / 1_000_000_000).toFixed(1) + ' tỷ'
  if (price >= 1_000_000) return (price / 1_000_000).toFixed(0) + ' tr'
  return price.toLocaleString() + ' đ'
}
</script>
