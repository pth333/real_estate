<template>
    <n-card title="Tiêu đề & mô tả" size="small" :segmented="{ content: true }">
        <template #header-extra>
            <div class="cursor-pointer" @click="collapsed = !collapsed">
                <n-icon class="text-gray-500 transition-transform duration-200" :class="{ 'rotate-180': !collapsed }">
                    <IconChevronDownOutline />
                </n-icon>
            </div>
        </template>

        <div v-show="!collapsed">
            <!-- Tiêu đề -->
            <div class="mb-4">
                <div class="flex items-center justify-between mb-1.5">
                    <label class="text-sm text-gray-700">
                        Tiêu đề <span class="text-red-500">*</span>
                    </label>
                    <ModalContentAI />
                </div>
                <n-input v-model:value="postStore.form.title" placeholder="Nhập tiêu đề tin đăng" :maxlength="99"
                    :status="postStore.errorsDescription.title ? 'error' : 'default'"
                    @update:value="clearError('title')" show-count />
                <div class="flex justify-between mt-1">
                    <span v-if="postStore.errorsDescription.title" class="text-xs text-red-500">{{
                        postStore.errorsDescription.title }}</span>
                </div>
            </div>

            <!-- Mô tả -->
            <div class="mb-4">
                <label class="text-sm text-gray-700 block mb-1.5">
                    Mô tả <span class="text-red-500">*</span>
                </label>
                <n-input v-model:value="postStore.form.description" type="textarea" placeholder="Nhập mô tả chi tiết"
                    :maxlength="3000" :rows="5" :status="postStore.errorsDescription.description ? 'error' : 'default'"
                    @update:value="clearError('description')" show-count />
                <div class="flex justify-between mt-1">
                    <span v-if="postStore.errorsDescription.description" class="text-xs text-red-500">{{
                        postStore.errorsDescription.description }}</span>
                </div>
            </div>

        </div>

    </n-card>
</template>

<script setup lang="ts">
import { useCreatePost } from '~/stores/create-post'

const postStore = useCreatePost()
const collapsed = ref(false)

// Xóa lỗi của 1 field khi người dùng sửa
const clearError = (field: keyof typeof postStore.errorsDescription) => {
    postStore.errorsDescription[field] = ''
}

// Validate phần Tiêu đề & mô tả (bắt buộc + độ dài), hiển thị lỗi dưới từng ô input
const validate = (): boolean => {
    const titleLen = postStore.form.title?.length ?? 0
    const descLen = postStore.form.description?.length ?? 0
    postStore.errorsDescription = {
        title: titleLen
            ? titleLen < 30
                ? `Tiêu đề tối thiểu 30 ký tự (hiện ${titleLen})`
                : ""
            : "Vui lòng nhập tiêu đề",
        description: descLen
            ? descLen < 30
                ? `Mô tả tối thiểu 30 ký tự (hiện ${descLen})`
                : ""
            : "Vui lòng nhập mô tả",
    }
    return Object.values(postStore.errorsDescription).every((msg) => msg === "")
}

defineExpose({ validate })
</script>
