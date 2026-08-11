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
        <n-input-number v-model:value="filterStore.filters.min_area" placeholder="Từ (m²)" :min="0" clearable
          class="flex-1" size="small" />
        <span class="text-gray-400">—</span>
        <n-input-number v-model:value="filterStore.filters.max_area" placeholder="Đến (m²)" :min="0" clearable
          class="flex-1" size="small" />
      </div>

      <!-- Preset ratios -->
      <div class="mb-3 grid grid-cols-2 gap-1.5">
        <button v-for="preset in presets" :key="preset.label" :class="[
          ' px-3 py-1.5 text-xs transition-colors',
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


interface AreaPreset {
  label: string
  range: Filter
}

const presets: AreaPreset[] = [
  { label: 'Dưới 30 m²', range: { max_area: 30 } },
  { label: '30 - 50 m²', range: { min_area: 30, max_area: 50 } },
  { label: '50 - 100 m²', range: { min_area: 50, max_area: 100 } },
  { label: '100 - 200 m²', range: { min_area: 100, max_area: 200 } },
  { label: 'Trên 200 m²', range: { min_area: 200 } },
];

const buttonLabel = computed(() => {
  const range = filterStore.filters;
  if (range.min_area != null && range.max_area != null) {
    return `${range.min_area} - ${range.max_area} m²`;
  }
  if (range.min_area != null) {
    return `Từ ${range.min_area} m²`;
  }
  if (range.max_area != null) {
    return `Đến ${range.max_area} m²`;
  }
  return 'Diện tích';
});

function isPresetActive(preset: AreaPreset): boolean {
  return (
    filterStore.filters.min_area === (preset.range.min_area ?? null) &&
    filterStore.filters.max_area === (preset.range.max_area ?? null)
  );
}

function selectPreset(preset: AreaPreset) {
  filterStore.filters.min_area = preset.range.min_area ?? undefined;
  filterStore.filters.max_area = preset.range.max_area ?? undefined;
}

function handleReset() {
  filterStore.filters.min_area = undefined
  filterStore.filters.max_area = undefined
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
