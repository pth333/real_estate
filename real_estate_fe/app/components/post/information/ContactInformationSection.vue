<template>
    <n-card title="Thông tin liên hệ" size="small" :segmented="{ content: true }">
        <template #header-extra>
            <div class="cursor-pointer" @click="collapsed = !collapsed">
                <n-icon class="text-gray-500 transition-transform duration-200" :class="{ 'rotate-180': !collapsed }">
                    <IconChevronDownOutline />
                </n-icon>
            </div>
        </template>

        <div v-show="!collapsed">
            <!-- Tên liên hệ -->
            <div class="mb-4">
                <label class="text-sm text-gray-700 block mb-1.5">
                    Tên liên hệ <span class="text-red-500">*</span>
                </label>
                <n-input v-model:value="postStore.form.contact_name" placeholder="Nhập tên liên hệ"
                    :status="postStore.errorsContact.contact_name ? 'error' : 'default'"
                    @update:value="clearError('contact_name')" />
                <span v-if="postStore.errorsContact.contact_name" class="text-xs text-red-500 mt-1 block">{{
                    postStore.errorsContact.contact_name }}</span>
            </div>

            <!-- Email -->
            <div class="mb-4">
                <label class="text-sm text-gray-700 block mb-1.5">
                    Email <span class="text-red-500">*</span>
                </label>
                <n-input v-model:value="postStore.form.contact_email" placeholder="Nhập email"
                    :status="postStore.errorsContact.contact_email ? 'error' : 'default'"
                    @update:value="clearError('contact_email')" />
                <span v-if="postStore.errorsContact.contact_email" class="text-xs text-red-500 mt-1 block">{{
                    postStore.errorsContact.contact_email }}</span>
            </div>

            <!-- Số điện thoại -->
            <div>
                <label class="text-sm text-gray-700 block mb-1.5">
                    Số điện thoại <span class="text-red-500">*</span>
                </label>
                <n-input v-model:value="postStore.form.contact_phone" placeholder="Nhập số điện thoại"
                    :status="postStore.errorsContact.contact_phone ? 'error' : 'default'"
                    @update:value="clearError('contact_phone')" />
                <span v-if="postStore.errorsContact.contact_phone" class="text-xs text-red-500 mt-1 block">{{
                    postStore.errorsContact.contact_phone }}</span>
            </div>
        </div>
    </n-card>
</template>

<script setup lang="ts">
import { useCreatePost } from '~/stores/create-post'

const postStore = useCreatePost()
const collapsed = ref(false)

// Xóa lỗi của 1 field khi người dùng bắt đầu nhập
const clearError = (field: keyof typeof postStore.errorsContact) => {
    postStore.errorsContact[field] = ''
}

// Validate phần Thông tin liên hệ, hiển thị lỗi dưới từng ô input
const validate = (): boolean => {
    postStore.errorsContact = {
        contact_name: postStore.form.contact_name ? "" : "Vui lòng nhập tên liên hệ",
        contact_email: postStore.form.contact_email ? "" : "Vui lòng nhập email",
        contact_phone: postStore.form.contact_phone ? "" : "Vui lòng nhập số điện thoại",
    }
    return Object.values(postStore.errorsContact).every((msg) => msg === "")
}

defineExpose({ validate })
</script>
