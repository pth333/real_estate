<template>
  <div class="flex h-screen flex-col overflow-hidden bg-gray-50">
    <DashboardTopbar title="Quản trị escrow" subtitle="Giám sát dòng tiền đặt cọc &amp; xử lý tranh chấp"
      :icon="IconShieldCheck" />

    <div class="flex flex-1 overflow-hidden">
      <DashboardSidebar :active-key="activeKey" :groups="menuGroups" />

      <main class="flex flex-1 flex-col overflow-y-auto p-4 lg:p-6">
        <div class="mb-5 flex shrink-0 flex-col gap-1">
          <nav class="flex items-center gap-1.5 text-xs text-gray-400">
            <span>Quản trị</span>
            <span>/</span>
            <span class="font-medium text-emerald-600">{{ currentTitle }}</span>
          </nav>
          <h1 class="text-xl font-bold tracking-tight text-gray-900">{{ currentTitle }}</h1>
        </div>

        <div class="flex min-h-0 flex-1 flex-col">
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DashboardMenuGroup } from '~/types/dashboard'
import DashboardTopbar from '~/components/dashboard/DashboardTopbar.vue'
import DashboardSidebar from '~/components/dashboard/DashboardSidebar.vue'
import IconShieldCheck from '~/icons/IconShieldCheck.vue'
import IconWallet from '~/icons/IconWallet.vue'
import IconInfo from '~/icons/IconInfo.vue'
import IconUser from '~/icons/IconUser.vue'

const route = useRoute()

const menuGroups: DashboardMenuGroup[] = [
  {
    label: 'Escrow',
    items: [
      {
        key: 'escrow',
        label: 'Tổng quan escrow',
        pageTitle: 'Tổng quan quỹ escrow',
        path: '/admin',
        icon: IconWallet,
        match: ['/admin'],
      },
      {
        key: 'disputes',
        label: 'Tranh chấp',
        pageTitle: 'Tranh chấp cần xử lý',
        path: '/admin/disputes',
        icon: IconInfo,
        match: ['disputes'],
      },
    ],
  },
  {
    label: 'Hệ thống',
    items: [
      {
        key: 'users',
        label: 'Người dùng & phân quyền',
        pageTitle: 'Người dùng & phân quyền',
        path: '/admin/users',
        icon: IconUser,
        match: ['/admin/users'],
      },
    ],
  },
]

const allItems = computed(() => menuGroups.flatMap((group) => group.items))

// Danh sách menu là hằng số nên luôn có mục đầu để làm mặc định
const fallbackItem = menuGroups[0]?.items[0]

// Mục "Tổng quan escrow" khớp /admin, nhưng /admin/users phải thuộc mục người dùng
// → findActiveItem chọn đoạn khớp dài nhất
const activeItem = computed(() => findActiveItem(allItems.value, route.path) ?? fallbackItem)

const activeKey = computed(() => activeItem.value?.key ?? '')
const currentTitle = computed(() => activeItem.value?.pageTitle ?? 'Quản trị escrow')
</script>
