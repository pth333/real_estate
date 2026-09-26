<!-- components/project/ProjectListings.vue -->
<template>
    <div class="mt-4 rounded-xl bg-white p-4 md:mt-6 md:p-5">
        <!-- Header -->
        <div class="mb-4 flex flex-wrap items-center justify-between gap-2">
            <div class="flex items-center gap-2">
                <h3 class="text-lg font-semibold text-gray-800">
                    Tin mua bán tại {{ project.name }}
                </h3>
            </div>

            <!-- Prev / Next -->
            <div class="flex items-center gap-2">
                <button
                    class="flex h-8 w-8 items-center justify-center rounded border border-gray-300 hover:border-gray-500 transition-colors disabled:opacity-30"
                    :disabled="currentIndex === 0" @click="prev">
                    <IconChevronLeft class="h-4 w-4" />
                </button>
                <button
                    class="flex h-8 w-8 items-center justify-center rounded border border-gray-300 hover:border-gray-500 transition-colors disabled:opacity-30"
                    :disabled="currentIndex + perPage >= items.length" @click="next">
                    <IconChevronRight class="h-4 w-4" />
                </button>
            </div>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
            <SkeletonCard :count="3" type="card" />
        </div>

        <!-- Empty -->
        <n-empty v-else-if="items.length === 0" description="Chưa có tin đăng nào cho dự án này" class="py-10" />

        <!-- Cards: 1 cột mobile, 2 cột tablet, 3 cột desktop -->
        <div v-else class="grid grid-cols-1 gap-4 overflow-hidden md:grid-cols-2 lg:grid-cols-3">
            <NuxtLink v-for="item in visibleItems" :key="item.id" :to="`${item.slug}`"
                class="cursor-pointer rounded-xl overflow-hidden border border-gray-100 shadow-sm hover:shadow-md transition-shadow duration-200 bg-white">

                <!-- Ảnh -->
                <div class="relative aspect-4/3 overflow-hidden bg-gray-100">
                    <img :src="item.images?.[0]?.url ?? '/placeholder.jpg'" :alt="item.title"
                        class="w-full h-full object-cover hover:scale-105 transition-transform duration-300" />
                </div>

                <!-- Nội dung -->
                <div class="p-3 flex flex-col gap-2">
                    <!-- Tiêu đề -->
                    <div class="text-sm font-semibold text-gray-800 leading-snug line-clamp-2">
                        {{ item.title }}
                    </div>

                    <!-- Giá + diện tích -->
                    <div class="flex items-center gap-2 text-sm">
                        <span class="text-red-500 font-bold">{{ formatPrice(item.price_vnd) }}</span>
                        <span class="text-gray-300">·</span>
                        <span class="text-gray-600">{{ item.acreage }} m²</span>
                    </div>

                    <!-- Địa chỉ -->
                    <div class="flex items-center gap-1 text-xs text-gray-500">
                        <IconMapPin class="h-3 w-3 shrink-0" />
                        <span class="line-clamp-1">{{ [item.district, item.city].filter(Boolean).join(', ') }}</span>
                    </div>

                    <!-- Footer -->
                    <div class="flex items-center justify-between pt-1 border-t border-gray-100">
                        <span class="text-xs text-gray-400">{{ formatDate(item.created_at) }}</span>
                        <button class="flex h-8 w-8 items-center justify-center rounded-full border transition-colors"
                            :class="item.is_favorite
                                ? 'border-red-500 text-red-500'
                                : 'border-gray-200 text-gray-400 hover:border-red-400 hover:text-red-500'"
                            @click.stop="toggleFavorite(item.id)">
                            <IconHeart class="h-4 w-4" :class="item.is_favorite ? 'fill-red-500 text-red-500' : ''" />
                        </button>
                    </div>
                </div>
            </NuxtLink>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { ProjectDetail } from '~/types/project';
import type { RealEstateResponse } from '~/types/real_estate';
import { useProjectService } from '~/services/project.service';

const props = defineProps<{
    project: ProjectDetail
}>()

const projectService = useProjectService()
const favorite = useFavorite()

const loading = ref(false)
const items = ref<RealEstateResponse[]>([])
const currentIndex = ref(0)
const perPage = 3

const visibleItems = computed(() =>
    items.value.slice(currentIndex.value, currentIndex.value + perPage)
)

const fetchListings = async () => {
    loading.value = true
    try {
        const res = await projectService.getListings(props.project.id)
        items.value = res || []
    } catch (err) {
        console.error('Lỗi khi tải tin dự án:', err)
    } finally {
        loading.value = false
    }
}

onMounted(() => fetchListings())

function prev() {
    currentIndex.value = Math.max(0, currentIndex.value - perPage)
}

function next() {
    currentIndex.value = Math.min(
        items.value.length - perPage,
        currentIndex.value + perPage
    )
}

function toggleFavorite(id: number) {
    const item = items.value.find(i => i.id === id)
    if (item) {
        favorite.toggle(id).then((next) => {
            if (next !== null) item.is_favorite = next
        })
    }
}
</script>