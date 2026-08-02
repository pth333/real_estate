<template>
    <n-card title="Thông tin chính" size="small" :segmented="{ content: true }">
        <template #header-extra>
            <div class="cursor-pointer" @click="collapsed = !collapsed">
                <n-icon class="text-gray-500 transition-transform duration-200" :class="{ 'rotate-180': !collapsed }">
                    <IconChevronDownOutline />
                </n-icon>
            </div>
        </template>

        <div v-show="!collapsed">
            <div class="mb-4">
                <label class="text-sm text-gray-700 block mb-1.5">
                    Loại bất động sản <span class="text-red-500">*</span>
                </label>
                <n-select v-model:value="postStore.form.real_estate_type" placeholder="Chọn loại BĐS"
                    :options="realEstateTypes" :status="postStore.errorsMainInfo.real_estate_type ? 'error' : 'default'"
                    @update:value="clearError('real_estate_type')" />
                <span v-if="postStore.errorsMainInfo.real_estate_type" class="text-xs text-red-500 mt-1 block">{{
                    postStore.errorsMainInfo.real_estate_type }}</span>
            </div>

            <div class="mb-4">
                <label class="text-sm text-gray-700 block mb-1.5">
                    Diện tích <span class="text-red-500">*</span>
                </label>
                <n-input-number v-model:value="postStore.form.area" placeholder="Nhập diện tích" :min="0" class="w-full"
                    :status="postStore.errorsMainInfo.area ? 'error' : 'default'" @update:value="clearError('area')">
                    <template #suffix>m²</template>
                </n-input-number>
                <span v-if="postStore.errorsMainInfo.area" class="text-xs text-red-500 mt-1 block">{{
                    postStore.errorsMainInfo.area
                }}</span>
            </div>
            <div class="flex gap-2 items-start">
                <!-- Mức giá -->
                <div class="flex flex-col flex-1">
                    <label class="text-sm text-gray-700 mb-1.5">
                        Mức giá <span class="text-red-500">*</span>
                    </label>
                    <n-input :value="price.displayed" placeholder="Nhập mức giá" class="w-full"
                        :status="postStore.errorsMainInfo.price ? 'error' : 'default'"
                        @update:value="(v: any) => { price.onInput(v); clearError('price') }" />
                    <span v-if="postStore.errorsMainInfo.price" class="text-xs text-red-500 mt-1 block">{{
                        postStore.errorsMainInfo.price
                    }}</span>
                    <span v-if="postStore.form.price && postStore.form.area" class="text-xs mt-1">
                        {{ formatTotalValue(postStore.form.price, postStore.form.area, postStore.form.unit) }}
                    </span>
                </div>

                <!-- Đơn vị -->
                <div class="flex flex-col w-36">
                    <label class="text-sm text-gray-700 mb-1.5">
                        Đơn vị <span class="text-red-500">*</span>
                    </label>
                    <n-select v-model:value="postStore.form.unit" placeholder="Chọn đơn vị" :options="unitOptions"
                        :status="postStore.errorsMainInfo.unit ? 'error' : 'default'"
                        @update:value="clearError('unit')" />
                    <span v-if="postStore.errorsMainInfo.unit" class="text-xs text-red-500 mt-1 block">{{
                        postStore.errorsMainInfo.unit
                    }}</span>
                </div>

            </div>
        </div>
    </n-card>
</template>
<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { useCreatePost } from '~/stores/create-post'
import type { OptionTypeRealestate } from '~/types/real_estate'
import { useNumberInput } from '~/composables/useNumberInput'

const postStore = useCreatePost()
const { $api } = useNuxtApp()
const collapsed = ref(false)
const realEstateTypes = ref<SelectOption[]>([])

// Ô nhập giá — tự format dấu chấm hàng nghìn, ghi số thực vào store
const price = useNumberInput(toRef(postStore.form, 'price'))

const unitOptions = ref<SelectOption[]>([
    { label: 'VNĐ', value: 'vnd' },
    { label: 'USD', value: 'usd' },
    { label: 'EUR', value: 'eur' },
])

// Định dạng số tiền theo đơn vị (triệu VNĐ, tỷ VNĐ, USD, EUR...)
const formatCurrency = (value: number, unit: string) => {
    const currencyMap: Record<string, string> = {
        vnd: 'VND',
        usd: 'USD',
        eur: 'EUR',
    }
    const currency = currencyMap[unit] ?? 'VND'
    return new Intl.NumberFormat('vi-VN', { style: 'currency', currency, maximumFractionDigits: 2 }).format(value)
}

// Hiển thị tổng giá trị và giá trên m² theo format:
// Tổng trị giá <x> (~<số tiền>/m²)
// price: giá thuần 1 đơn vị (VNĐ hoặc USD/EUR)
const formatTotalValue = (price: number, area: number, unit: string | null | undefined) => {
    const { formatPrice, formatPricePerM2 } = usePriceFormatter()

    if (unit === 'vnd') {
        const total = price * area
        return `Tổng trị giá ${formatPrice(total)} (~${formatPricePerM2(price)})`
    }

    const currencyUnit = unit ?? 'vnd'
    const total = price * area
    return `Tổng trị giá ${formatCurrency(total, currencyUnit)} (~${formatCurrency(price, currencyUnit)})`
}

// Xóa lỗi của 1 field khi người dùng bắt đầu nhập/sửa lại
const clearError = (field: keyof typeof postStore.errorsMainInfo) => {
    postStore.errorsMainInfo[field] = ''
}

// Validate phần Thông tin chính, hiển thị lỗi dưới từng ô input
const validate = (): boolean => {
    postStore.errorsMainInfo = {
        real_estate_type: postStore.form.real_estate_type
            ? ""
            : "Vui lòng chọn loại bất động sản",
        area: postStore.form.area ? "" : "Vui lòng nhập diện tích",
        price: postStore.form.price ? "" : "Vui lòng nhập mức giá",
        unit: postStore.form.unit ? "" : "Vui lòng chọn đơn vị giá",
    }
    return Object.values(postStore.errorsMainInfo).every((msg) => msg === "")
}

defineExpose({ validate })

const fetchRealEstateTypes = async () => {
    try {

        const response = await $api.get<{ data: OptionTypeRealestate[] }>('/real-estate/list/types')
        realEstateTypes.value = response.data.map((item) => {
            const normalizedName = item.name.replace('Bán', '').trim()
            const capitalizedName = normalizedName.charAt(0).toUpperCase() + normalizedName.slice(1)

            return {
                label: capitalizedName,
                value: item.id,
            }
        })
    } catch (error) {
        console.error('Error fetching real estate types:', error)
    }
}


onMounted(() => {
    fetchRealEstateTypes()
})
</script>