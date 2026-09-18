<template>
  <section class="py-8">
    <div class="container mx-auto px-24">
      <!-- Header -->
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-xl font-bold text-gray-900">Bất động sản dành cho bạn</h2>
        <div class="flex items-center gap-3 text-sm">
          <NuxtLink :to="realEstateNewest[0] ? `/${realEstateNewest[0].Slug}` : '#'"
            class="text-gray-600 hover:text-red-500 transition-colors">
            Tin nhà đất bán mới nhất
          </NuxtLink>
          <span class="text-gray-300">|</span>
          <NuxtLink :to="realEstateNewest[1] ? `/${realEstateNewest[1].Slug}` : '#'"
            class="text-gray-600 hover:text-red-500 transition-colors">
            Tin nhà đất cho thuê mới nhất
          </NuxtLink>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <SkeletonCard :count="8" type="card" />
      </div>

      <!-- Grid -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <NuxtLink v-for="item in visibleItems" :key="item.id" :to="`/${item.slug}`">
          <n-card hoverable size="small" content-style="padding: 12px;"
            class="cursor-pointer overflow-hidden rounded-lg shadow-sm">
            <template #cover>
              <div class="relative h-48 overflow-hidden bg-gray-100 group">
                <div v-if="item.image_urls?.length > 1"
                  class="absolute bottom-2 left-2 z-10 flex items-center gap-1 bg-black/50 text-white text-xs px-1.5 py-0.5 rounded">
                  <IconImage class="h-3 w-3" />
                  <span>{{ item.image_urls.length }}</span>
                </div>
                <img :src="item.thumbnail" :alt="item.title" loading="lazy"
                  class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105" />
              </div>
            </template>

            <!-- Nội dung -->
            <div class="flex flex-col gap-1.5">
              <!-- Badge xác thực + tiêu đề -->
              <div class="text-sm font-medium text-gray-800 leading-snug line-clamp-2 min-h-[40px]">
                <span v-if="item.verified" class="inline-flex items-center gap-1 text-green-600 font-semibold mr-1">
                  <IconShieldCheck class="h-3.5 w-3.5" />
                  XÁC THỰC
                </span>
                {{ item.title }}
              </div>

              <!-- Giá + diện tích -->
              <div class="flex items-center gap-2 text-sm">
                <span class="text-red-500 font-semibold">{{ item.price }}</span>
                <span class="text-gray-300">·</span>
                <span class="text-gray-600">{{ item.area }}</span>
              </div>

              <div v-if="item.bedrooms || item.bathrooms" class="flex items-center gap-2 text-xs text-gray-500">
                <template v-if="item.bedrooms">
                  <IconBed class="h-3 w-3" />
                  <span>{{ item.bedrooms }} PN</span>
                </template>
                <span v-if="item.bedrooms && item.bathrooms" class="text-gray-300">·</span>
                <template v-if="item.bathrooms">
                  <IconBath class="h-3 w-3" />
                  <span>{{ item.bathrooms }} WC</span>
                </template>
              </div>

              <!-- Địa chỉ -->
              <div class="flex items-center gap-1 text-xs text-gray-500">
                <IconMapPin class="h-3 w-3 shrink-0" />
                <span class="truncate">{{ item.location }}</span>
              </div>

              <!-- Footer: ngày đăng + yêu thích -->
              <div class="flex items-center justify-end mt-1" @click.stop.prevent>
                <button
                  class="flex h-8 w-8 items-center justify-center rounded-full border border-gray-200 transition-colors hover:border-red-400 hover:text-red-500"
                  :class="item.is_favorite ? 'border-red-500 text-red-500' : ''" @click="toggleFavorite(item.id)">
                  <IconHeart class="h-4 w-4"
                    :class="item.is_favorite ? 'fill-red-500 text-red-500' : 'text-gray-400'" />
                </button>
              </div>
            </div>
          </n-card>
        </NuxtLink>
      </div>

      <!-- Mở rộng -->
      <div class="flex justify-center mt-6">
        <n-button secondary class="px-6 h-10 text-gray-700 font-medium" @click="expanded = !expanded">
          {{ expanded ? 'Thu gọn' : 'Mở rộng' }}
          <template #icon>
            <IconChevronDown class="transition-transform duration-300" :class="{ 'rotate-180': expanded }" />
          </template>
        </n-button>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { RealEstateResponse } from '~/types/real_estate';

const expanded = ref(false);
const loading = ref(true);
const items = ref<any[]>([]);
const { $api } = useNuxtApp();
const favorite = useFavorite();

const menuStore = useMenuStore();

const realEstateNewest = computed(() => {
  return menuStore.menu?.categories ?? [];
});

// Fetch dữ liệu gợi ý từ API
const fetchRecommendations = async () => {
  loading.value = true;
  try {
    const res = await $api.get<{ data: RealEstateResponse[] }>('/real-estate/recommend', {
      params: { limit: 12 }
    });
    items.value = res.data || [];
  } catch (err) {
    console.error("Lỗi khi tải gợi ý BĐS:", err);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchRecommendations();
});

const visibleItems = computed(() => {
  return items.value.map((item) => ({
    ...item,
    thumbnail: item.image_urls?.[0] || '',
    price: formatPrice(item.price),
    area: formatAcreage(item.acreage),
    location: [item.district, item.city].join(', '),
  })).slice(0, expanded.value ? undefined : 8)
})

function toggleFavorite(id: number) {
  const item = items.value.find(i => i.id === id);
  if (item) {
    favorite.toggle(id).then((next) => {
      if (next !== null) item.is_favorite = next;
    });
  }
}
</script>
