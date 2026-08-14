<template>
  <div class="mx-auto max-w-[1200px] px-6 py-8">
    <!-- Header -->
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Danh sách Dự án</h1>
        <p class="text-sm text-gray-500 mt-1">Các dự án bất động sản nổi bật thuộc danh mục này</p>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="i in 6" :key="i" class="border border-gray-100 rounded-lg p-4 animate-pulse">
        <div class="aspect-[16/10] bg-gray-200 rounded-md mb-4"></div>
        <div class="h-5 bg-gray-200 rounded w-2/3 mb-2.5"></div>
        <div class="h-4 bg-gray-200 rounded w-full mb-2"></div>
        <div class="h-4 bg-gray-200 rounded w-1/2"></div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="projects.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
      <n-empty description="Không tìm thấy dự án nào thuộc danh mục này">
        <template #extra>
          <n-button type="primary" @click="navigateTo('/')">Quay lại trang chủ</n-button>
        </template>
      </n-empty>
    </div>

    <!-- Project Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="project in projects" :key="project.id"
        class="bg-white border border-gray-100 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-300 overflow-hidden cursor-pointer flex flex-col group"
        @click="goToProject(project)">
        <!-- Thumbnail -->
        <div class="relative aspect-[16/10] overflow-hidden bg-gray-50">
          <img src="https://placehold.co/600x400/e2e8f0/94a3b8?text=Project" :alt="project.name"
            class="w-full h-full object-cover group-hover:scale-102 transition-transform duration-300" />
          <!-- Status Badge -->
          <div class="absolute top-3 left-3">
            <span class="text-xs font-semibold px-2.5 py-1 rounded bg-white shadow-sm border"
              :class="statusClass(project.status)">
              {{ formatStatus(project.status) }}
            </span>
          </div>
        </div>

        <!-- Info content -->
        <div class="p-4 flex flex-col flex-grow gap-2">
          <!-- Name -->
          <h2 class="font-bold text-gray-900 text-base line-clamp-1 group-hover:text-emerald-600 transition-colors">
            {{ project.name }}
          </h2>

          <!-- Location -->
          <div class="flex items-center gap-1.5 text-sm text-gray-500">
            <span class="truncate">{{ project.full_address || 'Địa chỉ đang cập nhật' }}</span>
          </div>

          <!-- Project Specs -->
          <div class="flex items-center gap-4 text-xs text-gray-400 mt-1 pt-2 border-t border-gray-50">
            <span v-if="project.total_area_ha">Quy mô: <strong>{{ project.total_area_ha }} ha</strong></span>
            <span v-if="project.total_units">Số căn: <strong>{{ project.total_units }}</strong></span>
          </div>

          <!-- Price range -->
          <div class="mt-auto pt-2 flex items-center justify-between">
            <span class="text-xs text-gray-400">Mức giá khoảng</span>
            <span class="text-sm font-bold text-emerald-600">
              {{ formatPriceRange(project) }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
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

const statusClass = (status?: string) => {
  const formatted = formatStatus(status)
  if (formatted === 'Đang mở bán') {
    return 'border-green-100 text-green-600 bg-green-50'
  }
  if (formatted === 'Sắp mở bán') {
    return 'border-orange-100 text-orange-600 bg-orange-50'
  }
  return 'border-gray-100 text-gray-600 bg-gray-50'
}

const props = defineProps<{
  categorySlug: string
}>()

const { $api } = useNuxtApp()
const projects = ref<Project[]>([])
const loading = ref(false)

const fetchProjects = async () => {
  if (!props.categorySlug) return
  loading.value = true
  try {
    const res = await $api.get<{ data: Project[] }>(`/real-estate/project-category/${props.categorySlug}`)
    projects.value = res.data || []
  } catch (error) {
    console.error("Lỗi khi tải danh sách dự án:", error)
    projects.value = []
  } finally {
    loading.value = false
  }
}

const goToProject = (project: Project) => {
  const slug = project.slug
  navigateTo(`/${slug}-pj${project.id}`)

  // Trigger tăng lượt xem dự án ở backend
  $api.post(`/real-estate/project/view/${project.id}`).catch(err => {
    console.error("Lỗi khi tăng lượt xem dự án:", err)
  })
}

const formatPriceRange = (project: Project) => {
  if (project.price_min && project.price_max) {
    return `${project.price_min} - ${project.price_max} triệu/m²`
  } else if (project.price_min) {
    return `Từ ${project.price_min} triệu/m²`
  }
  return 'Đang cập nhật'
}

onMounted(() => {
  fetchProjects()
})

watch(() => props.categorySlug, () => {
  fetchProjects()
})
</script>
