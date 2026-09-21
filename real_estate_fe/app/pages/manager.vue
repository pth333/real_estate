<template>
  <div class="flex h-screen flex-col overflow-hidden bg-gray-50">
    <DashboardTopbar title="Khu vực quản lý" subtitle="Tin đăng, dự án và đơn đặt cọc của bạn" :icon="IconCreateOutline" />

    <div class="flex flex-1 overflow-hidden">
      <DashboardSidebar :active-key="activeKey" :groups="menuGroups" :cta="primaryAction" />

      <main class="flex flex-1 flex-col overflow-y-auto p-4 lg:p-6">
        <!-- Tiêu đề trang: breadcrumb + tên trang -->
        <div class="mb-5 flex shrink-0 flex-col gap-1">
          <nav class="flex items-center gap-1.5 text-xs text-gray-400">
            <span>Quản lý</span>
            <span>/</span>
            <span class="font-medium text-emerald-600">{{ currentTitle }}</span>
          </nav>
          <h1 class="text-xl font-bold tracking-tight text-gray-900">{{ currentTitle }}</h1>
        </div>

        <div class="flex min-h-0 flex-1 flex-col">
          <NuxtPage />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import type { DashboardMenuGroup } from '~/types/dashboard'
import DashboardTopbar from '~/components/dashboard/DashboardTopbar.vue'
import DashboardSidebar from '~/components/dashboard/DashboardSidebar.vue'
import IconCreateOutline from '~/icons/IconCreateOutline.vue'
import IconBuilding from '~/icons/IconBuilding.vue'
import IconWallet from '~/icons/IconWallet.vue'
import IconUser from '~/icons/IconUser.vue'
import IconHeart from '~/icons/IconHeart.vue'
import IconAddOutline from '~/icons/IconAddOutline.vue'

// Alias cả folder quản lý sang /nguoi-ban và tắt layout mặc định của website
definePageMeta({
  alias: '/nguoi-ban',
  layout: 'empty',
  requiresAuth: true,
})

const route = useRoute()

const primaryAction: { label: string; path: string; icon: Component } = {
  label: 'Đăng tin mới',
  path: '/nguoi-ban/dang-tin',
  icon: IconAddOutline,
}

/**
 * Cấu hình menu tập trung: mỗi mục khai báo cả nhãn menu, tiêu đề trang và các đoạn
 * path để nhận biết mục đang mở — tránh phải sửa 3 chỗ khi thêm trang mới.
 */
const menuGroups: DashboardMenuGroup[] = [
  {
    label: 'Quản lý',
    items: [
      {
        key: 'posts',
        label: 'Tin đăng',
        pageTitle: 'Danh sách bài đăng của bạn',
        path: '/nguoi-ban/quan-ly-tin-dang',
        icon: IconCreateOutline,
        match: ['quan-ly-tin-dang', 'posts'],
      },
      {
        key: 'projects',
        label: 'Dự án',
        pageTitle: 'Danh sách dự án',
        path: '/nguoi-ban/quan-ly-du-an',
        icon: IconBuilding,
        match: ['quan-ly-du-an', 'projects'],
      },
      {
        // Mục ẩn: chỉ để nhận biết trang tạo/sửa dự án và tô sáng mục "Dự án"
        key: 'project-form',
        label: 'Dự án',
        pageTitle: 'Tạo dự án mới',
        path: '/nguoi-ban/tao-du-an',
        icon: IconBuilding,
        hidden: true,
        match: ['tao-du-an'],
      },
      {
        key: 'deposits',
        label: 'Đơn đặt cọc',
        pageTitle: 'Đơn đặt cọc xem nhà',
        path: '/nguoi-ban/quan-ly-dat-coc',
        icon: IconWallet,
        // Chỉ môi giới (có quyền xem đơn được giao) mới thấy mục này
        permission: 'broker.deposit.list',
        match: ['quan-ly-dat-coc', 'deposits'],
      },
    ],
  },
  {
    label: 'Khách hàng',
    items: [
      {
        key: 'customers',
        label: 'Khách hàng',
        pageTitle: 'Danh sách khách hàng đăng ký',
        path: '/nguoi-ban/quan-ly-khach-hang',
        icon: IconUser,
        match: ['quan-ly-khach-hang', 'customers'],
      },
      {
        key: 'favorites',
        label: 'Yêu thích',
        pageTitle: 'Danh mục bất động sản yêu thích',
        path: '/nguoi-ban/quan-ly-yeu-thich',
        icon: IconHeart,
        match: ['quan-ly-yeu-thich', 'favorites'],
      },
    ],
  },
]

const allItems = computed(() => menuGroups.flatMap((group) => group.items))

// Danh sách menu là hằng số nên luôn có mục đầu để làm mặc định
const fallbackItem = menuGroups[0]?.items[0]

// Suy ra mục đang mở từ URL (mục ẩn vẫn tham gia để ra đúng tiêu đề)
const activeItem = computed(() => findActiveItem(allItems.value, route.path) ?? fallbackItem)

/** Trang con có thể ghi đè tiêu đề (VD Tạo dự án mới / Chỉnh sửa dự án) */
const { titleOverride, clearTitle } = useDashboardTitle()
watch(() => route.path, () => clearTitle())

const activeKey = computed(() => activeItem.value?.key ?? '')
const currentTitle = computed(
  () => titleOverride.value ?? activeItem.value?.pageTitle ?? 'Khu vực quản lý',
)

useHead({
  title: computed(() => `${currentTitle.value} | NhàViệt`),
})
</script>
