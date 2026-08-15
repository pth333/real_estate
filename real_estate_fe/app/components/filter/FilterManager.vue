<template>
    <div>
        <!-- Nút mở filter -->
        <div class="flex items-center gap-2">
            <n-button type="primary" ghost size="small" @click="openModal">
                <template #icon>
                    <IconFilter class="h-4 w-4" />
                </template>
                Bộ lọc
            </n-button>

            <FilterCategory />

            <!-- Filter price range -->
            <FilterPriceRange />

            <!-- Filter diện tích -->
            <FilterAreaRange />
        </div>


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
