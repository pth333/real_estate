<template>
  <div class="grid grid-cols-1 lg:grid-cols-1 gap-3">
    <CandleChart />
    <StockTable />
  </div>
</template>

<script setup lang="ts">
import CandleChart from '../charts/CandleChart.vue'
import StockTable from '../tables/StockTable.vue'
</script>
