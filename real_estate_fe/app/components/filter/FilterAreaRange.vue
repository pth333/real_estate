<template>
  <!-- Trigger button + popover -->
  <n-popover placement="bottom-start" trigger="click" :show="filterStore.showAreaPopover"
    @update:show="filterStore.showAreaPopover = $event" :style="{ padding: 0 }">
    <template #trigger>
      <n-button ghost size="small">
        <template #icon>
          <IconChevronDownOutline class="h-4 w-4" />
        </template>
        {{ buttonLabel }}
      </n-button>
    </template>

    <div class="w-80 p-4">
      <h3 class="mb-3 text-sm font-semibold text-gray-700">Diện tích</h3>

      <!-- Inputs -->
      <div class="mb-3 flex items-center gap-2">
        <n-input-number v-model:value="filterStore.filters.min_acreage" placeholder="Từ (m²)" :min="0" clearable
          class="flex-1" size="small" />
        <span class="text-gray-400">—</span>
        <n-input-number v-model:value="filterStore.filters.max_acreage" placeholder="Đến (m²)" :min="0" clearable
          class="flex-1" size="small" />
      </div>

      <!-- Preset ratios -->
      <div class="mb-3 grid grid-cols-2 gap-1.5">
        <button v-for="preset in presets" :key="preset.label" :class="[
          'rounded-md px-3 py-1.5 text-xs transition-colors',
          isPresetActive(preset)
            ? 'bg-emerald-600 text-white'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200',
        ]" @click="selectPreset(preset)">
          {{ preset.label }}
        </button>
      </div>

      <!-- Actions -->
      <div class="flex justify-end gap-2 border-t pt-3">
        <n-button quaternary size="tiny" @click="handleReset">Đặt lại</n-button>
        <n-button type="primary" size="tiny" @click="handleApply">Áp dụng</n-button>
      </div>
    </div>
  </n-popover>
</template>

<script setup lang="ts">
import { useFilterStore } from '~/stores/filter';
import { useRealEstateStore } from '~/stores/real_estate';
import type { Filter } from "~/types/real_estate";

const filterStore = useFilterStore();
const realEstateStore = useRealEstateStore();
const route = useRoute();

// const localMin = ref<number | null>(null);
// const localMax = ref<number | null>(null);

// Mở popover → đồng bộ local state từ store
// function syncLocalFromStore() {
//   localMin.value = filterStore.filters?.min_acreage ?? null;
//   localMax.value = filterStore.filters?.max_acreage ?? null;
// }

// watch(() => filterStore.showAreaPopover, (val) => {
//   if (val) syncLocalFromStore();
// });

interface AreaPreset {
  label: string
  range: Filter
}

const presets: AreaPreset[] = [
  { label: 'Dưới 30 m²', range: { max_acreage: 30 } },
  { label: '30 - 50 m²', range: { min_acreage: 30, max_acreage: 50 } },
  { label: '50 - 100 m²', range: { min_acreage: 50, max_acreage: 100 } },
  { label: '100 - 200 m²', range: { min_acreage: 100, max_acreage: 200 } },
  { label: 'Trên 200 m²', range: { min_acreage: 200 } },
];

const buttonLabel = computed(() => {
  const range = filterStore.filters;
  if (range.min_acreage != null && range.max_acreage != null) {
    return `${range.min_acreage} - ${range.max_acreage} m²`;
  }
  if (range.min_acreage != null) {
    return `Từ ${range.min_acreage} m²`;
  }
  if (range.max_acreage != null) {
    return `Đến ${range.max_acreage} m²`;
  }
  return 'Diện tích';
});

function isPresetActive(preset: AreaPreset): boolean {
  return (
    filterStore.filters.min_acreage === (preset.range.min_acreage ?? null) &&
    filterStore.filters.max_acreage === (preset.range.max_acreage ?? null)
  );
}

function selectPreset(preset: AreaPreset) {
  filterStore.filters.min_acreage = preset.range.min_acreage ?? undefined;
  filterStore.filters.max_acreage = preset.range.max_acreage ?? undefined;
}

function handleReset() {
  filterStore.filters.min_acreage = undefined
  filterStore.filters.max_acreage = undefined
}

function handleApply() {
  filterStore.showAreaPopover = false;

  // Build URL SEO server-driven → trigger fetch. Location = tên city (segment [location]).
  const catSlug: string = realEstateStore.categorySlug || (Array.isArray(route.params.category) ? route.params.category[0] ?? "" : route.params.category ?? "");
  const url = buildListUrl(
    catSlug,
    filterStore.filters,
    filterStore.cityOptions,
  );
  navigateTo(url);
}
</script>
