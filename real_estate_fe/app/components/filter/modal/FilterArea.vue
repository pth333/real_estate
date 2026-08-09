<template>
  <div class="flex flex-col gap-4">

    <!-- Input thủ công -->
    <div class="flex items-center gap-3">
      <n-input-number
        v-model:value="filterStore.filters.min_area"
        placeholder="Từ"
        :min="0"
        :show-button="false"
        class="flex-1"
      >
        <template #suffix>m²</template>
      </n-input-number>
      <span class="text-gray-400">→</span>
      <n-input-number
        v-model:value="filterStore.filters.max_area"
        placeholder="Đến"
        :min="0"
        :show-button="false"
        class="flex-1"
      >
        <template #suffix>m²</template>
      </n-input-number>
    </div>

    <!-- Radio ranges -->
    <div class="flex flex-col">
      <label
        v-for="opt in areaRanges"
        :key="opt.label"
        class="flex items-center justify-between py-3 border-b border-gray-100 cursor-pointer"
        @click="selectRange(opt)"
      >
        <span class="text-sm">{{ opt.label }}</span>
        <n-radio
          :checked="isSelected(opt)"
          @update:checked="selectRange(opt)"
        />
      </label>
    </div>

    <!-- Nút Áp dụng -->
    <div class="mt-2 flex justify-end">
      <n-button type="primary" @click="onApply">
        Áp dụng
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useFilterStore } from '~/stores/filter';

const filterStore = useFilterStore()

const areaRanges = [
  { label: 'Tất cả diện tích', min: undefined, max: undefined },
  { label: 'Dưới 30 m²', min: undefined, max: 30 },
  { label: '30 - 50 m²', min: 30, max: 50 },
  { label: '50 - 100 m²', min: 50, max: 100 },
  { label: '100 - 200 m²', min: 100, max: 200 },
  { label: 'Trên 200 m²', min: 200, max: undefined },
]

function isSelected(opt: typeof areaRanges[0]) {
  if (opt.min === undefined && opt.max === undefined) {
    return !filterStore.filters.min_area && !filterStore.filters.max_area
  }
  return filterStore.filters.min_area === opt.min && filterStore.filters.max_area === opt.max
}

function selectRange(opt: typeof areaRanges[0]) {
  filterStore.filters.min_area = opt.min
  filterStore.filters.max_area = opt.max
}

const onApply = () => {
  filterStore.screen = 'main'
}
</script>