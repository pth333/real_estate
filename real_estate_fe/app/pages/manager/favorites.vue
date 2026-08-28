<template>
  <div class="flex-1 bg-white p-6 rounded-lg border border-gray-200 flex flex-col gap-4 h-full min-h-0">
    <div class="flex items-center justify-between flex-shrink-0">
      <div>
        <p class="text-sm text-gray-500 mb-1">Quản lý yêu thích</p>
        <h2 class="text-lg font-bold text-gray-800">Danh mục bất động sản yêu thích</h2>
      </div>
      <span class="inline-flex items-center rounded-full bg-red-50 px-3 py-1 text-xs font-medium text-red-500">
        {{ total }} tin đăng
      </span>
    </div>

    <!-- Data Table danh sách yêu thích -->
    <div class="flex-1 min-h-0 flex flex-col">
      <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false"
        flex-height class="flex-1 min-h-0" :scroll-x="640" />
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex justify-end flex-shrink-0">
      <Pagination :current-page="managerStore.favoritesPage" :total-pages="totalPages" @page-change="goToPage" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { NImage, NTooltip, NButton, NIcon, NSpace, type DataTableColumns } from 'naive-ui'
import type { RealEstateResponse } from '~/types/real_estate'
import { formatPrice, formatDate } from '~/utils/format'
import { useManagerStore } from '~/stores/manager'
import IconEyeOutline from '~/icons/IconEyeOutline.vue'
import IconHeart from '~/icons/IconHeart.vue'

definePageMeta({
  alias: [
    '/nguoi-ban/quan-ly-yeu-thich'
  ],
})

const { $api } = useNuxtApp()
const favorite = useFavorite()
const managerStore = useManagerStore()
const loading = ref(false)

// Data + phân trang được cache trong managerStore
const items = computed(() => managerStore.favorites)
const total = computed(() => managerStore.favoritesTotal)
const totalPages = computed(() => Math.ceil(managerStore.favoritesTotal / managerStore.favoritesSize) || 1)

const fetchFavoritesData = async () => {
  loading.value = true
  try {
    await managerStore.fetchFavorites({
      page: managerStore.favoritesPage,
      size: managerStore.favoritesSize,
    })
  } catch {
    // store đã hiện message lỗi
  } finally {
    loading.value = false
  }
}

function goToPage(page: number) {
  if (page < 1 || page > totalPages.value) return
  managerStore.favoritesPage = page
  fetchFavoritesData()
}

onMounted(fetchFavoritesData)

// Xem chi tiết tin đăng
const handleView = (row: RealEstateResponse) => {
  const slug = row.slug || `-rs${row.id}`
  navigateTo(`/${slug}`)
}

// Bỏ yêu thích → gọi API rồi xoá khỏi danh sách (trong store)
const handleRemoveFavorite = async (row: RealEstateResponse) => {
  const next = await favorite.toggle(row.id)
  if (next === false) {
    managerStore.favorites = managerStore.favorites.filter((e) => e.id !== row.id)
    managerStore.favoritesTotal = Math.max(0, managerStore.favoritesTotal - 1)
    window.message?.success('Đã bỏ yêu thích')
  }
}

const columns: DataTableColumns<RealEstateResponse> = [
  {
    title: 'Bất động sản',
    key: 'property',
    width: 300,
    render(row) {
      return h('div', { class: 'flex items-center gap-3 py-1' }, [
        h(NImage, {
          src: row.image_urls?.[0] || 'https://picsum.photos/200/150?random=default',
          alt: 'BĐS Thumbnail',
          width: 60,
          height: 45,
          class: 'rounded object-cover border border-gray-100 flex-shrink-0',
          previewDisabled: false,
        }),
        h('div', { class: 'flex flex-col gap-1 min-w-0 flex-1' }, [
          h(
            NTooltip,
            { trigger: 'hover', placement: 'top' },
            {
              trigger: () =>
                h('span', { class: 'font-medium text-gray-800 text-sm truncate block w-full cursor-default' }, row.title),
              default: () => row.title,
            }
          ),
          h('span', { class: 'text-xs text-gray-400 truncate' }, [row.district, row.city].filter(Boolean).join(', ')),
        ]),
      ])
    },
  },
  {
    title: 'Giá',
    key: 'price',
    width: 120,
    render(row) {
      return h('span', { class: 'font-semibold text-red-500 text-sm whitespace-nowrap' }, formatPrice(row.price_vnd))
    },
  },
  {
    title: 'Diện tích',
    key: 'area',
    width: 100,
    render(row) {
      return h('span', { class: 'text-gray-700 text-sm whitespace-nowrap' }, `${row.acreage.toFixed(0)} m²`)
    },
  },
  {
    title: 'Ngày đăng',
    key: 'created_at',
    width: 120,
    render(row) {
      return h('span', { class: 'text-gray-500 text-sm whitespace-nowrap' }, formatDate(row.created_at))
    },
  },
  {
    title: 'Hành động',
    key: 'actions',
    width: 120,
    align: 'center',
    render(row) {
      return h(
        NSpace,
        { justify: 'center', size: 'small', align: 'center' },
        {
          default: () => [
            h(NButton, { size: 'small', quaternary: true, type: 'info', onClick: () => handleView(row) },
              { icon: () => h(NIcon, null, { default: () => h(IconEyeOutline) }) }
            ),
            // Bỏ yêu thích
            h(NButton, { size: 'small', quaternary: true, type: 'error', onClick: () => handleRemoveFavorite(row) },
              { icon: () => h(NIcon, null, { default: () => h(IconHeart, { class: 'fill-red-500 text-red-500' }) }) }
            ),
          ],
        }
      )
    },
  },
]
</script>
