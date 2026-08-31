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
            <n-button v-if="!postStore.isTabUpload()" size="large" type="error" @click="nextPage">
                Tiếp tục
            </n-button>
            <n-button v-if="postStore.isTabUpload()" size="large" type="error" @click="submitCreatePost">
                Đăng tin
            </n-button>
        </div>

        <ModalOTPAuthentication v-model:show="showOTPModal" />

        <!-- Modal hỏi tiếp tục bản nháp cũ -->
        <DraftModal v-model:show="showDraftModal" />

    </div>
</template>
<script setup lang="ts">
import { useCreatePost } from '~/stores/create-post'
import { useManagerStore } from '~/stores/manager'
import { InformationRealestate, type CreatePostResponse, type RealEstateResponse, type UpdatePostResponse } from '~/types/real_estate'
definePageMeta({
    alias: "/nguoi-ban/dang-tin",
})

useHead({
    title: "Đăng tin bất động sản",
})

const { phoneVerified, verifiedPhone, showOTPModal } = usePhoneVerification()
const { $api } = useNuxtApp()
const managerStore = useManagerStore()

const route = useRoute()
const isEdit = computed(() => route.query.id)

const postStore = useCreatePost()
const uploadComponent = ref()

const showDraftModal = ref(false)

const handleExit = () => {
    postStore.saveCurrentDraft()
    navigateTo('/')
}

const nextPage = () => {
    if (postStore.isTabInformation()) {
        if (!postStore.validateInformation()) {
            window.message?.warning('Vui lòng nhập đầy đủ thông tin')
            return
        }
        postStore.form.tab = 'upload'
    } else if (postStore.isTabUpload()) {
        if (!uploadComponent.value.validateImageCount()) return
        postStore.form.tab = 'review'
    }
    // } else if (postStore.form.isTabReview()) {
    // }
}

const prevPage = () => {
    if (postStore.isTabUpload()) {
        postStore.form.tab = 'information'
    } else if (postStore.isTabReview()) {
        postStore.form.tab = 'upload'
    }
}
const isSubmitting = ref(false)

const submitCreatePost = async () => {
    if (!uploadComponent.value.validateImageCount()) return
    isSubmitting.value = true
    try {
        let res;

        if (isEdit.value) {
            res = await $api.put<UpdatePostResponse>(
                `/manager/update-post/${isEdit.value}`,
                postStore.payload
            )
        } else {
            res = await $api.post<CreatePostResponse>(
                '/manager/create-post',
                postStore.payload,
            )
        }

        if (res.success) {
            window.message?.success(isEdit.value ? 'Cập nhật tin đăng thành công' : 'Đăng tin thành công')
            postStore.clearCurrentDraft()
            postStore.resetForm()
            // Đánh dấu cache bài viết cũ → trang quản lý bài viết sẽ fetch lại
            managerStore.invalidatePosts()
            navigateTo('/nguoi-ban/quan-ly-tin-dang')
        }
    } catch (err: unknown) {
    } finally {
        isSubmitting.value = false
    }
};

// Tải thông tin chi tiết bài đăng cũ và đưa vào form store bằng Class Object
const loadingPostDetail = async () => {
    try {
        const res = await $api.get<{ data: RealEstateResponse }>(`/real-estate/detail/${isEdit.value}`)
        if (res) {
            const editForm = InformationRealestate.fromResponse(res.data)
            postStore.form = editForm
            // Đẩy danh sách ảnh vào upload component để hiển thị preview
            if (editForm.images && editForm.images.length > 0) {
                nextTick(() => {
                    if (uploadComponent.value) {
                        uploadComponent.value.setFileList(editForm.images)
                    }
                })
            }
        }
    } catch (error) {
        window.message?.error("Không thể tải thông tin chi tiết bài viết cần sửa")
    }
}

watch(() => isEdit.value, (newId) => {
    if (newId) {
        postStore.resetForm()
        loadingPostDetail()
    } else {
        postStore.resetForm()
        // Có bản nháp cũ → hỏi người dùng có tiếp tục không
        if (postStore.loadCurrentDraft()) {
            showDraftModal.value = true
        }
        if (!phoneVerified.value) {
            showOTPModal.value = true
            return
        }
        postStore.form.contact_phone = verifiedPhone.value || ""
    }
}, { immediate: true })
</script>