<template>
  <div class="flex flex-col gap-5">
    <n-result v-if="!isAdmin" status="403" title="Không có quyền truy cập"
      description="Chỉ tài khoản quản trị viên mới xem được trang này" />

    <template v-else>
      <!-- Thẻ số liệu escrow -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div v-for="card in cards" :key="card.label"
          class="bg-white rounded-lg border border-gray-200 p-4 flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span class="text-xs text-gray-400">{{ card.label }}</span>
            <n-icon :class="card.iconClass" size="18">
              <component :is="card.icon" />
            </n-icon>
          </div>
          <span class="text-xl font-bold" :class="card.valueClass">{{ card.value }}</span>
          <span class="text-xs text-gray-400">{{ card.hint }}</span>
        </div>
      </div>

      <!-- Đơn đang chờ duyệt tài liệu mua nhà — phải duyệt mới trừ tồn kho dự án -->
      <n-alert v-if="pendingApprovalCount > 0" type="warning" :bordered="false" class="cursor-pointer"
        @click="filterPendingApproval">
        Có <strong>{{ pendingApprovalCount }}</strong> đơn đang chờ duyệt tài liệu mua nhà. Duyệt xong hệ thống mới
        hoàn tiền cọc cho khách và trừ 1 căn vào tồn kho của dự án.
        <template #icon>
          <n-icon>
            <IconShieldCheck />
          </n-icon>
        </template>
      </n-alert>

      <!-- Trạng thái đơn -->
      <div class="bg-white rounded-lg border border-gray-200 p-5 flex flex-col gap-4">
        <div class="flex items-center justify-between">
          <h2 class="text-base font-semibold text-gray-800">Đơn đặt cọc theo trạng thái</h2>
          <n-button text type="primary" @click="openDisputes">Xem tranh chấp đang mở ({{ summary?.open_disputes ?? 0 }})</n-button>
        </div>
        <n-spin :show="loading">
          <div class="flex flex-wrap gap-2">
            <div v-for="item in statusCounts" :key="item.status"
              class="flex items-center gap-2 rounded-lg border border-gray-100 px-3 py-2">
              <DepositStatusTag :status="item.status" />
              <span class="text-sm font-semibold text-gray-800">{{ item.count }}</span>
            </div>
          </div>
        </n-spin>
      </div>

      <!-- Danh sách đơn gần đây -->
      <div class="bg-white rounded-lg border border-gray-200 p-5 flex flex-col gap-4 h-[520px] min-h-0">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
          <h2 class="text-base font-semibold text-gray-800">Tất cả đơn đặt cọc</h2>
          <div class="w-full md:w-64">
            <n-select v-model:value="statusFilter" :options="statusOptions" placeholder="Tất cả trạng thái" clearable />
          </div>
        </div>

        <DepositTable :items="deposits" :loading="tableLoading" role="ADMIN" @view="openDetail" />

        <div class="flex justify-end">
          <Pagination :current-page="page" :total-pages="totalPages" @page-change="goToPage" />
        </div>
      </div>
    </template>

    <DepositDetailModal v-model:show="showDetail" :deposit-id="selectedId" role="ADMIN" @changed="fetchAll" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import type { Deposit, DepositStatus, EscrowSummary } from '~/types/deposit'
import { DEPOSIT_STATUS_META } from '~/types/deposit'
import { useDepositService } from '~/services/deposit.service'
import { useAuthStore } from '~/stores/auth'
import { formatVnd } from '~/utils/deposit'
import DepositStatusTag from '~/components/deposit/DepositStatusTag.vue'
import DepositTable from '~/components/deposit/DepositTable.vue'
import DepositDetailModal from '~/components/deposit/DepositDetailModal.vue'
import IconShieldCheck from '~/icons/IconShieldCheck.vue'
import IconWallet from '~/icons/IconWallet.vue'
import IconCheck from '~/icons/IconCheck.vue'
import IconXCircle from '~/icons/IconXCircle.vue'

definePageMeta({
  layout: 'admin',
  // Chặn ở tầng route: gõ tay /admin cũng không vào được nếu không phải ADMIN
  requiresRole: ['ADMIN'],
})

const authStore = useAuthStore()
const depositService = useDepositService()

const isAdmin = computed(() => authStore.user?.roles?.includes('ADMIN') ?? false)

const summary = ref<EscrowSummary | null>(null)
const deposits = ref<Deposit[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const loading = ref(false)
const tableLoading = ref(false)
const statusFilter = ref<DepositStatus | null>(null)
const showDetail = ref(false)
const selectedId = ref<number | null>(null)

const totalPages = computed(() => Math.ceil(total.value / size.value) || 1)

const statusOptions = Object.entries(DEPOSIT_STATUS_META).map(([value, meta]) => ({
  label: meta.label,
  value,
}))

const statusCounts = computed(() => {
  const counts = summary.value?.status_counts ?? {}
  return Object.entries(DEPOSIT_STATUS_META)
    .map(([status]) => ({ status: status as DepositStatus, count: counts[status] ?? 0 }))
    .filter((item) => item.count > 0)
})

// Số đơn đang chờ admin duyệt tài liệu mua nhà
const pendingApprovalCount = computed(
  () => summary.value?.status_counts?.['PENDING_PURCHASE_APPROVAL'] ?? 0,
)

function filterPendingApproval() {
  statusFilter.value = 'PENDING_PURCHASE_APPROVAL'
}

const cards = computed(() => [
  {
    label: 'Tiền đang giữ (escrow)',
    value: formatVnd(summary.value?.holding_amount ?? 0),
    hint: 'Tiền cọc chưa release cho bên nào',
    icon: IconWallet,
    iconClass: 'text-amber-500',
    valueClass: 'text-amber-600',
  },
  {
    label: 'Đã hoàn cho khách',
    value: formatVnd(summary.value?.refunded_amount ?? 0),
    hint: 'Tổng tiền đã refund',
    icon: IconCheck,
    iconClass: 'text-emerald-500',
    valueClass: 'text-emerald-600',
  },
  {
    label: 'Đã chuyển cho môi giới',
    value: formatVnd(summary.value?.transferred_amount ?? 0),
    hint: 'Cọc + phí môi giới đã release',
    icon: IconShieldCheck,
    iconClass: 'text-blue-500',
    valueClass: 'text-blue-600',
  },
  {
    label: 'Phạt thu từ môi giới',
    value: formatVnd(summary.value?.penalty_amount ?? 0),
    hint: 'Trừ vào tiền bảo lãnh',
    icon: IconXCircle,
    iconClass: 'text-red-500',
    valueClass: 'text-red-500',
  },
])

async function fetchSummary() {
  loading.value = true
  try {
    summary.value = await depositService.getEscrowSummary()
  } catch {
    summary.value = null
  } finally {
    loading.value = false
  }
}

async function fetchDeposits() {
  tableLoading.value = true
  try {
    const result = await depositService.getAllDeposits({
      status: statusFilter.value ?? undefined,
      page: page.value,
      size: size.value,
    })
    deposits.value = result.items
    total.value = result.total
  } catch {
    deposits.value = []
    total.value = 0
  } finally {
    tableLoading.value = false
  }
}

async function fetchAll() {
  await Promise.all([fetchSummary(), fetchDeposits()])
}

function goToPage(nextPage: number) {
  if (nextPage < 1 || nextPage > totalPages.value) return
  page.value = nextPage
  fetchDeposits()
}

function openDetail(id: number) {
  selectedId.value = id
  showDetail.value = true
}

function openDisputes() {
  navigateTo('/admin/disputes')
}

watch(statusFilter, () => {
  page.value = 1
  fetchDeposits()
})

onMounted(() => {
  if (isAdmin.value) fetchAll()
})
</script>
