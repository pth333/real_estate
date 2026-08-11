<template>
  <!-- Lọc theo loại bất động sản (Bang filter — hiển thị mọi trang) -->
  <n-popover placement="bottom-start" trigger="click" :show="showPopover"
    @update:show="showPopover = $event" :style="{ padding: 0 }">
    <template #trigger>
      <n-button ghost size="small">
        <template #icon>
          <IconChevronDownOutline class="h-4 w-4" />
        </template>
        {{ buttonLabel }}
      </n-button>
    </template>

    <!-- Danh sách category lấy từ bảng categories (GET /real-estate/list/types) -->
    <div class="w-72 p-3">
      <n-input v-model:value="keyword" placeholder="Tìm loại bất động sản..." clearable size="small"
        class="mb-2" />
      <div class="max-h-72 overflow-y-auto">  
        <button v-for="cat in filteredCategories" :key="cat.id"
          class="flex w-full items-center justify-between px-3 py-2 text-left text-sm transition-colors"
          :class="isActive(cat.id) ? 'bg-emerald-600 text-white' : 'text-gray-600 hover:bg-gray-100'"
          @click="toggleCategory(cat.id)">
          <span>{{ cat.name }}</span>
          <span v-if="isActive(cat.id)">✓</span>
        </button>
        <p v-if="filteredCategories.length === 0" class="px-3 py-2 text-sm text-gray-400">
          Không tìm thấy loại BĐS
        </p>
      </div>
    </div>
  </n-popover>
</template>

<script setup lang="ts">
import type { OptionTypeRealestate } from '~/types/real_estate';
import { useFilterStore } from '~/stores/filter';
import { useRealEstateStore } from '~/stores/real_estate';
import { buildListUrl } from '~/utils/slug';

const { $api } = useNuxtApp();
const filterStore = useFilterStore();
const realEstateStore = useRealEstateStore();
const route = useRoute();

const showPopover = ref(false);
const keyword = ref('');
const categories = ref<OptionTypeRealestate[]>([]);

const filteredCategories = computed(() => {
  const kw = keyword.value.trim().toLowerCase();
  if (!kw) return categories.value;
  return categories.value.filter((cat) => cat.name.toLowerCase().includes(kw));
});

const activeCategoryId = computed(() => filterStore.filters.category_id ?? null);

const buttonLabel = computed(() => {
  const id = activeCategoryId.value;
  const found = categories.value.find((c) => c.id === id);
  return found?.name ?? 'Loại BĐS';
});

function isActive(id: number): boolean {
  return activeCategoryId.value === id;
}

// Chọn/bỏ chọn 1 loại → đổi query string (server-driven)
function toggleCategory(id: number) {
  if (activeCategoryId.value === id) {
    filterStore.filters.category_id = undefined;
  } else {
    filterStore.filters.category_id = id;
  }
  showPopover.value = false;
  navigateList();
}

// Build URL server-driven: segment (category-city + giá + diện tích) + query
// Lấy slug category thuần để build URL (category + city là segment đầu tiên)
function navigateList() {
  const catSlug = getCategorySlug();
  navigateTo(buildListUrl(catSlug, filterStore.filters, filterStore.cityOptions));
}

function getCategorySlug(): string {
  const cat = route.params.category;
  const raw = Array.isArray(cat) ? cat[0] ?? '' : cat ?? '';
  return realEstateStore.categorySlug || raw;
}

// Đồng bộ active từ URL (lỡ category_id thay đổi ngoài component)
watch(
  () => route.query.category_id,
  (v) => {
    const raw = Array.isArray(v) ? v[0] : v;
    const id = raw ? Number(raw) : null;
    if (id && !Number.isNaN(id)) {
      filterStore.filters.category_id = id;
    } else {
      filterStore.filters.category_id = undefined;
    }
  },
  { immediate: true },
);

onMounted(async () => {
  try {
    const res = await $api.get<{ data: OptionTypeRealestate[] }>('/real-estate/list/types');
    categories.value = res.data || [];
  } catch (err) {
    window.message?.error(err instanceof Error ? err.message : 'Lỗi tải loại BĐS');
  }
});
</script>