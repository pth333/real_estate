<template>
  <n-layout class="bg-transparent">
    <n-layout-content class="mx-auto max-w-[1200px] px-6 py-6 bg-transparent">
      <!-- Breadcrumb Naive UI chuyên nghiệp -->
      <n-breadcrumb class="mb-4">
        <n-breadcrumb-item @click="navigateTo('/')">Trang chủ</n-breadcrumb-item>
        <n-breadcrumb-item>Dự án</n-breadcrumb-item>
        <n-breadcrumb-item>Dự án BĐS Toàn Quốc</n-breadcrumb-item>
      </n-breadcrumb>

      <!-- Tiêu đề + sắp xếp -->
      <n-space justify="space-between" align="end" class="mb-5">
        <n-space vertical :size="4">
          <n-h1 class="text-2xl! font-bold! m-0! !text-gray-900">Dự án toàn quốc</n-h1>
          <n-text depth="3" class="text-sm">
            Hiện đang có <strong class="text-gray-700">{{ projects.length.toLocaleString('vi-VN') }}</strong> dự án
          </n-text>
        </n-space>

        <!-- Sort selector Naive UI xịn sò -->
        <n-space align="center">
          <n-text depth="3" class="text-sm">Sắp xếp:</n-text>
          <n-select
            v-model:value="sortOrder"
            :options="sortOptions"
            class="w-40"
            size="small"
          />
        </n-space>
      </n-space>

      <!-- Grid Layout Naive UI cố định 4 cột đồng bộ -->
      <n-grid :cols="4" :x-gap="24" :y-gap="24">
        <!-- Cột trái (Chiếm 3/4 cột) -->
        <n-grid-item :span="3" class="min-w-0">
          <n-spin :show="loading">
            <!-- Empty State Naive UI -->
            <n-empty v-if="!loading && projects.length === 0" description="Không tìm thấy dự án nào thuộc danh mục này" class="py-20">
              <template #extra>
                <n-button type="primary" @click="navigateTo('/')">Quay lại trang chủ</n-button>
              </template>
            </n-empty>

            <!-- Danh sách dự án bằng n-space vertical và n-card -->
            <n-space v-else vertical :size="16">
              <n-card
                v-for="project in pagedProjects"
                :key="project.id"
                hoverable
                content-style="padding: 0;"
                class="overflow-hidden cursor-pointer group"
                @click="goToProject(project)"
              >
                <n-grid :cols="12" class="h-44">
                  <!-- Thumbnail dự án -->
                  <n-grid-item :span="4" class="relative overflow-hidden bg-gray-100">
                    <img
                      src="https://placehold.co/440x296/e2e8f0/94a3b8?text=Project"
                      :alt="project.name"
                      class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                    />
                    <!-- Trạng thái tag Naive UI -->
                    <div class="absolute top-2.5 left-2.5">
                      <n-tag
                        :type="statusTagType(project.status)"
                        size="small"
                        round
                        class="shadow-sm font-semibold"
                      >
                        {{ formatStatus(project.status) }}
                      </n-tag>
                    </div>
                  </n-grid-item>

                  <!-- Thông tin dự án -->
                  <n-grid-item :span="8" class="p-4 flex flex-col justify-between min-w-0">
                    <n-space vertical :size="6" class="min-w-0">
                      <!-- Tên dự án -->
                      <n-h2 class="!text-base !font-bold !m-0 !text-gray-900 group-hover:text-emerald-600 transition-colors truncate">
                        {{ project.name }}
                      </n-h2>

                      <!-- Quy mô / Diện tích / Số căn -->
                      <n-space :size="16" align="center" class="text-sm">
                        <n-text depth="3" v-if="project.total_area_ha" class="flex items-center gap-1">
                          <span class="text-gray-400">◱</span>
                          <strong class="text-gray-700">{{ project.total_area_ha }} ha</strong>
                        </n-text>
                        <n-text depth="3" v-if="project.total_units" class="flex items-center gap-1">
                          <span class="text-gray-400">⊞</span>
                          <strong class="text-gray-700">{{ project.total_units }}</strong> căn
                        </n-text>
                      </n-space>

                      <!-- Địa chỉ -->
                      <n-text depth="3" class="text-xs truncate block">
                        {{ project.full_address || 'Địa chỉ đang cập nhật' }}
                      </n-text>

                      <!-- Mô tả ngắn -->
                      <n-text depth="3" class="text-xs line-clamp-2 leading-relaxed">
                        {{ project.description || `${project.name} là dự án bất động sản tọa lạc tại ${project.full_address || 'vị trí đang cập nhật'}.` }}
                      </n-text>
                    </n-space>

                    <!-- Khoảng giá bán -->
                    <n-space justify="space-between" align="center" class="border-t border-gray-50 pt-2">
                      <n-text depth="3" class="text-xs">Mức giá khoảng</n-text>
                      <n-text type="success" class="!font-bold !text-sm text-emerald-600">
                        {{ formatPriceRange(project) }}
                      </n-text>
                    </n-space>
                  </n-grid-item>
                </n-grid>
              </n-card>

              <!-- Phân trang dùng component có sẵn Pagination.vue -->
              <Pagination
                v-if="pageCount > 1"
                :current-page="page"
                :total-pages="pageCount"
                :total-items="projects.length"
                :page-size="pageSize"
                @update:current-page="handlePageChange"
              />
            </n-space>
          </n-spin>
        </n-grid-item>

        <!-- Cột phải: Sidebar (Chiếm 1/4 cột) -->
        <n-grid-item :span="1">
          <n-card
            title="Đánh giá dự án"
            header-style="border-b: 1px solid #f3f4f6; padding: 12px 16px;"
            content-style="padding: 0;"
            class="overflow-hidden"
          >
            <template #header-extra>
              <n-button text type="primary" class="text-xs font-semibold" @click="navigateTo('/reviews')">
                Xem tất cả →
              </n-button>
            </template>

            <!-- Sidebar list sử dụng n-space vertical -->
            <n-space vertical :size="0" class="divide-y divide-gray-100">
              <div
                v-for="(item, i) in sidebarReviews"
                :key="i"
                class="relative h-40 overflow-hidden cursor-pointer group"
              >
                <img
                  :src="item.img"
                  :alt="item.title"
                  class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                />
                <div class="absolute inset-0 from-black/70 via-black/10 to-transparent"></div>
                <div class="absolute bottom-0 left-0 right-0 p-3">
                  <n-text class="text-white text-xs font-semibold line-clamp-2 leading-snug">
                    {{ item.title }}
                  </n-text>
                </div>
              </div>
            </n-space>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-layout-content>
  </n-layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'

interface Project {
  id: number
  name: string
  slug: string
  status: string
  full_address: string
  description?: string
  total_area_ha?: number
  total_units?: number
  price_min?: number
  price_max?: number
}

const props = defineProps<{
  categorySlug: string
}>()

const { $api } = useNuxtApp()
const projects = ref<Project[]>([])
const loading = ref(false)
const sortOrder = ref('newest')

// Phân trang state
const page = ref(1)
const pageSize = ref(6) // Hiển thị 6 dự án mỗi trang

// Options sắp xếp cho n-select của Naive UI
const sortOptions = [
  { label: 'Mới nhất', value: 'newest' },
  { label: 'Cũ nhất', value: 'oldest' },
  { label: 'Giá tăng dần', value: 'price_asc' },
  { label: 'Giá giảm dần', value: 'price_desc' },
]

// Sidebar mock data — replace with real API sau
const sidebarReviews = ref([
  {
    title: 'Những Dự Án Căn Hộ View Sông Sài Gòn Tại TP.HCM',
    img: 'https://placehold.co/560x320/94a3b8/ffffff?text=Review+1',
  },
  {
    title: 'Top 10 Dự Án Đáng Đầu Tư Nhất Năm 2025',
    img: 'https://placehold.co/560x320/7c9eb2/ffffff?text=Review+2',
  },
  {
    title: 'Bất động sản Hà Nội: Xu hướng và cơ hội',
    img: 'https://placehold.co/560x320/64748b/ffffff?text=Review+3',
  },
])

const formatStatus = (status?: string | boolean): string => {
  if (!status) return 'Đang cập nhật'
  const s = String(status).toLowerCase().trim()
  if (s === 'active') {
    return 'Đang mở bán'
  }
  if (s === 'inactive') {
    return 'Sắp mở bán'
  }
  return status as string
}

// Map trạng thái sang kiểu màu Naive UI (default, primary, info, success, warning, error)
const statusTagType = (status?: string) => {
  const formatted = formatStatus(status)
  if (formatted === 'Đang mở bán') return 'success'
  if (formatted === 'Sắp mở bán') return 'warning'
  return 'default'
}

const formatPriceRange = (project: Project) => {
  const min = project.price_min
  const max = project.price_max
  if (!min && !max) return 'Liên hệ'

  const toBillion = (val: number) => {
    return (val / 1000000000).toFixed(1).replace('.0', '') + ' tỷ'
  }

  if (min && max) {
    return `${toBillion(min)} - ${toBillion(max)}`
  }
  if (min) return `Từ ${toBillion(min)}`
  return `Đến ${toBillion(max!)}`
}

const sortedProjects = computed(() => {
  const list = [...projects.value]
  if (sortOrder.value === 'price_asc') {
    return list.sort((a, b) => (a.price_min ?? 0) - (b.price_min ?? 0))
  }
  if (sortOrder.value === 'price_desc') {
    return list.sort((a, b) => (b.price_min ?? 0) - (a.price_min ?? 0))
  }
  if (sortOrder.value === 'oldest') return list.sort((a, b) => a.id - b.id)
  return list.sort((a, b) => b.id - a.id)
})

// Dự án của trang hiện tại sau khi phân trang
const pagedProjects = computed(() => {
  const start = (page.value - 1) * pageSize.value
  const end = start + pageSize.value
  return sortedProjects.value.slice(start, end)
})

// Tổng số trang
const pageCount = computed(() => {
  return Math.ceil(projects.value.length / pageSize.value)
})

// Cuộn mượt mà lên đầu trang khi chuyển trang
const handlePageChange = (val: number) => {
  page.value = val
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

const fetchProjects = async () => {
  if (!props.categorySlug) return
  loading.value = true
  try {
    const res = await $api.get<{ data: Project[] }>(`/real-estate/project-category/${props.categorySlug}`)
    projects.value = res.data || []
  } catch (error) {
    console.error('Lỗi khi tải danh sách dự án:', error)
    projects.value = []
  } finally {
    loading.value = false
  }
}

const goToProject = (project: Project) => {
  navigateTo(`/${project.slug}-pj${project.id}`)
}

onMounted(() => {
  fetchProjects()
})

// Reset trang về 1 khi đổi danh mục hoặc đổi kiểu sắp xếp
watch([() => props.categorySlug, sortOrder], () => {
  page.value = 1
})

watch(() => props.categorySlug, () => {
  fetchProjects()
})
</script>
