<template>
  <div>
    <div class="mb-3 flex items-center justify-between">
      <h2 class="text-lg font-semibold">Danh sách Bất Động Sản</h2>
      <span v-if="realEstateStore.total" class="text-sm text-gray-500">
        Trang {{ realEstateStore.page }}/{{ pageCount }} - Tổng: {{ realEstateStore.total }}
      </span>
    </div>
    <ag-grid-vue
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
    field: 'Title',
    flex: 2,
    minWidth: 250,
    cellRenderer: (params: { data: { SourceURL: string; Title: string } }) => {
      return `<a href="${params.data.SourceURL}" target="_blank" class="text-blue-600 hover:underline">${params.data.Title}</a>`
    },
  },
  {
    headerName: 'Giá (tỷ)',
    field: 'PriceVND',
    width: 130,
    valueFormatter: (params: { value: number }) => {
      if (params.value == null) return ''
      return (params.value / 1_000_000_000).toFixed(2)
    },
    filter: 'agNumberColumnFilter',
  },
  {
    headerName: 'Diện tích (m²)',
    field: 'Acreage',
    width: 140,
    valueFormatter: (params: { value: number }) => {
      if (params.value == null) return ''
      return params.value.toLocaleString() + ' m²'
    },
    filter: 'agNumberColumnFilter',
  },
  {
    headerName: 'Giá / m²',
    field: 'PricePerM2',
    width: 130,
    valueFormatter: (params: { value: number }) => {
      if (params.value == null) return ''
      return (params.value / 1_000_000).toFixed(1) + ' tr'
    },
  },
  {
    headerName: 'Địa chỉ',
    field: 'Address',
    flex: 1,
    minWidth: 200,
  },
  {
    headerName: 'Quận',
    field: 'District',
    width: 120,
  },
  {
    headerName: 'Nguồn',
    field: 'Source',
    width: 150,
  },
  {
    headerName: 'Ngày crawl',
    field: 'CrawledAt',
    width: 160,
    valueFormatter: (params: { value: string }) => {
      if (!params.value) return ''
      return new Date(params.value).toLocaleDateString('vi-VN')
    },
  },
]
</script>
