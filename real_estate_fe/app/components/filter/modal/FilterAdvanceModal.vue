<template>
  <n-modal :show="showModalAll" preset="card" :style="{ maxWidth: '520px' }" mask-closable
    @update:show="emit('update:showModalAll', $event)">

    <!-- Header động -->
    <template #header>
      <div v-if="filterStore.screen === 'main'" class="text-base font-semibold">Bộ lọc nâng cao</div>
      <div v-else class="flex items-center gap-2">
        <n-button text @click="backToMain">
          <template #icon>
            <IconChevronRight class="h-4 w-4 rotate-180" />
          </template>
        </n-button>
        <span class="text-base font-semibold">Khu vực</span>
      </div>
    </template>

    <!-- Nội dung chính -->
    <div v-if="filterStore.screen === 'main'">
      <n-form label-placement="top" class="mt-2">
        <n-form-item label="Khu vực">
          <n-button text @click="openLocation">
            <template #icon>
              <IconChevronRight class="h-4 w-4" />
            </template>
            {{ locationLabel }}
          </n-button>
        </n-form-item>

      </n-form>
    </div>

    <!-- Màn hình chọn khu vực -->
    <div v-else-if="filterStore.screen === 'location'">
      <FilterLocation />
    </div>

    <!-- Footer actions -->
    <template #footer>
      <div v-if="filterStore.screen === 'main'" class="flex justify-end gap-2">
        <n-button quaternary @click="handleReset">Đặt lại</n-button>
        <n-button type="primary" @click="handleApply">Áp dụng</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { useFilterStore } from '~/stores/filter';
import { useRealEstateStore } from '~/stores/real_estate';
import { buildCategoryPath } from '~/utils/slug';

const props = defineProps<{
  showModalAll: boolean
}>()
const emit = defineEmits<{
  (e: 'update:showModalAll', value: boolean): void
}>()
const filterStore = useFilterStore();
const realEstateStore = useRealEstateStore()
const route = useRoute();

const locationLabel = computed(() => {
  const loc = filterStore.filterLocation
  const city = filterStore.cityOptions.find(item => item.value === loc.city)
  const district = filterStore.districtOptions.find(item => item.value === loc.district)
  const ward = filterStore.wardOptions.find(item => item.value === loc.ward)

  const parts: string[] = [];
  if (city) parts.push(city.label as string);
  if (district) parts.push(district.label as string);
  if (ward) parts.push(ward.label as string);

  if (parts.length === 0) return 'Chọn khu vực';
  return parts.join(', ');
});

function openLocation() {
  filterStore.screen = 'location';
}

function backToMain() {
  filterStore.screen = 'main';
}


function handleReset() {
  // filterStore.filters.price_range = {};
  filterStore.filterLocation = {};
  console.log(realEstateStore.categorySlug)
}

const handleApply = () => {
  // Commit draft location vào filters actual
  filterStore.filters.location = { ...filterStore.filterLocation };
  emit('update:showModalAll', false);

  // Reset URL về trang 1 để trigger fetch lại — vì filters đã thay đổi
  realEstateStore.currentPage = 0;
  const catSlug: string = realEstateStore.categorySlug || (Array.isArray(route.params.category) ? route.params.category[0] ?? "" : route.params.category ?? "");
  const url = buildCategoryPath(
    catSlug,
    filterStore.filterLocation,
    filterStore.cityOptions,
    filterStore.filterPriceRange,
    filterStore.filterAreaRange,
    1,
  );
  console.log(url)
  navigateTo(url);
};

// Reset screen khi mở modal
watch(() => props.showModalAll, (val) => {
  if (val) filterStore.screen = 'main';
});
</script>
