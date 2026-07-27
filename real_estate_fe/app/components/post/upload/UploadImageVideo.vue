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
        <div v-if="allFiles.length > 0" class="grid grid-cols-4 gap-2.5">
            <div v-for="(item, index) in allFiles" :key="item.id"
                class="relative aspect-[4/3] rounded-lg overflow-hidden bg-gray-100 group">

                <!-- Badge ảnh đại diện -->
                <span v-if="index === 0 && item.fileType === 'image'"
                    class="absolute top-1 left-1 z-10 bg-red-500 text-white text-[10px] font-medium px-1.5 py-0.5 rounded">
                    Ảnh đại diện
                </span>

                <!-- Preview -->
                <img v-if="item.fileType === 'image'" :src="item.previewUrl" :alt="item.file.name"
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

                <!-- Spinner khi đang lấy presign / confirming -->
                <div v-if="item.status === 'gettingPresign' || item.status === 'confirming'"
                    class="absolute inset-0 bg-black/30 flex items-center justify-center">
                    <n-spin size="small" />
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

        <p v-if="showError" class="text-xs text-red-500 -mt-2">
            Vui lòng tải lên ít nhất 3 ảnh.
        </p>

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
import type { ConfirmResponse, FileItem, PresignResponse } from '~/types/uploadmedia'



const { $api } = useNuxtApp()
const { showToast } = useToast()

const imageInputRef = ref<HTMLInputElement | null>(null)
const videoInputRef = ref<HTMLInputElement | null>(null)

const isDragging = ref(false)
const imageFiles = ref<FileItem[]>([])
const videoFiles = ref<FileItem[]>([])
const showError = ref(false)

const allFiles = computed(() => [...imageFiles.value, ...videoFiles.value])
const totalCount = computed(() => allFiles.value.length)

const uploadedImages = computed(() =>
    imageFiles.value.filter((f) => f.status === 'done' && f.imageId != null)
)
const allUploaded = computed(() =>
    [...imageFiles.value, ...videoFiles.value].filter((f) => f.status === 'done' && f.imageId != null)
)

let fileIdCounter = 0
function generateId(): string {
    return `file_${Date.now()}_${++fileIdCounter}`
}

const MAX_FILES = 25

function triggerImageInput() {
    imageInputRef.value?.click()
}

function triggerVideoInput() {
    videoInputRef.value?.click()
}

function addFiles(fileList: FileList, type: 'image' | 'video') {
    showError.value = false

    const remaining = MAX_FILES - totalCount.value
    if (remaining <= 0) return

    const files = Array.from(fileList).slice(0, remaining)
    const newItems: FileItem[] = files.map((file) => ({
        id: generateId(),
        file,
        fileType: type,
        previewUrl: URL.createObjectURL(file),
        status: 'pending' as const,
        progress: 0,
    }))

    if (type === 'image') {
        imageFiles.value.push(...newItems)
    } else {
        videoFiles.value.push(...newItems)
    }

    // Step 1-5: Upload từng file theo luồng presigned URL
    newItems.forEach((item) => {
        uploadFile(item)
    })
}

async function uploadFile(item: FileItem) {
    try {
        // Step 1: Lấy presigned URL
        item.status = 'gettingPresign'
        const presignRes = await $api.post<PresignResponse>('/upload/presign', {
            filename: item.file.name,
            content_type: item.file.type,
        })

        if (!presignRes.success || !presignRes.data) {
            throw new Error(presignRes.message || 'Không thể lấy presigned URL')
        }

        const { upload_url, key, expires_at } = presignRes.data
        item.key = key
        item.uploadUrl = upload_url
        item.expiresAt = new Date(expires_at).getTime()

        // Step 4: Upload file trực tiếp lên R2 (không qua backend)
        item.status = 'uploading'
        await uploadToR2(item.file, upload_url, (progress) => {
            item.progress = progress
        })

        // Step 5: Confirm upload → backend kiểm tra file trên R2 rồi lưu DB
        item.status = 'confirming'
        const confirmRes = await $api.post<ConfirmResponse>('/upload/confirm', { key })

        if (!confirmRes.success || !confirmRes.data) {
            throw new Error(confirmRes.message || 'Xác nhận upload thất bại')
        }

        item.imageId = confirmRes.data.image_id
        item.publicUrl = confirmRes.data.public_url
        item.status = 'done'
        item.progress = 100
    } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : 'Upload thất bại'
        item.status = 'error'
        item.errorMessage = msg
        showToast('error', msg)
        console.error('Upload error:', err)
    }
}

async function uploadToR2(file: File, url: string, onProgress: (pct: number) => void): Promise<void> {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest()
        xhr.open('PUT', url, true)
        xhr.setRequestHeader('Content-Type', file.type)

        xhr.upload.onprogress = (e) => {
            if (e.lengthComputable) {
                onProgress(Math.round((e.loaded / e.total) * 100))
            }
        }

        xhr.onload = () => {
            if (xhr.status >= 200 && xhr.status < 300) {
                onProgress(100)
                resolve()
            } else {
                reject(new Error(`Upload R2 thất bại: ${xhr.status}`))
            }
        }

        xhr.onerror = () => reject(new Error('Lỗi mạng khi upload lên R2'))
        xhr.ontimeout = () => reject(new Error('Upload lên R2 đã hết thời gian'))

        xhr.send(file)
    })
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
        addFiles(target.files, 'image')
    }
    target.value = ''
}

function onVideoInputChange(e: Event) {
    const target = e.target as HTMLInputElement
    if (target.files && target.files.length > 0) {
        addFiles(target.files, 'video')
    }
    target.value = ''
}

function removeFile(id: string) {
    imageFiles.value = imageFiles.value.filter((f) => f.id !== id)
    videoFiles.value = videoFiles.value.filter((f) => f.id !== id)
}

function validateImageCount(): boolean {
    if (imageFiles.value.filter((f) => f.status === 'done').length < 3) {
        showError.value = true
        return false
    }
    showError.value = false
    return true
}

defineExpose({ validateImageCount, uploadedImages, allUploaded })

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
