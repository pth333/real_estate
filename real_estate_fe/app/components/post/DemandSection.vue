<template>
    <n-card title="Nhu cầu" size="small" :segmented="{ content: true }">
        <template #header-extra>
            <n-icon class="text-gray-500">
                <IconChevronDownOutline />
            </n-icon>
        </template>

        <div class="grid grid-cols-2 gap-3">
            <!-- Bán: active -->
            <div @click="postStore.listingType = 'sell'" :class="[
                'flex items-center gap-2.5 px-5 py-3.5 rounded-lg cursor-pointer text-sm border transition-colors',
                postStore.listingType === 'sell'
                    ? 'border-gray-900 text-gray-900 font-semibold'
                    : 'border-gray-300 text-gray-600 hover:border-red-500 hover:text-red-500'
            ]">
                <n-icon size="20">
                    <IconPricetagOutline />
                </n-icon>
                <span>Bán</span>
            </div>
            <!-- Cho thuê -->
            <div @click="postStore.listingType = 'rent'" :class="[
                'flex items-center gap-2.5 px-5 py-3.5 rounded-lg cursor-pointer text-sm border transition-colors',
                postStore.listingType === 'rent'
                    ? 'border-gray-900 text-gray-900 font-semibold'
                    : 'border-gray-300 text-gray-600 hover:border-red-500 hover:text-red-500'
            ]">
                <n-icon size="20">
                    <IconKeyOutline />
                </n-icon>
                <span>Cho thuê</span>
            </div>
        </div>
    </n-card>
</template>
<script setup lang="ts">
import { useCreatePost } from '~/stores/create-post'
const postStore = useCreatePost()


</script>