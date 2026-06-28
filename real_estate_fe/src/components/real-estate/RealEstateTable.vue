<template>
  <div>
    <div class="mb-3 flex items-center justify-between">
      <h2 class="text-lg font-semibold">Danh sách Bất Động Sản</h2>
      <span v-if="store.total" class="text-sm text-gray-500">
        Tổng: {{ store.total }} kết quả
      </span>
    </div>

    <ag-grid-vue
      class="ag-theme-alpine w-full"
      style="height: 500px"
      :column-defs="columnDefs"
      :row-data="store.items"
      :pagination="true"
      :pagination-page-size="store.pageSize"
      :loading="store.loading"
      :default-col-def="defaultColDef"
      suppress-cell-focus
      @grid-ready="onGridReady"
    />
  </div>
</template>

<script setup lang="ts">
import { AgGridVue } from 'ag-grid-vue3'
import type { ColDef, GridApi, GridReadyEvent } from 'ag-grid-community'
import { useRealEstateStore } from '@/stores/real_estate'

const store = useRealEstateStore()

let gridApi: GridApi | null = null

function onGridReady(params: GridReadyEvent) {
  gridApi = params.api
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
