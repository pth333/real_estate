<template>
  <section class="py-8">
    <div class="container mx-auto px-24">
      <!-- Header -->
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-xl font-bold text-gray-900">Bất động sản dành cho bạn</h2>
        <div class="flex items-center gap-3 text-sm">
          <a href="#" class="text-gray-600 hover:text-red-500 transition-colors">Tin nhà đất bán mới nhất</a>
          <span class="text-gray-300">|</span>
          <a href="#" class="text-gray-600 hover:text-red-500 transition-colors">Tin nhà đất cho thuê mới nhất</a>
        </div>
      </div>

      <!-- Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <n-card
          v-for="item in visibleItems"
          :key="item.id"
          hoverable
          size="small"
          content-style="padding: 12px;"
          class="cursor-pointer overflow-hidden rounded-lg shadow-sm"
        >
          <template #cover>
            <div class="relative aspect-4/3 overflow-hidden bg-gray-100">
              <img
                :src="item.thumbnail"
                :alt="item.title"
                class="w-full h-full object-cover hover:scale-105 transition-transform duration-300"
              />
            </div>
          </template>

          <!-- Nội dung -->
          <div class="flex flex-col gap-1.5">
            <!-- Badge xác thực + tiêu đề -->
            <div class="text-sm font-medium text-gray-800 leading-snug line-clamp-2">
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

            <!-- Địa chỉ -->
            <div class="flex items-center gap-1 text-xs text-gray-500">
              <IconMapPin class="h-3 w-3 shrink-0" />
              <span class="truncate">{{ item.location }}</span>
            </div>

            <!-- Footer: ngày đăng + yêu thích -->
            <div class="flex items-center justify-between mt-1">
              <span class="text-xs text-gray-400">{{ item.postedAt }}</span>
              <button
                class="flex h-8 w-8 items-center justify-center rounded-full border border-gray-200 transition-colors hover:border-red-400 hover:text-red-500"
                :class="item.isFavorite ? 'border-red-500 text-red-500' : ''"
                @click.stop="toggleFavorite(item.id)"
              >
                <IconHeart
                  class="h-4 w-4"
                  :class="item.isFavorite ? 'fill-red-500 text-red-500' : 'text-gray-400'"
                />
              </button>
            </div>
          </div>
        </n-card>
      </div>

      <!-- Mở rộng -->
      <div class="flex justify-center mt-6">
        <n-button
          secondary
          class="px-6 h-10 text-gray-700 font-medium"
          @click="expanded = !expanded"
        >
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
const loading = ref(false);
const items = ref<any[]>([]);
const { $api } = useNuxtApp();
const favorite = useFavorite();

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

// Hiển thị 8 khi thu gọn, tất cả khi mở rộng
const visibleItems = computed(() => expanded.value ? items.value : items.value.slice(0, 8));

function toggleFavorite(id: number) {
  const item = items.value.find(i => i.id === id);
  if (item) {
    favorite.toggleWithConfirm(id, item.isFavorite).then((next) => {
      if (next !== null) item.isFavorite = next;
    });
  }
}
</script>
