<template>
  <div class="rounded-lg bg-white p-4 shadow-sm">
    <h3 class="mb-3 text-sm font-semibold text-gray-600">
      Phân bố giá / m²
    </h3>
    <Bar
      v-if="chartData"
      :data="chartData"
      :options="chartOptions"
      :height="200"
    />
    <p v-else class="py-8 text-center text-sm text-gray-400">Đang tải...</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
} from 'chart.js'
import { useRealEstateStore } from '@/stores/real_estate'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip)

const store = useRealEstateStore()

// Nhóm giá theo khoảng (tỷ)
const chartData = computed(() => {
  const items = store.items
  if (!items.length) return null

  // Tạo bucket: <1 tỷ, 1-2 tỷ, 2-3 tỷ, 3-5 tỷ, 5-10 tỷ, >=10 tỷ
  const buckets: Record<string, number> = {
    '< 1 tỷ': 0,
    '1-2 tỷ': 0,
    '2-3 tỷ': 0,
    '3-5 tỷ': 0,
    '5-10 tỷ': 0,
    '≥ 10 tỷ': 0,
  }

  items.forEach((item) => {
    const priceInBillion = item.PriceVND / 1_000_000_000
    if (priceInBillion < 1) buckets['< 1 tỷ']++
    else if (priceInBillion < 2) buckets['1-2 tỷ']++
    else if (priceInBillion < 3) buckets['2-3 tỷ']++
    else if (priceInBillion < 5) buckets['3-5 tỷ']++
    else if (priceInBillion < 10) buckets['5-10 tỷ']++
    else buckets['≥ 10 tỷ']++
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
        label: (ctx: { parsed: { y: number } }) => `${ctx.parsed.y} BĐS`,
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
