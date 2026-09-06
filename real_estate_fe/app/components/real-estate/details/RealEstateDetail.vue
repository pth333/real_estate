<template>
    <div class="mx-auto max-w-285 px-3 py-4">
        <!-- Loading -->
        <SkeletonCard v-if="realEstateDetailStore.loading" />

        <!-- 404 -->
        <div v-else-if="!realEstateDetailStore.listing" class="px-6 py-16 text-center">
            <p class="text-base text-gray-400">Không tìm thấy tin đăng</p>
            <n-button class="mt-4" type="primary" @click="navigateTo('/')">Quay lại</n-button>
        </div>

        <!-- Detail -->
        <div v-else class="grid grid-cols-1 md:grid-cols-4 gap-5">

            <!-- Cột trái -->
            <div class="md:col-span-3 space-y-2">
                <RealEstateGallery />

                <!-- Breadcrumb -->
                <nav class="text-xs text-emerald-600 px-1">
                    <span>Bán</span>
                    <span class="mx-1 text-gray-400">/</span>
                    <span>{{ realEstateDetailStore.listing.city }}</span>
                    <span v-if="realEstateDetailStore.listing.district">
                        <span class="mx-1 text-gray-400">/</span>
                        <span>{{ realEstateDetailStore.listing.district }}</span>
                    </span>
                </nav>

                <RealEstateTitlePrice />
                <RealEstateDescription />
                <RealEstateFeatures />
                <!-- <InformationProject /> -->
                <RealEstateDetailForYou />

                <!-- Tiện ích quanh bds -->
                <NearbyAmenities />
            </div>

            <ContactRealestateDetail />

        </div>
    </div>
</template>

<script setup lang="ts">
import { useRealEstateDetail } from '~/stores/detail/real_estate_detail';
import { useTracking } from '~/composables/useTracking';
const props = defineProps<{ id: number }>();

const realEstateDetailStore = useRealEstateDetail()
const { trackView, cleanupTracking } = useTracking()

onMounted(() => {
    realEstateDetailStore.fetchDetail(props.id);
    trackView(props.id);
});

const detailTitle = computed(() => realEstateDetailStore.listing?.title || "Chi tiết bất động sản")
useHead({
    title: detailTitle,
})

onBeforeUnmount(() => {
    cleanupTracking();
});
</script>
