<template>
  <!-- Trigger button + popover -->
  <n-popover placement="bottom-start" trigger="click" :show="filterStore.showPopover"
    @update:show="filterStore.showPopover = $event" :style="{ padding: 0 }">
    <template #trigger>
      <n-button ghost size="small">
        <template #icon>
          <IconChevronDown class="h-4 w-4" />
        </template>
        {{ buttonLabel }}
      </n-button>
    </template>

    <div class="w-80 p-4">
      <h3 class="mb-3 text-sm font-semibold text-gray-700">Khoảng giá</h3>

      <!-- Inputs -->
      <div class="mb-3 flex items-center gap-2">
        <n-input-number v-model:value="localMin" placeholder="Từ (VNĐ)" :min="0" clearable class="flex-1"
          size="small" />
        <span class="text-gray-400">—</span>
        <n-input-number v-model:value="localMax" placeholder="Đến (VNĐ)" :min="0" clearable class="flex-1"
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
import type { FilterPriceRange } from '~/types/real_estate';
import { useFilterStore } from '~/stores/filter';

const filterStore = useFilterStore();

const localMin = ref<number | null>(null);
const localMax = ref<number | null>(null);

// Mở popover → đồng bộ local state từ store
function syncLocalFromStore() {
  localMin.value = filterStore.filterPriceRange?.min_price ?? null;
  localMax.value = filterStore.filterPriceRange?.max_price ?? null;
}

watch(() => filterStore.showPopover, (val) => {
  if (val) syncLocalFromStore();
});

interface PricePreset {
  label: string
  range: FilterPriceRange
}

const presets: PricePreset[] = [
  { label: 'Dưới 500 triệu', range: { max_price: 500_000_000 } },
  { label: '500 tr - 1 tỷ', range: { min_price: 500_000_000, max_price: 1_000_000_000 } },
  { label: '1 tỷ - 3 tỷ', range: { min_price: 1_000_000_000, max_price: 3_000_000_000 } },
  { label: '3 tỷ - 5 tỷ', range: { min_price: 3_000_000_000, max_price: 5_000_000_000 } },
  { label: '5 tỷ - 10 tỷ', range: { min_price: 5_000_000_000, max_price: 10_000_000_000 } },
  { label: 'Trên 10 tỷ', range: { min_price: 10_000_000_000 } },
];

const { formatPrice } = usePriceFormatter();

const buttonLabel = computed(() => {
  const range = filterStore.filterPriceRange;
  if (range.min_price != null && range.max_price != null) {
    return `${formatPrice(range.min_price)} - ${formatPrice(range.max_price)}`;
  }
  if (range.min_price != null) {
    return `Từ ${formatPrice(range.min_price)}`;
  }
  if (range.max_price != null) {
    return `Đến ${formatPrice(range.max_price)}`;
  }
  return 'Khoảng giá';
});

function isPresetActive(preset: PricePreset): boolean {
  return (
    localMin.value === (preset.range.min_price ?? null) &&
    localMax.value === (preset.range.max_price ?? null)
  );
}

function selectPreset(preset: PricePreset) {
  localMin.value = preset.range.min_price ?? null;
  localMax.value = preset.range.max_price ?? null;
}

function handleReset() {
  localMin.value = null;
  localMax.value = null;
  filterStore.filterPriceRange = {};
}

function handleApply() {
  filterStore.filterPriceRange = {
    min_price: localMin.value ?? undefined,
    max_price: localMax.value ?? undefined,
  };
  filterStore.showPopover = false;

  // Commit lên filters.price_range để đồng nhất
  filterStore.filters.price_range = { ...filterStore.filterPriceRange };
  filterStore.filters.location = { ...filterStore.filterLocation };

  // Reset URL về trang 1 → trigger fetch
  const route = useRoute();
  const slug = route.params.slug;
  const slugStr = slug ? (Array.isArray(slug) ? slug[0] : slug) : '';
  navigateTo(`/${slugStr}`);
}
</script>
