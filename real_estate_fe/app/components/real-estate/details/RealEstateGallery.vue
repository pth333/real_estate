<template>
    <div class="overflow-hidden rounded-lg border border-gray-200 bg-gray-100 shadow-sm">
        <!-- Ảnh chính: thấp ở mobile, cao dần theo breakpoint -->
        <div class="relative">
            <img v-if="mainImage" :src="mainImage" :alt="store.listing?.title"
                class="h-56 w-full object-cover md:h-80 lg:h-105" />
            <div v-else class="flex h-56 w-full items-center justify-center bg-gray-200 text-gray-400 md:h-80 lg:h-105">
                Không có ảnh
            </div>
            <!-- Nút chuyển ảnh -->
            <template v-if="allImages.length > 1">
                <button
                    class="absolute left-2 top-1/2 flex h-8 w-8 -translate-y-1/2 items-center justify-center rounded-full bg-white/80 shadow hover:bg-white md:left-3 md:h-9 md:w-9"
                    @click="prevImage">
                    <IconChevronLeft class="h-5 w-5 text-gray-700" />
                </button>
                <button
                    class="absolute right-2 top-1/2 flex h-8 w-8 -translate-y-1/2 items-center justify-center rounded-full bg-white/80 shadow hover:bg-white md:right-3 md:h-9 md:w-9"
                    @click="nextImage">
                    <IconChevronRight class="h-5 w-5 text-gray-700" />
                </button>
            </template>
            <!-- Đếm ảnh -->
            <span v-if="allImages.length > 0"
                class="absolute bottom-2 right-2 rounded bg-black/50 px-2 py-0.5 text-[11px] text-white md:bottom-3 md:right-3 md:text-xs">
                {{ activeImageIndex + 1 }} / {{ allImages.length }}
            </span>
        </div>

        <!-- Thumbnail: cuộn ngang khi màn hình hẹp -->
        <div v-if="allImages.length > 1"
            class="flex gap-2 overflow-x-auto bg-gray-800 p-2 [scrollbar-width:thin]">
            <button v-for="(img, i) in allImages" :key="i"
                class="h-12 w-16 shrink-0 overflow-hidden rounded border-2 transition md:h-16 md:w-24"
                :class="activeImageIndex === i ? 'border-red-500' : 'border-transparent opacity-70 hover:opacity-100'"
                @click="activeImageIndex = i">
                <img :src="img" :alt="`Ảnh ${i + 1}`" class="h-full w-full object-cover" />
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useRealEstateDetail } from '~/stores/detail/real_estate_detail';

const store = useRealEstateDetail();

const activeImageIndex = ref(0);

const allImages = computed(() => (store.listing?.images ?? []).map(img => img.url));
const mainImage = computed(() => allImages.value[activeImageIndex.value] || '');

function nextImage() {
    if (allImages.value.length > 1) {
        activeImageIndex.value = (activeImageIndex.value + 1) % allImages.value.length;
    }
}

function prevImage() {
    if (allImages.value.length > 1) {
        activeImageIndex.value = (activeImageIndex.value - 1 + allImages.value.length) % allImages.value.length;
    }
}
</script>
