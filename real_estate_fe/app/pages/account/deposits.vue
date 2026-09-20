<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-6xl mx-auto px-4 py-6 flex flex-col gap-4">
      <!-- Tiêu đề + bộ lọc -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
        <div class="flex flex-col">
          <h1 class="text-xl font-bold text-gray-900">Đơn đặt cọc xem nhà của tôi</h1>
          <p class="text-sm text-gray-500">
            Tiền cọc được platform giữ an toàn và chỉ release khi buổi xem nhà có kết quả rõ ràng
          </p>
        </div>
        <div class="w-full md:w-64">
          <n-select v-model:value="statusFilter" :options="statusOptions" placeholder="Tất cả trạng thái" clearable />
        </div>
      </div>

      <!-- Thống kê nhanh -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div v-for="card in summaryCards" :key="card.label"
          class="bg-white rounded-lg border border-gray-200 p-3 flex flex-col gap-1">
          <span class="text-xs text-gray-400">{{ card.label }}</span>
          <span class="text-lg font-bold" :class="card.className">{{ card.value }}</span>
        </div>
      </div>

      <!-- Danh sách -->
      <div class="bg-white rounded-lg border border-gray-200 p-4 flex flex-col gap-4 h-[calc(100vh-320px)] min-h-[420px]">
        <DepositTable :items="deposits" :loading="loading" role="CUSTOMER" @view="openDetail" />
        <div class="flex justify-end flex-shrink-0">
          <Pagination :current-page="page" :total-pages="totalPages" @page-change="goToPage" />
        </div>
      </div>
    </div>

    <DepositDetailModal v-model:show="showDetail" :deposit-id="selectedId" role="CUSTOMER" @changed="fetchData" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import type { Deposit, DepositStatus } from '~/types/deposit'
import { useDepositService } from '~/services/deposit.service'
import { DEPOSIT_STATUS_META } from '~/types/deposit'
import { formatVnd } from '~/utils/deposit'
import DepositTable from '~/components/deposit/DepositTable.vue'
import DepositDetailModal from '~/components/deposit/DepositDetailModal.vue'

// Chặn ở tầng route: chỉ user có quyền xem đơn của mình mới vào được.
// Tên file đặt tiếng Anh; giữ alias tiếng Việt cho URL người dùng đã quen.
definePageMeta({
  requiresPermission: ['deposit.view.own'],
  alias: ['/tai-khoan/dat-coc'],
})

const depositService = useDepositService()

const deposits = ref<Deposit[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const loading = ref(false)
const statusFilter = ref<DepositStatus | null>(null)
const showDetail = ref(false)
const selectedId = ref<number | null>(null)

const totalPages = computed(() => Math.ceil(total.value / size.value) || 1)

const statusOptions = Object.entries(DEPOSIT_STATUS_META).map(([value, meta]) => ({
  label: meta.label,
  value,
}))

const summaryCards = computed(() => {
  const holding = deposits.value
    .filter((item) => ['PENDING', 'BROKER_CONFIRMED', 'CHECKED_IN', 'DISPUTE'].includes(item.status))
    .reduce((sum, item) => sum + item.amount, 0)
  const refunded = deposits.value.reduce((sum, item) => sum + (item.refund_amount ?? 0), 0)
  const waiting = deposits.value.filter((item) => item.status === 'PENDING').length

  return [
    { label: 'Tổng số đơn', value: String(total.value), className: 'text-gray-900' },
    { label: 'Tiền đang được giữ', value: formatVnd(holding), className: 'text-amber-600' },
    { label: 'Đã hoàn cho bạn', value: formatVnd(refunded), className: 'text-emerald-600' },
    { label: 'Chờ môi giới xác nhận', value: String(waiting), className: 'text-blue-600' },
  ]
})

async function fetchData() {
  loading.value = true
  try {
    const result = await depositService.getMyDeposits({
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
    loading.value = false
  }
}

function goToPage(nextPage: number) {
  if (nextPage < 1 || nextPage > totalPages.value) return
  page.value = nextPage
  fetchData()
}

function openDetail(id: number) {
  selectedId.value = id
  showDetail.value = true
}

watch(statusFilter, () => {
  page.value = 1
  fetchData()
})

onMounted(fetchData)
</script>
