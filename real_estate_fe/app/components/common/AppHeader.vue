<template>
  <header class="sticky top-0 z-50 border-b border-gray-100 bg-white/95 backdrop-blur">
    <div class="mx-auto flex h-16 max-w-7xl items-center gap-2 px-4 md:gap-4 lg:px-6">
      <!-- Logo: bấm về trang chủ -->
      <NuxtLink to="/" class="flex shrink-0 items-center gap-2">
        <img :src="LOGO_PATH" :alt="BRAND_NAME" class="h-9 w-9" :style="{ filter: LOGO_FILTER }" />
        <span class="text-base font-bold tracking-tight text-gray-800 md:text-lg">
          Nhà<span class="text-emerald-600">{{ BRAND_NAME_ACCENT }}</span>
        </span>
      </NuxtLink>

      <!-- Menu danh mục ngang: chỉ hiện ở desktop vì lg mới đủ chỗ -->
      <div class="hidden min-w-0 flex-1 lg:block">
        <CategoryMenu />
      </div>

      <!-- Nhóm tài khoản + CTA.
           Phân cấp thị giác: Đăng nhập (nhẹ) < Đăng ký (viền) < Đăng tin (đặc).
           Dùng đúng type của Naive UI thay vì đè màu bằng Tailwind để hover/padding chuẩn. -->
      <nav class="ml-auto flex shrink-0 items-center gap-1 md:gap-2">
        <template v-if="auth.isAuthenticated">
          <NotificationBell />
          <UserProfileMenu />
        </template>
        <!-- Khách vãng lai: từ tablet trở lên mới đủ chỗ, mobile nằm trong drawer -->
        <div v-else class="hidden items-center gap-2 md:flex">
          <n-button quaternary @click="goToLogin">Đăng nhập</n-button>
          <n-button type="primary" ghost @click="goToRegister">Đăng ký</n-button>
        </div>

        <!-- Vạch ngăn giữa nhóm tài khoản và hành động chính -->
        <span class="mx-1 hidden h-5 w-px bg-gray-200 md:block" aria-hidden="true" />

        <!-- Mobile: nút Đăng tin nằm trong drawer cho gọn header -->
        <div class="hidden items-center md:flex">
          <n-button type="primary" @click="goToCreatePost">
            <template #icon>
              <n-icon>
                <IconAddOutline />
              </n-icon>
            </template>
            <!-- Tablet chỉ còn icon, desktop mới hiện chữ -->
            <span class="hidden lg:inline">Đăng tin</span>
          </n-button>
        </div>

        <!-- Hamburger (mobile/tablet) -->
        <div class="flex items-center lg:hidden">
          <AppMobileMenu />
        </div>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { BRAND_NAME, BRAND_NAME_ACCENT, LOGO_FILTER, LOGO_PATH } from '~/utils/brand'
import IconAddOutline from '~/icons/IconAddOutline.vue'
import AppMobileMenu from '~/components/common/AppMobileMenu.vue'

const auth = useAuthStore()

// Dùng alias tiếng Việt cho URL người dùng nhìn thấy
const goToLogin = () => navigateTo('/dang-nhap')
const goToRegister = () => navigateTo('/dang-ky')
const goToCreatePost = () => navigateTo('/nguoi-ban/dang-tin')
</script>
