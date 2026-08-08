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
        <n-input-number v-model:value="localMin" placeholder="Từ (m²)" :min="0" clearable class="flex-1"
          size="small" />
        <span class="text-gray-400">—</span>
        <n-input-number v-model:value="localMax" placeholder="Đến (m²)" :min="0" clearable class="flex-1"
          size="small" />
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
import type { FilterAreaRange } from '~/types/real_estate';
import { useFilterStore } from '~/stores/filter';
import { useRealEstateStore } from '~/stores/real_estate';
import { buildCategoryPath } from '~/utils/slug';

const filterStore = useFilterStore();
const realEstateStore = useRealEstateStore();
const route = useRoute();

const localMin = ref<number | null>(null);
const localMax = ref<number | null>(null);

// Mở popover → đồng bộ local state từ store
function syncLocalFromStore() {
  localMin.value = filterStore.filterAreaRange?.min_acreage ?? null;
  localMax.value = filterStore.filterAreaRange?.max_acreage ?? null;
}

watch(() => filterStore.showAreaPopover, (val) => {
  if (val) syncLocalFromStore();
});

interface AreaPreset {
  label: string
  range: FilterAreaRange
}

const presets: AreaPreset[] = [
  { label: 'Dưới 30 m²', range: { max_acreage: 30 } },
  { label: '30 - 50 m²', range: { min_acreage: 30, max_acreage: 50 } },
  { label: '50 - 100 m²', range: { min_acreage: 50, max_acreage: 100 } },
  { label: '100 - 200 m²', range: { min_acreage: 100, max_acreage: 200 } },
  { label: 'Trên 200 m²', range: { min_acreage: 200 } },
];

const buttonLabel = computed(() => {
  const range = filterStore.filterAreaRange;
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
    localMin.value === (preset.range.min_acreage ?? null) &&
    localMax.value === (preset.range.max_acreage ?? null)
  );
}

function selectPreset(preset: AreaPreset) {
  localMin.value = preset.range.min_acreage ?? null;
  localMax.value = preset.range.max_acreage ?? null;
}

function handleReset() {
  localMin.value = null;
  localMax.value = null;
  filterStore.filterAreaRange = {};
}

function handleApply() {
  filterStore.filterAreaRange = {
    min_acreage: localMin.value ?? undefined,
    max_acreage: localMax.value ?? undefined,
  };
  filterStore.showAreaPopover = false;

  // Commit lên filters để đồng nhất
  filterStore.filters.min_acreage = filterStore.filterAreaRange.min_acreage;
  filterStore.filters.max_acreage = filterStore.filterAreaRange.max_acreage;
  filterStore.filters.location = { ...filterStore.filterLocation };

  // Build URL SEO server-driven → trigger fetch. Location = tên city (segment [location]).
  const catSlug: string = realEstateStore.categorySlug || (Array.isArray(route.params.category) ? route.params.category[0] ?? "" : route.params.category ?? "");
  const url = buildCategoryPath(
    catSlug,
    filterStore.filterLocation,
    filterStore.cityOptions,
    filterStore.filterPriceRange,
    filterStore.filterAreaRange,
    1,
  );
  navigateTo(url);
}
</script>
