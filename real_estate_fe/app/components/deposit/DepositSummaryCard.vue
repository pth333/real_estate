<template>
  <div class="flex flex-col gap-4">
    <!-- Bất động sản -->
    <div class="flex items-center gap-3 rounded-lg border border-gray-100 p-3">
      <n-image v-if="deposit.real_estate_thumbnail" :src="deposit.real_estate_thumbnail" width="72" height="54"
        class="rounded object-cover border border-gray-100 flex-shrink-0" />
      <div class="flex flex-col gap-1 min-w-0">
        <span class="font-semibold text-gray-800 text-sm truncate">{{ deposit.real_estate_title || 'Bất động sản' }}</span>
        <span class="text-xs text-gray-500 truncate">
          {{ deposit.real_estate_address || 'Chưa cập nhật địa chỉ' }}
        </span>
        <span class="text-xs text-emerald-600 font-medium">
          Lịch xem: {{ formatSlot(deposit) }}
        </span>
      </div>
    </div>

    <!-- Hai bên tham gia -->
    <div class="grid grid-cols-2 gap-3">
      <div class="rounded-lg border border-gray-100 p-3 flex flex-col gap-1">
        <span class="text-xs font-semibold uppercase tracking-wide text-gray-400">Khách hàng</span>
        <span class="text-sm font-medium text-gray-800">{{ deposit.customer_name || '—' }}</span>
        <span class="text-xs text-gray-500">{{ deposit.customer_phone || '—' }}</span>
        <span class="text-xs text-gray-500 truncate">{{ deposit.customer_email || '—' }}</span>
      </div>
      <div class="rounded-lg border border-gray-100 p-3 flex flex-col gap-1">
        <span class="text-xs font-semibold uppercase tracking-wide text-gray-400">Môi giới</span>
        <span class="text-sm font-medium text-gray-800">{{ deposit.broker_name || '—' }}</span>
        <span class="text-xs text-gray-500">{{ deposit.broker_phone || '—' }}</span>
        <span class="text-xs text-gray-500 truncate">{{ deposit.broker_email || '—' }}</span>
      </div>
    </div>

    <!-- Dòng tiền -->
    <n-descriptions :column="2" size="small" bordered label-placement="left">
      <n-descriptions-item label="Phí môi giới">
        <span class="font-semibold text-gray-800">{{ formatVnd(deposit.broker_fee) }}</span>
      </n-descriptions-item>
      <!-- Cột chừa sẵn cho tiền đặt cọc — nghiệp vụ cọc làm sau nên hiện chưa có số -->
      <n-descriptions-item label="Đặt cọc">
        <span class="text-gray-400">—</span>
      </n-descriptions-item>
      <n-descriptions-item label="Đã hoàn cho khách">
        <span :class="deposit.refund_amount ? 'text-emerald-600 font-semibold' : 'text-gray-400'">
          {{ deposit.refund_amount === null ? 'Chưa tất toán' : formatVnd(deposit.refund_amount) }}
        </span>
      </n-descriptions-item>
      <n-descriptions-item label="Phạt môi giới">
        <span :class="deposit.penalty_amount ? 'text-red-500 font-semibold' : 'text-gray-400'">
          {{ deposit.penalty_amount === null ? '—' : formatVnd(deposit.penalty_amount) }}
        </span>
      </n-descriptions-item>
    </n-descriptions>

    <!-- Thanh toán -->
    <n-descriptions :column="2" size="small" bordered label-placement="left">
      <n-descriptions-item label="Cổng thanh toán">
        {{ PAYMENT_METHOD_LABEL[deposit.payment_method] || deposit.payment_method || '—' }}
      </n-descriptions-item>
      <n-descriptions-item label="Mã giao dịch">
        <span class="text-xs">{{ deposit.payment_ref || '—' }}</span>
      </n-descriptions-item>
      <n-descriptions-item label="Thời điểm thanh toán">
        {{ deposit.paid_at ? fromNow(deposit.paid_at) : 'Chưa thanh toán' }}
      </n-descriptions-item>
      <n-descriptions-item label="Xác nhận có mặt">
        <div class="flex flex-col gap-0.5 text-xs">
          <span>Môi giới: {{ checkinSideLabel(deposit.broker_checkin, deposit.broker_checkin_at) }}</span>
          <span>Khách: {{ checkinSideLabel(deposit.customer_checkin, deposit.customer_checkin_at) }}</span>
          <!-- Bằng chứng vị trí: chỉ là căn cứ đối chiếu, KHÔNG tự chuyển trạng thái đã check-in -->
          <span v-if="deposit.checkin_matched" class="font-semibold text-emerald-600">
            {{ deposit.status === 'CHECKED_IN' ? 'Hai bên đã gặp nhau' : 'Vị trí 2 bên đã trùng nhau' }}
            (cách {{ formatDistance(deposit.checkin_distance_meters) }})
            <template v-if="deposit.status !== 'CHECKED_IN'">
              — cần khách nhập mã OTP để xác nhận đã check-in
            </template>
          </span>
          <span v-else-if="deposit.checkin_distance_meters !== null" class="text-amber-600">
            Vị trí 2 bên lệch {{ formatDistance(deposit.checkin_distance_meters) }}
          </span>
        </div>
      </n-descriptions-item>
    </n-descriptions>

    <!-- Báo cáo kết quả + bằng chứng từng bên -->
    <div v-if="deposit.broker_report || deposit.customer_report" class="rounded-lg border border-gray-100 p-3 flex flex-col gap-3">
      <span class="text-xs font-semibold uppercase tracking-wide text-gray-400">Báo cáo kết quả</span>

      <div class="grid grid-cols-2 gap-3 text-sm">
        <div class="flex flex-col gap-1">
          <span class="text-gray-600">
            Môi giới: <strong class="text-gray-800">{{ reportLabel(deposit.broker_report, isCheckedIn) }}</strong>
          </span>
          <div v-if="deposit.broker_report_evidence?.length" class="flex flex-wrap gap-1.5">
            <n-image v-for="(url, index) in deposit.broker_report_evidence" :key="`b-${index}`" :src="url" width="56"
              height="56" class="rounded object-cover border border-gray-100" />
          </div>
          <span v-else-if="deposit.broker_report" class="text-xs text-gray-400">Không kèm bằng chứng</span>
        </div>

        <div class="flex flex-col gap-1">
          <span class="text-gray-600">
            Khách: <strong class="text-gray-800">{{ reportLabel(deposit.customer_report, isCheckedIn) }}</strong>
          </span>
          <div v-if="deposit.customer_report_evidence?.length" class="flex flex-wrap gap-1.5">
            <n-image v-for="(url, index) in deposit.customer_report_evidence" :key="`c-${index}`" :src="url" width="56"
              height="56" class="rounded object-cover border border-gray-100" />
          </div>
          <span v-else-if="deposit.customer_report" class="text-xs text-gray-400">Không kèm bằng chứng</span>
          <span v-if="deposit.customer_report === 'BOUGHT'" class="text-xs text-emerald-700">
            Tài liệu: {{ PURCHASE_PROOF_LABEL[deposit.customer_purchase_proof] || '—' }}
          </span>
        </div>
      </div>

      <span v-if="deposit.report_deadline" class="text-xs text-gray-400">
        Hạn giữ phí: {{ formatDate(deposit.report_deadline) }} — khách đặt cọc mua BĐS trước mốc này thì được hoàn 100% phí, quá hạn thì phí thuộc về môi giới.
      </span>
    </div>

    <!-- Thời điểm khách đặt cọc mua BĐS (nếu đã đặt cọc) -->
    <div v-if="deposit.purchase_deposit_at" class="rounded-lg bg-emerald-50 border border-emerald-100 p-3 text-sm text-emerald-700">
      Khách đã đặt cọc mua bất động sản lúc {{ fromNow(deposit.purchase_deposit_at) }} — phí môi giới được hoàn 100% cho khách.
    </div>

    <!-- Lý do môi giới từ chối -->
    <div v-if="deposit.reject_reason" class="rounded-lg bg-red-50 border border-red-100 p-3 text-sm text-red-700">
      <strong>Lý do:</strong> {{ deposit.reject_reason }}
    </div>

    <!-- Tranh chấp -->
    <div v-if="deposit.dispute" class="rounded-lg bg-amber-50 border border-amber-100 p-3 flex flex-col gap-2">
      <div class="flex items-center justify-between gap-2">
        <span class="text-sm font-semibold text-amber-800">Tranh chấp #{{ deposit.dispute.id }}</span>
        <n-tag size="small" :type="disputeTagType" :bordered="false">
          {{ disputeStatusLabel }}
        </n-tag>
      </div>
      <p class="text-sm text-amber-800">
        <strong>Bên mở:</strong> {{ disputeRaisedByLabel }} — {{ deposit.dispute.reason }}
      </p>
      <p v-if="deposit.dispute.evidence_deadline" class="text-xs text-amber-700">
        Hạn upload bằng chứng: {{ formatDate(deposit.dispute.evidence_deadline) }}
      </p>
      <div v-if="deposit.dispute.evidence_urls?.length" class="flex flex-wrap gap-2">
        <n-image v-for="(url, index) in deposit.dispute.evidence_urls" :key="index" :src="url" width="70" height="70"
          class="rounded object-cover border border-amber-100" />
      </div>
      <p v-if="deposit.dispute.resolution" class="text-sm text-amber-900">
        <strong>Quyết định của admin:</strong> {{ disputeResolutionLabel }}
      </p>
    </div>

    <!-- Đánh giá -->
    <div v-if="deposit.rating" class="rounded-lg border border-gray-100 p-3 flex flex-col gap-1">
      <span class="text-xs font-semibold uppercase tracking-wide text-gray-400">Đánh giá môi giới</span>
      <n-rate readonly :value="deposit.rating.rating" size="small" />
      <span v-if="deposit.rating.comment" class="text-sm text-gray-600">{{ deposit.rating.comment }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Deposit } from '~/types/deposit'
import { PAYMENT_METHOD_LABEL, PURCHASE_PROOF_LABEL } from '~/types/deposit'
import { formatDate, fromNow } from '~/utils/format'
import { formatSlot, formatVnd, reportLabel, formatDistance } from '~/utils/deposit'
import { NDescriptions, NDescriptionsItem, NImage, NRate, NTag } from 'naive-ui'

const props = defineProps<{ deposit: Deposit }>()

const isCheckedIn = computed(
  () => props.deposit.status === 'CHECKED_IN' || props.deposit.status.startsWith('VISITED'),
)

/** "đã xác nhận 10:05" / "chưa xác nhận" cho từng bên */
function checkinSideLabel(done: boolean, at: string): string {
  if (!done) return 'chưa xác nhận'
  if (!at) return 'đã xác nhận'
  const time = new Date(at).toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' })
  return `đã xác nhận lúc ${time}`
}

const disputeTagType = computed(() => {
  switch (props.deposit.dispute?.status) {
    case 'RESOLVED':
      return 'success'
    case 'REVIEWING':
      return 'info'
    default:
      return 'warning'
  }
})

const disputeStatusLabel = computed(() => {
  switch (props.deposit.dispute?.status) {
    case 'RESOLVED':
      return 'Đã xử lý'
    case 'REVIEWING':
      return 'Admin đang xem xét'
    default:
      return 'Chờ admin tiếp nhận'
  }
})

const disputeRaisedByLabel = computed(() => {
  switch (props.deposit.dispute?.raised_by) {
    case 'CUSTOMER':
      return 'Khách hàng'
    case 'BROKER':
      return 'Môi giới'
    default:
      return 'Hệ thống'
  }
})

const disputeResolutionLabel = computed(() => {
  switch (props.deposit.dispute?.resolution) {
    case 'REFUND_CUSTOMER':
      return 'Hoàn tiền cho khách'
    case 'TRANSFER_BROKER':
      return 'Chuyển tiền cho môi giới'
    case 'SPLIT':
      return 'Chia tiền cho cả 2 bên'
    default:
      return '—'
  }
})
</script>
