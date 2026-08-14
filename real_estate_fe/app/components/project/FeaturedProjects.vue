<template>
    <section class="py-8">
        <div class="container mx-auto px-24">
            <!-- Header -->
            <div class="flex items-center justify-between mb-5">
                <h2 class="text-xl font-bold text-gray-900">Dự án bất động sản nổi bật</h2>
                <a href="#" class="flex items-center gap-1 text-emerald-600 text-sm font-medium hover:underline">
                    Xem thêm
                    <IconArrowRight class="h-4 w-4" />
                </a>
            </div>

            <!-- Loading State -->
            <div v-if="loading" class="grid grid-cols-4 gap-4">
                <div v-for="i in 4" :key="i" class="border border-gray-100 rounded-lg p-4 animate-pulse">
                    <div class="aspect-[4/3] bg-gray-200 rounded-md mb-4"></div>
                    <div class="h-4 bg-gray-200 rounded w-2/3 mb-2"></div>
                    <div class="h-3 bg-gray-200 rounded w-1/2"></div>
                </div>
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
                    <div v-for="item in visibleItems" :key="item.id"
                        class="bg-white border border-gray-100 rounded-lg shadow-sm hover:shadow-md overflow-hidden cursor-pointer group flex flex-col" @click="goToProject(item)">
                        <!-- Ảnh -->
                        <div class="relative aspect-[4/3] overflow-hidden bg-gray-100 rounded-t-lg">
                            <img :src="item.thumbnail" :alt="item.name"
                                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                            <!-- Số ảnh -->
                            <div
                                class="absolute bottom-2 right-2 flex items-center gap-1 bg-black/50 text-white text-xs px-1.5 py-0.5 rounded">
                                <IconImage class="h-3 w-3" />
                                5
                            </div>
                        </div>

                        <!-- Nội dung -->
                        <div class="p-4 flex flex-col gap-1.5 flex-grow">
                            <!-- Badge trạng thái -->
                            <div class="flex items-center gap-2 flex-wrap">
                                <span class="text-[10px] font-semibold px-2 py-0.5 border rounded" :class="statusClass(item.status)">
                                    {{ formatStatus(item.status) }}
                                </span>
                            </div>

                            <!-- Tên dự án -->
                            <div
                                class="font-semibold text-gray-900 text-sm truncate group-hover:text-emerald-600 transition-colors">
                                {{ item.name }}
                            </div>

                            <!-- Quy mô diện tích -->
                            <div class="text-xs text-gray-500">
                                Quy mô: {{ item.total_area_ha ? item.total_area_ha + ' ha' : 'Đang cập nhật' }}
                            </div>

                            <!-- Địa chỉ -->
                            <div class="text-xs text-gray-400 truncate">{{ item.full_address || 'Địa chỉ đang cập nhật'
                                }}</div>
                        </div>
                    </div>
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
interface Project {
    id: number
    name: string
    slug: string
    status: string
    full_address: string
    total_area_ha?: number
    total_units?: number
    price_min?: number
    price_max?: number
    thumbnail?: string
}

const formatStatus = (status?: string | boolean): string => {
    if (!status) return 'Chưa cập nhật'
    const s = String(status).toLowerCase().trim()
    if (s === 'active' || s === 'true' || s === '1' || s === 'đang mở bán' || s === 'dang mo ban') {
        return 'Đang mở bán'
    }
    if (s === 'inactive' || s === 'false' || s === '0' || s === 'sắp mở bán' || s === 'sap mo ban') {
        return 'Sắp mở bán'
    }
    return status as string
}

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

const { $api } = useNuxtApp()
const projects = ref<Project[]>([])
const loading = ref(false)

const fetchFeaturedProjects = async () => {
    loading.value = true
    try {
        const res = await $api.get<{ data: Project[] }>("/real-estate/project/featured?limit=12")
        projects.value = (res.data || []).map((p, index) => ({
            ...p,
            thumbnail: `https://placehold.co/400x300/e2e8f0/94a3b8?text=Project+${index + 1}`
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

const goToProject = (project: Project) => {
    // Sinh SEO URL driven chuyển hướng chi tiết dự án có định dạng dạng `/slug-du-an-pj{id}`
    const slug = project.slug
    navigateTo(`/${slug}-pj${project.id}`)
}

onMounted(() => {
    fetchFeaturedProjects()
})
</script>
