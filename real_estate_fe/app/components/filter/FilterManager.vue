<template>
    <div>
        <!-- Nút mở filter + các chip lọc nhanh.
             Mobile/tablet: cuộn ngang, tràn ra sát mép màn hình cho dễ vuốt.
             Desktop (lg): dàn hàng ngang và tự xuống dòng. -->
        <div
            class="-mx-4 flex items-center gap-2 overflow-x-auto px-4 pb-1 md:mx-0 md:flex-wrap md:overflow-visible md:px-0 md:pb-0 [scrollbar-width:none] [&>*]:shrink-0 [&::-webkit-scrollbar]:hidden">
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
