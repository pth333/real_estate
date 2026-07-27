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
            <div class="flex flex-col gap-4" v-if="postStore.tab === 'information'">
                <DemandSection />

                <AddressSection />
                <!-- Section: Thông tin chính -->
                <MainInformationSection />
            </div>
            <!-- Bước 2. Upload ảnh, video -->
            <div v-else-if="postStore.tab === 'upload'">
                <UploadImageVideo />
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
const { showToast } = useToast();
const { validate } = useFormValidation()
const showOTPModal = ref(false)

const validateInformation = () => {
    const error = validate({
        province: { value: postStore.province, label: 'tỉnh/thành' },
        ward: { value: postStore.ward, label: 'phường/xã' },
        detail_address: { value: postStore.detail_address?.trim(), label: 'địa chỉ chi tiết' },
    })
    if (error) {
        showToast('error', `Vui lòng nhập ${error}`)
    }
    return error
}

const nextPage = () => {
    if (postStore.tab === 'information') {
        if (validateInformation()) return
        postStore.tab = 'upload'
    }
    if (postStore.tab === 'upload') { }

}

const prevPage = () => {
    if (postStore.tab === 'upload') {
        postStore.tab = 'information'
    } else if (postStore.tab === 'review') {
        postStore.tab = 'upload'
    }
}

onMounted(() => {
    if (!phoneVerified.value) {
        showOTPModal.value = true
    }
})

</script>