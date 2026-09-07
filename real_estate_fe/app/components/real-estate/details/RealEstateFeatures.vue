<template>
    <div class="bg-white rounded-lg py-4">
        <h2 class="mb-4 border-b border-gray-200 pb-2 text-base font-bold text-gray-800">
            Đặc điểm bất động sản
        </h2>

        <div class="grid grid-cols-1 gap-x-8 sm:grid-cols-2">
            <template v-for="attr in attrs" :key="attr.label">
                <div v-if="attr.value" class="flex items-center justify-between border-b border-gray-100 py-2">
                    <span class="flex items-center gap-2 text-sm text-gray-500">
                        <component :is="attr.icon" class="h-4 w-4 shrink-0" />
                        {{ attr.label }}
                    </span>
                    <span class="text-sm font-medium text-gray-800">{{ attr.value }}</span>
                </div>
            </template>
        </div>
    </div>
</template>

<script setup lang="ts">
import { formatPrice, formatPricePerM2 } from '~/utils/format';
import { useRealEstateDetail } from '~/stores/detail/real_estate_detail';
import IconPrice from '~/icons/IconPrice.vue';
import IconArea from '~/icons/IconArea.vue';
import IconBed from '~/icons/IconBed.vue';
import IconBath from '~/icons/IconBath.vue';
import IconBuilding from '~/icons/IconBuilding.vue';
import IconCompass from '~/icons/IconCompass.vue';
import IconBalcony from '~/icons/IconBalcony.vue';
import IconShieldCheck from '~/icons/IconShieldCheck.vue';
import IconSofa from '~/icons/IconSofa.vue';
import IconZap from '~/icons/IconZap.vue';
import IconDroplet from '~/icons/IconDroplet.vue';
import IconWifi from '~/icons/IconWifi.vue';
import IconSparkles from '~/icons/IconSparkles.vue';

const store = useRealEstateDetail();

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

const formatPriceNumber = (n: number) => n.toLocaleString('vi-VN');

const attrs = computed(() => {
    const l = store.listing;
    if (!l) return [];
    return [
        { label: 'Khoảng giá', icon: IconPrice, value: formatPrice(l.price_vnd) },
        { label: 'Giá/m²', icon: IconPrice, value: formatPricePerM2(l.price_per_m2) },
        { label: 'Diện tích', icon: IconArea, value: `${l.acreage} m²` },
        { label: 'Số phòng ngủ', icon: IconBed, value: l.bedrooms ? `${l.bedrooms} phòng` : '' },
        { label: 'Số phòng tắm, vệ sinh', icon: IconBath, value: l.bathrooms ? `${l.bathrooms} phòng` : '' },
        { label: 'Số tầng', icon: IconBuilding, value: l.floors ? `${l.floors} tầng` : '' },
        { label: 'Hướng nhà', icon: IconCompass, value: l.house_direction ? DIRECTION_LABELS[l.house_direction] ?? l.house_direction : '' },
        { label: 'Hướng ban công', icon: IconBalcony, value: l.balcony_direction ? DIRECTION_LABELS[l.balcony_direction] ?? l.balcony_direction : '' },
        { label: 'Pháp lý', icon: IconShieldCheck, value: l.legal_docs ? LEGAL_DOC_LABELS[l.legal_docs] ?? l.legal_docs : '' },
        { label: 'Nội thất', icon: IconSofa, value: l.interior ? INTERIOR_LABELS[l.interior] ?? l.interior : '' },
        { label: 'Giá điện', icon: IconZap, value: l.price_electricity ? `${formatPriceNumber(l.price_electricity)} đ/kWh` : '' },
        { label: 'Giá nước', icon: IconDroplet, value: l.price_water ? `${formatPriceNumber(l.price_water)} đ/m³` : '' },
        { label: 'Giá internet', icon: IconWifi, value: l.price_internet ? `${formatPriceNumber(l.price_internet)} đ/tháng` : '' },
        { label: 'Tiện ích', icon: IconSparkles, value: l.amenities?.length ? l.amenities.map((a) => AMENITY_LABELS[a] ?? a).join(', ') : '' },
    ];
});
</script>
