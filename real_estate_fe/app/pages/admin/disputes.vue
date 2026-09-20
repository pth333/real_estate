<template>
  <div class="flex flex-col gap-5">
    <n-result v-if="!isAdmin" status="403" title="Không có quyền truy cập"
      description="Chỉ tài khoản quản trị viên mới xử lý được tranh chấp" />

    <template v-else>
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
        <div>
          <h1 class="text-lg font-bold text-gray-900">Tranh chấp cần xử lý</h1>
          <p class="text-sm text-gray-500">
            Tiền cọc của các đơn dưới đây đang bị tạm giữ. Chỉ release sau khi admin ra quyết định.
          </p>
        </div>
        <div class="w-full md:w-56">
          <n-select v-model:value="statusFilter" :options="statusOptions" placeholder="Tất cả trạng thái" clearable />
        </div>
      </div>

      <div class="bg-white rounded-lg border border-gray-200 p-4 flex flex-col gap-4 min-h-[420px]">
        <n-spin :show="loading">
          <n-empty v-if="!disputes.length && !loading" description="Không có tranh chấp nào" class="py-16" />

          <div v-else class="flex flex-col gap-3">
            <div v-for="dispute in disputes" :key="dispute.id"
              class="rounded-lg border border-gray-200 p-4 flex flex-col gap-3">
              <!-- Dòng tiêu đề -->
              <div class="flex flex-wrap items-center justify-between gap-2">
                <div class="flex items-center gap-2">
                  <span class="font-semibold text-gray-800">Tranh chấp #{{ dispute.id }}</span>
                  <n-tag size="small" :type="disputeStatusType(dispute.status)" :bordered="false">
                    {{ disputeStatusLabel(dispute.status) }}
                  </n-tag>
                  <n-tag size="small" :bordered="false" type="default">
                    Bên mở: {{ raisedByLabel(dispute.raised_by) }}
                  </n-tag>
                </div>
                <span class="text-xs text-gray-400">{{ fromNow(dispute.created_at) }}</span>
              </div>

              <!-- Đơn liên quan -->
              <div v-if="dispute.deposit" class="grid grid-cols-2 md:grid-cols-4 gap-3 text-sm">
                <div class="flex flex-col">
                  <span class="text-xs text-gray-400">Đơn đặt cọc</span>
                  <span class="text-gray-800 font-medium">#{{ dispute.deposit.id }}</span>
                </div>
                <div class="flex flex-col">
                  <span class="text-xs text-gray-400">Bất động sản</span>
                  <span class="text-gray-800 truncate">{{ dispute.deposit.real_estate_title || '—' }}</span>
                </div>
                <div class="flex flex-col">
                  <span class="text-xs text-gray-400">Tiền cọc</span>
                  <span class="text-gray-800 font-medium">{{ formatVnd(dispute.deposit.amount) }}</span>
                </div>
                <div class="flex flex-col">
                  <span class="text-xs text-gray-400">Lịch xem</span>
                  <span class="text-gray-800">{{ formatSlot(dispute.deposit) }}</span>
                </div>
              </div>

              <!-- Lý do -->
              <p class="text-sm text-gray-700">
                <strong>Lý do:</strong> {{ dispute.reason }}
              </p>

              <!-- Báo cáo + bằng chứng 2 bên (căn cứ để admin phán) -->
              <div v-if="dispute.deposit?.broker_report || dispute.deposit?.customer_report"
                class="rounded-lg bg-gray-50 border border-gray-100 p-3 flex flex-col gap-3">
                <span class="text-xs font-semibold uppercase tracking-wide text-gray-400">
                  Báo cáo &amp; bằng chứng hai bên
                </span>
                <div class="grid grid-cols-2 gap-3">
                  <div class="flex flex-col gap-1.5">
                    <span class="text-sm text-gray-700">
                      Môi giới: <strong>{{ reportLabel(dispute.deposit.broker_report, reportPhaseCheckedIn(dispute.deposit)) }}</strong>
                    </span>
                    <div v-if="dispute.deposit.broker_report_evidence?.length" class="flex flex-wrap gap-1.5">
                      <n-image v-for="(url, index) in dispute.deposit.broker_report_evidence" :key="`be-${index}`"
                        :src="url" width="64" height="64" class="rounded object-cover border border-gray-200" />
                    </div>
                    <span v-else class="text-xs text-gray-400">Không kèm bằng chứng</span>
                  </div>
                  <div class="flex flex-col gap-1.5">
                    <span class="text-sm text-gray-700">
                      Khách: <strong>{{ reportLabel(dispute.deposit.customer_report, reportPhaseCheckedIn(dispute.deposit)) }}</strong>
                    </span>
                    <div v-if="dispute.deposit.customer_report_evidence?.length" class="flex flex-wrap gap-1.5">
                      <n-image v-for="(url, index) in dispute.deposit.customer_report_evidence" :key="`ce-${index}`"
                        :src="url" width="64" height="64" class="rounded object-cover border border-gray-200" />
                    </div>
                    <span v-else class="text-xs text-gray-400">Không kèm bằng chứng</span>
                    <span v-if="dispute.deposit.customer_report === 'BOUGHT'" class="text-xs text-emerald-700">
                      Tài liệu khách xuất trình:
                      <strong>{{ PURCHASE_PROOF_LABEL[dispute.deposit.customer_purchase_proof] || 'chưa chọn' }}</strong>
                    </span>
                  </div>
                </div>
              </div>

              <!-- Bằng chứng -->
              <div v-if="dispute.evidence_urls?.length" class="flex flex-wrap gap-2">
                <n-image v-for="(url, index) in dispute.evidence_urls" :key="index" :src="url" width="72" height="72"
                  class="rounded object-cover border border-gray-100" />
              </div>
              <span v-else class="text-xs text-gray-400">Chưa có bằng chứng nào được gửi lên</span>

              <span v-if="dispute.evidence_deadline" class="text-xs text-amber-600">
                Hạn upload bằng chứng: {{ formatDate(dispute.evidence_deadline) }}
              </span>

              <!-- Hành động -->
              <div class="flex justify-end gap-2">
                <n-button size="small" @click="openDetail(dispute.deposit_id)">Xem đơn đặt cọc</n-button>
                <n-button v-if="dispute.status !== 'RESOLVED'" size="small" type="primary" @click="openResolve(dispute)">
                  Ra quyết định
                </n-button>
                <n-tag v-else size="small" type="success" :bordered="false">
                  {{ resolutionLabel(dispute.resolution) }}
                </n-tag>
              </div>
            </div>
          </div>
        </n-spin>
      </div>

      <div class="flex justify-end">
        <Pagination :current-page="page" :total-pages="totalPages" @page-change="goToPage" />
      </div>
    </template>

    <!-- Modal ra quyết định -->
    <n-modal v-model:show="showResolve" :mask-closable="false">
      <div class="w-[520px] max-w-[94vw] bg-white rounded-xl p-6 flex flex-col gap-4">
        <span class="font-semibold text-gray-800">Quyết định xử lý tranh chấp #{{ resolving?.id }}</span>

        <div v-if="resolving?.deposit" class="rounded-lg bg-gray-50 p-3 text-sm text-gray-600">
          Tiền cọc đang giữ: <strong>{{ formatVnd(resolving.deposit.amount) }}</strong>
        </div>

        <n-radio-group v-model:value="resolution" class="flex flex-col gap-2">
          <n-radio value="REFUND_CUSTOMER">
            <div class="flex flex-col">
              <span class="text-sm text-gray-800">Hoàn toàn bộ cho khách</span>
              <span class="text-xs text-gray-400">Khách nhận lại 100% tiền cọc</span>
            </div>
          </n-radio>
          <n-radio value="TRANSFER_BROKER">
            <div class="flex flex-col">
              <span class="text-sm text-gray-800">Chuyển toàn bộ cho môi giới</span>
              <span class="text-xs text-gray-400">Môi giới nhận 100% tiền cọc</span>
            </div>
          </n-radio>
          <n-radio value="SPLIT">
            <div class="flex flex-col">
              <span class="text-sm text-gray-800">Chia tiền cho cả 2 bên</span>
              <span class="text-xs text-gray-400">Nhập % hoàn cho khách, phần còn lại chuyển môi giới</span>
            </div>
          </n-radio>
        </n-radio-group>

        <n-form-item v-if="resolution === 'SPLIT'" label="% tiền cọc hoàn cho khách" :show-feedback="false">
          <n-slider v-model:value="splitPercent" :min="0" :max="100" :step="5" :marks="{ 0: '0%', 50: '50%', 100: '100%' }" />
        </n-form-item>

        <n-input v-model:value="note" type="textarea" :rows="3" placeholder="Ghi chú gửi 2 bên (không bắt buộc)" />

        <div class="flex justify-end gap-2">
          <n-button @click="showResolve = false">Huỷ</n-button>
          <n-button type="primary" :loading="processing" @click="submitResolve">Xác nhận &amp; release tiền</n-button>
        </div>
      </div>
    </n-modal>

    <DepositDetailModal v-model:show="showDetail" :deposit-id="selectedId" role="ADMIN" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import type { Deposit, DepositDispute, DisputeResolution } from '~/types/deposit'
import { PURCHASE_PROOF_LABEL } from '~/types/deposit'
import { useDepositService } from '~/services/deposit.service'
import { useAuthStore } from '~/stores/auth'
import { formatSlot, formatVnd, reportLabel } from '~/utils/deposit'
import { formatDate, fromNow } from '~/utils/format'
import DepositDetailModal from '~/components/deposit/DepositDetailModal.vue'

definePageMeta({
  layout: 'admin',
})

const authStore = useAuthStore()
const depositService = useDepositService()

const isAdmin = computed(() => authStore.user?.role === 'ADMIN')

const disputes = ref<DepositDispute[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const loading = ref(false)
const statusFilter = ref<string | null>(null)

const showDetail = ref(false)
const selectedId = ref<number | null>(null)

const showResolve = ref(false)
const resolving = ref<DepositDispute | null>(null)
const resolution = ref<DisputeResolution>('REFUND_CUSTOMER')
const splitPercent = ref(50)
const note = ref('')
const processing = ref(false)

const totalPages = computed(() => Math.ceil(total.value / size.value) || 1)

const statusOptions = [
  { label: 'Chờ tiếp nhận', value: 'OPEN' },
  { label: 'Đang xem xét', value: 'REVIEWING' },
  { label: 'Đã xử lý', value: 'RESOLVED' },
]

function disputeStatusType(status: string) {
  if (status === 'RESOLVED') return 'success'
  if (status === 'REVIEWING') return 'info'
  return 'warning'
}

function disputeStatusLabel(status: string) {
  if (status === 'RESOLVED') return 'Đã xử lý'
  if (status === 'REVIEWING') return 'Đang xem xét'
  return 'Chờ tiếp nhận'
}

function raisedByLabel(raisedBy: string) {
  if (raisedBy === 'CUSTOMER') return 'Khách hàng'
  if (raisedBy === 'BROKER') return 'Môi giới'
  return 'Hệ thống'
}

function resolutionLabel(value: string) {
  if (value === 'REFUND_CUSTOMER') return 'Hoàn tiền cho khách'
  if (value === 'TRANSFER_BROKER') return 'Chuyển tiền cho môi giới'
  if (value === 'SPLIT') return 'Chia tiền cho cả 2 bên'
  return value
}

// Giai đoạn báo cáo: đã check-in thì báo mua/không mua, chưa thì báo điểm danh
function reportPhaseCheckedIn(deposit: Deposit): boolean {
  return deposit.status === 'CHECKED_IN' || deposit.status.startsWith('VISITED')
}

async function fetchDisputes() {
  loading.value = true
  try {
    const result = await depositService.getDisputes({
      status: statusFilter.value ?? undefined,
      page: page.value,
      size: size.value,
    })
    disputes.value = result.items
    total.value = result.total
  } catch {
    disputes.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function goToPage(nextPage: number) {
  if (nextPage < 1 || nextPage > totalPages.value) return
  page.value = nextPage
  fetchDisputes()
}

function openDetail(depositId: number) {
  selectedId.value = depositId
  showDetail.value = true
}

function openResolve(dispute: DepositDispute) {
  resolving.value = dispute
  resolution.value = 'REFUND_CUSTOMER'
  splitPercent.value = 50
  note.value = ''
  showResolve.value = true
}

async function submitResolve() {
  if (!resolving.value) return
  processing.value = true
  try {
    await depositService.resolveDispute(resolving.value.id, {
      resolution: resolution.value,
      note: note.value,
      split_customer_percent: resolution.value === 'SPLIT' ? splitPercent.value : 0,
    })
    window.message?.success('Đã release tiền theo quyết định và thông báo cho 2 bên')
    showResolve.value = false
    await fetchDisputes()
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Xử lý tranh chấp thất bại'
    window.message?.error(message)
  } finally {
    processing.value = false
  }
}

watch(statusFilter, () => {
  page.value = 1
  fetchDisputes()
})

onMounted(() => {
  if (isAdmin.value) fetchDisputes()
})
</script>
