<template>
  <header class="flex h-16 shrink-0 items-center justify-between gap-4 border-b border-gray-100 bg-white px-4 lg:px-6">
    <!-- Trái: thương hiệu + tên khu vực -->
    <div class="flex min-w-0 items-center gap-3">
      <NuxtLink to="/" class="flex shrink-0 items-center gap-2">
        <img :src="LOGO_PATH" :alt="BRAND_NAME" class="h-8 w-8" :style="{ filter: LOGO_FILTER }" />
        <span class="hidden text-base font-bold tracking-tight text-gray-800 sm:block">
          Nhà<span class="text-emerald-600">{{ BRAND_NAME_ACCENT }}</span>
        </span>
      </NuxtLink>

      <span class="hidden h-6 w-px bg-gray-200 sm:block" aria-hidden="true" />

      <div class="flex min-w-0 items-center gap-2.5">
        <span class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-emerald-600">
          <n-icon size="18">
            <component :is="icon" />
          </n-icon>
        </span>
        <div class="flex min-w-0 flex-col leading-tight">
          <span class="truncate text-sm font-semibold text-gray-900">{{ title }}</span>
          <span class="truncate text-xs text-gray-400">{{ subtitle }}</span>
        </div>
      </div>
    </div>

    <!-- Phải: về trang chủ + menu tài khoản -->
    <div class="flex shrink-0 items-center gap-2">
      <NuxtLink to="/"
        class="hidden items-center gap-1.5 text-sm text-gray-500 transition hover:text-emerald-600 sm:inline-flex">
        <n-icon size="16">
          <IconChevronLeft />
        </n-icon>
        Về trang chủ
      </NuxtLink>
      <span class="hidden h-5 w-px bg-gray-200 sm:block" aria-hidden="true" />
      <UserProfileMenu />
    </div>
  </header>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import { BRAND_NAME, BRAND_NAME_ACCENT, LOGO_FILTER, LOGO_PATH } from '~/utils/brand'
import IconChevronLeft from '~/icons/IconChevronLeft.vue'

defineProps<{
  /** Tên khu vực, VD "Quản lý tin đăng" */
  title: string
  /** Mô tả ngắn dưới tên khu vực */
  subtitle: string
  /** Icon đại diện khu vực */
  icon: Component
}>()
</script>
