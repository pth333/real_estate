<template>
    <div class="flex flex-col gap-4">

        <!-- Input thủ công -->
        <div class="flex items-center gap-3">
            <n-input-number v-model:value="filterStore.filters.min_price" placeholder="Từ" :min="0" :show-button="false"
                class="flex-1">
                <template #suffix>đ</template>
            </n-input-number>
            <span class="text-gray-400">→</span>
            <n-input-number v-model:value="filterStore.filters.max_price" placeholder="Đến" :min="0"
                :show-button="false" class="flex-1">
                <template #suffix>đ</template>
            </n-input-number>
        </div>

        <!-- Radio ranges -->
        <div class="flex flex-col">
            <label v-for="opt in priceRanges" :key="opt.label"
                class="flex items-center justify-between py-3 border-b border-gray-100 cursor-pointer"
                @click="selectRange(opt)">
                <span class="text-sm">{{ opt.label }}</span>
                <n-radio :checked="isSelected(opt)" @update:checked="selectRange(opt)" />
            </label>
        </div>

        <!-- Nút Áp dụng -->
        <div class="mt-2 flex justify-end">
            <n-button type="primary" @click="onApply">
                Áp dụng
            </n-button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useFilterStore } from '~/stores/filter';

const filterStore = useFilterStore()

const TY = 1_000_000_000

const priceRanges = [
    { label: 'Tất cả khoảng giá', min: undefined, max: undefined },
    { label: 'Dưới 1 tỷ', min: undefined, max: 1 * TY },
    { label: '1 - 3 tỷ', min: 1 * TY, max: 3 * TY },
    { label: '3 - 5 tỷ', min: 3 * TY, max: 5 * TY },
    { label: '5 - 10 tỷ', min: 5 * TY, max: 10 * TY },
    { label: 'Trên 10 tỷ', min: 10 * TY, max: undefined },
]

function isSelected(opt: typeof priceRanges[0]) {
    if (opt.min === undefined && opt.max === undefined) {
        return !filterStore.filters.min_price && !filterStore.filters.max_price
    }
    return filterStore.filters.min_price === opt.min && filterStore.filters.max_price === opt.max
}

function selectRange(opt: typeof priceRanges[0]) {
    filterStore.filters.min_price = opt.min
    filterStore.filters.max_price = opt.max
}

const onApply = () => {
    filterStore.screen = 'main'
}
</script>