<template>
  <div class="flex min-h-screen bg-white">
    <!-- ── Cột form ── -->
    <div class="flex w-full flex-col px-6 py-8 sm:px-10 lg:w-1/2 lg:px-16">
      <!-- Thương hiệu: bấm về trang chủ (khách vãng lai vẫn xem web được nên luôn có đường về) -->
      <NuxtLink to="/" class="flex shrink-0 items-center gap-2 self-start">
        <img :src="LOGO_PATH" :alt="BRAND_NAME" class="h-9 w-9" :style="{ filter: LOGO_FILTER }" />
        <span class="text-lg font-bold tracking-tight text-gray-800">
          Nhà<span class="text-emerald-600">{{ BRAND_NAME_ACCENT }}</span>
        </span>
      </NuxtLink>

      <!-- Nội dung form căn giữa theo chiều dọc -->
      <div class="flex flex-1 items-center justify-center py-10">
        <div class="w-full max-w-[400px]">
          <h1 class="text-2xl font-bold tracking-tight text-gray-900">{{ title }}</h1>
          <p class="mt-1.5 text-sm text-gray-500">{{ subtitle }}</p>

          <div class="mt-8">
            <slot />
          </div>

          <div class="mt-6">
            <slot name="footer" />
          </div>
        </div>
      </div>

      <!-- Đường về trang chủ -->
      <div class="shrink-0">
        <NuxtLink to="/"
          class="inline-flex items-center gap-1.5 text-sm text-gray-500 transition hover:text-emerald-600">
          <n-icon size="16">
            <IconChevronLeft />
          </n-icon>
          Về trang chủ
        </NuxtLink>
      </div>
    </div>

    <!-- ── Cột giới thiệu nền tảng (ẩn trên mobile) ── -->
    <div
      class="relative hidden overflow-hidden bg-linear-to-br from-emerald-600 to-emerald-800 lg:flex lg:w-1/2 lg:flex-col lg:justify-center lg:px-16">
      <!-- Hoạ tiết nền -->
      <svg class="pointer-events-none absolute inset-0 h-full w-full opacity-[0.12]" aria-hidden="true">
        <defs>
          <pattern id="auth-waves" x="0" y="0" width="140" height="140" patternUnits="userSpaceOnUse">
            <path d="M10 60c18-10 42-10 60 0s42 10 60 0" fill="none" stroke="#ffffff" stroke-width="1" />
            <path d="M10 90c18-10 42-10 60 0s42 10 60 0" fill="none" stroke="#ffffff" stroke-width="1" />
            <path d="M10 30c18-10 42-10 60 0s42 10 60 0" fill="none" stroke="#ffffff" stroke-width="1" />
          </pattern>
        </defs>
        <rect width="100%" height="100%" fill="url(#auth-waves)" />
      </svg>

      <div class="relative max-w-md">
        <h2 class="text-2xl font-bold leading-snug text-white">
          Tìm nhà an tâm,<br />đặt cọc minh bạch
        </h2>

        <!-- Đúng các giá trị nền tảng đang cung cấp cho luồng đặt cọc -->
        <ul class="mt-8 flex flex-col gap-5">
          <li v-for="item in highlights" :key="item.title" class="flex items-start gap-3">
            <span class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-white/15 text-white">
              <n-icon size="18">
                <component :is="item.icon" />
              </n-icon>
            </span>
            <span class="flex flex-col">
              <span class="text-sm font-semibold text-white">{{ item.title }}</span>
              <span class="text-sm text-emerald-50/80">{{ item.description }}</span>
            </span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import {
  BRAND_NAME,
  BRAND_NAME_ACCENT,
  LOGO_FILTER,
  LOGO_PATH,
} from '~/utils/brand'
import IconChevronLeft from '~/icons/IconChevronLeft.vue'
import IconShieldCheck from '~/icons/IconShieldCheck.vue'
import IconWallet from '~/icons/IconWallet.vue'
import IconKeyOutline from '~/icons/IconKeyOutline.vue'

defineProps<{
  /** Tiêu đề trang, VD "Đăng nhập" */
  title: string
  /** Mô tả ngắn dưới tiêu đề */
  subtitle: string
}>()

const highlights: { title: string; description: string; icon: Component }[] = [
  {
    title: 'Tiền cọc được nền tảng giữ hộ',
    description: 'Không chuyển thẳng cho môi giới, chỉ giải ngân khi buổi xem có kết quả rõ ràng.',
    icon: IconShieldCheck,
  },
  {
    title: 'Hoàn tiền tự động theo thoả thuận',
    description: 'Môi giới từ chối hoặc không phản hồi trong 24 giờ là hoàn lại 100% tiền cọc.',
    icon: IconWallet,
  },
  {
    title: 'Xác nhận 2 chiều bằng OTP tại chỗ',
    description: 'Khách và môi giới cùng xác nhận đã gặp nhau, không bên nào tự quyết được.',
    icon: IconKeyOutline,
  },
]
</script>
