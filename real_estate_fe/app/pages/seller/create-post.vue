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
                <n-button @click="handleExit">Thoát</n-button>
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
            <div class="flex flex-col gap-4" v-show="postStore.form.tab === 'information'">
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
            <div v-show="postStore.form.tab === 'upload'">
                <UploadImageVideo ref="uploadComponent" />
            </div>
        </div>

        <!-- Footer -->
        <div class="sticky bottom-0 flex justify-between px-6 py-4 border-t border-gray-200 bg-white">
            <n-button size="large" @click="prevPage">
                Quay lại
            </n-button>
            <n-button v-if="!postStore.form.isTabUpload()" size="large" type="error" @click="nextPage">
                Tiếp tục
            </n-button>
            <n-button v-if="postStore.form.isTabUpload()" size="large" type="error" @click="submitCreatePost">
                Đăng tin
            </n-button>
        </div>

        <ModalOTPAuthentication v-model:show="showOTPModal" />

        <!-- Modal hỏi tiếp tục bản nháp cũ -->
        <n-modal v-model:show="showDraftModal" preset="dialog" type="warning" title="Phát hiện bản nháp"
            :mask-closable="false" @click="onDraftModalClose">
            <template #default>
                <p>
                    Bạn có bản nháp tin đăng chưa hoàn thành. Bạn có muốn tiếp tục
                    sử dụng bản nháp trước đó không?
                </p>
            </template>
            <template #action>
                <n-button size="large" @click="discardDraft">
                    Không, tạo tin mới
                </n-button>
                <n-button size="large" type="error" @click="continueDraft">
                    Tiếp tục bản nháp
                </n-button>
            </template>
        </n-modal>

    </div>
</template>
<script setup lang="ts">
import { useCreatePost } from '~/stores/create-post'
definePageMeta({
    alias: "/nguoi-ban/dang-tin",
})
const { phoneVerified, verifiedPhone, showOTPModal } = usePhoneVerification()
const { $api } = useNuxtApp()

const postStore = useCreatePost()
const uploadComponent = ref()

const showDraftModal = ref(false)

const handleExit = () => {
    postStore.saveCurrentDraft()
    navigateTo('/')
}

const continueDraft = () => {
    postStore.applyDraft()
    showDraftModal.value = false
}

const discardDraft = () => {
    postStore.clearCurrentDraft()
    postStore.resetForm()
    showDraftModal.value = false
}

const onDraftModalClose = () => {
    showDraftModal.value = false
}
const nextPage = () => {
    if (postStore.form.isTabInformation()) {
        if (!postStore.validateInformation()) {
            window.message?.warning('Vui lòng nhập đầy đủ thông tin')
            return
        }
        postStore.form.tab = 'upload'
    } else if (postStore.form.isTabUpload()) {
        if (!uploadComponent.value.validateImageCount()) return
        postStore.form.tab = 'review'
    }
    // } else if (postStore.form.isTabReview()) {
    // }
}

const prevPage = () => {
    if (postStore.form.isTabUpload()) {
        postStore.form.tab = 'information'
    } else if (postStore.form.isTabReview()) {
        postStore.form.tab = 'upload'
    }
}

const isSubmitting = ref(false)

const submitCreatePost = async () => {
    isSubmitting.value = true
    try {
        const res = await $api.post<{ success: boolean; message?: string; data?: { id: number } }>(
            '/real-estate/create-post',
            postStore.payload,
        )
        if (res.success) {
            window.message?.success('Đăng tin thành công')
            postStore.clearCurrentDraft()
            postStore.resetForm()
            navigateTo('/')
        } else {
            window.message?.error(res.message || 'Đăng tin thất bại')
        }
    } catch (err: unknown) {
        const msg = (err as { data?: { message?: string } })?.data?.message
        window.message?.error(msg || 'Đăng tin thất bại, vui lòng thử lại')
    } finally {
        isSubmitting.value = false
    }
};


// Tự động điền số điện thoại khi đã xác thực thành công hoặc khi số xác thực thay đổi
watch(verifiedPhone, (newPhone) => {
    if (newPhone && !postStore.form.contact_phone) {
        postStore.form.contact_phone = newPhone
    }
}, { immediate: true })

onMounted(() => {
    if (!phoneVerified.value) {
        showOTPModal.value = true
        return
    }
    // Điền số điện thoại mặc định từ localStorage nếu chưa có
    if (verifiedPhone.value && !postStore.form.contact_phone) {
        postStore.form.contact_phone = verifiedPhone.value
    }
    // Có bản nháp cũ → hỏi người dùng có tiếp tục không
    if (postStore.loadCurrentDraft()) {
        showDraftModal.value = true
    }
})

</script>