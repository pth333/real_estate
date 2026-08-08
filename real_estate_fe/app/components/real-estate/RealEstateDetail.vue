<template>
  <div class="mx-auto max-w-[1000px] px-4 py-6">
    <!-- Loading -->
    <SkeletonCard v-if="loading" />

    <!-- 404 -->
    <div v-else-if="!listing" class="px-6 py-16 text-center">
      <p class="text-base text-gray-400">Không tìm thấy tin đăng</p>
      <n-button class="mt-4" type="primary" @click="navigateTo('/')">Quay lại</n-button>
    </div>

    <!-- Detail -->
    <div v-else class="space-y-4">

      <!-- Gallery -->
      <div class="overflow-hidden  bg-gray-100">
        <!-- Main image -->
        <div class="relative">
          <img v-if="mainImage" :src="mainImage" :alt="listing.title" class="h-[420px] w-full object-cover" />
          <div v-else class="flex h-[420px] w-full items-center justify-center bg-gray-200 text-gray-400">
            Không có ảnh
          </div>
          <!-- Nav arrow -->
          <button v-if="allImages.length > 1"
            class="absolute right-3 top-1/2 -translate-y-1/2 flex h-9 w-9 items-center justify-center rounded bg-white/80 shadow hover:bg-white"
            @click="nextImage">
            <IconChevronRight class="h-5 w-5 text-gray-700" />
          </button>
          <!-- Counter -->
          <span v-if="allImages.length > 0"
            class="absolute bottom-3 right-3 rounded bg-black/50 px-2 py-0.5 text-xs text-white">
            {{ activeImageIndex + 1 }} / {{ allImages.length }}
          </span>
        </div>

        <!-- Thumbnails -->
        <div v-if="allImages.length > 1" class="flex gap-2 bg-gray-800 p-2">
          <button v-for="(img, i) in allImages" :key="i"
            class="h-16 w-24 flex-shrink-0 overflow-hidden rounded border-2 transition"
            :class="activeImageIndex === i ? 'border-red-500' : 'border-transparent opacity-70 hover:opacity-100'"
            @click="activeImageIndex = i">
            <img :src="img" :alt="`Ảnh ${i + 1}`" class="h-full w-full object-cover" />
          </button>
        </div>
      </div>

      <!-- Breadcrumb -->
      <nav class="text-xs text-blue-600">
        <span>Bán</span>
        <span class="mx-1 text-gray-400">/</span>
        <span>{{ listing.city }}</span>
        <span v-if="listing.district">
          <span class="mx-1 text-gray-400">/</span>
          <span>{{ listing.district }}</span>
        </span>
      </nav>

      <!-- Tiêu đề + Giá -->
      <div class="border border-gray-200 bg-white p-5">
        <h1 class="text-xl font-bold leading-snug text-gray-800">{{ listing.title }}</h1>

        <!-- Địa chỉ -->
        <p class="mt-2 flex items-start gap-1 text-sm">
          <IconLocationOutline class="mt-0.5 h-4 w-4 flex-shrink-0 text-gray-500" />
          <span class="text-gray-600">{{ listing.address }}</span>
        </p>

        <!-- Giá / Diện tích / Phòng ngủ -->
        <div class="mt-4 flex flex-wrap items-end gap-6">
          <div>
            <p class="text-xs text-gray-500">Khoảng giá</p>
            <p class="text-2xl font-bold text-gray-800">{{ formatPrice(listing.price_vnd) }}</p>
            <p class="text-xs text-gray-500">{{ formatPricePerM2(listing.price_per_m2) }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-500">Diện tích</p>
            <p class="text-lg font-bold text-gray-800">{{ listing.acreage.toFixed(0) }} m²</p>
          </div>
          <div v-if="listing.bedrooms">
            <p class="text-xs text-gray-500">Phòng ngủ</p>
            <p class="text-lg font-bold text-gray-800">{{ listing.bedrooms }} PN</p>
          </div>
          <!-- Action icons -->
          <div class="ml-auto flex items-center gap-3 text-gray-400">
            <button class="hover:text-blue-500" @click="handleShare">
              <IconShare />
            </button>
            <button class="hover:text-red-500">
              <IconHeart />
            </button>
          </div>
        </div>
      </div>

      <!-- Thông tin mô tả -->
      <div v-if="listing.description" class="border border-gray-200 bg-white p-5">
        <h2 class="mb-3 border-b border-gray-200 pb-2 text-base font-bold text-gray-800">
          Thông tin mô tả
        </h2>
        <p class="whitespace-pre-line text-sm leading-relaxed text-gray-700">{{ listing.description }}</p>

        <!-- Agent contact -->
        <div v-if="listing.agent_phone" class="mt-4 flex flex-wrap items-center gap-2 text-sm text-gray-600">
          <span>Liên hệ ngay để biết thêm thông tin chi tiết:</span>
          <span class="font-medium text-gray-800">{{ maskedPhone }}</span>
          <n-button size="small" type="primary" @click="showPhone = true">Hiện số</n-button>
          <span v-if="listing.agent_name"> - {{ listing.agent_name }}</span>
        </div>
      </div>

      <!-- Đặc điểm bất động sản -->
      <div class=" border border-gray-200 bg-white p-5">
        <h2 class="mb-4 border-b border-gray-200 pb-2 text-base font-bold text-gray-800">
          Đặc điểm bất động sản
        </h2>

        <div class="grid grid-cols-1 gap-x-8 sm:grid-cols-2">
          <template v-for="attr in attrs" :key="attr.label">
            <div v-if="attr.value" class="flex items-center justify-between border-b border-gray-100 py-2.5">
              <span class="flex items-center gap-2 text-sm text-gray-500">
                <component :is="attr.icon" class="h-4 w-4 flex-shrink-0" />
                {{ attr.label }}
              </span>
              <span class="text-sm font-medium text-gray-800">{{ attr.value }}</span>
            </div>
          </template>
        </div>
      </div>

      <!-- Liên hệ -->
      <div class=" border border-gray-200 bg-white p-5">
        <h2 class="mb-3 border-b border-gray-200 pb-2 text-base font-bold text-gray-800">Liên hệ</h2>
        <div class="flex items-center gap-3">
          <div
            class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-blue-500 text-lg font-bold text-white">
            {{ agentInitial }}
          </div>
          <div>
            <p class="font-semibold text-gray-800">{{ listing.agent_name || 'Người đăng' }}</p>
            <p v-if="listing.agent_phone" class="text-sm text-gray-500">{{ listing.agent_phone }}</p>
          </div>
          <n-button v-if="listing.agent_phone" type="primary" class="ml-auto" @click="handleCall">
            Gọi {{ listing.agent_phone }}
          </n-button>
        </div>
      </div>

      <!-- Footer meta -->
      <div class="grid grid-cols-2 gap-3  border border-gray-200 bg-white p-5 sm:grid-cols-3">
        <div>
          <p class="text-xs text-gray-500">Ngày đăng</p>
          <p class="text-sm font-medium text-gray-800">{{ formatDate(listing.created_at) }}</p>
        </div>
        <div v-if="listing.badge">
          <p class="text-xs text-gray-500">Loại tin</p>
          <p class="text-sm font-medium text-gray-800">{{ listing.badge }}</p>
        </div>
        <div>
          <p class="text-xs text-gray-500">Mã tin</p>
          <p class="text-sm font-medium text-gray-800">{{ listing.id }}</p>
        </div>
        <div v-if="listing.source">
          <p class="text-xs text-gray-500">Nguồn</p>
          <a v-if="listing.source_url" :href="listing.source_url" target="_blank"
            class="text-sm font-medium text-blue-600 hover:underline">
            {{ listing.source }}
          </a>
          <p v-else class="text-sm font-medium text-gray-800">{{ listing.source }}</p>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import IconArea from '~/icons/IconArea.vue';
import IconBath from '~/icons/IconBath.vue';
import IconBed from '~/icons/IconBed.vue';
import IconPrice from '~/icons/IconPrice.vue';
import type { RealEstateResponse } from '~/types/real_estate';
import { formatPrice, formatPricePerM2, formatDate } from '~/utils/format';

const listing = ref<RealEstateResponse | null>(null);
const loading = ref(false);
const activeImageIndex = ref(0);
const showPhone = ref(false);

const allImages = computed(() => listing.value?.images ?? []);
const mainImage = computed(() => allImages.value[activeImageIndex.value] || '');

const agentInitial = computed(() => {
  const name = listing.value?.agent_name || 'Q';
  return name.charAt(0).toUpperCase();
});

const maskedPhone = computed(() => {
  const phone = listing.value?.agent_phone || '';
  if (showPhone.value || !phone) return phone;
  return phone.replace(/.{3}$/, '***');
});

function nextImage() {
  if (allImages.value.length > 1) {
    activeImageIndex.value = (activeImageIndex.value + 1) % allImages.value.length;
  }
}

// Map giá trị lưu DB (value form) → label hiển thị
const DIRECTION_LABELS: Record<string, string> = {
  dong: 'Đông', tay: 'Tây', nam: 'Nam', bac: 'Bắc',
  dong_bac: 'Đông Bắc', tay_bac: 'Tây Bắc', dong_nam: 'Đông Nam', tay_nam: 'Tây Nam',
};
const LEGAL_DOC_LABELS: Record<string, string> = {
  so_do: 'Sổ đỏ/ Sổ hồng', hop_dong_mua_ban: 'Hợp đồng mua bán', dang_cho_so: 'Đang chờ sổ',
};
const INTERIOR_LABELS: Record<string, string> = {
  day_du: 'Đầy đủ nội thất', co_ban: 'Cơ bản', chua_co: 'Chưa có nội thất',
};
const AMENITY_LABELS: Record<string, string> = {
  camera: 'Camera', bao_ve: 'Bảo vệ', pccc: 'PCCC',
};

// Chỉ dùng các field có trong RealEstateResponse
const attrs = computed(() => {
  const l = listing.value;
  if (!l) return [];
  return [
    { label: 'Khoảng giá', icon: IconPrice, value: formatPrice(l.price_vnd) },
    { label: 'Giá/m²', icon: IconPrice, value: formatPricePerM2(l.price_per_m2) },
    { label: 'Diện tích', icon: IconArea, value: `${l.acreage} m²` },
    { label: 'Số phòng ngủ', icon: IconBed, value: l.bedrooms ? `${l.bedrooms} phòng` : '' },
    { label: 'Số phòng tắm, vệ sinh', icon: IconBath, value: l.bathrooms ? `${l.bathrooms} phòng` : '' },
    { label: 'Số tầng', icon: IconArea, value: l.floors ? `${l.floors} tầng` : '' },
    { label: 'Hướng nhà', icon: IconArea, value: l.house_direction ? DIRECTION_LABELS[l.house_direction] ?? l.house_direction : '' },
    { label: 'Hướng ban công', icon: IconArea, value: l.balcony_direction ? DIRECTION_LABELS[l.balcony_direction] ?? l.balcony_direction : '' },
    { label: 'Pháp lý', icon: IconArea, value: l.legal_docs ? LEGAL_DOC_LABELS[l.legal_docs] ?? l.legal_docs : '' },
    { label: 'Nội thất', icon: IconArea, value: l.interior ? INTERIOR_LABELS[l.interior] ?? l.interior : '' },
    { label: 'Giá điện', icon: IconArea, value: l.price_electricity ? `${formatNumber(l.price_electricity)} đ/kWh` : '' },
    { label: 'Giá nước', icon: IconArea, value: l.price_water ? `${formatNumber(l.price_water)} đ/m³` : '' },
    { label: 'Giá internet', icon: IconArea, value: l.price_internet ? `${formatNumber(l.price_internet)} đ/tháng` : '' },
    { label: 'Tiện ích', icon: IconArea, value: l.amenities?.length ? l.amenities.map((a) => AMENITY_LABELS[a] ?? a).join(', ') : '' },
  ];
});

function formatNumber(n: number): string {
  return n.toLocaleString('vi-VN');
}

// Chia sẻ: copy link hiện tại vào clipboard
function handleShare() {
  const url = window.location.href;
  navigator.clipboard?.writeText(url);
  window.message?.success('Đã copy link chia sẻ');
}

function handleCall() {
  if (listing.value?.agent_phone) {
    window.open(`tel:${listing.value.agent_phone}`, '_self');
  }
}

async function fetchDetail(id: number) {
  const { $api } = useNuxtApp();
  loading.value = true;
  try {
    const res = await $api.get<{ data: RealEstateResponse }>(`/real-estate/detail/${id}`);
    listing.value = res?.data ?? null;
  } catch {
    listing.value = null;
  } finally {
    loading.value = false;
  }
}

const props = defineProps<{ id: number }>();

onMounted(() => {
  fetchDetail(props.id);
});
</script>