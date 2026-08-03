<template>
    <div class="flex flex-col gap-4">
        <!-- Header -->
        <div class="flex items-center justify-between">
            <p class="text-sm font-semibold text-gray-900">
                Hình ảnh &amp; video <span class="text-red-500">*</span>
            </p>
            <n-button size="small" secondary @click="triggerVideoInput">
                <template #icon>
                    <n-icon>
                        <IconPlayCircleOutline />
                    </n-icon>
                </template>
                Thêm Video
            </n-button>
        </div>

        <!-- Validations -->
        <div class="flex flex-col gap-1.5">
            <div class="flex items-center gap-2">
                <n-icon size="16" :color="imageFiles.length >= 3 ? '#22c55e' : '#6b7280'">
                    <IconCheckmarkOutline v-if="imageFiles.length >= 3" />
                    <IconRadioButtonOff v-else />
                </n-icon>
                <span class="text-sm" :class="imageFiles.length >= 3 ? 'text-green-600' : 'text-gray-600'">
                    Đăng tối thiểu 3 ảnh
                </span>
            </div>
            <div class="flex items-start gap-2">
                <n-icon size="16" color="#6b7280" class="mt-0.5 shrink-0">
                    <IconRadioButtonOff />
                </n-icon>
                <span class="text-sm text-gray-600">
                    Gợi ý: Thêm video giúp tin đăng chân thật và cuốn hút hơn.
                    <br />
                    <span class="text-amber-500">🎁 Nhận ngay 1 lượt đẩy Tin Thường.</span>
                    <a href="#" class="text-red-500 underline ml-1 text-xs">Chi tiết</a>
                </span>
            </div>
        </div>

        <!-- Dropzone -->
        <div class="flex items-center gap-3 rounded-lg border-2 border-dashed border-gray-300 bg-gray-50 px-4 py-5 cursor-pointer transition-colors hover:border-red-400 hover:bg-red-50/20"
            :class="{ 'border-red-400 bg-red-50/20': isDragging }" @dragover.prevent="isDragging = true"
            @dragleave.prevent="isDragging = false" @drop.prevent="onDrop" @click="triggerImageInput">
            <div class="w-9 h-9 rounded-full bg-gray-200 flex items-center justify-center shrink-0">
                <n-icon size="20" color="#9ca3af">
                    <IconAddOutline />
                </n-icon>
            </div>
            <div class="flex-1">
                <p class="text-sm font-semibold text-gray-800 mb-0.5">
                    Kéo thả hình ảnh và video
                    <span class="text-gray-400 font-normal">({{ totalCount }}/25)</span>
                </p>
                <p class="text-xs text-gray-400 mb-1">Hỗ trợ các định dạng PNG, JPG, JPEG, GIF, MP4, MOV</p>
                <span class="text-xs text-red-500 cursor-pointer hover:underline" @click.stop="triggerImageInput">
                    Tải lên từ thiết bị
                </span>
            </div>
        </div>

        <!-- Hidden inputs -->
        <input ref="imageInputRef" type="file" accept="image/png,image/jpeg,image/jpg,image/gif" multiple class="hidden"
            @change="onImageInputChange" />
        <input ref="videoInputRef" type="file" accept="video/mp4,video/quicktime" class="hidden"
            @change="onVideoInputChange" />

        <!-- Preview grid -->
        <div v-if="files.length > 0" class="grid grid-cols-4 gap-2.5">
            <div v-for="(item, index) in files" :key="item.id"
                class="relative aspect-[4/3] rounded-lg overflow-hidden bg-gray-100 group">

                <!-- Badge ảnh đại diện -->
                <span v-if="index === 0 && item.fileType === 'image'"
                    class="absolute top-1 left-1 z-10 bg-red-500 text-white text-[10px] font-medium px-1.5 py-0.5 rounded">
                    Ảnh đại diện
                </span>

                <!-- Preview -->
                <img v-if="item.fileType === 'image'" :src="item.previewUrl" :alt="item.file.name"
                    class="w-full h-full object-cover" />
                <img v-else-if="item.thumbnailUrl" :src="item.thumbnailUrl" alt="Video thumbnail"
                    class="w-full h-full object-cover" />
                <div v-else class="w-full h-full flex items-center justify-center bg-gray-100">
                    <n-icon size="28" color="#9ca3af">
                        <IconVideoOutline />
                    </n-icon>
                </div>

                <!-- Progress bar -->
                <div v-if="item.status === 'uploading'" class="absolute inset-x-0 bottom-0 h-1 bg-gray-200">
                    <div class="h-full bg-red-500 transition-all" :style="{ width: item.progress + '%' }" />
                </div>

                <!-- Done icon -->
                <div v-if="item.status === 'done'"
                    class="absolute top-1 right-1 w-4 h-4 rounded-full bg-green-500 flex items-center justify-center">
                    <n-icon size="10" color="#fff">
                        <IconCheckmarkOutline />
                    </n-icon>
                </div>

                <!-- Error icon -->
                <div v-if="item.status === 'error'"
                    class="absolute top-1 right-1 w-4 h-4 rounded-full bg-red-500 flex items-center justify-center">
                    <n-icon size="10" color="#fff">
                        <IconCloseOutline />
                    </n-icon>
                </div>


                <!-- Nút xoá -->
                <button
                    class="absolute top-1 right-1 w-5 h-5 rounded-full bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity hover:bg-black/70"
                    @click.stop="removeFile(item.id)">
                    <n-icon size="10" color="#fff">
                        <IconCloseOutline />
                    </n-icon>
                </button>
            </div>
        </div>



        <!-- Tiêu chuẩn -->
        <div>
            <p class="text-sm font-semibold text-gray-800 mb-3">Tiêu chuẩn hình ảnh/video</p>
            <div class="grid grid-cols-2 gap-3">
                <div>
                    <p class="text-xs font-medium text-gray-600 mb-2">Video</p>
                    <n-card :bordered="true" size="small" content-style="padding: 0;">
                        <table class="w-full text-xs">
                            <tbody>
                                <tr v-for="row in videoStandards" :key="row.label"
                                    class="border-b border-gray-100 last:border-b-0">
                                    <td class="px-3 py-2 text-gray-500 w-1/2">{{ row.label }}</td>
                                    <td class="px-3 py-2 text-gray-800">
                                        {{ row.value }}
                                        <n-tag v-if="row.isNew" size="tiny" type="error" :bordered="false"
                                            class="ml-1">Mới</n-tag>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </n-card>
                </div>
                <div>
                    <p class="text-xs font-medium text-gray-600 mb-2">Hình ảnh</p>
                    <n-card :bordered="true" size="small" content-style="padding: 0;">
                        <table class="w-full text-xs">
                            <tbody>
                                <tr v-for="row in imageStandards" :key="row.label"
                                    class="border-b border-gray-100 last:border-b-0">
                                    <td class="px-3 py-2 text-gray-500 w-1/2">{{ row.label }}</td>
                                    <td class="px-3 py-2 text-gray-800">{{ row.value }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </n-card>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { FileItem } from '~/types/uploadmedia'
import { uploadFile } from '~/composables/upload'
import { useCreatePost } from '~/stores/create-post'

const { validateImage, validateVideo } = useValidate()

const postStore = useCreatePost()

const imageInputRef = ref<HTMLInputElement | null>(null)
const videoInputRef = ref<HTMLInputElement | null>(null)

const isDragging = ref(false)
const files = ref<FileItem[]>([])

const imageFiles = computed(() => files.value.filter(f => f.fileType === 'image'))
const totalCount = computed(() => files.value.length)

let fileIdCounter = 0
function generateId(): string {
    return `file_${Date.now()}_${++fileIdCounter}`
}

const MAX_FILES = 25
const MAX_VIDEOS = 3
function triggerImageInput() {
    imageInputRef.value?.click()
}

function triggerVideoInput() {
    videoInputRef.value?.click()
}

function addFiles(fileList: FileList, type: 'image' | 'video') {

    const currentVideos = files.value.filter(f => f.fileType === 'video').length
    const remainingTotal = MAX_FILES - totalCount.value
    const remainingVideo = MAX_VIDEOS - currentVideos

    if (remainingTotal <= 0) return

    if (type === 'video' && remainingVideo <= 0) {
        window.message?.warning(`Chỉ được tải lên tối đa ${MAX_VIDEOS} video`)
        return
    }

    const limit = type === 'video'
        ? Math.min(remainingTotal, remainingVideo)
        : remainingTotal

    const newItems: FileItem[] = Array.from(fileList).slice(0, limit).map((file) => ({
        id: generateId(),
        file,
        fileType: type,
        previewUrl: URL.createObjectURL(file),
        status: 'pending' as const,
        progress: 0,
    }))

    files.value.push(...newItems)
    newItems.forEach(({ id }) => {
        uploadFile(files.value.find(f => f.id === id)!).then(() => {
            syncImageIdsToStore()
        })
    })
}

// Đồng bộ id các ảnh/video đã upload thành công lên store để gửi khi Đăng tin
function syncImageIdsToStore() {
    postStore.form.image_ids = files.value
        .filter((f) => f.status === 'done' && f.imageId)
        .map((f) => f.imageId!)
}

function onDrop(e: DragEvent) {
    isDragging.value = false
    const files = e.dataTransfer?.files
    if (files && files.length > 0) {
        for (const file of Array.from(files)) {
            const dt = new DataTransfer()
            dt.items.add(file)
            if (file.type.startsWith('image/')) {
                addFiles(dt.files, 'image')
            } else if (file.type.startsWith('video/')) {
                addFiles(dt.files, 'video')
            }
        }
    }
}

function onImageInputChange(e: Event) {
    const target = e.target as HTMLInputElement
    if (target.files && target.files.length > 0) {
        const validFiles: File[] = []
        for (const file of Array.from(target.files)) {
            const result = validateImage(file)
            if (!result.valid) {
                window.message?.warning(`Ảnh "${file.name}": ${result.message}`)
                continue
            }
            validFiles.push(file)
        }
        if (validFiles.length > 0) {
            const dt = new DataTransfer()
            validFiles.forEach(f => dt.items.add(f))
            addFiles(dt.files, 'image')
        }
    }
    target.value = ''
}

function onVideoInputChange(e: Event) {
    const target = e.target as HTMLInputElement
    if (target.files && target.files.length > 0) {
        const validFiles: File[] = []
        for (const file of Array.from(target.files)) {
            const result = validateVideo(file)
            if (!result.valid) {
                window.message?.warning(`Video "${file.name}": ${result.message}`)
                continue
            }
            validFiles.push(file)
        }
        if (validFiles.length > 0) {
            const dt = new DataTransfer()
            validFiles.forEach(f => dt.items.add(f))
            addFiles(dt.files, 'video')
        }
    }
    target.value = ''
}

function removeFile(id: string) {
    const item = files.value.find(f => f.id === id)
    if (item) URL.revokeObjectURL(item.previewUrl)
    files.value = files.value.filter(f => f.id !== id)
    syncImageIdsToStore()
}

function validateImageCount(): boolean {
    if (imageFiles.value.filter((f) => f.status === 'done').length < 3) {
        window.message?.warning('Vui lòng đăng ít nhất 3 ảnh')
        return false
    }
    return true
}

defineExpose({ validateImageCount })

const videoStandards = [
    { label: 'Thời lượng', value: 'Tối thiểu 10 giây, tối đa 5 phút' },
    { label: 'Độ phân giải', value: 'Tối thiểu 720p' },
    { label: 'Kích thước', value: 'Tối đa 200MB', isNew: true },
    { label: 'Định dạng', value: 'MP4, MOV' },
]

const imageStandards = [
    { label: 'Kích thước', value: 'Tối thiểu 400 x 300px' },
    { label: 'Dung lượng', value: 'Tối đa 10MB' },
    { label: 'Định dạng', value: 'PNG, JPG, JPEG, GIF' },
    { label: 'Tỷ lệ', value: '4:3 hoặc vuông' },
]
</script>
