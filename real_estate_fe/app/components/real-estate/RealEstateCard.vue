<template>
    <div class="grid grid-cols-1 gap-4">
        <div v-for="estate in realEstates" :key="estate.id"
            class="group relative overflow-hidden rounded-lg border border-gray-100 bg-white shadow-sm transition-shadow duration-300 hover:shadow-md">
            <div v-if="estate.badge"
                class="absolute left-3 top-3 z-10 rounded bg-red-600 px-3 py-1 text-xs font-bold uppercase text-white">
                {{ estate.badge }}
            </div>

            <div class="grid h-80 grid-cols-[2fr_1fr] gap-0.5 overflow-hidden rounded-t-lg bg-gray-100">
                <div class="overflow-hidden">
                    <img :src="mainImage(estate)" :alt="estate.title" class="h-full w-full object-cover"
                        @error="handleImageError" />
                </div>
                <div class="flex h-full min-h-0 flex-col gap-0.5">
                    <div v-if="thumbnails(estate)[0]" class="relative min-h-0 flex-1 overflow-hidden bg-gray-200">
                        <img :src="thumbnails(estate)[0]" :alt="`${estate.title} 2`"
                            class="absolute inset-0 h-full w-full object-cover" @error="handleImageError" />
                    </div>
                    <div v-if="thumbnails(estate)[1] || remainingImagesCount(estate) > 0"
                        class="flex min-h-0 flex-1 gap-0.5">
                        <div v-if="thumbnails(estate)[1]" class="relative min-h-0 flex-1 overflow-hidden bg-gray-200">
                            <img :src="thumbnails(estate)[1]" :alt="`${estate.title} 3`"
                                class="absolute inset-0 h-full w-full object-cover" @error="handleImageError" />
                        </div>
                        <div v-if="remainingImagesCount(estate) > 0"
                            class="relative min-h-0 flex-1 overflow-hidden bg-gray-200">
                            <img v-if="overlayImage(estate)" :src="overlayImage(estate)" :alt="`${estate.title} more`"
                                class="absolute inset-0 h-full w-full object-cover" @error="handleImageError" />
                            <div
                                class="absolute inset-0 flex items-center justify-center bg-black/60 text-xl font-bold text-white">
                                +{{ remainingImagesCount(estate) }}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="space-y-3 p-4">
                <h3 class="cursor-pointer text-base font-semibold uppercase leading-tight text-gray-800 line-clamp-2 hover:text-emerald-600"
                    @click="goToDetail(estate)">
                    {{ estate.title }}
                </h3>

                <div class="flex flex-wrap items-center gap-3 text-sm text-gray-500">
                    <span class="text-lg font-bold text-red-600">{{ formatPrice(estate.price_vnd) }}</span>
                    <span>{{ estate.acreage.toFixed(1) }} m²</span>
                    <span>{{ formatPricePerM2(estate.price_per_m2) }}</span>
                    <div class="flex gap-3">
                        <span v-if="estate.bedrooms" class="flex items-center gap-1">
                            <IconBed /> {{ estate.bedrooms }}
                        </span>
                        <span v-if="estate.bathrooms" class="flex items-center gap-1">
                            <IconBath /> {{ estate.bathrooms }}
                        </span>
                    </div>
                    <span class="text-gray-500">{{ fullLocation(estate) }}</span>
                </div>

                <p v-if="estate.description" class="line-clamp-2 text-sm leading-relaxed text-gray-500">
                    {{ truncatedDescription(estate) }}
                </p>

                <div class="flex items-center justify-between border-t border-gray-200 pt-3">
                    <div class="flex items-center gap-2">
                        <div
                            class="flex h-10 w-10 items-center justify-center rounded-full bg-emerald-500 text-lg font-bold text-white">
                            {{ agentInitial(estate) }}
                        </div>
                        <div>
                            <div class="text-sm font-semibold text-gray-800">{{ estate.agent_name || 'Người đăng' }}
                            </div>
                            <div class="text-xs text-gray-400">{{ postTime(estate) }}</div>
                        </div>
                    </div>

                    <div class="flex gap-2">
                        <button v-if="estate.agent_phone"
                            class="flex cursor-pointer items-center gap-1.5 rounded-md bg-emerald-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600"
                            @click="handleCall(estate)">
                            <IconPhone /> {{ formattedPhone(estate) }}
                        </button>
                        <button
                            class="flex h-8 w-8 items-center justify-center rounded-full border border-gray-200 transition-colors hover:border-red-400 hover:text-red-500"
                            :class="estate.is_favorite ? 'border-red-500 text-red-500' : ''"
                            @click.stop="handleToggleFavorite(estate)">
                            <IconHeart class="h-4 w-4"
                                :class="estate.is_favorite ? 'fill-red-500 text-red-500' : 'text-gray-400'" />
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { RealEstateResponse } from '~/types/real_estate'
import { buildDetailUrl } from '~/utils/slug'

const props = defineProps<{
    realEstates: RealEstateResponse[]
}>()

const emit = defineEmits<{
    call: [phone: string]
}>()

const favorite = useFavorite()
const DEFAULT_IMAGE = '/placeholder.jpg'

const images = (estate: RealEstateResponse) => estate.image_urls || []
const mainImage = (estate: RealEstateResponse) => images(estate)[0] || DEFAULT_IMAGE
const thumbnails = (estate: RealEstateResponse) => images(estate).slice(1, 3)
const overlayImage = (estate: RealEstateResponse) => images(estate)[3] || DEFAULT_IMAGE
const remainingImagesCount = (estate: RealEstateResponse) => Math.max(0, images(estate).length - 3)

const handleImageError = (event: Event) => {
    const image = event.target as HTMLImageElement
    image.src = DEFAULT_IMAGE
}

const fullLocation = (estate: RealEstateResponse) => [estate.district, estate.city].filter(Boolean).join(', ')
const truncatedDescription = (estate: RealEstateResponse) => {
    const description = estate.description || ''
    return description.length > 150 ? `${description.substring(0, 150)}...` : description
}
const agentInitial = (estate: RealEstateResponse) => (estate.agent_name || 'Q').charAt(0).toUpperCase()
const postTime = (estate: RealEstateResponse) => fromNow(estate.created_at)

const formattedPhone = (estate: RealEstateResponse) => {
    const phone = estate.agent_phone || ''
    return phone.length >= 10 ? `${phone.substring(0, 4)} ${phone.substring(4, 7)} *** - Hiện số` : phone
}

const handleCall = (estate: RealEstateResponse) => {
    if (estate.agent_phone) emit('call', estate.agent_phone)
}

const handleToggleFavorite = async (estate: RealEstateResponse) => {
    const next = await favorite.toggle(estate.id)
    if (next !== null) estate.is_favorite = next
}

const goToDetail = (estate: RealEstateResponse) => {
    // Dùng helper để luôn ra đường dẫn tuyệt đối /<slug>-rs<id> (xem utils/slug.ts)
    const slug = estate.slug || `-rs${estate.id}`
    navigateTo(buildDetailUrl(slug))
}
</script>
