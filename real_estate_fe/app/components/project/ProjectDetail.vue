<template>
  <div class="container mx-auto px-24 py-8">
    <!-- Breadcrumb -->
    <div class="mb-5 text-sm text-gray-500">
      <span class="cursor-pointer hover:text-emerald-600" @click="navigateTo('/')">Trang chủ</span>
      <span class="mx-2">/</span>
      <span class="cursor-pointer hover:text-emerald-600" @click="navigateTo('/du-an')">Dự án</span>
      <span class="mx-2">/</span>
      <span class="text-gray-900 font-medium">{{ project?.name || 'Chi tiết dự án' }}</span>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="animate-pulse flex flex-col gap-6">
      <div class="h-96 bg-gray-200 rounded-lg w-full"></div>
      <div class="h-8 bg-gray-200 rounded w-1/3"></div>
      <div class="h-4 bg-gray-200 rounded w-1/2"></div>
      <div class="h-24 bg-gray-200 rounded w-full"></div>
    </div>

    <!-- Error/Not Found State -->
    <div v-else-if="!project" class="text-center py-16 text-gray-500">
      <p class="text-lg">Không tìm thấy thông tin dự án này hoặc đã xảy ra lỗi.</p>
      <n-button type="primary" class="mt-4" @click="navigateTo('/')">Quay lại trang chủ</n-button>
    </div>

    <!-- Main Content -->
    <div v-else class="grid grid-cols-3 gap-8">
      <!-- Cột trái: Thông tin chính -->
      <div class="col-span-2 flex flex-col gap-6">
        <!-- Gallery ảnh (placeholder) -->
        <div class="relative h-[450px] overflow-hidden rounded-xl bg-gray-100 shadow-sm border border-gray-100">
          <img :src="project.thumbnail" :alt="project.name" class="w-full h-full object-cover" />
          <div
            class="absolute bottom-4 right-4 bg-black/60 text-white text-xs px-3 py-1.5 rounded-full flex items-center gap-1.5 font-medium">
            <IconImage class="h-3.5 w-3.5" />
            Xem tất cả hình ảnh
          </div>
        </div>

        <!-- Tiêu đề, địa chỉ và trạng thái -->
        <div>
          <div class="flex items-center gap-3 mb-2.5">
            <span class="text-xs font-semibold px-2.5 py-1 border rounded-md" :class="statusClass(project.status)">
              {{ formatStatus(project.status) }}
            </span>
            <span class="text-xs text-gray-500 flex items-center gap-1">
              <IconEye class="h-3.5 w-3.5" />
              Lượt xem: {{ project.view_count || 0 }}
            </span>
          </div>
          <h1 class="text-2xl font-bold text-gray-900 mb-2">{{ project.name }}</h1>
          <p class="text-sm text-gray-500 flex items-start gap-1.5">
            <IconMapPin class="h-4 w-4 text-gray-400 mt-0.5 shrink-0" />
            {{ project.full_address || 'Địa chỉ đang cập nhật' }}
          </p>
        </div>

        <hr class="border-gray-100" />

        <!-- Thông tin cơ bản: Diện tích, quy mô -->
        <div class="grid grid-cols-3 gap-4 py-2">
          <div class="bg-gray-50/50 p-4 rounded-xl border border-gray-100/50 flex flex-col gap-1">
            <span class="text-xs text-gray-400 font-medium">QUY MÔ</span>
            <span class="text-base font-bold text-gray-800">
              {{ project.total_area_ha ? project.total_area_ha + ' ha' : 'Đang cập nhật' }}
            </span>
          </div>
          <div class="bg-gray-50/50 p-4 rounded-xl border border-gray-100/50 flex flex-col gap-1">
            <span class="text-xs text-gray-400 font-medium">SỐ CĂN HỘ / NỀN</span>
            <span class="text-base font-bold text-gray-800">
              {{ project.total_units ? project.total_units + ' căn' : 'Đang cập nhật' }}
            </span>
          </div>
          <div class="bg-gray-50/50 p-4 rounded-xl border border-gray-100/50 flex flex-col gap-1">
            <span class="text-xs text-gray-400 font-medium">KHOẢNG GIÁ</span>
            <span class="text-base font-bold text-emerald-600">
              {{ formatPriceRange(project.price_min, project.price_max) }}
            </span>
          </div>
        </div>

        <hr class="border-gray-100" />

        <!-- Giới thiệu dự án -->
        <div>
          <h3 class="text-lg font-bold text-gray-900 mb-3">Thông tin chi tiết</h3>
          <div class="text-sm text-gray-600 leading-relaxed space-y-3">
            <p>
              Dự án <strong>{{ project.name }}</strong> tọa lạc tại vị trí đắc địa thuộc khu vực {{ project.full_address
              }}.
              Với tổng quy mô đầu tư phát triển lên đến {{ project.total_area_ha ? project.total_area_ha + ' ha' :
                'nhiều ha' }},
              dự án hứa hẹn sẽ mang đến không gian sống đẳng cấp, tiện nghi cùng cơ hội đầu tư sinh lời vượt trội cho
              quý khách hàng.
            </p>
            <p>
              Được quy hoạch bài bản đồng bộ với tổng số lượng sản phẩm khoảng {{ project.total_units ?
                project.total_units + ' căn hộ/nhà phố' : 'nhiều sản phẩm đa dạng' }},
              thiết kế hiện đại chuẩn xanh, tối ưu hóa công năng và ánh sáng tự nhiên.
            </p>
          </div>
        </div>
      </div>

      <!-- Cột phải: Form liên hệ / Tư vấn -->
      <div class="col-span-1 flex flex-col gap-6">
        <div class="bg-white p-6 rounded-xl border border-gray-100 shadow-sm sticky top-6">
          <div class="flex items-center gap-3 mb-4">
            <div
              class="w-12 h-12 rounded-full bg-emerald-50 flex items-center justify-center text-emerald-600 font-bold text-lg">
              S
            </div>
            <div>
              <div class="font-bold text-gray-900 text-sm">Phòng Kinh Doanh</div>
              <div class="text-xs text-gray-400">Hỗ trợ tư vấn dự án 24/7</div>
            </div>
          </div>

          <div class="flex flex-col gap-3">
            <n-button type="primary" size="large"
              class="w-full !bg-emerald-600 hover:!bg-emerald-700 font-semibold flex items-center gap-2">
              <IconPhone class="h-4 w-4" />
              Yêu cầu gọi lại tư vấn
            </n-button>
            <n-button ghost size="large" class="w-full text-gray-700 border-gray-300 font-medium">
              Tải bảng giá & tài liệu pháp lý
            </n-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const props = defineProps<{
  id: number
}>()

interface ProjectDetail {
  id: number
  name: string
  slug: string
  status: string
  full_address: string
  total_area_ha?: number
  total_units?: number
  price_min?: number
  price_max?: number
  view_count?: number
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
const project = ref<ProjectDetail | null>(null)
const loading = ref(true)

const fetchProjectDetail = async () => {
  loading.value = true
  try {
    // Gọi API lấy chi tiết dự án chính xác theo ID
    const res = await $api.get<{ data: any }>("/real-estate/project/detail/" + props.id)
    if (res.data) {
      project.value = {
        ...res.data,
        thumbnail: res.data.thumbnail || "https://placehold.co/600x400/e2e8f0/94a3b8?text=Project"
      }
    }
  } catch (error) {
    console.error("Lỗi khi tải chi tiết dự án:", error)
  } finally {
    loading.value = false
  }
}

function formatPriceRange(min?: number, max?: number) {
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

onMounted(async () => {
  // Trigger tăng lượt xem dự án ở backend ngay khi người dùng xem chi tiết
  try {
    await $api.post("/real-estate/project/view/" + props.id)
  } catch (err) {
    console.error("Lỗi khi tăng lượt xem dự án:", err)
  }

  // Nạp thông tin chi tiết dự án (sẽ chứa view_count mới nhất vừa được tăng)
  fetchProjectDetail()
})
</script>
