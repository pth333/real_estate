<template>
    <div class="md:col-span-1">
        <div class="border border-gray-200 bg-white rounded-lg p-4 shadow-sm sticky top-6 space-y-4">
            <h3 class="text-base font-bold text-gray-800 border-b border-gray-100 pb-2">Thông tin liên hệ</h3>

            <!-- Người bán -->
            <div class="flex items-center gap-2.5">
                <n-avatar round :size="44" color="#ecfdf5"
                    style="color: #059669; font-weight: 700; border: 1px solid #a7f3d0;">
                    {{ agentInitial }}
                </n-avatar>
                <div>
                    <p class="font-bold text-gray-900 text-sm leading-tight">{{
                        realEstateDetailStore.listing?.agent_name
                        || 'Người đăng' }}
                    </p>
                    <p class="text-[11px] text-gray-400 mt-0.5">Thành viên môi giới uy tín</p>
                </div>
            </div>

            <!-- SĐT -->
            <div v-if="realEstateDetailStore.listing?.agent_phone" class="space-y-2.5">
                <n-descriptions :column="1" size="small" bordered label-placement="left">
                    <n-descriptions-item label="Số điện thoại">
                        <span class="font-bold tracking-wide">{{ realEstateDetailStore.maskedPhone }}</span>
                    </n-descriptions-item>
                </n-descriptions>

                <n-button v-if="!realEstateDetailStore.showPhone" type="warning" ghost block
                    @click="realEstateDetailStore.showPhone = true">
                    Hiện số điện thoại
                </n-button>
            </div>

            <!-- Action buttons -->
            <n-space vertical>
                <n-button v-if="realEstateDetailStore.listing?.agent_phone" type="primary" block @click="handleCall">
                    <template #icon>
                        <IconPhone />
                    </template>
                    Yêu cầu gọi lại tư vấn
                </n-button>
                <n-button block ghost @click="realEstateDetailStore.handleShare">
                    Chia sẻ tin đăng này
                </n-button>
            </n-space>
        </div>
    </div>
</template>
<script setup lang="ts">
import { useRealEstateDetail } from '~/stores/detail/real_estate_detail';
const realEstateDetailStore = useRealEstateDetail()


const agentInitial = computed(() => {
    const name = realEstateDetailStore.listing?.agent_name || 'Q';
    return name.charAt(0).toUpperCase();
});

function handleCall() {
    if (realEstateDetailStore.listing?.agent_phone) {
        window.open(`tel:${realEstateDetailStore.listing.agent_phone}`, '_self');
    }
}


</script>