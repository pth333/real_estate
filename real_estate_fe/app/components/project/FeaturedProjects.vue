<template>
    <section class="py-8">
        <div class="container mx-auto px-24">
            <!-- Header -->
            <div class="flex items-center justify-between mb-5">
                <h2 class="text-xl font-bold text-gray-900">Dự án bất động sản nổi bật</h2>
                <NuxtLink :to="thirdCategorySlug ? `/${thirdCategorySlug}` : '#'"
                    class="flex items-center gap-1 text-emerald-600 text-sm font-medium hover:underline">
                    Xem thêm
                    <IconArrowRight class="h-4 w-4" />
                </NuxtLink>
            </div>

            <!-- Loading State -->
            <div v-if="loading" class="grid grid-cols-4 gap-4">
                <SkeletonCard :count="4" type="project" />
            </div>

            <!-- Slider wrapper -->
            <div v-else-if="projects.length > 0" class="relative">
                <!-- Nút prev -->
                <button
                    class="absolute -left-5 top-1/2 -translate-y-1/2 z-10 w-9 h-9 bg-white border border-gray-200 shadow rounded-md flex items-center justify-center hover:shadow-md transition-shadow"
                    @click="prev" :disabled="currentIndex === 0"
                    :class="{ 'opacity-50 cursor-not-allowed': currentIndex === 0 }">
                    <IconChevronLeft class="h-4 w-4 text-gray-600" />
                </button>

                <!-- Cards -->
                <div class="grid grid-cols-4 gap-4 overflow-hidden">

                    <NuxtLink v-for="item in visibleItems" :key="item.id" :to="`${item.slug}`"
                        class="bg-white border border-gray-100 rounded-lg shadow-sm hover:shadow-md overflow-hidden cursor-pointer group flex flex-col no-underline">
                        <!-- Ảnh -->
                        <div class="relative aspect-4/3 overflow-hidden bg-gray-100 rounded-t-lg">
                            <img :src="item.thumbnail" :alt="item.name"
                                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                            
                        </div>

                        <!-- Nội dung -->
                        <div class="p-3 flex flex-col gap-1.5 grow">
                            <div class="flex items-center flex-wrap">
                                <n-tag :type="statusTagType(item.status)" size="small" 
                                    class="shadow-sm font-semibold">
                                    {{ formatStatus(item.status) }}
                                </n-tag>
                            </div>
                            <div
                                class="font-semibold text-gray-900 text-sm truncate group-hover:text-emerald-600 transition-colors">
                                {{ item.name }}
                            </div>
                            <div class="text-xs text-gray-500">
                                Quy mô: {{ item.total_area_ha ? item.total_area_ha + ' ha' : 'Đang cập nhật' }}
                            </div>
                            <div class="text-xs text-gray-400 truncate">
                                {{ item.full_address || 'Địa chỉ đang cập nhật' }}
                            </div>
                        </div>
                    </NuxtLink>
                </div>

                <!-- Nút next -->
                <button
                    class="absolute -right-5 top-1/2 -translate-y-1/2 z-10 w-9 h-9 bg-white border border-gray-200 shadow rounded-md flex items-center justify-center hover:shadow-md transition-shadow"
                    @click="next" :disabled="currentIndex + pageSize >= projects.length"
                    :class="{ 'opacity-50 cursor-not-allowed': currentIndex + pageSize >= projects.length }">
                    <IconChevronRight class="h-4 w-4 text-gray-600" />
                </button>
            </div>

            <!-- Empty State -->
            <div v-else class="text-center py-8 text-gray-500 text-sm">
                Chưa có dự án nổi bật nào được ghi nhận.
            </div>
        </div>
    </section>
</template>

<script setup lang="ts">
import { useMenuStore } from '~/stores/menu';
import type { ProjectSummary } from '~/types/project';
import { useProjectService } from '~/services/project.service';

const menuStore = useMenuStore();

const thirdCategorySlug = computed(() => {
    return menuStore.menu?.categories?.[2]?.Slug;
});
// const formatStatus = (status?: string | boolean): string => {
//     if (!status) return 'Chưa cập nhật'
//     const s = String(status).toLowerCase().trim()
//     if (s === 'active') {
//         return 'Đang mở bán'
//     }
//     if (s === 'inactive') {
//         return 'Sắp mở bán'
//     }
//     return status as string
// }

function statusClass(status?: string) {
    const formatted = formatStatus(status)
    if (formatted === 'Đang mở bán') {
        return 'border-green-400 text-green-600 bg-green-50'
    }
    if (formatted === 'Sắp mở bán') {
        return 'border-red-300 text-red-500 bg-red-50'
    }
    return 'border-gray-300 text-gray-500 bg-gray-50'
}

const projectService = useProjectService()
const projects = ref<ProjectSummary[]>([])
const loading = ref(true)

const fetchFeaturedProjects = async () => {
    loading.value = true
    try {
        const res = await projectService.getFeatured(12)
        projects.value = (res || []).map((p, index) => ({
            ...p,
            // Ưu tiên ảnh từ API, fallback placeholder
            thumbnail: p.thumbnail
        }))
    } catch (error) {
        console.error("Lỗi khi tải danh sách dự án nổi bật:", error)
        projects.value = []
    } finally {
        loading.value = false
    }
}

const currentIndex = ref(0)
const pageSize = 4

const visibleItems = computed(() => projects.value.slice(currentIndex.value, currentIndex.value + pageSize))

function prev() {
    if (currentIndex.value > 0) currentIndex.value -= pageSize
}

function next() {
    if (currentIndex.value + pageSize < projects.value.length) currentIndex.value += pageSize
}

const goToProject = (project: ProjectSummary) => {
    // Sinh SEO URL driven chuyển hướng chi tiết dự án có định dạng dạng `/slug-du-an-pj{id}`
    const slug = project.slug
    navigateTo(`/${slug}-pj${project.id}`)
}

onMounted(() => {
    fetchFeaturedProjects()
})
</script>
