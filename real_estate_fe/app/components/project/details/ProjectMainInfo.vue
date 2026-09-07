<template>
  <n-space v-if="project" vertical :size="24">
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
    <n-card title="Thông tin chi tiết" header-style="border-bottom: 1px solid #f3f4f6; padding: 12px 16px;" class="rounded-xl border border-gray-100 shadow-sm">
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
</template>

<script setup lang="ts">
import { useProjectDetail } from '~/stores/detail/project_detail';

const store = useProjectDetail();
const { project } = storeToRefs(store);
</script>
