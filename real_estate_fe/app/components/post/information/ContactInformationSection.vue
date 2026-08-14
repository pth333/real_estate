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

            <!-- Số điện thoại (Dưới dạng select xổ xuống chọn hoặc thêm mới) -->
            <div>
                <label class="text-sm text-gray-700 block mb-1.5">
                    Số điện thoại <span class="text-red-500">*</span>
                </label>
                <n-select v-model:value="postStore.form.contact_phone" :options="phoneOptions"
                    placeholder="Chọn số điện thoại liên hệ"
                    :status="postStore.errorsContact.contact_phone ? 'error' : 'default'"
                    @update:value="handlePhoneSelect" />
                <span v-if="postStore.errorsContact.contact_phone" class="text-xs text-red-500 mt-1 block">{{
                    postStore.errorsContact.contact_phone }}</span>
            </div>
        </div>
    </n-card>
</template>

<script setup lang="ts">
import { computed, nextTick, ref } from 'vue'
import { useCreatePost } from '~/stores/create-post'
import { usePhoneVerification } from '~/composables/usePhoneVerification'

const postStore = useCreatePost()
const collapsed = ref(false)
const { verifiedPhones, showOTPModal } = usePhoneVerification()

// Định nghĩa danh sách option bao gồm các số đã xác thực và lựa chọn thêm mới
const phoneOptions = computed(() => {
    const list = verifiedPhones.value.map(phone => ({
        label: phone,
        value: phone
    }))
    list.push({
        label: '+ Thêm số điện thoại mới',
        value: 'add_new_phone'
    })
    return list
})

// Xử lý khi thay đổi giá trị trong select
const handlePhoneSelect = (value: string) => {
    if (value === 'add_new_phone') {
        // Mở modal OTP để nhập số điện thoại mới và OTP
        showOTPModal.value = true

        // Khôi phục lại giá trị hiển thị cũ để không hiển thị chữ 'add_new_phone' trên ô select
        nextTick(() => {
            const previousValue = verifiedPhones.value.includes(postStore.form.contact_phone)
                ? postStore.form.contact_phone
                : (verifiedPhones.value[0] || '')
            postStore.form.contact_phone = previousValue
        })
    } else {
        clearError('contact_phone')
    }
}

// Xóa lỗi của 1 field khi người dùng bắt đầu nhập
const clearError = (field: keyof typeof postStore.errorsContact) => {
    postStore.errorsContact[field] = ''
}
</script>
