<template>
    <div>
        <!-- Nút mở filter -->
        <div class="flex items-center gap-2">
            <n-button type="primary" ghost size="small" @click="openModal">
                <template #icon>
                    <IconFilter class="h-4 w-4" />
                </template>
                Bộ lọc
                <!-- <span v-if="activeFilterCount > 0"
                class="ml-1.5 flex h-5 w-5 items-center justify-center rounded-full bg-blue-500 text-[11px] text-white">
                {{ activeFilterCount }}
            </span> -->
            </n-button>

            <!-- Filter price range -->
            <FilterPriceRange />

            <!-- Filter diện tích -->
            <FilterAreaRange />
        </div>
        <!-- Badge filter đang active -->
        <!-- <div v-if="hasActiveFilters" class="mt-2 flex flex-wrap items-center gap-2">
            <span class="text-xs text-gray-400">Đang lọc:</span>
            <n-tag v-if="filterStore.filters.min_price != null" closable size="small"
                @close="removeFilter('min_price')">
                Từ {{ formatFilterPrice(filterStore.filters.min_price) }}
            </n-tag>
            <n-tag v-if="filterStore.filters.max_price != null" closable size="small"
                @close="removeFilter('max_price')">
                Đến {{ formatFilterPrice(filterStore.filters.max_price) }}
            </n-tag>
            <n-tag v-if="filterStore.filters.district" closable size="small" @close="removeFilter('district')">
                {{ filterStore.filters.district }}
            </n-tag>
            <n-button size="tiny" text type="warning" @click="resetFilters">
                Xoá tất cả
            </n-button>
        </div> -->

        <!-- Modal filter -->
        <FilterAdvanceModal v-model:showModalAll="filterStore.showModalAll" />
    </div>
</template>

<script setup lang="ts">
import { useFilterStore } from '~/stores/filter';

const filterStore = useFilterStore();
const openModal = () => {
    filterStore.showModalAll = true
    filterStore.showPopover = false
}
</script>
