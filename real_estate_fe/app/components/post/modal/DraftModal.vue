<template>
    <n-modal :show="show" preset="dialog" type="warning" title="Phát hiện bản nháp" :mask-closable="false"
        @update:show="onShowChange">
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
</template>
<script setup lang="ts">
import { useCreatePost } from '~/stores/create-post'
const { verifiedPhone } = usePhoneVerification()

const props = defineProps<{
    show: boolean
}>()

const emit = defineEmits<{
    'update:show': [value: boolean]
}>()

const postStore = useCreatePost()

const onShowChange = (value: boolean) => {
    emit('update:show', value)
}

// Tiếp tục bản nháp: áp draft vào form rồi đóng modal
const continueDraft = () => {
    postStore.applyDraft()
    emit('update:show', false)
}

// Hủy bản nháp: xóa draft + reset form rồi đóng modal
const discardDraft = () => {
    postStore.clearCurrentDraft()
    postStore.resetForm()
    postStore.form.contact_phone = verifiedPhone.value || ""
    emit('update:show', false)
}
</script>
