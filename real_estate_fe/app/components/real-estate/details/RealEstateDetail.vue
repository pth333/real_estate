<template>
    <!-- pb lớn ở mobile/tablet để thanh CTA cố định dưới màn hình không che nội dung -->
    <div class="mx-auto max-w-285 px-3 py-4 pb-24 md:px-4 md:py-5 lg:pb-5">
        <!-- Loading -->
        <SkeletonCard v-if="realEstateDetailStore.loading" />

        <!-- Detail: desktop chia 2 cột (nội dung + liên hệ), mobile/tablet xếp dọc -->
        <div v-else class="grid grid-cols-1 gap-4 md:gap-5 lg:grid-cols-4">

            <!-- Cột trái -->
            <div class="space-y-2 lg:col-span-3">
                <RealEstateGallery />

                <!-- Breadcrumb -->
                <nav class="text-xs text-emerald-600 px-1">
                    <span>Bán</span>
                    <span class="mx-1 text-gray-400">/</span>
                    <span>{{ realEstateDetailStore.listing?.city }}</span>
                    <span v-if="realEstateDetailStore.listing?.district">
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
