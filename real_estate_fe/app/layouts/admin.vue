<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- Header -->
    <div class="h-[65px] bg-white border-b border-gray-200 px-6 flex items-center justify-between flex-shrink-0">
      <div class="flex items-center gap-3">
        <n-avatar round size="medium" class="bg-amber-50 text-amber-600">
          <template #icon>
            <n-icon>
              <IconShieldCheck />
            </n-icon>
          </template>
        </n-avatar>
        <div class="flex flex-col leading-tight">
          <span class="text-lg font-bold text-gray-900">Quản trị escrow</span>
          <span class="text-xs text-gray-400">Giám sát dòng tiền đặt cọc &amp; xử lý tranh chấp</span>
        </div>
      </div>
      <n-button text class="text-emerald-600" @click="navigateTo('/')">Quay lại trang chủ</n-button>
    </div>

    <div class="flex flex-1 overflow-hidden">
      <!-- Sidebar -->
      <div class="w-60 border-r border-gray-200 bg-white py-4 px-4 flex-shrink-0">
        <n-menu v-model:value="activeKey" :options="menuOptions" :indent="18" @update:value="handleMenuSelect" />
      </div>

      <!-- Nội dung -->
      <div class="flex-1 p-6 overflow-y-auto">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { h, ref, type Component } from 'vue'
import { NIcon, type MenuOption } from 'naive-ui'
import IconShieldCheck from '~/icons/IconShieldCheck.vue'
import IconWallet from '~/icons/IconWallet.vue'
import IconInfo from '~/icons/IconInfo.vue'
import IconUser from '~/icons/IconUser.vue'

const route = useRoute()

// Xác định mục menu đang mở theo đường dẫn hiện tại
function resolveActiveKey(path: string): string {
  if (path.includes('/admin/users')) return 'users'
  if (path.includes('disputes')) return 'disputes'
  return 'escrow'
}

const activeKey = ref<string>(resolveActiveKey(route.path))

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions: MenuOption[] = [
  { label: 'Tổng quan escrow', key: 'escrow', icon: renderIcon(IconWallet) },
  { label: 'Tranh chấp', key: 'disputes', icon: renderIcon(IconInfo) },
  { label: 'Người dùng & phân quyền', key: 'users', icon: renderIcon(IconUser) },
]

/** Đường dẫn tương ứng với từng mục menu */
const MENU_PATHS: Record<string, string> = {
  escrow: '/admin',
  disputes: '/admin/disputes',
  users: '/admin/users',
}

function handleMenuSelect(key: string) {
  activeKey.value = key
  navigateTo(MENU_PATHS[key] ?? '/admin')
}
</script>
