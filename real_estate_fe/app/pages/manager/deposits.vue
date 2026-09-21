<template>
  <div class="flex h-full min-h-0 flex-col gap-4 rounded-xl border border-gray-200 bg-white p-5 lg:p-6">
    <!-- Tiêu đề trang đã có ở khu vực quản lý, ở đây chỉ giữ bộ lọc -->
    <div class="flex flex-shrink-0 flex-col justify-between gap-3 md:flex-row md:items-center">
      <span class="text-sm text-gray-500">Lọc theo trạng thái để xử lý nhanh từng nhóm đơn</span>
      <div class="w-full md:w-64">
        <n-select v-model:value="statusFilter" :options="statusOptions" placeholder="Tất cả trạng thái" clearable />
      </div>
    </div>

    <!-- Nhắc việc cần làm -->
    <n-alert v-if="pendingCount > 0" type="warning" :bordered="false" class="flex-shrink-0">
      Bạn có <strong>{{ pendingCount }}</strong> đơn chờ xác nhận. Quá 24 giờ không phản hồi, hệ thống sẽ tự động
      từ chối và hoàn 100% tiền cọc cho khách.
    </n-alert>

    <!-- Bảng -->
    <DepositTable :items="deposits" :loading="loading" role="BROKER" @view="openDetail" />

    <div class="flex justify-end flex-shrink-0">
      <Pagination :current-page="page" :total-pages="totalPages" @page-change="goToPage" />
    </div>

    <DepositDetailModal v-model:show="showDetail" :deposit-id="selectedId" role="BROKER" @changed="fetchData" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import type { Deposit, DepositStatus } from '~/types/deposit'
import { DEPOSIT_STATUS_META } from '~/types/deposit'
import { useDepositService } from '~/services/deposit.service'
import DepositTable from '~/components/deposit/DepositTable.vue'
import DepositDetailModal from '~/components/deposit/DepositDetailModal.vue'

definePageMeta({
  alias: ['/nguoi-ban/quan-ly-dat-coc'],
  // Chỉ user có quyền môi giới mới xem được danh sách đơn được giao
  requiresPermission: ['broker.deposit.list'],
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
const pendingCount = computed(() => deposits.value.filter((item) => item.status === 'PENDING').length)

const statusOptions = Object.entries(DEPOSIT_STATUS_META).map(([value, meta]) => ({
  label: meta.label,
  value,
}))

async function fetchData() {
  loading.value = true
  try {
    const result = await depositService.getBrokerDeposits({
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
