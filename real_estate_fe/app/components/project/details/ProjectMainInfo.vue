<template>
  <div v-if="project" class="flex flex-col gap-6">
    <!-- Gallery -->
    <div class="relative h-90 cursor-pointer overflow-hidden rounded-xl" @click="openLightbox(0)">
      <div class="grid h-full gap-1" :style="gridStyle">
        <!-- Ảnh lớn đầu tiên -->
        <div class="row-span-2 overflow-hidden">
          <img :src="images[0]?.url" :alt="images[0]?.file_name" class="h-full w-full object-cover" />
        </div>

        <!-- Các ảnh còn lại dynamic -->
        <div v-for="(img, idx) in previewImages" :key="img.id" class="relative overflow-hidden">
          <img :src="img.thumbnail_url || img.url" :alt="img.file_name" class="h-full w-full object-cover" />

          <!-- Badge +N ở ảnh cuối -->
          <div
            v-if="idx === previewImages.length - 1 && remainingCount > 0"
            class="absolute inset-0 flex items-center justify-center bg-black/50 text-lg font-semibold text-white"
          >
            +{{ remainingCount }}
          </div>
        </div>
      </div>

      <!-- Nút xem tất cả -->
      <div class="absolute bottom-3 right-3">
        <n-button size="small" round secondary @click.stop="openLightbox(0)">
          <template #icon>
            <IconImage />
          </template>
          Xem tất cả {{ images.length }} ảnh
        </n-button>
      </div>
    </div>

    <!-- Title block -->
    <div class="flex flex-col gap-2">
      <div class="flex items-center gap-3">
        <n-tag :type="statusTagType(project.status)" size="small" round class="font-semibold shadow-sm">
          {{ formatStatus(project.status) }}
        </n-tag>
        <span class="flex items-center gap-1 text-xs text-gray-400">
          <IconEye class="h-3.5 w-3.5" />
          Lượt xem: {{ project.view_count || 0 }}
        </span>
      </div>

      <h1 class="m-0 text-2xl font-bold leading-snug text-gray-900">{{ project.name }}</h1>
      <span class="flex items-start gap-1 text-sm text-gray-400">
        <IconMapPin class="mt-0.5 h-4 w-4 shrink-0" />
        {{ project.full_address || 'Địa chỉ đang cập nhật' }}
      </span>
    </div>

    <!-- Stats grid -->
    <div class="grid grid-cols-3 gap-4">
      <div class="rounded-xl border border-gray-100/50 bg-gray-50/50 p-3">
        <div class="flex flex-col gap-1">
          <span class="text-[10px] font-bold uppercase tracking-wider text-gray-400">QUY MÔ</span>
          <span class="text-base font-bold text-gray-800">
            {{ project.total_area_ha ? `${project.total_area_ha} ha` : 'Đang cập nhật' }}
          </span>
        </div>
      </div>

      <div class="rounded-xl border border-gray-100/50 bg-gray-50/50 p-3">
        <div class="flex flex-col gap-1">
          <span class="text-[10px] font-bold uppercase tracking-wider text-gray-400">SỐ CĂN HỘ / NỀN</span>
          <span class="text-base font-bold text-gray-800">
            {{ project.total_units ? `${project.total_units} căn` : 'Đang cập nhật' }}
          </span>
          <!-- Tồn kho: chỉ trừ khi admin duyệt tài liệu mua nhà của khách -->
          <span v-if="remainingUnits !== null"
            :class="remainingUnits === 0 ? 'text-xs font-semibold text-red-500' : 'text-xs font-medium text-emerald-600'">
            {{ remainingUnits === 0 ? 'Đã hết căn' : `Còn ${remainingUnits} căn` }}
          </span>
        </div>
      </div>

      <div class="rounded-xl border border-gray-100/50 bg-gray-50/50 p-3">
        <div class="flex flex-col gap-1">
          <span class="text-[10px] font-bold uppercase tracking-wider text-gray-400">KHOẢNG GIÁ</span>
          <span class="text-base font-bold text-emerald-600">
            {{ formatPriceRange(project.price_min, project.price_max) }}
          </span>
        </div>
      </div>
    </div>

    <!-- Detail card -->
    <div class="rounded-xl border border-gray-100 shadow-sm">
      <div class="border-b border-gray-100 px-4 py-3 font-semibold text-gray-800">Thông tin chi tiết</div>
      <div class="flex flex-col gap-3 p-4 text-sm leading-relaxed text-gray-600">
        <p class="m-0">
          Dự án <strong>{{ project.name }}</strong> tọa lạc tại vị trí đắc địa thuộc khu vực {{ project.full_address }}.
          Với tổng quy mô đầu tư phát triển lên đến {{ project.total_area_ha ? `${project.total_area_ha} ha` : 'nhiều ha' }},
          dự án hứa hẹn sẽ mang đến không gian sống đẳng cấp, tiện nghi cùng cơ hội đầu tư sinh lời vượt trội cho quý khách hàng.
        </p>
        <p class="m-0">
          Được quy hoạch bài bản đồng bộ với tổng số lượng sản phẩm khoảng {{ project.total_units ? `${project.total_units} căn hộ/nhà phố` : 'nhiều sản phẩm đa dạng' }},
          thiết kế hiện đại chuẩn xanh, tối ưu hóa công năng và ánh sáng tự nhiên.
        </p>
      </div>
    </div>

    <n-modal
      v-model:show="showLightbox"
      :mask-closable="true"
      :closable="true"
      preset="card"
      class="!m-0 !max-w-[100vw] !rounded-none !border-0 !p-0"
      style="width: 100vw; height: 100vh; background: rgba(0, 0, 0, 0.9);"
    >
      <div
        class="flex h-full flex-col bg-black/90"
        tabindex="0"
        @keydown.left.prevent="prev"
        @keydown.right.prevent="next"
      >
        <!-- Header -->
        <div class="flex items-center justify-between px-6 py-4 text-white">
          <span class="text-sm text-white/70">{{ currentIndex + 1 }} / {{ images.length }}</span>
          <button
            class="rounded-full p-2 text-white/70 transition hover:bg-white/10 hover:text-white"
            @click="showLightbox = false"
          >
            <IconX class="h-5 w-5" />
          </button>
        </div>

        <!-- Main image -->
        <div class="relative flex flex-1 items-center justify-center px-10 py-2">
          <button
            class="absolute left-4 flex h-11 w-11 items-center justify-center rounded-full bg-white/10 text-white transition hover:bg-white/20"
            @click="prev"
          >
            <IconChevronLeft class="h-6 w-6" />
          </button>

          <img
            :src="images[currentIndex]?.url"
            class="max-h-[72vh] max-w-[78vw] rounded-lg object-contain shadow-2xl"
          />

          <button
            class="absolute right-4 flex h-11 w-11 items-center justify-center rounded-full bg-white/10 text-white transition hover:bg-white/20"
            @click="next"
          >
            <IconChevronRight class="h-6 w-6" />
          </button>
        </div>

        <!-- Thumbnail strip -->
        <div class="thumbnail-strip flex justify-center gap-2 overflow-x-auto px-4 pb-5 pt-2">
          <div
            v-for="(img, idx) in images"
            :key="img.id"
            class="h-16 w-16 shrink-0 cursor-pointer overflow-hidden rounded-md border border-white/10 transition-all"
            :class="idx === currentIndex ? 'scale-105 border-white ring-2 ring-white' : 'opacity-60 hover:opacity-90'"
            @click="currentIndex = idx"
          >
            <img :src="img.thumbnail_url || img.url" class="h-full w-full object-cover" />
          </div>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import type { ProjectDetail } from '~/types/project';
import { formatPriceRange, formatStatus, statusTagType } from '~/utils/format';

const MAX_PREVIEW = 5

const props = defineProps<{
  project: ProjectDetail;
}>();

const images = computed(() => props.project.images ?? [])
// Tồn kho còn lại = total_units - sold_units; null khi dự án không giới hạn số căn
const remainingUnits = computed(() => {
  if (props.project.total_units == null) return null
  return Math.max(0, props.project.total_units - (props.project.sold_units ?? 0))
})
const previewImages = computed(() => images.value.slice(1, MAX_PREVIEW))
const remainingCount = computed(() => Math.max(0, images.value.length - MAX_PREVIEW))
const gridStyle = computed(() => {
  const cols = Math.min(Math.ceil(previewImages.value.length / 2), 2)
  return {
    gridTemplateColumns: `2fr ${Array(cols).fill('1fr').join(' ')}`,
  }
})

const showLightbox = ref(false)
const currentIndex = ref(0)
const lightboxRef = ref<HTMLElement | null>(null)



function openLightbox(idx: number) {
  if (!images.value.length) return
  currentIndex.value = idx
  showLightbox.value = true
}

function prev() {
  if (!images.value.length) return
  currentIndex.value = (currentIndex.value - 1 + images.value.length) % images.value.length
}

function next() {
  if (!images.value.length) return
  currentIndex.value = (currentIndex.value + 1) % images.value.length
}
</script>

