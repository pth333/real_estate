<template>
  <div>
    <div class="mb-3 flex items-center justify-between">
      <h2 class="text-lg font-semibold">Danh sách Bất Động Sản</h2>
      <span v-if="realEstateStore.total" class="text-sm text-gray-500">
        Trang {{ realEstateStore.page }}/{{ pageCount }} - Tổng: {{ realEstateStore.total }}
      </span>
    </div>

    <!-- Loading overlay -->
    <div v-if="realEstateStore.loading" class="flex items-center justify-center rounded-lg bg-white py-20 shadow-sm">
      <div class="flex flex-col items-center gap-2">
        <div class="size-8 animate-spin rounded-full border-4 border-blue-200 border-t-blue-600"></div>
        <span class="text-sm text-gray-400">Đang tải dữ liệu...</span>
      </div>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="rounded-lg bg-red-50 p-6 text-center text-sm text-red-600">
      {{ error }}
    </div>

    <ag-grid-vue
      v-else
      :theme="themeBalham"
      style="height: 500px"
      :column-defs="columnDefs"
      :row-data="realEstateStore.items"
      :pagination="false"
      :default-col-def="defaultColDef"
      suppress-cell-focus
    />
    <!-- Pagination thay thế n-pagination của Naive UI -->
    <div class="flex items-center justify-end gap-2 py-3">
      <span class="text-sm text-gray-500">Số dòng/trang:</span>
      <select
        :value="realEstateStore.pageSize"
        @change="updatePageSize(Number(($event.target as HTMLSelectElement).value))"
        class="border rounded px-2 py-1 text-sm"
      >
        <option :value="20">20</option>
        <option :value="50">50</option>
        <option :value="100">100</option>
      </select>
      <div class="flex items-center gap-1">
        <button
          v-for="p in pageCount"
          :key="p"
          :class="[
            'px-3 py-1 text-sm rounded border',
            p === realEstateStore.page
              ? 'bg-blue-500 text-white border-blue-500'
              : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-100'
          ]"
          @click="updatePage(p)"
        >
          {{ p }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { AgGridVue } from 'ag-grid-vue3'
import { themeBalham } from 'ag-grid-community'
import type { ColDef } from 'ag-grid-community'

const realEstateStore = useRealEstateStore()
const error = ref<string | null>(null)

const pageCount = computed(() => Math.ceil(realEstateStore.total / realEstateStore.pageSize))

function updatePage(page: number) {
  realEstateStore.setPage(page)
}

function updatePageSize(pageSize: number) {
  realEstateStore.setPageSize(pageSize)
}

const defaultColDef: ColDef = {
  sortable: true,
  resizable: true,
  filter: true,
}

const columnDefs: ColDef[] = [
  {
    headerName: 'Tiêu đề',
    field: 'title',
    flex: 2,
    minWidth: 250,
    cellRenderer: (params: { data: { source_url: string; title: string } }) => {
      return `<a href="${params.data.source_url}" target="_blank" class="text-blue-600 hover:underline">${params.data.title}</a>`
    },
  },
  {
    headerName: 'Giá (tỷ)',
    field: 'price_vnd',
    width: 130,
    valueFormatter: (params: { value: number }) => {
      if (params.value == null) return ''
      return (params.value / 1_000_000_000).toFixed(2)
    },
    filter: 'agNumberColumnFilter',
  },
  {
    headerName: 'Diện tích (m²)',
    field: 'acreage',
    width: 140,
    valueFormatter: (params: { value: number }) => {
      if (params.value == null) return ''
      return params.value.toLocaleString() + ' m²'
    },
    filter: 'agNumberColumnFilter',
  },
  {
    headerName: 'Giá / m²',
    field: 'price_per_m2',
    width: 130,
    valueFormatter: (params: { value: number }) => {
      if (params.value == null) return ''
      return (params.value / 1_000_000).toFixed(1) + ' tr'
    },
  },
  {
    headerName: 'Địa chỉ',
    field: 'address',
    flex: 1,
    minWidth: 200,
  },
  {
    headerName: 'Quận',
    field: 'district',
    width: 120,
  },
  {
    headerName: 'Nguồn',
    field: 'source',
    width: 150,
  },
  {
    headerName: 'Ngày crawl',
    field: 'created_at',
    width: 160,
    valueFormatter: (params: { value: string }) => {
      if (!params.value) return ''
      return new Date(params.value).toLocaleDateString('vi-VN')
    },
  },
]
</script>
