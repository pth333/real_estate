<template>
  <div class="rounded-lg bg-white p-4 shadow-sm">
    <h3 class="mb-3 text-sm font-semibold text-gray-600">Phân bố giá / m²</h3>
    <Bar v-if="chartData" :data="chartData" :options="chartOptions" :height="200" />
    <p v-else class="py-8 text-center text-sm text-gray-400">Đang tải...</p>
  </div>
</template>

<script setup lang="ts">
import { Bar } from 'vue-chartjs'

import type { RealEstateResponse } from '~/types/real_estate'

const store = useRealEstateStore()

type PriceBucket = '< 1 tỷ' | '1-2 tỷ' | '2-3 tỷ' | '3-5 tỷ' | '5-10 tỷ' | '≥ 10 tỷ'

const chartData = computed(() => {
  const items = store.items as RealEstateResponse[]
  if (items && !items.length) return null

  const buckets: Record<PriceBucket, number> = {
    '< 1 tỷ': 0,
    '1-2 tỷ': 0,
    '2-3 tỷ': 0,
    '3-5 tỷ': 0,
    '5-10 tỷ': 0,
    '≥ 10 tỷ': 0,
  }
  if (!items) return null
  items.forEach((item) => {
    const priceInBillion = item.price_vnd / 1_000_000_000
    let key: PriceBucket
    if (priceInBillion < 1) key = '< 1 tỷ'
    else if (priceInBillion < 2) key = '1-2 tỷ'
    else if (priceInBillion < 3) key = '2-3 tỷ'
    else if (priceInBillion < 5) key = '3-5 tỷ'
    else if (priceInBillion < 10) key = '5-10 tỷ'
    else key = '≥ 10 tỷ'
    buckets[key]++
  })

  return {
    labels: Object.keys(buckets),
    datasets: [
      {
        label: 'Số lượng',
        backgroundColor: '#3b82f6',
        borderRadius: 4,
        data: Object.values(buckets),
      },
    ],
  }
})

const chartOptions = {
  responsive: true,
  plugins: {
    tooltip: {
      callbacks: {
        label: (ctx: { parsed: { y: number | null } }) => `${ctx.parsed.y} BĐS`,
      },
    },
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: { stepSize: 1 },
    },
  },
}
</script>
