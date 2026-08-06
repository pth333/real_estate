<template>
  <div class="group relative rounded-lg bg-white shadow-sm transition-shadow duration-300 hover:shadow-md">
    <!-- Badge VIP -->
    <div v-if="estate.badge"
      class="absolute left-3 top-3 z-10 rounded bg-red-600 px-3 py-1 text-xs font-bold uppercase text-white">
      {{ estate.badge }}
    </div>

    <!-- Grid ảnh: cột trái ảnh chính lớn, cột phải grid 2x2 -->
    <div class="grid h-80 grid-cols-[2fr_1fr] gap-0.5 bg-gray-100 overflow-hidden">
      <div class="overflow-hidden">
        <img :src="mainImage" :alt="estate.title" class="h-full w-full object-cover" @error="handleImageError" />
      </div>

      <div class="flex h-full min-h-0 flex-col gap-0.5">
        <!-- Ảnh 2: hàng trên (chiếm nửa chiều cao) -->
        <div v-if="thumbnails[0]" class="relative min-h-0 flex-1 overflow-hidden bg-gray-200">
          <img :src="thumbnails[0]" :alt="`${estate.title} 2`" class="absolute inset-0 h-full w-full object-cover"
            @error="handleImageError" />
        </div>
        <!-- Hàng dưới: ảnh 3 (trái) + +N (phải), chia đôi -->
        <div v-if="thumbnails[1] || remainingImagesCount > 0" class="flex min-h-0 flex-1 gap-0.5">
          <div v-if="thumbnails[1]" class="relative min-h-0 flex-1 overflow-hidden bg-gray-200">
            <img :src="thumbnails[1]" :alt="`${estate.title} 3`" class="absolute inset-0 h-full w-full object-cover"
              @error="handleImageError" />
          </div>
          <div v-if="remainingImagesCount > 0" class="relative min-h-0 flex-1 overflow-hidden bg-gray-200">
            <img v-if="overlayImage" :src="overlayImage" :alt="`${estate.title} more`"
              class="absolute inset-0 h-full w-full object-cover" @error="handleImageError" />
            <div class="absolute inset-0 flex items-center justify-center bg-black/60 text-xl font-bold text-white">
              +{{ remainingImagesCount }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Thông tin chính -->
    <div class="space-y-3 p-4">
      <h3 class="line-clamp-2 text-base font-semibold uppercase leading-tight text-gray-800">
        {{ estate.title }}
      </h3>

      <div class="flex flex-wrap items-center gap-3 text-sm text-gray-500">
        <span class="text-lg font-bold text-red-600">{{ formattedPrice }}</span>
        <span>{{ formattedArea }}</span>
        <span>{{ formattedPricePerM2 }}</span>

        <div class="flex gap-3">
          <span v-if="estate.bedrooms" class="flex items-center gap-1">
            <IconBed /> {{ estate.bedrooms }}
          </span>
          <span v-if="estate.bathrooms" class="flex items-center gap-1">
            <IconBath /> {{ estate.bathrooms }}
          </span>
        </div>

        <span class="text-gray-500">{{ fullLocation }}</span>
      </div>

      <p v-if="estate.description" class="line-clamp-2 text-sm leading-relaxed text-gray-500">
        {{ truncatedDescription }}
      </p>

      <!-- Footer -->
      <div class="flex items-center justify-between border-t border-gray-200 pt-3">
        <div class="flex items-center gap-2">
          <div class="flex h-10 w-10 items-center justify-center rounded-full bg-blue-500 text-lg font-bold text-white">
            {{ agentInitial }}
          </div>
          <div>
            <div class="text-sm font-semibold text-gray-800">{{ estate.agent_name || 'Người đăng' }}</div>
            <div class="text-xs text-gray-400">{{ postTime }}</div>
          </div>
        </div>

        <div class="flex gap-2">
          <button v-if="estate.agent_phone"
            class="flex cursor-pointer items-center gap-1.5 rounded bg-emerald-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600"
            @click="handleCall">
            <IconPhone /> {{ formattedPhone }}
          </button>
          <button
            class="flex h-10 w-10 cursor-pointer items-center justify-center rounded border-2 border-gray-300 bg-transparent transition hover:border-red-500"
            :class="{ 'border-red-500 bg-red-500 text-white': isFavorite }" @click="handleToggleFavorite">
            <IconHeart />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { RealEstateResponse } from '~/types/real_estate'
import { useImageGrid } from '~/composables/useImageGrid'


const props = defineProps({
  estate: {
    type: Object as PropType<RealEstateResponse>,
    required: true,
  }
})


const emit = defineEmits<{
  call: [phone: string]
  toggleFavorite: [id: number]
}>()

const { mainImage, thumbnails, overlayImage, remainingImagesCount, handleImageError } = useImageGrid(computed(() => props.estate.images))

const isFavorite = ref(props.estate.is_favorite || false)

const formattedPrice = computed(() => formatPrice(props.estate.price_vnd))

const formattedArea = computed(() => {
  return `${props.estate.acreage.toFixed(1)} m²`
})

const formattedPricePerM2 = computed(() => formatPricePerM2(props.estate.price_per_m2))

const fullLocation = computed(() => {
  const parts = [props.estate.district, props.estate.city].filter(Boolean)
  return parts.join(', ')
})

const truncatedDescription = computed(() => {
  const desc = props.estate.description || ''
  return desc.length > 150 ? desc.substring(0, 150) + '...' : desc
})

const agentInitial = computed(() => {
  const name = props.estate.agent_name || 'Q'
  return name.charAt(0).toUpperCase()
})

const postTime = computed(() => fromNow(props.estate.created_at))

const formattedPhone = computed(() => {
  const phone = props.estate.agent_phone || ''
  if (phone.length >= 10) {
    return `${phone.substring(0, 4)} ${phone.substring(4, 7)} *** - Hiện số`
  }
  return phone
})

const handleCall = () => {
  if (props.estate.agent_phone) {
    emit('call', props.estate.agent_phone)
  }
}

const handleToggleFavorite = () => {
  isFavorite.value = !isFavorite.value
  emit('toggleFavorite', props.estate.id)
}
</script>
