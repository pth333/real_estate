<template>
    <div class="mb-8">
        <n-card :bordered="true" embedded size="small">
            <div class="flex flex-wrap items-end gap-4">
                <n-form-item label="Giá tối thiểu" :show-feedback="false" label-placement="top">
                    <n-input v-model:value="filterStore.filters.min_price" placeholder="VD: 1,000,000,000" :min="0"
                        clearable style="width: 180px" />
                </n-form-item>
                <n-form-item label="Giá tối đa" :show-feedback="false" label-placement="top">
                    <n-input v-model:value="filterStore.filters.max_price" placeholder="VD: 5,000,000,000" :min="0"
                        clearable style="width: 180px" />
                </n-form-item>
                <n-form-item label="Quận/Huyện" :show-feedback="false" label-placement="top">
                    <n-input v-model:value="filterStore.filters.district" placeholder="VD: Quận 1" clearable
                        style="width: 180px" />
                </n-form-item>
                <n-button tertiary type="warning" @click="resetFilters">
                    Reset
                </n-button>
                <n-button type="primary" @click="$emit('apply-filters')">
                    Áp dụng
                </n-button>
            </div>
        </n-card>
    </div>
</template>

<script setup lang="ts">
import { useFilterStore } from '~/stores/filter';
const filterStore = useFilterStore();
const resetFilters = () => {
    filterStore.filters = {
        min_price: undefined,
        max_price: undefined,
        district: undefined,
    };
};

</script>
