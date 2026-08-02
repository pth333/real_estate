<template>
    <div class="min-h-screen bg-white flex flex-col font-sans">

        <!-- Header -->
        <div class="flex justify-between items-center px-6 py-3.5 border-b border-gray-200">
            <span class="text-lg font-semibold text-gray-900">Tạo tin đăng</span>
            <div class="flex items-center gap-2.5">
                <n-button text>
                    <template #icon>
                        <n-icon>
                            <IconEyeOutline />
                        </n-icon>
                    </template>
                    Xem trước
                </n-button>
                <n-button>Thoát</n-button>
            </div>
        </div>

        <!-- Progress -->
        <div class="px-6 pb-3 border-b border-gray-200">
            <span class="text-xs text-red-500 font-medium block mb-1.5 mt-1.5">{{ postStore.currentStepLabel }}</span>
            <n-progress type="line" :percentage="postStore.currentStepProgress" :show-indicator="false" color="#ef4444"
                rail-color="#f3f4f6" />
        </div>

        <div class="flex-1 px-6 py-6 w-full max-w-2xl mx-auto flex flex-col gap-4">

            <!-- Section: Nhu cầu -->
            <!-- Bước 1. Thông tin BĐS -->
            <!-- Section: Địa chỉ -->
            <div class="flex flex-col gap-4" v-show="postStore.tab === 'information'">
                <DemandSection />

                <AddressSection />
                <!-- Section: Thông tin chính -->
                <MainInformationSection />
                <!-- Section: Thông tin khác -->
                <InformationOther />
                <!-- Section: Thông tin liên hệ -->
                <ContactInformationSection />
                <!-- Section: Tiêu đề & mô tả -->
                <DescriptionSection />
            </div>
            <!-- Bước 2. Upload ảnh, video -->
            <div v-show="postStore.tab === 'upload'">
                <UploadImageVideo ref="uploadComponent" />
            </div>
        </div>

        <!-- Footer -->
        <div class="sticky bottom-0 flex justify-between px-6 py-4 border-t border-gray-200 bg-white">
            <n-button size="large" @click="prevPage">
                Quay lại
            </n-button>
            <n-button size="large" type="error" @click="nextPage">
                Tiếp tục
            </n-button>
            <n-button v-show="postStore.tab === 'review'" size="large" type="error" @click="submitCreatePost">
                Đăng tin
            </n-button>
        </div>

        <ModalOTPAuthentication v-model:show="showOTPModal" />

    </div>
</template>
<script setup lang="ts">
import { useCreatePost } from '~/stores/create-post'
definePageMeta({
    alias: "/nguoi-ban/dang-tin",
})
const { phoneVerified } = usePhoneVerification()
const postStore = useCreatePost()
const showOTPModal = ref(false)
const uploadComponent = ref()

const nextPage = () => {

    if (postStore.tab === 'information') {
        if (!postStore.validateInformation()) {
            window.message?.warning('Vui lòng nhập đầy đủ thông tin')
            return
        }
        postStore.tab = 'upload'
    } else if (postStore.tab === 'upload') {
        if (!uploadComponent.value.validateImageCount()) return
        postStore.tab = 'review'
    } else if (postStore.tab === 'review') {
    }
}

const prevPage = () => {
    if (postStore.tab === 'upload') {
        postStore.tab = 'information'
    } else if (postStore.tab === 'review') {
        postStore.tab = 'upload'
    }
}

const submitCreatePost = () => {
};

onMounted(() => {
    if (!phoneVerified.value) {
        showOTPModal.value = true
    }
})

</script>