<template>
    <n-card title="Địa chỉ" size="small" :segmented="{ content: true }">
        <template #header-extra>
            <n-icon class="text-gray-500 rotate-180">
                <IconChevronDownOutline />
            </n-icon>
        </template>
        <n-form-item label="Khu vực">
            <n-input readonly :value="locationLabel" placeholder="Chọn khu vực..." class="cursor-pointer"
                @click="openModal">
                <template #suffix>
                    <IconChevronRight class="h-4 w-4" />
                </template>
            </n-input>
        </n-form-item>
    </n-card>
    <n-modal v-model:show="showLocationModal" title="Nhập địa chỉ" style="width: 600px; max-height: 800px;">
        <n-card title="Nhập địa chỉ">
            <n-form-item label="Tỉnh / Thành phố" path="province">
                <n-select v-model:value="postStore.province" placeholder="Chọn tỉnh/thành phố" clearable filterable
                    :options="provinceOption" @update:value="onWardChange" />
            </n-form-item>

            <n-form-item label="Phường / Xã" path="ward">
                <n-select v-model:value="postStore.ward" placeholder="Chọn phường/xã" clearable filterable
                    :options="wardOptions" :loading="loadingWard" />
            </n-form-item>

            <!-- <n-form-item label="Đường / Phố" path="street">
                <n-select placeholder="Chọn đường/phố" clearable filterable />
            </n-form-item> -->

            <n-form-item label="Địa chỉ chi tiết" path="detail">
                <n-input v-model:value="postStore.detail_address" placeholder="Nhập số nhà, khu phố, ngõ hẻm..."
                    clearable />
            </n-form-item>
            <template #action>
                <div class="flex justify-end">
                    <n-button type="primary" @click="applyAddress">Áp dụng</n-button>
                </div>
            </template>
        </n-card>

    </n-modal>
</template>
<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import type { CityOption, WardOption } from '~/types/real_estate'
import { useCreatePost } from '~/stores/create-post'

const { $api } = useNuxtApp()
const postStore = useCreatePost()
const showLocationModal = ref(false)
const openModal = () => {
    showLocationModal.value = true
}
const applyAddress = () => {
    showLocationModal.value = false
}

const loadingWard = ref(false)
const provinceOption = ref<SelectOption[]>([])
const wardOptions = ref<SelectOption[]>([])

const locationLabel = computed(() => {
    const city = provinceOption.value.find(item => item.value === postStore.province)
    const ward = wardOptions.value.find(item => item.value === postStore.ward)

    const parts: string[] = [];
    if (city) parts.push(city.label as string);
    if (ward) parts.push(ward.label as string);
    if (postStore.detail_address) parts.push(postStore.detail_address);
    if (parts.length === 0) return 'Chọn khu vực';
    return parts.join(', ');
});

const fetchListProvice = async () => {
    try {
        const res = await $api.get<{ data: CityOption[] }>("/real-estate/list/city")
        provinceOption.value = res.data.map((item: CityOption) => ({
            label: item.name,
            value: item.code
        }))
    } catch (error) {
        console.error("Lỗi khi tải danh sách tỉnh/thành phố:", error)
        provinceOption.value = []
    }
}

const onWardChange = async (provinceCode: string | null) => {
    postStore.ward = null
    wardOptions.value = []

    if (!provinceCode) return

    try {
        loadingWard.value = true
        const res = await $api.get<{ data: WardOption[] }>(`/real-estate/list/ward`, {
            params: { code: provinceCode }
        })
        wardOptions.value = res.data.map(item => ({
            label: item.name,
            value: item.code
        }))
    } catch (error) {
        console.error("Lỗi khi tải danh sách phường/xã:", error)
        wardOptions.value = []
    } finally {
        loadingWard.value = false
    }
}
onMounted(() => {
    fetchListProvice()
})
</script>
