<template>
  <div class="rounded-lg bg-white p-4 shadow-sm">
    <h3 class="mb-3 text-sm font-semibold text-gray-600">Thống kê theo quận</h3>
    <Bar v-if="chartData" :data="chartData" :options="chartOptions" :height="200" />
    <p v-else class="py-8 text-center text-sm text-gray-400">Đang tải...</p>
  </div>
</template>

<script setup lang="ts">
import { Bar } from 'vue-chartjs'

const store = useRealEstateStore()

const chartData = computed(() => {
  const items = store.items
  if (!items.length) return null

  const districtMap: Record<string, number> = {}
  items.forEach((item) => {
    const d = item.District || 'Không xác định'
    districtMap[d] = (districtMap[d] || 0) + 1
  })

  const sorted = Object.entries(districtMap)
    .sort(([, a], [, b]) => b - a)
    .slice(0, 10)

  return {
    labels: sorted.map(([name]) => name),
    datasets: [
      {
        label: 'Số lượng',
        backgroundColor: '#10b981',
        borderRadius: 4,
        data: sorted.map(([, count]) => count),
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
