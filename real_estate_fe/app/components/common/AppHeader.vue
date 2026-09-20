<template>
  <header class="sticky top-0 z-50 border-b border-gray-100 bg-white/95 backdrop-blur">
    <div class="mx-auto flex h-16 max-w-7xl items-center gap-4 px-4 lg:px-6">
      <!-- Logo -->
      <NuxtLink to="/" class="flex shrink-0 items-center gap-2">
        <img src="/logo.svg" alt="NhàViệt" class="h-9 w-9" :style="{ filter: LOGO_FILTER }" />
        <span class="text-lg font-bold tracking-tight text-gray-800">
          Nhà<span class="text-emerald-600">Việt</span>
        </span>
      </NuxtLink>

      <!-- Menu danh mục (chiếm phần giữa, tự co giãn) -->
      <div class="min-w-0 flex-1">
        <CategoryMenu />
      </div>

      <!-- Nhóm tài khoản + CTA.
           Phân cấp thị giác: Đăng nhập (nhẹ) < Đăng ký (viền) < Đăng tin (đặc).
           Dùng đúng type của Naive UI thay vì đè màu bằng Tailwind để hover/padding chuẩn. -->
      <nav class="flex shrink-0 items-center gap-2">
        <template v-if="auth.isAuthenticated">
          <NotificationBell />
          <UserProfileMenu />
        </template>
        <template v-else>
          <n-button quaternary @click="goToLogin">Đăng nhập</n-button>
          <n-button type="primary" ghost @click="goToRegister">Đăng ký</n-button>
        </template>

        <!-- Vạch ngăn giữa nhóm tài khoản và hành động chính -->
        <span class="mx-1 hidden h-5 w-px bg-gray-200 sm:block" aria-hidden="true" />

        <n-button type="primary" @click="goToCreatePost">
          <template #icon>
            <n-icon>
              <IconAddOutline />
            </n-icon>
          </template>
          Đăng tin
        </n-button>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import IconAddOutline from '~/icons/IconAddOutline.vue'

/** Đổi màu logo.svg sang tông emerald của thương hiệu */
const LOGO_FILTER =
  'invert(48%) sepia(79%) saturate(476%) hue-rotate(114deg) brightness(95%) contrast(95%)'

const auth = useAuthStore()

// Dùng alias tiếng Việt cho URL người dùng nhìn thấy
const goToLogin = () => navigateTo('/dang-nhap')
const goToRegister = () => navigateTo('/dang-ky')
const goToCreatePost = () => navigateTo('/nguoi-ban/dang-tin')
</script>
