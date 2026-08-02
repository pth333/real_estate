<template>
    <div>
        <!-- Nút mở modal, đặt bên phải ô tiêu đề -->
        <n-button size="medium" type="primary" ghost @click="showModal = true">
            <template #icon>
                <n-icon>
                    <IconCreateOutline />
                </n-icon>
            </template>
            Tạo với AI
        </n-button>

        <!-- Modal: tự quản lý mọi logic bên trong -->
        <n-modal v-model:show="showModal" title="Tạo tiêu đề & mô tả với AI" style="width: 560px">
            <n-card title="Tạo tiêu đề & mô tả với AI">
                <!-- Chọn văn phong -->
                <div class="mb-4">
                    <label class="text-sm text-gray-700 block mb-1.5">Văn phong</label>
                    <n-radio-group v-model:value="tone">
                        <n-space>
                            <n-radio value="lich_su">Lịch sự</n-radio>
                            <n-radio value="tre_trung">Trẻ trung</n-radio>
                        </n-space>
                    </n-radio-group>
                </div>

                <!-- Nút tạo nội dung -->
                <div class="mb-4">
                    <n-button type="primary" :loading="generating" @click="generateContent">
                        <template #icon>
                            <n-icon>
                                <IconCreateOutline />
                            </n-icon>
                        </template>
                        Tạo nội dung
                    </n-button>
                </div>

                <!-- Nội dung AI tạo ra — hiển thị trong input để xem trước / chỉnh sửa -->
                <div v-if="contentAI.title || contentAI.description">
                    <div class="mb-4">
                        <label class="text-sm text-gray-700 block mb-1.5">Tiêu đề</label>
                        <n-input v-model:value="contentAI.title" placeholder="AI sẽ đề xuất tiêu đề" :maxlength="99"
                            show-count />
                    </div>
                    <div>
                        <label class="text-sm text-gray-700 block mb-1.5">Mô tả</label>
                        <n-input v-model:value="contentAI.description" type="textarea" placeholder="AI sẽ đề xuất mô tả"
                            :maxlength="3000" :rows="5" show-count />
                    </div>
                </div>
                <p v-else class="text-sm text-gray-400">
                    Bấm <b>"Tạo nội dung"</b> để AI đề xuất tiêu đề và mô tả.
                </p>

                <template #action>
                    <div class="flex justify-end gap-2">
                        <n-button @click="showModal = false">Hủy</n-button>
                        <n-button type="primary" :disabled="!contentAI.title && !contentAI.description"
                            @click="applyContent">
                            Áp dụng
                        </n-button>
                    </div>
                </template>
            </n-card>
        </n-modal>
    </div>
</template>

<script setup lang="ts">
import { useCreatePost } from '~/stores/create-post'

interface AIContent {
    title: string
    description: string
}
const { $api } = useNuxtApp()
const postStore = useCreatePost()
const showModal = ref(false)
const tone = ref<'lich_su' | 'tre_trung'>('lich_su')
const generating = ref(false)
const contentAI = ref<AIContent>({
    title: '',
    description: '',
})

const generateContent = async () => {
    generating.value = true
    try {
        const response = await $api.post<{ data: AIContent }>('/ai/generate-content', {
            tone: tone.value,
        })
        contentAI.value = response.data
    } finally {
        generating.value = false
    }
}

const applyContent = () => {
    postStore.form.title = contentAI.value?.title || ''
    postStore.form.description = contentAI.value?.description || ''
    // Xóa lỗi (nếu có) vì đã có nội dung
    postStore.errorsDescription.title = ''
    postStore.errorsDescription.description = ''
    showModal.value = false
}
</script>
