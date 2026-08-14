<template>
    <n-card title="Địa chỉ" size="small" :segmented="{ content: true }">
        <template #header-extra>
            <div class="cursor-pointer" @click="collapsed = !collapsed">
                <n-icon class="text-gray-500 transition-transform duration-200" :class="{ 'rotate-180': !collapsed }">
                    <IconChevronDownOutline />
                </n-icon>
            </div>
        </template>
        <div v-show="!collapsed">
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
        </div>
    </n-card>
    <n-modal v-model:show="showLocationModal" title="Nhập địa chỉ" style="width: 600px; max-height: 800px;">
        <n-card title="Nhập địa chỉ">
            <n-form-item label="Tỉnh / Thành phố" path="province" :feedback="postStore.errorsAddress.province"
                :validation-status="postStore.errorsAddress.province ? 'error' : undefined">
                <n-select v-model:value="postStore.form.province" placeholder="Chọn tỉnh/thành phố" clearable filterable
                    :options="provinceOption" @update:value="onWardChange" />
            </n-form-item>

            <n-form-item label="Phường / Xã" path="ward" :feedback="postStore.errorsAddress.ward"
                :validation-status="postStore.errorsAddress.ward ? 'error' : undefined">
                <n-select v-model:value="postStore.form.ward" placeholder="Chọn phường/xã" clearable filterable
                    :options="wardOptions" :loading="loadingWard" @update:value="clearError('ward')"
                    :disabled="isDisabledWard" />
            </n-form-item>

            <n-form-item label="Địa chỉ chi tiết" path="detail" :feedback="postStore.errorsAddress.detail_address"
                :validation-status="postStore.errorsAddress.detail_address ? 'error' : undefined">
                <n-input v-model:value="postStore.form.detail_address" placeholder="Nhập số nhà, khu phố, ngõ hẻm..."
                    clearable @update:value="clearError('detail_address')" />
            </n-form-item>

            <n-form-item label="Dự án" path="project_id">
                <n-select v-model:value="postStore.form.project_id" placeholder="Chọn dự án" clearable filterable
                    :options="projectOptions" :disabled="isDisabledPJ"/>
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
import type { CityOption, WardOption, ProjectOption } from '~/types/real_estate'
import { useCreatePost } from '~/stores/create-post'

const { $api } = useNuxtApp()
const postStore = useCreatePost()
const collapsed = ref(false)
const showLocationModal = ref(false)
const isDisabledWard = computed(() => !postStore.form.province)
const isDisabledPJ = computed(() => !postStore.form.province || !postStore.form.ward)


const openModal = () => {
    showLocationModal.value = true
}

// Áp dụng địa chỉ: nếu thiếu field bắt buộc thì hiện lỗi ngay trong modal, không tắt
const applyAddress = () => {
    // Validate qua store để hiển thị lỗi dưới từng ô
    const valid = postStore.validateAddress()
    if (!valid) return // Giữ modal mở để người dùng sửa

    showLocationModal.value = false
}

const loadingWard = ref(false)
const provinceOption = ref<SelectOption[]>([])
const wardOptions = ref<SelectOption[]>([])
const projectOptions = ref<SelectOption[]>([])

// Có lỗi địa chỉ hay không (để hiện status error dưới ô Khu vực)
const hasAddressError = computed(() =>
    Object.values(postStore.errorsAddress).some((message) => message !== '')
)

// Vô hiệu hóa ô Chọn phường/xã khi chưa chọn tỉnh/thành phố

// Xóa lỗi của 1 field địa chỉ khi người dùng sửa
const clearError = (field: keyof typeof postStore.errorsAddress) => {
    postStore.errorsAddress[field] = ''
}

const locationLabel = computed(() => {
    const city = provinceOption.value.find(item => item.value === postStore.form.province)
    const ward = wardOptions.value.find(item => item.value === postStore.form.ward)
    const project = projectOptions.value.find(item => item.value === postStore.form.project_id)

    const parts: string[] = [];
    if (city) parts.push(city.label as string);
    if (ward) parts.push(ward.label as string);
    if (project) parts.push(project.label as string);
    if (postStore.form.detail_address) parts.push(postStore.form.detail_address);
    if (parts.length === 0) return 'Chọn khu vực';
    return parts.join(', ');
});

const fetchListProvice = async () => {
    try {
        const res = await $api.get("/real-estate/list/city") as { data: CityOption[] }
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
    postStore.form.ward = null
    wardOptions.value = []
    clearError('province')

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



const fetchListProject = async () => {
    try {
        const params: Record<string, string> = {}
        if (postStore.form.province) {
            params.province = postStore.form.province
        }
        if (postStore.form.ward) {
            params.ward = postStore.form.ward
        }

        const res = await $api.get<{ data: ProjectOption[] }>("/real-estate/list/project", {
            params
        })
        projectOptions.value = res.data.map((item: ProjectOption) => ({
            label: item.name,
            value: item.id
        }))
    } catch (error) {
        console.error("Lỗi khi tải danh sách dự án:", error)
        projectOptions.value = []
    }
}

watch(
    () => [postStore.form.province, postStore.form.ward],
    async ([newProvince, newWard], [oldProvince, oldWard]) => {
        // Reset dự án đã chọn nếu tỉnh hoặc xã thay đổi
        if (newProvince !== oldProvince || newWard !== oldWard) {
            postStore.form.project_id = null
        }
        await fetchListProject()
    }
)

onMounted(() => {
    fetchListProvice()
    // fetchListProject()
})
</script>
