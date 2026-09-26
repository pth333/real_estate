<template>
    <div class="rounded-lg bg-white px-3 py-3 md:px-4 md:py-4">
        <h1 class="text-lg font-bold leading-snug text-gray-800 md:text-xl lg:text-2xl">{{ store.listing?.title }}</h1>

        <!-- Địa chỉ -->
        <p class="mt-2 flex items-start gap-1 text-sm">
            <IconLocationOutline class="mt-0.5 h-4 w-4 shrink-0 text-gray-500" />
            <span class="text-gray-600">{{ locationFull }}</span>
        </p>

        <!-- Giá / Diện tích / Phòng ngủ -->
        <div class="mt-3 flex flex-wrap items-end gap-x-5 gap-y-3 md:mt-4 md:gap-x-6">
            <div>
                <p class="text-xs text-gray-500">Khoảng giá</p>
                <p class="text-xl font-bold text-gray-800 md:text-2xl">{{ formatPrice(store.listing?.price_vnd) }}</p>
                <p class="text-xs text-gray-500">{{ formatPricePerM2(store.listing?.price_per_m2) }}</p>
            </div>
            <div>
                <p class="text-xs text-gray-500">Diện tích</p>
                <p class="text-base font-bold text-gray-800 md:text-lg">{{ store.listing?.acreage.toFixed(0) }} m²</p>
            </div>
            <div v-if="store.listing?.bedrooms">
                <p class="text-xs text-gray-500">Phòng ngủ</p>
                <p class="text-base font-bold text-gray-800 md:text-lg">{{ store.listing?.bedrooms }} PN</p>
            </div>
            <!-- Action icons -->
            <div class="ml-auto flex items-center gap-3 text-gray-400">
                <button
                    class="flex h-9 w-9 items-center justify-center rounded-full border border-gray-200 transition-colors hover:border-emerald-400 hover:text-emerald-500"
                    @click="store.handleShare()">
                    <IconShare class="h-5 w-5" />
                </button>
                <button
                    class="flex h-9 w-9 items-center justify-center rounded-full border border-gray-200 transition-colors hover:border-red-400 hover:text-red-500"
                    :class="store.listing?.is_favorite ? 'border-red-500 text-red-500' : ''"
                    @click="toggleFavoriteDetail">
                    <IconHeart class="h-5 w-5"
                        :class="store.listing?.is_favorite ? 'fill-red-500 text-red-500' : 'text-gray-400'" />
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { formatPrice, formatPricePerM2 } from '~/utils/format';
import { useRealEstateDetail } from '~/stores/detail/real_estate_detail';

const store = useRealEstateDetail();
const favorite = useFavorite();

const locationFull = computed(() => {
    const l = store.listing;
    if (!l) return '';
    const parts = [l.address, l.district, l.city].filter(p => p && p.trim() !== '');
    return parts.join(', ');
});

const toggleFavoriteDetail = async () => {
    const id = store.listing?.id;
    if (!id) return;
    const next = await favorite.toggle(id);
    if (next !== null && store.listing) {
        store.listing.is_favorite = next;
    }
};
</script>
