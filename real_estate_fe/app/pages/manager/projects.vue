<template>
  <div class="flex-1 bg-white p-6 rounded-lg border border-gray-200 flex flex-col gap-4 h-full min-h-0">
    <!-- Header: Tìm kiếm + nút tạo dự án -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 flex-shrink-0">
      <div class="w-full md:w-80">
        <n-input v-model:value="searchQuery" placeholder="Tìm kiếm theo tên dự án..." clearable @input="handleSearch">
          <template #prefix>
            <n-icon>
              <IconSearch />
            </n-icon>
          </template>
        </n-input>
      </div>

      <n-button type="error" @click="goToCreateProject">
        <template #icon>
          <n-icon>
            <IconAddOutline />
          </n-icon>
        </template>
        Tạo dự án
      </n-button>
    </div>

    <!-- Data Table danh sách dự án -->
    <div class="flex-1 min-h-0 flex flex-col">
      <n-data-table :columns="columns" :data="projects" :loading="loading" :bordered="false" :single-line="false"
        flex-height class="flex-1 min-h-0" :scroll-x="760" />
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex justify-end flex-shrink-0">
      <Pagination :current-page="managerStore.projectsPage" :total-pages="totalPages" @page-change="goToPage" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { NImage, NTooltip, NButton, NIcon, NSpace, NTag, type DataTableColumns } from 'naive-ui'
import { formatPrice } from '~/utils/format'
import type { ManagerProject } from '~/types/manager'
import { useManagerStore } from '~/stores/manager'
import IconSearch from '~/icons/IconSearch.vue'
import IconAddOutline from '~/icons/IconAddOutline.vue'
import IconEyeOutline from '~/icons/IconEyeOutline.vue'
import IconCreateOutline from '~/icons/IconCreateOutline.vue'

definePageMeta({
  alias: [
    '/nguoi-ban/quan-ly-du-an'
  ],
})

const managerStore = useManagerStore()
const loading = ref(false)

// Data + phân trang được cache trong managerStore
const projects = computed(() => managerStore.projects)
const total = computed(() => managerStore.projectsTotal)
const totalPages = computed(() => Math.ceil(managerStore.projectsTotal / managerStore.projectsSize) || 1)
const searchQuery = computed({
  get: () => managerStore.projectsSearch,
  set: (val: string) => { managerStore.projectsSearch = val },
})

const fetchProjectsData = async () => {
  loading.value = true
  try {
    await managerStore.fetchProjects({
      page: managerStore.projectsPage,
      size: managerStore.projectsSize,
      search: managerStore.projectsSearch,
    })
  } catch {
    // store đã hiện message lỗi
  } finally {
    loading.value = false
  }
}

function goToPage(page: number) {
  if (page < 1 || page > totalPages.value) return
  managerStore.projectsPage = page
  fetchProjectsData()
}

onMounted(fetchProjectsData)

let searchTimeout: any = null
const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    managerStore.projectsPage = 1
    fetchProjectsData()
  }, 300)
}

const goToCreateProject = () => {
  navigateTo('/nguoi-ban/tao-du-an')
}

const handleViewProject = (row: ManagerProject) => {
  navigateTo(`/${row.slug}-pj${row.id}`)
}

const handleEditProject = (row: ManagerProject) => {
  navigateTo(`/nguoi-ban/tao-du-an?id=${row.id}`)
}

const STATUS_LABELS: Record<string, { label: string; type: 'success' | 'info' | 'warning' | 'default' }> = {
  active: { label: 'Đang mở bán', type: 'success' },
  upcoming: { label: 'Sắp mở bán', type: 'info' },
  handed_over: { label: 'Đã bàn giao', type: 'default' },
  paused: { label: 'Tạm dừng', type: 'warning' },
}

const statusTag = (status: string) => STATUS_LABELS[status] || { label: status, type: 'default' as const }

const columns: DataTableColumns<ManagerProject> = [
  {
    title: 'Dự án',
    key: 'name',
    width: 280,
    render(row) {
      return h('div', { class: 'flex items-center gap-3 py-1' }, [
        h(NImage, {
          src: row.thumbnail || 'https://placehold.co/200x150/e2e8f0/94a3b8?text=Project',
          alt: 'Dự án',
          width: 60,
          height: 45,
          class: 'rounded object-cover border border-gray-100 flex-shrink-0',
          previewDisabled: true,
        }),
        h('div', { class: 'flex flex-col gap-1 min-w-0 flex-1' }, [
          h(
            NTooltip,
            { trigger: 'hover', placement: 'top' },
            {
              trigger: () =>
                h('span', { class: 'font-medium text-gray-800 text-sm truncate block w-full cursor-default' }, row.name),
              default: () => row.name,
            }
          ),
          h('span', { class: 'text-xs text-gray-400 truncate' }, row.alternative_name || row.slug),
        ]),
      ])
    },
  },
  {
    title: 'Trạng thái',
    key: 'status',
    width: 120,
    render(row) {
      const tag = statusTag(row.status)
      return h(NTag, { type: tag.type, size: 'small', round: true }, { default: () => tag.label })
    },
  },
  {
    title: 'Vị trí',
    key: 'full_address',
    width: 200,
    render(row) {
      return h('span', { class: 'text-gray-600 text-sm truncate block' }, row.full_address || '—')
    },
  },
  {
    title: 'Quy mô',
    key: 'scale',
    width: 130,
    render(row) {
      const parts: string[] = []
      if (row.total_units != null) parts.push(`${row.total_units} căn`)
      if (row.total_area_ha != null) parts.push(`${row.total_area_ha} ha`)
      return h('span', { class: 'text-gray-700 text-sm whitespace-nowrap' }, parts.join(' · ') || '—')
    },
  },
  {
    title: 'Giá',
    key: 'price',
    width: 130,
    render(row) {
      if (row.price_min && row.price_max) {
        return h('span', { class: 'font-semibold text-red-500 text-sm whitespace-nowrap' },
          `${formatPrice(row.price_min)} - ${formatPrice(row.price_max)}`)
      }
      return h('span', { class: 'text-gray-400 text-sm' }, '—')
    },
  },
  {
    title: 'Ngày tạo',
    key: 'created_at',
    width: 110,
    render(row) {
      return h('span', { class: 'text-gray-500 text-sm whitespace-nowrap' }, row.created_at)
    },
  },
  {
    title: 'Hành động',
    key: 'actions',
    width: 100,
    align: 'center',
    render(row) {
      return h(
        NSpace,
        { justify: 'center', size: 'small', align: 'center' },
        {
          default: () => [
            // Sửa dự án → trang tạo/chỉnh sửa dự án
            h(NButton, { size: 'small', quaternary: true, type: 'info', onClick: () => handleEditProject(row) },
              { icon: () => h(NIcon, null, { default: () => h(IconCreateOutline) }) }
            ),
            h(NButton, { size: 'small', quaternary: true, type: 'info', onClick: () => handleViewProject(row) },
              { icon: () => h(NIcon, null, { default: () => h(IconEyeOutline) }) }
            ),
          ],
        }
      )
    },
  },
]
</script>
