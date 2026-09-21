<template>
  <aside class="flex h-full w-64 shrink-0 flex-col border-r border-gray-100 bg-white">
    <div class="flex flex-1 flex-col gap-6 overflow-y-auto px-3 py-4">
      <!-- Hành động chính của khu vực (nếu có) -->
      <div v-if="cta" class="px-1">
        <n-button type="primary" block size="large" class="rounded-lg" @click="navigateTo(cta.path)">
          <template #icon>
            <n-icon>
              <component :is="cta.icon" />
            </n-icon>
          </template>
          {{ cta.label }}
        </n-button>
      </div>

      <!-- Các nhóm menu -->
      <div v-for="group in visibleGroups" :key="group.label" class="flex flex-col gap-1.5">
        <p class="px-3 text-[11px] font-semibold uppercase tracking-wider text-gray-400">
          {{ group.label }}
        </p>
        <n-menu :value="activeKey" :options="group.options" :indent="16" @update:value="handleSelect" />
      </div>

      <!-- Không còn mục nào để hiển thị -->
      <n-empty v-if="!visibleGroups.length && !cta" size="small" description="Không có mục nào" class="mt-6" />
    </div>

    <!-- Thông tin người đang đăng nhập -->
    <div class="shrink-0 border-t border-gray-100 px-4 py-3.5">
      <div class="flex items-center gap-3">
        <n-avatar round :size="36" class="bg-emerald-50 text-emerald-600">
          {{ initial }}
        </n-avatar>
        <div class="flex min-w-0 flex-col">
          <span class="truncate text-sm font-semibold text-gray-800">{{ displayName }}</span>
          <span class="truncate text-xs text-gray-400">{{ auth.user?.email || "Chưa đăng nhập" }}</span>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { h, type Component } from 'vue'
import { NIcon, type MenuOption } from 'naive-ui'
import type { DashboardMenuGroup } from '~/types/dashboard'
import { useAuthStore } from '~/stores/auth'

const props = defineProps<{
  /** Mục đang mở */
  activeKey: string
  /** Danh sách nhóm menu */
  groups: DashboardMenuGroup[]
  /** Nút hành động chính ở đầu sidebar */
  cta?: { label: string; path: string; icon: Component }
}>()

const auth = useAuthStore()

const displayName = computed(() => auth.user?.name || 'Người dùng')
const initial = computed(() => displayName.value.charAt(0).toUpperCase())

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

/** Bỏ mục ẩn (chỉ dùng để nhận biết trang) và mục mà user không có quyền */
const visibleGroups = computed(() =>
  props.groups
    .map((group) => ({
      label: group.label,
      options: group.items
        .filter((item) => !item.hidden)
        .filter((item) => !item.permission || auth.can(item.permission))
        .map<MenuOption>((item) => ({
          label: item.label,
          key: item.key,
          icon: renderIcon(item.icon),
        })),
    }))
    .filter((group) => group.options.length > 0),
)

// Điều hướng theo key: tra trong toàn bộ items (không phụ thuộc nhóm nào đang hiện)
function handleSelect(key: string) {
  const item = props.groups
    .flatMap((group) => group.items)
    .find((candidate) => candidate.key === key)
  if (item) navigateTo(item.path)
}
</script>
