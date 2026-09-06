<template>
    <div v-if="store.listing?.description" class="bg-white rounded-lg py-4">
        <h2 class="mb-3 border-b border-gray-200 pb-2 text-base font-bold text-gray-800">
            Thông tin mô tả
        </h2>
        <p class="whitespace-pre-line text-sm leading-relaxed text-gray-700">{{ store.listing?.description }}</p>

        <!-- Liên hệ trên mobile -->
        <div v-if="store.listing?.agent_phone"
            class="mt-4 flex flex-wrap items-center gap-2 text-sm text-gray-600 sm:hidden">
            <span>Liên hệ:</span>
            <span class="font-medium text-gray-800">{{ store.maskedPhone }}</span>
            <n-button size="small" type="primary" @click="store.showPhone = true">Hiện số</n-button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useRealEstateDetail } from '~/stores/detail/real_estate_detail';

const store = useRealEstateDetail();
</script>
