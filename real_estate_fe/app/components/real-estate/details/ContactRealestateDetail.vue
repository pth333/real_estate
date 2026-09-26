<template>
    <div class="lg:col-span-1">
        <!-- Khối liên hệ đầy đủ: mobile/tablet nằm dưới nội dung (theo thứ tự DOM),
             desktop dính ở cột phải. -->
        <div class="space-y-4 rounded-lg border border-gray-200 bg-white p-4 shadow-sm lg:sticky lg:top-6">
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
                <n-button type="primary" block @click="handleOpenBooking">
                    <template #icon>
                        <IconWallet />
                    </template>
                    Đặt cọc giữ lịch xem nhà
                </n-button>
                <n-button v-if="realEstateDetailStore.listing?.agent_phone" block ghost @click="handleCall">
                    <template #icon>
                        <IconPhone />
                    </template>
                    Yêu cầu gọi lại tư vấn
                </n-button>
                <n-button block ghost @click="realEstateDetailStore.handleShare">
                    Chia sẻ tin đăng này
                </n-button>
            </n-space>

            <!-- Modal đặt cọc: mức cọc do hệ thống tính theo giá BĐS, tiền platform giữ -->
            <DepositBookingModal v-if="listingId" v-model:show="showBookingModal" :real-estate-id="listingId"
                :estate-title="realEstateDetailStore.listing?.title" />
        </div>

        <!-- Thanh CTA cố định đáy màn hình (chỉ mobile/tablet): luôn thấy nút gọi & đặt cọc khi cuộn -->
        <Teleport to="body">
            <div class="fixed inset-x-0 bottom-0 z-40 flex items-center gap-2 border-t border-gray-200 bg-white/95 px-3 pt-2 backdrop-blur lg:hidden"
                style="padding-bottom: max(0.5rem, env(safe-area-inset-bottom));">
                <n-button v-if="realEstateDetailStore.listing?.agent_phone" type="primary" ghost class="flex-1"
                    @click="handleCall">
                    <template #icon>
                        <IconPhone />
                    </template>
                    Gọi tư vấn
                </n-button>
                <n-button type="primary" class="flex-1" @click="handleOpenBooking">
                    <template #icon>
                        <IconWallet />
                    </template>
                    Đặt cọc
                </n-button>
            </div>
        </Teleport>
    </div>
</template>
<script setup lang="ts">
import { ref } from 'vue';
import { useRealEstateDetail } from '~/stores/detail/real_estate_detail';
import { useAuthStore } from '~/stores/auth';
import IconWallet from '~/icons/IconWallet.vue';
import DepositBookingModal from '~/components/deposit/DepositBookingModal.vue';

const realEstateDetailStore = useRealEstateDetail()
const authStore = useAuthStore()
const route = useRoute()

// Trạng thái mở modal đặt cọc
const showBookingModal = ref(false);

// ID BĐS đang xem — bắt buộc phải có mới đặt cọc được
const listingId = computed(() => realEstateDetailStore.listing?.id ?? 0);

/**
 * Trang chi tiết BĐS là trang CÔNG KHAI nên khách vãng lai vẫn xem được.
 * Nhưng đặt cọc thì cần tài khoản → nhắc đăng nhập và quay lại đúng tin này.
 */
function handleOpenBooking() {
    if (!authStore.isAuthenticated) {
        window.message?.info('Vui lòng đăng nhập để đặt cọc giữ lịch xem nhà');
        // navigateTo({ path: '/dang-nhap', query: { redirect: route.fullPath } });
        return;
    }
    showBookingModal.value = true;
}

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