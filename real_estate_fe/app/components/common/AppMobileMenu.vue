<template>
  <!-- Nút hamburger: chỉ dùng dưới desktop (lg là layout chính) -->
  <n-button quaternary circle size="large" aria-label="Mở danh mục" @click="showDrawer = true">
    <template #icon>
      <n-icon>
        <IconMenu />
      </n-icon>
    </template>
  </n-button>

  <!-- Drawer danh mục + hành động tài khoản cho mobile/tablet -->
  <n-drawer v-model:show="showDrawer" width="min(300px, 85vw)" placement="left">
    <n-drawer-content title="Danh mục bất động sản" closable>
      <div class="flex flex-col gap-4">
        <n-tree :data="treeData" block-line block-node selectable default-expand-all :selected-keys="selectedKeys"
          class="text-sm" @update:selected-keys="handleSelect" />

        <n-divider class="my-0!" />

        <!-- Khách vãng lai: 2 nút này chỉ nằm ở header từ tablet trở lên -->
        <div v-if="!auth.isAuthenticated" class="flex flex-col gap-2">
          <n-button block @click="goToLogin">Đăng nhập</n-button>
          <n-button block type="primary" ghost @click="goToRegister">Đăng ký</n-button>
        </div>

        <n-button block type="primary" @click="goToCreatePost">
          <template #icon>
            <n-icon>
              <IconAddOutline />
            </n-icon>
          </template>
          Đăng tin
        </n-button>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import type { TreeOption } from 'naive-ui'
import type { Category } from '~/types/menu'
import { useMenuStore } from '~/stores/menu'
import { useAuthStore } from '~/stores/auth'
import { useRealEstateStore } from '~/stores/real_estate'

type CategoryTreeOption = TreeOption & { slug?: string }

const menuStore = useMenuStore()
const auth = useAuthStore()
const realEstateStore = useRealEstateStore()

const showDrawer = ref(false)

const treeData = computed<CategoryTreeOption[]>(() =>
  (menuStore.menu?.categories ?? []).map(toTreeOption),
)

function toTreeOption(cat: Category): CategoryTreeOption {
  return {
    key: cat.Slug ?? String(cat.ID),
    label: cat.Name,
    slug: cat.Slug,
    children: (cat.children ?? []).map(toTreeOption),
  }
}

// Key của node chính là slug danh mục nên so khớp trực tiếp được
const selectedKeys = computed<Array<string | number>>(() =>
  realEstateStore.categorySlug ? [realEstateStore.categorySlug] : [],
)

function findSlugByKey(options: CategoryTreeOption[], key: string | number): string | null {
  for (const opt of options) {
    if (opt.key === key) return opt.slug ?? null
    const child = findSlugByKey((opt.children ?? []) as CategoryTreeOption[], key)
    if (child) return child
  }
  return null
}

function handleSelect(keys: Array<string | number>) {
  const key = keys[0]
  if (key == null) return

  const slug = findSlugByKey(treeData.value, key)
  if (!slug) return

  realEstateStore.categorySlug = slug
  showDrawer.value = false
  navigateTo(`/${slug}`)
}

// Dùng alias tiếng Việt cho URL người dùng nhìn thấy
const goToLogin = () => { showDrawer.value = false; navigateTo('/dang-nhap') }
const goToRegister = () => { showDrawer.value = false; navigateTo('/dang-ky') }
const goToCreatePost = () => { showDrawer.value = false; navigateTo('/nguoi-ban/dang-tin') }
</script>
