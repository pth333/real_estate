<template>
    <n-card title="Địa chỉ" size="small" :segmented="{ content: true }">
        <template #header-extra>
            <div class="cursor-pointer" @click="collapsed = !collapsed">
                <n-icon class="text-gray-500 transition-transform duration-200" :class="{ 'rotate-180': !collapsed }">
                    <IconChevronDownOutline />
                </n-icon>
            </div>
        </template>
        <div v-show="!collapsed" class="flex flex-col gap-4">
            <n-form-item label="Khu vực"
                :feedback="postStore.errorsAddress.province || postStore.errorsAddress.detail_address"
                :validation-status="hasAddressError ? 'error' : undefined">
                <n-input readonly :value="locationLabel" placeholder="Chọn khu vực..." class="cursor-pointer"
                    @click="openModal">
                    <template #suffix>
                        <IconChevronRight class="h-4 w-4" />
                    </template>
                </n-input>
            </n-form-item>

            <!-- Bản đồ nhỏ hiển thị sau khi đã chọn toạ độ -->
            <div v-if="staticMapSrc" class="rounded-lg overflow-hidden border border-gray-200">
                <img :src="staticMapSrc" alt="Vị trí trên bản đồ" class="h-44 w-full object-cover" />
            </div>
        </div>
    </n-card>

    <!-- Modal chọn vị trí — logic + options nằm trong AddressModal -->
    <AddressModal v-model:show="showLocationModal" v-model:location-label="locationLabel" />
</template>
<script setup lang="ts">
import { useCreatePost } from '~/stores/create-post'
import { useGeoapify } from '~/composables/useGeoapify'

const postStore = useCreatePost()
const { staticMapUrl } = useGeoapify()

const collapsed = ref(false)
const showLocationModal = ref(false)
const locationLabel = ref('')

// URL ảnh static map nhỏ dưới ô "Chọn khu vực" (tự cập nhật theo toạ độ)
const staticMapSrc = computed(() => {
    if (!postStore.form.latitude || !postStore.form.longitude) return ''
    return staticMapUrl(postStore.form.latitude, postStore.form.longitude, {
        width: 600,
        height: 400,
        zoom: 16,
    })
})

const hasAddressError = computed(() =>
    Object.values(postStore.errorsAddress).some((message) => message !== '')
)

const openModal = () => {
    showLocationModal.value = true
}
</script>
