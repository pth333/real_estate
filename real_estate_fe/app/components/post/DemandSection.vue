<template>
    <n-card title="Nhu cầu" size="small" :segmented="{ content: true }">
        <template #header-extra>
            <div class="cursor-pointer" @click="collapsed = !collapsed">
                <n-icon class="text-gray-500 transition-transform duration-200" :class="{ 'rotate-180': !collapsed }">
                    <IconChevronDownOutline />
                </n-icon>
            </div>
        </template>

        <div v-show="!collapsed">
            <div class="grid grid-cols-2 gap-3">
                <!-- Bán: active -->
                <div @click="postStore.form.listingType = 'sell'" :class="[
                    'flex items-center gap-2.5 px-5 py-3.5 rounded-lg cursor-pointer text-sm border transition-colors',
                    postStore.form.listingType === 'sell'
                        ? 'border-gray-900 text-gray-900 font-semibold'
                        : 'border-gray-300 text-gray-600 hover:border-red-500 hover:text-red-500'
                ]">
                    <n-icon size="20">
                        <IconPricetagOutline />
                    </n-icon>
                    <span>Bán</span>
                </div>
                <!-- Cho thuê -->
                <div @click="postStore.form.listingType = 'rent'" :class="[
                    'flex items-center gap-2.5 px-5 py-3.5 rounded-lg cursor-pointer text-sm border transition-colors',
                    postStore.form.listingType === 'rent'
                        ? 'border-gray-900 text-gray-900 font-semibold'
                        : 'border-gray-300 text-gray-600 hover:border-red-500 hover:text-red-500'
                ]">
                    <n-icon size="20">
                        <IconKeyOutline />
                    </n-icon>
                    <span>Cho thuê</span>
                </div>
            </div>
        </div>
    </n-card>
</template>
<script setup lang="ts">
import { useCreatePost } from '~/stores/create-post'
const postStore = useCreatePost()
const collapsed = ref(false)

</script>