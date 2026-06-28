<template>
  <div class="flex flex-col gap-4">
    <!-- Thống kê tổng quan -->
    <SummaryCards />

    <!-- Charts -->
    <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
      <PriceChart />
      <DistrictChart />
    </div>

    <!-- Bảng dữ liệu -->
    <RealEstateTable />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRealEstateStore } from '@/stores/real_estate'
import SummaryCards from '@/components/real-estate/SummaryCards.vue'
import PriceChart from '@/components/real-estate/PriceChart.vue'
import DistrictChart from '@/components/real-estate/DistrictChart.vue'
import RealEstateTable from '@/components/real-estate/RealEstateTable.vue'

const store = useRealEstateStore()

onMounted(() => {
  store.fetchList()
  store.fetchSummary()
})
</script>
