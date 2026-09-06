<template>
    <div class="overflow-hidden rounded-lg border border-gray-200 bg-gray-100 shadow-sm">
        <!-- Ảnh chính -->
        <div class="relative">
            <img v-if="mainImage" :src="mainImage" :alt="store.listing?.title"
                class="h-105 w-full object-cover" />
            <div v-else class="flex h-105 w-full items-center justify-center bg-gray-200 text-gray-400">
                Không có ảnh
            </div>
            <!-- Nút chuyển ảnh -->
            <template v-if="allImages.length > 1">
                <button
                    class="absolute left-3 top-1/2 -translate-y-1/2 flex h-9 w-9 items-center justify-center bg-white/80 shadow hover:bg-white rounded-full"
                    @click="prevImage">
                    <IconChevronLeft class="h-5 w-5 text-gray-700" />
                </button>
                <button
                    class="absolute right-3 top-1/2 -translate-y-1/2 flex h-9 w-9 items-center justify-center bg-white/80 shadow hover:bg-white rounded-full"
                    @click="nextImage">
                    <IconChevronRight class="h-5 w-5 text-gray-700" />
                </button>
            </template>
            <!-- Đếm ảnh -->
            <span v-if="allImages.length > 0"
                class="absolute bottom-3 right-3 bg-black/50 px-2 py-0.5 text-xs text-white rounded">
                {{ activeImageIndex + 1 }} / {{ allImages.length }}
            </span>
        </div>

        <!-- Thumbnail -->
        <div v-if="allImages.length > 1" class="flex gap-2 bg-gray-800 p-2">
            <button v-for="(img, i) in allImages" :key="i"
                class="h-16 w-24 shrink-0 overflow-hidden border-2 transition rounded"
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
