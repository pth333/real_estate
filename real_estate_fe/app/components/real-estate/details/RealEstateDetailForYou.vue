<template>
    <div class="container mx-auto py-4">
        <h2 class="mb-4 border-b border-gray-200 pb-2 text-base font-bold text-gray-800">
            Bất dộng sản dành cho bạn
        </h2>
        <!-- Loading State -->
        <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            <SkeletonCard :count="6" type="card" />
        </div>

        <!-- Grid 3 cột -->
        <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            <div v-for="item in visibleItems" :key="item.id"
                class="cursor-pointer rounded-xl overflow-hidden border border-gray-100 shadow-sm hover:shadow-md transition-shadow duration-200 bg-white">

                <!-- Ảnh + Badge -->
                <div class="relative aspect-4/3 overflow-hidden bg-gray-100">
                    <img :src="item.image_urls?.[0] ?? '/placeholder.jpg'" :alt="item.title"
                        class="w-full h-full object-cover hover:scale-105 transition-transform duration-300" />
                    <!-- Badge VIP / nổi bật -->
                    <span v-if="item.badge"
                        class="absolute top-2 left-2 px-2 py-0.5 text-xs font-bold rounded text-white"
                        :class="item.badge === 'VIP' ? 'bg-yellow-400 text-yellow-900' : 'bg-blue-500'">
                        {{ item.badge }}
                    </span>
                </div>

                <!-- Nội dung -->
                <div class="p-3 flex flex-col gap-2">

                    <!-- Tiêu đề -->
                    <div class="text-sm font-semibold text-gray-800 leading-snug line-clamp-2">
                        {{ item.title }}
                    </div>

                    <!-- Giá + diện tích -->
                    <div class="flex items-center gap-2 text-sm">
                        <span class="text-red-500 font-bold">{{ formatPrice(item.price_vnd) }}</span>
                        <span class="text-gray-300">·</span>
                        <span class="text-gray-600">{{ item.acreage }} m²</span>
                    </div>

                    <!-- Phòng ngủ / WC / Nội thất -->
                    <div class="flex items-center gap-3 text-xs text-gray-500">
                        <span v-if="item.bedrooms" class="flex items-center gap-1">
                            <IconBed class="h-3.5 w-3.5" />
                            {{ item.bedrooms }}
                        </span>
                        <span v-if="item.bathrooms" class="flex items-center gap-1">
                            <IconBath class="h-3.5 w-3.5" />
                            {{ item.bathrooms }}
                        </span>
                        <span v-if="item.interior" class="flex items-center gap-1">
                            <IconArmchair class="h-3.5 w-3.5" />
                            {{ item.interior }}
                        </span>
                        <span v-if="item.house_direction" class="flex items-center gap-1">
                            <IconCompass class="h-3.5 w-3.5" />
                            {{ item.house_direction }}
                        </span>
                    </div>

                    <!-- Địa chỉ -->
                    <div class="flex items-start gap-1 text-xs text-gray-500">
                        <IconMapPin class="h-3 w-3 shrink-0 mt-0.5" />
                        <span class="line-clamp-1">{{ [item.address, item.district,
                        item.city].filter(Boolean).join(', ') }}</span>
                    </div>

                    <!-- Footer: ngày đăng + yêu thích -->
                    <div class="flex items-center justify-between pt-1 border-t border-gray-100">
                        <span class="text-xs text-gray-400">{{ formatDate(item.created_at) }}</span>
                        <button class="flex h-8 w-8 items-center justify-center rounded-full border transition-colors"
                            :class="item.is_favorite
                                ? 'border-red-500 text-red-500'
                                : 'border-gray-200 text-gray-400 hover:border-red-400 hover:text-red-500'"
                            @click.stop="toggleFavorite(item.id)">
                            <IconHeart class="h-4 w-4" :class="item.is_favorite ? 'fill-red-500 text-red-500' : ''" />
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Mở rộng / Thu gọn -->
        <div class="flex justify-center mt-6">
            <n-button secondary class="px-6 h-10 text-gray-700 font-medium" @click="expanded = !expanded">
                {{ expanded ? 'Thu gọn' : 'Xem thêm' }}
                <template #icon>
                    <IconChevronDown class="transition-transform duration-300" :class="{ 'rotate-180': expanded }" />
                </template>
            </n-button>
        </div>

    </div>
</template>

<script setup lang="ts">
import type { RealEstateResponse } from '~/types/real_estate';

const expanded = ref(false);
const loading = ref(false);
const items = ref<RealEstateResponse[]>([]);
const { $api } = useNuxtApp();
const favorite = useFavorite();

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

const visibleItems = computed(() =>
    expanded.value ? items.value : items.value.slice(0, 6)
);

function toggleFavorite(id: number) {
    const item = items.value.find(i => i.id === id);
    if (item) {
        favorite.toggle(id).then((next) => {
            if (next !== null) item.is_favorite = next;
        });
    }
}
</script>