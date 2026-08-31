<template>
    <div class="mx-auto max-w-[1140px] px-3 py-4">
        <!-- Loading -->
        <SkeletonCard v-if="realEstateDetailStore.loading" />

        <!-- 404 -->
        <div v-else-if="!realEstateDetailStore.listing" class="px-6 py-16 text-center">
            <p class="text-base text-gray-400">Không tìm thấy tin đăng</p>
            <n-button class="mt-4" type="primary" @click="navigateTo('/')">Quay lại</n-button>
        </div>

        <!-- Detail -->
        <div v-else class="grid grid-cols-1 md:grid-cols-4 gap-5">

            <!-- Cột trái: Hình ảnh và thông tin chi tiết (md:col-span-3 - chiếm 75% độ rộng) -->
            <div class="md:col-span-3 space-y-2">
                <!-- Gallery -->
                <div class="overflow-hidden rounded-lg border border-gray-200 bg-gray-100 shadow-sm">
                    <!-- Main image -->
                    <div class="relative">
                        <img v-if="mainImage" :src="mainImage" :alt="realEstateDetailStore.listing.title"
                            class="h-[420px] w-full object-cover" />
                        <div v-else class="flex h-[420px] w-full items-center justify-center bg-gray-200 text-gray-400">
                            Không có ảnh
                        </div>
                        <!-- Nav arrows -->
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
                        <!-- Counter -->
                        <span v-if="allImages.length > 0"
                            class="absolute bottom-3 right-3 bg-black/50 px-2 py-0.5 text-xs text-white rounded">
                            {{ activeImageIndex + 1 }} / {{ allImages.length }}
                        </span>
                    </div>

                    <!-- Thumbnails -->
                    <div v-if="allImages.length > 1" class="flex gap-2 bg-gray-800 p-2">
                        <button v-for="(img, i) in allImages" :key="i"
                            class="h-16 w-24 shrink-0 overflow-hidden border-2 transition rounded"
                            :class="activeImageIndex === i ? 'border-red-500' : 'border-transparent opacity-70 hover:opacity-100'"
                            @click="activeImageIndex = i">
                            <img :src="img" :alt="`Ảnh ${i + 1}`" class="h-full w-full object-cover" />
                        </button>
                    </div>
                </div>

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

                <!-- Tiêu đề + Giá -->
                <div class=" bg-white rounded-lg p-4">
                    <h1 class="text-xl font-bold leading-snug text-gray-800">{{ realEstateDetailStore.listing.title }}
                    </h1>

                    <!-- Địa chỉ -->
                    <p class="mt-2 flex items-start gap-1 text-sm">
                        <IconLocationOutline class="mt-0.5 h-4 w-4 shrink-0 text-gray-500" />
                        <span class="text-gray-600">{{ locationFull }}</span>
                    </p>

                    <!-- Giá / Diện tích / Phòng ngủ -->
                    <div class="mt-4 flex flex-wrap items-end gap-6">
                        <div>
                            <p class="text-xs text-gray-500">Khoảng giá</p>
                            <p class="text-2xl font-bold text-gray-800">{{
                                formatPrice(realEstateDetailStore.listing.price_vnd) }}</p>
                            <p class="text-xs text-gray-500">{{
                                formatPricePerM2(realEstateDetailStore.listing.price_per_m2) }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-gray-500">Diện tích</p>
                            <p class="text-lg font-bold text-gray-800">{{
                                realEstateDetailStore.listing.acreage.toFixed(0) }} m²</p>
                        </div>
                        <div v-if="realEstateDetailStore.listing.bedrooms">
                            <p class="text-xs text-gray-500">Phòng ngủ</p>
                            <p class="text-lg font-bold text-gray-800">{{ realEstateDetailStore.listing.bedrooms }} PN
                            </p>
                        </div>
                        <!-- Action icons -->
                        <div class="ml-auto flex items-center gap-3 text-gray-400">
                            <button
                                class="flex h-9 w-9 items-center justify-center rounded-full border border-gray-200 transition-colors hover:border-emerald-400 hover:text-emerald-500"
                                @click="realEstateDetailStore.handleShare">
                                <IconShare class="h-5 w-5" />
                            </button>
                            <button
                                class="flex h-9 w-9 items-center justify-center rounded-full border border-gray-200 transition-colors hover:border-red-400 hover:text-red-500"
                                :class="realEstateDetailStore.listing?.is_favorite ? 'border-red-500 text-red-500' : ''"
                                @click="toggleFavoriteDetail">
                                <IconHeart class="h-5 w-5"
                                    :class="realEstateDetailStore.listing?.is_favorite ? 'fill-red-500 text-red-500' : 'text-gray-400'" />
                            </button>
                        </div>
                    </div>
                </div>

                <!-- Thông tin mô tả -->
                <div v-if="realEstateDetailStore.listing.description" class=" bg-white rounded-lg p-4 ">
                    <h2 class="mb-3 border-b border-gray-200 pb-2 text-base font-bold text-gray-800">
                        Thông tin mô tả
                    </h2>
                    <p class="whitespace-pre-line text-sm leading-relaxed text-gray-700">{{
                        realEstateDetailStore.listing.description }}
                    </p>

                    <!-- Agent contact fallback for mobile -->
                    <div v-if="realEstateDetailStore.listing.agent_phone"
                        class="mt-4 flex flex-wrap items-center gap-2 text-sm text-gray-600 sm:hidden">
                        <span>Liên hệ:</span>
                        <span class="font-medium text-gray-800">{{ realEstateDetailStore.maskedPhone }}</span>
                        <n-button size="small" type="primary" @click="realEstateDetailStore.showPhone = true">Hiện
                            số</n-button>
                    </div>
                </div>

                <!-- Đặc điểm bất động sản  -->
                <div class=" bg-white rounded-lg p-4 ">
                    <h2 class="mb-4 border-b border-gray-200 pb-2 text-base font-bold text-gray-800">
                        Đặc điểm bất động sản
                    </h2>

                    <div class="grid grid-cols-1 gap-x-8 sm:grid-cols-2">
                        <template v-for="attr in attrs" :key="attr.label">
                            <div v-if="attr.value"
                                class="flex items-center justify-between border-b border-gray-100 py-2">
                                <span class="flex items-center gap-2 text-sm text-gray-500">
                                    <component :is="attr.icon" class="h-4 w-4 shrink-0" />
                                    {{ attr.label }}
                                </span>
                                <span class="text-sm font-medium text-gray-800">{{ attr.value }}</span>
                            </div>
                        </template>
                    </div>
                </div>

                <!--Tiện ích quanh bds  -->
                <NearbyAmenities />

                <!-- Footer meta -->
                <!-- <div
                    class="grid grid-cols-2 gap-3 bg-white rounded-lg p-4  sm:grid-cols-3">
                    <div>
                        <p class="text-xs text-gray-500">Ngày đăng</p>
                        <p class="text-sm font-medium text-gray-800">{{
                            formatDate(realEstateDetailStore.listing.created_at) }}</p>
                    </div>
                    <div v-if="realEstateDetailStore.listing.badge">
                        <p class="text-xs text-gray-500">Loại tin</p>
                        <p class="text-sm font-medium text-gray-800">{{ realEstateDetailStore.listing.badge }}</p>
                    </div>
                    <div>
                        <p class="text-xs text-gray-500">Mã tin</p>
                        <p class="text-sm font-medium text-gray-800">{{ realEstateDetailStore.listing.id }}</p>
                    </div>
                    <div v-if="realEstateDetailStore.listing.source">
                        <p class="text-xs text-gray-500">Nguồn</p>
                        <a v-if="realEstateDetailStore.listing.source_url"
                            :href="realEstateDetailStore.listing.source_url" target="_blank"
                            class="text-sm font-medium text-emerald-600 hover:underline">
                            {{ realEstateDetailStore.listing.source }}
                        </a>
                        <p v-else class="text-sm font-medium text-gray-800">{{ realEstateDetailStore.listing.source }}
                        </p>
                    </div>
                </div> -->
            </div>

            <ContactRealestateDetail />

        </div>
    </div>
</template>

<script setup lang="ts">
import { formatPrice, formatPricePerM2 } from '~/utils/format';
import { useRealEstateDetail } from '~/stores/detail/real_estate_detail';
import { useTracking } from '~/composables/useTracking';
const props = defineProps<{ id: number }>();

const realEstateDetailStore = useRealEstateDetail()
const { trackView, cleanupTracking } = useTracking()
const favorite = useFavorite()

const activeImageIndex = ref(0);

const locationFull = computed(() => {
    const l = realEstateDetailStore.listing;
    if (!l) return "";
    const parts = [l.address, l.district, l.city].filter(
        (p) => p && p.trim() !== "",
    );
    return parts.join(", ");
});

const toggleFavoriteDetail = async () => {
    const id = realEstateDetailStore.listing?.id
    if (!id) return
    const next = await favorite.toggle(id)
    if (next !== null && realEstateDetailStore.listing) {
        realEstateDetailStore.listing.is_favorite = next
    }
};

const allImages = computed(() => (realEstateDetailStore.listing?.images ?? []).map(img => img.url));
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

const attrs = computed(() => {
    const l = realEstateDetailStore.listing;
    if (!l) return [];
    return [
        { label: 'Khoảng giá', icon: 'IconPrice', value: formatPrice(l.price_vnd) },
        { label: 'Giá/m²', icon: 'IconPrice', value: formatPricePerM2(l.price_per_m2) },
        { label: 'Diện tích', icon: 'IconArea', value: `${l.acreage} m²` },
        { label: 'Số phòng ngủ', icon: 'IconBed', value: l.bedrooms ? `${l.bedrooms} phòng` : '' },
        { label: 'Số phòng tắm, vệ sinh', icon: 'IconBath', value: l.bathrooms ? `${l.bathrooms} phòng` : '' },
        { label: 'Số tầng', icon: 'IconBuilding', value: l.floors ? `${l.floors} tầng` : '' },
        { label: 'Hướng nhà', icon: 'IconCompass', value: l.house_direction ? DIRECTION_LABELS[l.house_direction] ?? l.house_direction : '' },
        { label: 'Hướng ban công', icon: 'IconBalcony', value: l.balcony_direction ? DIRECTION_LABELS[l.balcony_direction] ?? l.balcony_direction : '' },
        { label: 'Pháp lý', icon: 'IconShieldCheck', value: l.legal_docs ? LEGAL_DOC_LABELS[l.legal_docs] ?? l.legal_docs : '' },
        { label: 'Nội thất', icon: 'IconSofa', value: l.interior ? INTERIOR_LABELS[l.interior] ?? l.interior : '' },
        { label: 'Giá điện', icon: 'IconZap', value: l.price_electricity ? `${formatPriceNumber(l.price_electricity)} đ/kWh` : '' },
        { label: 'Giá nước', icon: 'IconDroplet', value: l.price_water ? `${formatPriceNumber(l.price_water)} đ/m³` : '' },
        { label: 'Giá internet', icon: 'IconWifi', value: l.price_internet ? `${formatPriceNumber(l.price_internet)} đ/tháng` : '' },
        { label: 'Tiện ích', icon: 'IconSparkles', value: l.amenities?.length ? l.amenities.map((a) => AMENITY_LABELS[a] ?? a).join(', ') : '' },
    ];
});

onMounted(() => {
    realEstateDetailStore.fetchDetail(props.id);
    trackView(props.id);
});

// Cập nhật tiêu đề trang cho Bất động sản
const detailTitle = computed(() => realEstateDetailStore.listing?.title || "Chi tiết bất động sản")
useHead({
    title: detailTitle,
})

// Gửi tracking khi Huỷ Component (Rời trang)
onBeforeUnmount(() => {
    cleanupTracking();
});
</script>
