<template>
  <n-modal :show="showModalAll" preset="card" :style="{ maxWidth: '520px' }" mask-closable content-style="max-height: 70vh; overflow-y: auto;"
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
        <span class="text-base font-semibold">{{ screenTitle }}</span>
      </div>
    </template>

    <!-- Nội dung chính -->
    <div v-if="filterStore.screen === 'main'" class="flex flex-col gap-5">

      <!-- Khu vực / Khoảng giá / Diện tích (dạng ô input) -->
      <div v-for="field in inputFields" :key="field.label">
        <div class="text-sm font-semibold mb-2">{{ field.label }}</div>
        <div class="flex items-center justify-between w-full px-3 py-2.5 border border-gray-300 rounded-md bg-white cursor-pointer hover:border-emerald-400 transition-colors"
          @click="filterStore.screen = field.screen">
          <span class="text-sm" :class="field.isDefault ? 'text-gray-400' : 'text-gray-800'">{{ field.value }}</span>
          <IconChevronRight class="h-4 w-4 text-gray-400 shrink-0" />
        </div>
      </div>


      <!-- Số phòng ngủ -->
      <div>
        <div class="text-sm font-semibold mb-3">Số phòng ngủ</div>
        <div class="flex gap-2 flex-wrap">
          <button v-for="opt in roomOptions" :key="opt.value"
            class="px-4 py-1.5 border rounded-md text-sm transition-colors" :class="filterStore.filters.bedrooms === opt.value
              ? 'bg-red-500 text-white border-red-500'
              : 'bg-white text-gray-700 border-gray-300 hover:border-emerald-400'"
            @click="toggleSingle('bedrooms', opt.value)">
            {{ opt.label }}
          </button>
        </div>
      </div>


      <!-- Số phòng tắm -->
      <div>
        <div class="text-sm font-semibold mb-3">Số phòng tắm, vệ sinh</div>
        <div class="flex gap-2 flex-wrap">
          <button v-for="opt in roomOptions" :key="opt.value"
            class="px-4 py-1.5 border rounded-md text-sm transition-colors" :class="filterStore.filters.bathrooms === opt.value
              ? 'bg-red-500 text-white border-red-500'
              : 'bg-white text-gray-700 border-gray-300 hover:border-emerald-400'"
            @click="toggleSingle('bathrooms', opt.value)">
            {{ opt.label }}
          </button>
        </div>
      </div>


      <!-- Hướng nhà -->
      <div>
        <div class="text-sm font-semibold mb-3">Hướng nhà</div>
        <div class="flex gap-2 flex-wrap">
          <button v-for="opt in directionOptions" :key="opt.value"
            class="px-3 py-1.5 border rounded-md text-sm transition-colors" :class="filterStore.filters.house_direction === opt.value
              ? 'bg-red-500 text-white border-red-500'
              : 'bg-white text-gray-700 border-gray-300 hover:border-emerald-400'"
            @click="toggleSingle('house_direction', opt.value)">
            {{ opt.label }}
          </button>
        </div>
      </div>


      <!-- Hướng ban công -->
      <div>
        <div class="text-sm font-semibold mb-3">Hướng ban công</div>
        <div class="flex gap-2 flex-wrap">
          <button v-for="opt in directionOptions" :key="opt.value"
            class="px-3 py-1.5 border rounded-md text-sm transition-colors" :class="filterStore.filters.balcony_direction === opt.value
              ? 'bg-red-500 text-white border-red-500'
              : 'bg-white text-gray-700 border-gray-300 hover:border-emerald-400'"
            @click="toggleSingle('balcony_direction', opt.value)">
            {{ opt.label }}
          </button>
        </div>
      </div>


      <!-- Pháp lý -->
      <div>
        <div class="text-sm font-semibold mb-3">Pháp lý</div>
        <div class="flex gap-2 flex-wrap">
          <button v-for="opt in legalOptions" :key="opt.value"
            class="px-3 py-1.5 border rounded-md text-sm transition-colors" :class="filterStore.filters.legal_docs === opt.value
              ? 'bg-red-500 text-white border-red-500'
              : 'bg-white text-gray-700 border-gray-300 hover:border-emerald-400'"
            @click="toggleSingle('legal_docs', opt.value)">
            {{ opt.label }}
          </button>
        </div>
      </div>


      <!-- Nội thất -->
      <div>
        <div class="text-sm font-semibold mb-3">Nội thất</div>
        <div class="flex gap-2 flex-wrap">
          <button v-for="opt in interiorOptions" :key="opt.value"
            class="px-3 py-1.5 border rounded-md text-sm transition-colors" :class="filterStore.filters.interior === opt.value
              ? 'bg-red-500 text-white border-red-500'
              : 'bg-white text-gray-700 border-gray-300 hover:border-emerald-400'"
            @click="toggleSingle('interior', opt.value)">
            {{ opt.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- Màn hình chọn khu vực -->
    <div v-else-if="filterStore.screen === 'location'">
      <FilterLocation />
    </div>

    <!-- Màn hình khoảng giá -->
    <div v-else-if="filterStore.screen === 'price'">
      <FilterPrice />
    </div>

    <!-- Màn hình diện tích -->
    <div v-else-if="filterStore.screen === 'area'">
      <FilterArea />
    </div>

    <!-- Footer actions -->
    <template #footer>
      <div v-if="filterStore.screen === 'main'" class="flex justify-end gap-2 mt-2">
        <n-button  @click="handleReset">Đặt lại</n-button>
        <n-button type="primary" @click="handleApply">Áp dụng</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { useFilterStore } from '~/stores/filter';
import { useRealEstateStore } from '~/stores/real_estate';

const props = defineProps<{
  showModalAll: boolean
}>()
const emit = defineEmits<{
  (e: 'update:showModalAll', value: boolean): void
}>()
const filterStore = useFilterStore();
const realEstateStore = useRealEstateStore()
const route = useRoute();

// --- Screen title ---
const screenTitle = computed(() => {
  const map: Record<string, string> = {
    location: 'Khu vực',
    price: 'Khoảng giá',
    area: 'Diện tích',
  }
  return map[filterStore.screen] ?? ''
})

// --- Options ---
const roomOptions = [
  { label: '1', value: 1 },
  { label: '2', value: 2 },
  { label: '3', value: 3 },
  { label: '4', value: 4 },
  { label: '5+', value: 5 },
]

const directionOptions = [
  { label: 'Bắc', value: 'bac' },
  { label: 'Đông Bắc', value: 'dong_bac' },
  { label: 'Đông', value: 'dong' },
  { label: 'Đông Nam', value: 'dong_nam' },
  { label: 'Nam', value: 'nam' },
  { label: 'Tây Nam', value: 'tay_nam' },
  { label: 'Tây', value: 'tay' },
  { label: 'Tây Bắc', value: 'tay_bac' },
]

const legalOptions = [
  { label: 'Sổ đỏ/Sổ hồng', value: 'so_do' },
  { label: 'Hợp đồng mua bán', value: 'hop_dong' },
  { label: 'Đang chờ sổ', value: 'cho_so' },
  { label: 'Khác', value: 'khac' },
]

const interiorOptions = [
  { label: 'Đầy đủ', value: 'day_du' },
  { label: 'Cơ bản', value: 'co_ban' },
  { label: 'Không nội thất', value: 'khong' },
  { label: 'Khác', value: 'khac' },
]

const TY = 1_000_000_000

// --- Display labels ---
const priceLabel = computed(() => {
  const { min_price, max_price } = filterStore.filters
  if (!min_price && !max_price) return 'Tất cả'
  if (!min_price) return `Dưới ${max_price! / TY} tỷ`
  if (!max_price) return `Trên ${min_price / TY} tỷ`
  return `${min_price / TY} - ${max_price / TY} tỷ`
})

const areaLabel = computed(() => {
  const { min_area, max_area } = filterStore.filters
  if (!min_area && !max_area) return 'Tất cả'
  if (!min_area) return `Dưới ${max_area} m²`
  if (!max_area) return `Trên ${min_area} m²`
  return `${min_area} - ${max_area} m²`
})

// --- Location label ---
const locationLabel = computed(() => {
  const loc = filterStore.filters
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

// --- Các ô chọn dạng input (Khu vực / Khoảng giá / Diện tích) ---
const inputFields = computed(() => [
  {
    label: 'Khu vực',
    screen: 'location' as const,
    value: locationLabel.value,
    isDefault: locationLabel.value === 'Chọn khu vực',
  },
  {
    label: 'Khoảng giá',
    screen: 'price' as const,
    value: priceLabel.value,
    isDefault: priceLabel.value === 'Tất cả',
  },
  {
    label: 'Diện tích',
    screen: 'area' as const,
    value: areaLabel.value,
    isDefault: areaLabel.value === 'Tất cả',
  },
]);

// --- Helpers ---
function toggleSingle(field: keyof typeof filterStore.filters, value: any) {
  if (filterStore.filters[field] === value) {
    (filterStore.filters as any)[field] = typeof value === 'number' ? null : ''
  } else {
    (filterStore.filters as any)[field] = value
  }
}

function backToMain() {
  filterStore.screen = 'main';
}

function handleReset() {
  filterStore.filters = {};
}

const handleApply = () => {
  emit('update:showModalAll', false);

  realEstateStore.currentPage = 0;
  const catSlug: string = realEstateStore.categorySlug || (Array.isArray(route.params.category) ? route.params.category[0] ?? "" : route.params.category ?? "");
  console.log(filterStore.filters.city)
  const url = buildListUrl(
    catSlug,
    filterStore.filters,
    filterStore.cityOptions,
  );

  navigateTo(url);
};

// Reset screen khi mở modal
watch(() => props.showModalAll, (val) => {
  if (val) filterStore.screen = 'main';
});
</script>