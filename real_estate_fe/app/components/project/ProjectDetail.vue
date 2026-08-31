<template>
  <n-layout class="bg-transparent">
    <n-layout-content class="mx-auto max-w-[1200px] px-6 py-6 bg-transparent">
      <!-- Breadcrumb Naive UI chuyên nghiệp -->
      <n-breadcrumb class="mb-5">
        <n-breadcrumb-item @click="navigateTo('/')">Trang chủ</n-breadcrumb-item>
        <n-breadcrumb-item @click="navigateTo('/du-an')">Dự án</n-breadcrumb-item>
        <n-breadcrumb-item>{{ project?.name || 'Chi tiết dự án' }}</n-breadcrumb-item>
      </n-breadcrumb>

      <!-- Loading State bằng n-spin cực xịn bọc toàn bộ nội dung để tránh layout giật lag -->
      <n-spin :show="loading" size="large">
        <template #description>
          Đang tải thông tin dự án...
        </template>

        <!-- Khi lỗi/không tìm thấy -->
        <n-empty v-if="!loading && !project" description="Không tìm thấy thông tin dự án này hoặc đã xảy ra lỗi." class="py-20">
          <template #extra>
            <n-button type="primary" @click="navigateTo('/')">Quay lại trang chủ</n-button>
          </template>
        </n-empty>

        <!-- Layout chính chia cột - Chỉ hiển thị khi project đã load thành công (khác null) -->
        <n-grid v-else-if="project" :cols="4" :x-gap="24" :y-gap="24" item-responsive>
          <!-- Cột trái: Thông tin chính (Chiếm 3/4) -->
          <n-grid-item :span="3" class="min-w-0">
            <n-space vertical :size="24">
              <!-- Gallery ảnh -->
              <n-card content-style="padding: 0;" class="overflow-hidden rounded-xl border border-gray-100 shadow-sm relative h-[450px]">
                <img :src="project.thumbnail" :alt="project.name" class="w-full h-full object-cover" />
                <div class="absolute bottom-4 right-4">
                  <n-button secondary strong round type="tertiary" size="small" class="bg-black/60! text-white! flex items-center gap-1">
                    <IconImage class="h-3.5 w-3.5" />
                    Xem tất cả hình ảnh
                  </n-button>
                </div>
              </n-card>

              <!-- Tiêu đề, địa chỉ và trạng thái -->
              <n-space vertical :size="8">
                <n-space align="center" :size="12">
                  <n-tag :type="statusTagType(project.status)" size="small" round class="font-semibold shadow-sm">
                    {{ formatStatus(project.status) }}
                  </n-tag>
                  <n-text depth="3" class="text-xs flex items-center gap-1">
                    <IconEye class="h-3.5 w-3.5 text-gray-400" />
                    Lượt xem: {{ project.view_count || 0 }}
                  </n-text>
                </n-space>

                <n-h1 class="text-2xl! font-bold! m-0! text-gray-900! leading-snug">{{ project.name }}</n-h1>
                <n-text depth="3" class="text-sm flex items-start gap-1">
                  <IconMapPin class="h-4 w-4 text-gray-400 mt-0.5 shrink-0" />
                  {{ project.full_address || 'Địa chỉ đang cập nhật' }}
                </n-text>
              </n-space>

              <!-- Thông tin cơ bản dạng Grid/Cards Naive UI -->
              <n-grid :cols="3" :x-gap="16" :y-gap="16" item-responsive class="w-full">
                <n-grid-item>
                  <n-card size="small" class="bg-gray-50/50 border border-gray-100/50 rounded-xl">
                    <n-space vertical :size="4">
                      <n-text depth="3" class="text-[10px] font-bold tracking-wider uppercase">QUY MÔ</n-text>
                      <n-text class="text-base font-bold text-gray-800">
                        {{ project.total_area_ha ? project.total_area_ha + ' ha' : 'Đang cập nhật' }}
                      </n-text>
                    </n-space>
                  </n-card>
                </n-grid-item>

                <n-grid-item>
                  <n-card size="small" class="bg-gray-50/50 border border-gray-100/50 rounded-xl">
                    <n-space vertical :size="4">
                      <n-text depth="3" class="text-[10px] font-bold tracking-wider uppercase">SỐ CĂN HỘ / NỀN</n-text>
                      <n-text class="text-base font-bold text-gray-800">
                        {{ project.total_units ? project.total_units + ' căn' : 'Đang cập nhật' }}
                      </n-text>
                    </n-space>
                  </n-card>
                </n-grid-item>

                <n-grid-item>
                  <n-card size="small" class="bg-gray-50/50 border border-gray-100/50 rounded-xl">
                    <n-space vertical :size="4">
                      <n-text depth="3" class="text-[10px] font-bold tracking-wider uppercase">KHOẢNG GIÁ</n-text>
                      <n-text type="success" class="text-base font-bold text-emerald-600">
                        {{ formatPriceRange(project.price_min, project.price_max) }}
                      </n-text>
                    </n-space>
                  </n-card>
                </n-grid-item>
              </n-grid>

              <!-- Giới thiệu dự án -->
              <n-card title="Thông tin chi tiết" header-style="border-b: 1px solid #f3f4f6; padding: 12px 16px;" class="rounded-xl border border-gray-100 shadow-sm">
                <n-space vertical :size="12" class="text-sm text-gray-600 leading-relaxed">
                  <n-text>
                    Dự án <strong>{{ project.name }}</strong> tọa lạc tại vị trí đắc địa thuộc khu vực {{ project.full_address }}.
                    Với tổng quy mô đầu tư phát triển lên đến {{ project.total_area_ha ? project.total_area_ha + ' ha' : 'nhiều ha' }},
                    dự án hứa hẹn sẽ mang đến không gian sống đẳng cấp, tiện nghi cùng cơ hội đầu tư sinh lời vượt trội cho quý khách hàng.
                  </n-text>
                  <n-text>
                    Được quy hoạch bài bản đồng bộ với tổng số lượng sản phẩm khoảng {{ project.total_units ? project.total_units + ' căn hộ/nhà phố' : 'nhiều sản phẩm đa dạng' }},
                    thiết kế hiện đại chuẩn xanh, tối ưu hóa công năng và ánh sáng tự nhiên.
                  </n-text>
                </n-space>
              </n-card>
            </n-space>
          </n-grid-item>

          <!-- Cột phải: Khung liên hệ tư vấn (Chiếm 1/4) -->
          <n-grid-item :span="1" class="relative">
            <!-- Sử dụng div với class sticky của tailwind, và bọc gọn gàng bên trong grid-item để không bị trôi -->
            <div class="sticky top-6">
              <n-card class="border border-gray-100 shadow-sm rounded-xl">
                <n-space vertical :size="16">
                  <!-- Avatar và thông tin liên hệ -->
                  <n-space align="center" :size="12">
                    <n-avatar
                      round
                      :size="48"
                      class="bg-emerald-50 text-emerald-600 font-bold border border-emerald-100"
                    >
                      S
                    </n-avatar>
                    <n-space vertical :size="2">
                      <n-text class="font-bold text-gray-900 text-sm">Phòng Kinh Doanh</n-text>
                      <n-text depth="3" class="text-xs">Hỗ trợ tư vấn dự án 24/7</n-text>
                    </n-space>
                  </n-space>

                  <!-- Các nút hành động Naive UI chuyên nghiệp -->
                  <n-space vertical :size="10">
                    <n-button
                      type="primary"
                      size="large"
                      class="w-full bg-emerald-600! hover:bg-emerald-700! font-semibold! flex items-center justify-center gap-2"
                    >
                      <IconPhone class="h-4 w-4 shrink-0" />
                      Yêu cầu gọi lại tư vấn
                    </n-button>
                    <n-button
                      ghost
                      size="large"
                      class="w-full text-gray-700 border-gray-300 font-medium"
                    >
                      Tải bảng giá & pháp lý
                    </n-button>
                  </n-space>
                </n-space>
              </n-card>
            </div>
          </n-grid-item>
        </n-grid>
      </n-spin>
    </n-layout-content>
  </n-layout>
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

// Trả về type màu chuẩn Naive UI cho tag (success, warning, error, default)
const statusTagType = (status?: string) => {
  const formatted = formatStatus(status)
  if (formatted === 'Đang mở bán') return 'success'
  if (formatted === 'Sắp mở bán') return 'warning'
  return 'default'
}

const { $api } = useNuxtApp()
const project = ref<ProjectDetail | null>(null)
const loading = ref(true)

const fetchProjectDetail = async () => {
  loading.value = true
  try {
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

const projectTitle = computed(() => project.value?.name || "Chi tiết dự án")
useHead({
  title: projectTitle,
})

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
