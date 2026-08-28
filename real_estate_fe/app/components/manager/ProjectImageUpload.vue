<template>
  <div class="flex flex-col gap-3">
    <div class="flex items-center justify-between">
      <p class="text-sm font-semibold text-gray-800">Hình ảnh dự án</p>
      <span class="text-xs text-gray-400">({{ totalCount }}/10)</span>
    </div>

    <!-- Dropzone -->
    <div
      class="flex items-center gap-3 border-2 border-dashed border-gray-300 bg-gray-50 px-4 py-4 cursor-pointer transition-colors hover:border-emerald-400 hover:bg-emerald-50/20"
      @click="triggerInput">
      <n-icon size="20" color="#9ca3af">
        <IconAddOutline />
      </n-icon>
      <div>
        <p class="text-sm font-medium text-gray-700">Tải ảnh dự án (PNG, JPG, JPEG)</p>
        <p class="text-xs text-gray-400 mt-0.5">Tối đa 10 ảnh</p>
      </div>
    </div>

    <input ref="inputRef" type="file" accept="image/png,image/jpeg,image/jpg,image/gif" multiple class="hidden"
      @change="onInputChange" />

    <!-- Preview: ảnh đã lưu (edit) + ảnh mới upload -->
    <div v-if="totalCount > 0" class="grid grid-cols-4 gap-2.5">
      <!-- Ảnh đã có của dự án -->
      <div v-for="img in existing" :key="'e' + img.id"
        class="relative aspect-[4/3] overflow-hidden bg-gray-100 group">
        <img :src="img.url" :alt="'Ảnh dự án ' + img.id" class="w-full h-full object-cover" />
        <button
          class="absolute top-1 right-1 z-10 w-5 h-5 bg-black/60 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity hover:bg-red-500"
          @click="removeExisting(img.id)">
          <IconCloseOutline class="h-3 w-3 text-white" />
        </button>
      </div>

      <!-- Ảnh mới đang upload -->
      <div v-for="item in newFiles" :key="item.id" class="relative aspect-[4/3] overflow-hidden bg-gray-100 group">
        <img :src="item.previewUrl" :alt="item.file.name" class="w-full h-full object-cover" />
        <div v-if="item.status !== 'done' && item.status !== 'error'"
          class="absolute inset-0 z-10 flex items-center justify-center bg-black/30">
          <n-spin size="small" />
        </div>
        <p v-if="item.status === 'error'"
          class="absolute inset-x-0 bottom-0 z-10 bg-red-500 text-[10px] text-white px-1 py-0.5 text-center">Lỗi upload</p>
        <button
          class="absolute top-1 right-1 z-10 w-5 h-5 bg-black/60 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity hover:bg-red-500"
          @click="removeNew(item)">
          <IconCloseOutline class="h-3 w-3 text-white" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FileItem } from '~/types/uploadmedia'
import { uploadFile } from '~/composables/upload'
import IconAddOutline from '~/icons/IconAddOutline.vue'
import IconCloseOutline from '~/icons/IconCloseOutline.vue'

// Ảnh dự án đã lưu (khi chỉnh sửa)
interface ProjectImage {
  id: number
  url: string
}

const props = defineProps<{
  existing?: ProjectImage[]
  ids: number[]
}>()

const emit = defineEmits<{
  'update:ids': [ids: number[]]
}>()

const { validateImage } = useValidate()

const inputRef = ref<HTMLInputElement | null>(null)
const newFiles = ref<FileItem[]>([])

const MAX_IMAGES = 10
const totalCount = computed(() => (props.existing?.length || 0) + newFiles.value.length)

function triggerInput() {
  inputRef.value?.click()
}

function onInputChange(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files) addFiles(input.files)
  input.value = ''
}

function addFiles(list: FileList) {
  const remaining = MAX_IMAGES - totalCount.value
  if (remaining <= 0) {
    window.message?.warning(`Chỉ được tải lên tối đa ${MAX_IMAGES} ảnh`)
    return
  }

  const items: FileItem[] = []
  Array.from(list).slice(0, remaining).forEach((file) => {
    const v = validateImage(file)
    if (!v.valid) {
      window.message?.warning(v.message)
      return
    }
    items.push({
      id: `file_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
      file,
      fileType: 'image',
      previewUrl: URL.createObjectURL(file),
      status: 'pending',
      progress: 0,
    })
  })

  newFiles.value.push(...items)
  items.forEach((item) => {
    uploadFile(item, 'project').then(() => {
      // Upload xong → thêm image_id vào danh sách ids
      if (item.status === 'done' && item.imageId) {
        emit('update:ids', [...props.ids, item.imageId])
      }
    })
  })
}

function removeNew(item: FileItem) {
  newFiles.value = newFiles.value.filter((f) => f.id !== item.id)
  if (item.imageId) {
    emit('update:ids', props.ids.filter((id) => id !== item.imageId))
  }
}

function removeExisting(id: number) {
  emit('update:ids', props.ids.filter((i) => i !== id))
}
</script>
