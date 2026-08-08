<template>
  <div class="mx-auto max-w-[1000px] px-6 py-6">
    <!-- Loading -->
    <SkeletonCard v-if="loading" />

    <!-- 404 -->
    <div v-else-if="!listing" class="px-6 py-16 text-center">
      <p class="text-base text-gray-400">Không tìm thấy tin đăng</p>
      <n-button class="mt-4" type="primary" @click="navigateTo('/')">Quay lại</n-button>
    </div>

    <!-- Detail -->
    <div v-else class="space-y-6">
      <!-- Hình chính -->
      <div class="overflow-hidden rounded-lg bg-gray-100">
        <img v-if="mainImage" :src="mainImage" :alt="listing.title"
          class="h-[420px] w-full object-cover" />
        <div v-for="(img, i) in otherImages" :key="i" class="sr-only">
          <img :src="img" :alt="`${listing.title} ${i + 2}`" />
        </div>
      </div>

      <!-- Thông tin chính -->
      <div class="rounded-lg border border-gray-200 bg-white p-6">
        <h1 class="text-2xl font-bold text-gray-800">{{ listing.title }}</h1>

        <div class="mt-3 flex flex-wrap items-center gap-3 text-sm text-gray-500">
          <span class="text-2xl font-bold text-red-600">{{ formatPrice(listing.price_vnd) }}</span>
          <span>{{ formatPricePerM2(listing.price_per_m2) }}</span>
          <span>{{ listing.acreage.toFixed(1) }} m²</span>
          <span v-if="listing.bedrooms" class="flex items-center gap-1"><IconBed /> {{ listing.bedrooms }}</span>
          <span v-if="listing.bathrooms" class="flex items-center gap-1"><IconBath /> {{ listing.bathrooms }}</span>
        </div>

        <p class="mt-3 flex items-center gap-1 text-sm text-gray-500">
          <IconLocationOutline class="h-4 w-4" />
          {{ [listing.district, listing.city].filter(Boolean).join(', ') }}
        </p>
      </div>

      <!-- Mô tả -->
      <div v-if="listing.description" class="rounded-lg border border-gray-200 bg-white p-4">
        <h2 class="mb-2 text-sm font-semibold uppercase text-gray-700">Mô tả</h2>
        <p class="whitespace-pre-line text-sm leading-relaxed text-gray-600">{{ listing.description }}</p>
      </div>

      <!-- Liên hệ -->
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <h2 class="mb-3 text-sm font-semibold uppercase text-gray-700">Liên hệ</h2>
        <div class="flex items-center gap-3">
          <div class="flex h-12 w-12 items-center justify-center rounded-full bg-blue-500 text-lg font-bold text-white">
            {{ agentInitial }}
          </div>
          <div>
            <p class="font-semibold text-gray-800">{{ listing.agent_name || 'Người đăng' }}</p>
            <p class="text-sm text-gray-500">{{ listing.agent_phone || '' }}</p>
          </div>
          <n-button v-if="listing.agent_phone" type="primary" class="ml-auto" @click="handleCall">
            Gọi {{ listing.agent_phone }}
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { RealEstateResponse } from '~/types/real_estate';

const listing = ref<RealEstateResponse | null>(null);
const loading = ref(false);

const mainImage = computed(() => listing.value?.images?.[0] || '');
const otherImages = computed(() => listing.value?.images?.slice(1) || []);
const agentInitial = computed(() => {
  const name = listing.value?.agent_name || 'Q';
  return name.charAt(0).toUpperCase();
});

function handleCall() {
  if (listing.value?.agent_phone) {
    window.open(`tel:${listing.value.agent_phone}`, '_self');
  }
}

async function fetchDetail(slug: string) {
  const { $api } = useNuxtApp();
  loading.value = true;
  try {
    const res = await $api.get<{ type: string; data: RealEstateResponse }>(`/real-estate/${slug}`);
    listing.value = res?.data ?? null;
  } catch {
    listing.value = null;
  } finally {
    loading.value = false;
  }
}

const route = useRoute();
watch(
  () => route.params.category,
  (val) => {
    const slug = Array.isArray(val) ? val[0] ?? '' : val ?? '';
    if (slug) fetchDetail(slug);
  },
  { immediate: true },
);
</script>