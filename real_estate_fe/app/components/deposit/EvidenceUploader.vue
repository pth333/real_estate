<template>
  <div class="flex flex-col gap-2">
    <n-upload :custom-request="handleUpload" :file-list="fileList" list-type="image-card" :max="6"
      @remove="handleRemove">
      Chọn ảnh bằng chứng
    </n-upload>
    <span class="text-xs text-gray-400">
      Hỗ trợ PNG/JPG/GIF, tối đa 6 ảnh. Ảnh sẽ được gửi kèm cho admin xem xét.
    </span>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NUpload, type UploadCustomRequestOptions, type UploadFileInfo } from 'naive-ui'
import { uploadImageToBackend } from '~/composables/upload'

const props = defineProps<{
  urls: string[]
}>()

const emit = defineEmits<{
  'update:urls': [value: string[]]
}>()

const fileList = ref<UploadFileInfo[]>([])

// Upload ảnh lên storage, lưu lại public_url để gửi cho backend làm bằng chứng
async function handleUpload({ file, onFinish, onError }: UploadCustomRequestOptions) {
  if (!file.file) return
  try {
    const result = await uploadImageToBackend(file.file, undefined)
    emit('update:urls', [...props.urls, result.public_url])
    onFinish()
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Upload ảnh thất bại'
    window.message?.error(message)
    onError()
  }
}

function handleRemove({ file }: { file: UploadFileInfo }) {
  const index = fileList.value.findIndex((item) => item.id === file.id)
  if (index === -1) return
  const nextUrls = [...props.urls]
  nextUrls.splice(index, 1)
  fileList.value.splice(index, 1)
  emit('update:urls', nextUrls)
}
</script>
