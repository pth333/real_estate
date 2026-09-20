<template>
  <div v-if="hasAnyAction" class="rounded-lg border border-emerald-100 bg-emerald-50/50 p-4 flex flex-col gap-3">
    <span class="text-xs font-semibold uppercase tracking-wide text-emerald-700">Thao tác của bạn</span>

    <n-space>
      <!-- Môi giới -->
      <n-button v-if="isBroker && deposit.can_confirm" type="primary" :loading="processing" @click="confirmDeposit">
        Xác nhận lịch
      </n-button>
      <n-button v-if="isBroker && deposit.can_confirm" type="error" ghost @click="openModal('reject')">
        Từ chối
      </n-button>
      <n-button v-if="isBroker && canGenerateOtp" type="warning" :loading="processing" @click="generateOtp">
        Sinh mã OTP check-in
      </n-button>

      <!-- Khách -->
      <n-button v-if="isCustomer && canCheckin" type="primary" @click="openModal('checkin')">
        Nhập OTP check-in
      </n-button>
      <n-button v-if="isCustomer && canRate" secondary @click="openModal('rating')">
        Đánh giá môi giới
      </n-button>

      <!-- Dùng chung -->
      <n-button v-if="canReport" type="info" ghost @click="openModal('report')">
        Báo cáo kết quả
      </n-button>
      <n-button v-if="canAddEvidence" warning ghost @click="openModal('evidence')">
        Gửi bằng chứng
      </n-button>
      <n-button v-if="canOpenDispute" type="error" ghost @click="openModal('dispute')">
        Mở tranh chấp
      </n-button>

      <!-- Admin duyệt tài liệu mua nhà -->
      <n-button v-if="canApprovePurchase" type="primary" @click="openModal('purchase')">
        Duyệt tài liệu mua nhà
      </n-button>
      <n-button v-if="canApprovePurchase" type="error" ghost @click="openModal('purchaseReject')">
        Từ chối tài liệu
      </n-button>
    </n-space>

    <!-- Từ chối lịch -->
    <n-modal v-model:show="modalVisible.reject" :mask-closable="false">
      <div class="w-[440px] bg-white rounded-xl p-6 flex flex-col gap-4">
        <span class="font-semibold text-gray-800">Từ chối lịch xem nhà</span>
        <p class="text-sm text-gray-500">
          Khách sẽ được hoàn 100% tiền cọc. Vui lòng cho khách biết lý do.
        </p>
        <n-input v-model:value="rejectReason" type="textarea" :rows="3" placeholder="VD: Bất động sản đã có khách cọc mua" />
        <div class="flex justify-end gap-2">
          <n-button @click="closeModal('reject')">Huỷ</n-button>
          <n-button type="error" :loading="processing" @click="rejectDeposit">Xác nhận từ chối</n-button>
        </div>
      </div>
    </n-modal>

    <!-- Hiển thị OTP cho môi giới -->
    <n-modal v-model:show="modalVisible.otp" :mask-closable="false">
      <div class="w-[400px] bg-white rounded-xl p-6 flex flex-col items-center gap-3">
        <span class="font-semibold text-gray-800">Mã OTP check-in</span>
        <p class="text-sm text-gray-500 text-center">
          Đọc mã này cho khách nhập vào ứng dụng để xác nhận cả 2 có mặt tại chỗ.
        </p>
        <span class="text-3xl font-bold tracking-[0.4em] text-emerald-600">{{ otpCode }}</span>
        <span class="text-xs text-gray-400">Hết hạn lúc {{ otpExpiresAt }}</span>
        <n-button class="mt-2" @click="closeModal('otp')">Đóng</n-button>
      </div>
    </n-modal>

    <!-- Khách nhập OTP -->
    <n-modal v-model:show="modalVisible.checkin" :mask-closable="false">
      <div class="w-[420px] bg-white rounded-xl p-6 flex flex-col gap-4">
        <span class="font-semibold text-gray-800">Check-in buổi xem nhà</span>
        <p class="text-sm text-gray-500">Nhập mã 6 số do môi giới hiển thị để xác nhận bạn đã gặp môi giới.</p>
        <n-input v-model:value="checkinOtp" placeholder="Nhập 6 số" :maxlength="6" size="large" />
        <div class="flex justify-end gap-2">
          <n-button @click="closeModal('checkin')">Huỷ</n-button>
          <n-button type="primary" :loading="processing" @click="submitCheckin">Xác nhận</n-button>
        </div>
      </div>
    </n-modal>

    <!-- Báo cáo kết quả -->
    <n-modal v-model:show="modalVisible.report" :mask-closable="false">
      <div class="w-[460px] bg-white rounded-xl p-6 flex flex-col gap-4">
        <span class="font-semibold text-gray-800">Báo cáo kết quả buổi xem</span>
        <n-radio-group v-model:value="reportValue" class="flex flex-col gap-2">
          <n-radio v-for="option in options" :key="option.value" :value="option.value" class="items-start">
            <div class="flex flex-col">
              <span class="text-sm text-gray-800">{{ option.label }}</span>
              <span class="text-xs text-gray-400">{{ option.description }}</span>
            </div>
          </n-radio>
        </n-radio-group>

        <!-- Báo cáo mua/không mua bắt buộc kèm bằng chứng -->
        <div v-if="evidenceRequired" class="flex flex-col gap-3">
          <!-- Khách khai đã mua phải chọn loại tài liệu chứng minh -->
          <div v-if="purchaseProofRequired" class="flex flex-col gap-2">
            <span class="text-sm font-medium text-gray-700">
              Tài liệu chứng minh đã mua <span class="text-red-500">*</span>
            </span>
            <n-select v-model:value="purchaseProofValue" :options="purchaseProofOptions"
              placeholder="Chọn loại tài liệu" />
            <span class="text-xs text-gray-400">{{ purchaseProofHint }}</span>
          </div>

          <div class="flex flex-col gap-2">
            <span class="text-sm font-medium text-gray-700">
              Ảnh tài liệu/bằng chứng <span class="text-red-500">*</span>
            </span>
            <EvidenceUploader v-model:urls="evidenceUrls" />
          </div>
        </div>

        <p class="text-xs text-amber-600">
          Báo cáo chỉ được gửi <strong>một lần</strong> và không sửa lại được. Nếu 2 bên báo cáo khác nhau,
          tiền cọc sẽ được tạm giữ và chuyển cho admin xác minh dựa trên bằng chứng.
          <strong>Bên không gửi báo cáo trong thời hạn coi như không chứng minh được và bị áp theo báo cáo của bên kia.</strong>
        </p>
        <div class="flex justify-end gap-2">
          <n-button @click="closeModal('report')">Huỷ</n-button>
          <n-button type="primary" :loading="processing" @click="submitReport">Gửi báo cáo</n-button>
        </div>
      </div>
    </n-modal>

    <!-- Mở tranh chấp -->
    <n-modal v-model:show="modalVisible.dispute" :mask-closable="false">
      <div class="w-[480px] bg-white rounded-xl p-6 flex flex-col gap-4">
        <span class="font-semibold text-gray-800">Mở tranh chấp</span>
        <p class="text-sm text-gray-500">
          Tiền cọc sẽ bị tạm giữ cho tới khi admin ra quyết định. Vui lòng mô tả rõ sự việc và gửi kèm bằng chứng.
        </p>
        <n-input v-model:value="disputeReason" type="textarea" :rows="3" placeholder="Mô tả sự việc" />
        <EvidenceUploader v-model:urls="evidenceUrls" />
        <div class="flex justify-end gap-2">
          <n-button @click="closeModal('dispute')">Huỷ</n-button>
          <n-button type="error" :loading="processing" @click="submitDispute">Mở tranh chấp</n-button>
        </div>
      </div>
    </n-modal>

    <!-- Gửi bổ sung bằng chứng -->
    <n-modal v-model:show="modalVisible.evidence" :mask-closable="false">
      <div class="w-[480px] bg-white rounded-xl p-6 flex flex-col gap-4">
        <span class="font-semibold text-gray-800">Gửi bằng chứng cho tranh chấp #{{ deposit.dispute?.id }}</span>
        <EvidenceUploader v-model:urls="evidenceUrls" />
        <div class="flex justify-end gap-2">
          <n-button @click="closeModal('evidence')">Huỷ</n-button>
          <n-button type="primary" :loading="processing" @click="submitEvidence">Gửi bằng chứng</n-button>
        </div>
      </div>
    </n-modal>

    <!-- Admin duyệt tài liệu mua nhà -->
    <n-modal v-model:show="modalVisible.purchase" :mask-closable="false">
      <div class="w-[520px] max-w-[94vw] bg-white rounded-xl p-6 flex flex-col gap-4">
        <span class="font-semibold text-gray-800">Duyệt tài liệu mua nhà</span>
        <p class="text-sm text-gray-500">
          Kiểm tra tài liệu khách xuất trình:
          <strong>{{ PURCHASE_PROOF_LABEL[deposit.customer_purchase_proof] || 'chưa chọn' }}</strong>
        </p>
        <n-alert type="warning" :bordered="false">
          Duyệt sẽ <strong>hoàn 100% tiền cọc</strong> cho khách và
          <strong>trừ 1 căn</strong> vào số căn đã bán của dự án. Hành động không hoàn tác được.
        </n-alert>
        <n-input v-model:value="purchaseNote" type="textarea" :rows="2"
          placeholder="Ghi chú duyệt (không bắt buộc)" />
        <div class="flex justify-end gap-2">
          <n-button @click="closeModal('purchase')">Huỷ</n-button>
          <n-button type="primary" :loading="processing" @click="submitPurchaseDecision(true)">
            Xác nhận duyệt
          </n-button>
        </div>
      </div>
    </n-modal>

    <!-- Admin từ chối tài liệu -->
    <n-modal v-model:show="modalVisible.purchaseReject" :mask-closable="false">
      <div class="w-[520px] max-w-[94vw] bg-white rounded-xl p-6 flex flex-col gap-4">
        <span class="font-semibold text-gray-800">Từ chối tài liệu mua nhà</span>
        <p class="text-sm text-gray-500">
          Đơn sẽ chuyển sang <strong>tranh chấp</strong> để xử lý tiếp, tiền cọc vẫn bị tạm giữ.
        </p>
        <n-input v-model:value="purchaseNote" type="textarea" :rows="3"
          placeholder="Lý do từ chối (bắt buộc, VD: hợp đồng không có chữ ký 2 bên)" />
        <div class="flex justify-end gap-2">
          <n-button @click="closeModal('purchaseReject')">Huỷ</n-button>
          <n-button type="error" :loading="processing" @click="submitPurchaseDecision(false)">
            Xác nhận từ chối
          </n-button>
        </div>
      </div>
    </n-modal>

    <!-- Đánh giá môi giới -->
    <n-modal v-model:show="modalVisible.rating" :mask-closable="false">
      <div class="w-[420px] bg-white rounded-xl p-6 flex flex-col gap-4">
        <span class="font-semibold text-gray-800">Đánh giá môi giới</span>
        <n-rate v-model:value="ratingValue" size="large" />
        <n-input v-model:value="ratingComment" type="textarea" :rows="3" placeholder="Chia sẻ trải nghiệm của bạn (không bắt buộc)" />
        <div class="flex justify-end gap-2">
          <n-button @click="closeModal('rating')">Huỷ</n-button>
          <n-button type="primary" :loading="processing" @click="submitRating">Gửi đánh giá</n-button>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { NButton, NInput, NModal, NRadio, NRadioGroup, NRate, NSpace } from 'naive-ui'
import type { Deposit, DepositActorRole } from '~/types/deposit'
import { PURCHASE_PROOF_LABEL, PURCHASE_PROOF_OPTIONS } from '~/types/deposit'
import { useDepositService } from '~/services/deposit.service'
import { reportOptions } from '~/utils/deposit'

const props = defineProps<{
  deposit: Deposit
  role: DepositActorRole
}>()

const emit = defineEmits<{
  changed: []
}>()

const depositService = useDepositService()

const isBroker = computed(() => props.role === 'BROKER')
const isCustomer = computed(() => props.role === 'CUSTOMER')
const isAdmin = computed(() => props.role === 'ADMIN')

const processing = ref(false)
const modalVisible = reactive<Record<string, boolean>>({})
const evidenceUrls = ref<string[]>([])

const rejectReason = ref('')
const otpCode = ref('')
const otpExpiresAt = ref('')
const checkinOtp = ref('')
const reportValue = ref('')
const purchaseProofValue = ref<string | null>(null)
const disputeReason = ref('')
const ratingValue = ref(5)
const ratingComment = ref('')
const purchaseNote = ref('')

const options = computed(() => reportOptions(props.deposit, isBroker.value))

// Báo cáo mua/không mua (đã check-in) bắt buộc kèm ảnh; báo điểm danh thì không
const evidenceRequired = computed(() => props.deposit.status === 'CHECKED_IN')

// Khách khai BOUGHT phải kèm tài liệu mua bán (bên duy nhất hưởng lợi khi khai BOUGHT)
const purchaseProofRequired = computed(
  () => isCustomer.value && evidenceRequired.value && reportValue.value === 'BOUGHT',
)

const purchaseProofOptions = PURCHASE_PROOF_OPTIONS.map((item) => ({
  label: item.label,
  value: item.value,
}))

const purchaseProofHint = computed(
  () => PURCHASE_PROOF_OPTIONS.find((item) => item.value === purchaseProofValue.value)?.hint ?? '',
)

const canGenerateOtp = computed(() => isBroker.value && props.deposit.can_checkin && !props.deposit.broker_checkin)
const canCheckin = computed(() => isCustomer.value && props.deposit.can_checkin)
const canReport = computed(() => {
  if (!isCustomer.value && !isBroker.value) return false
  if (!props.deposit.can_report) return false
  return isBroker.value ? !props.deposit.broker_report : !props.deposit.customer_report
})
const canOpenDispute = computed(() => {
  if (props.deposit.has_dispute || !isCustomer.value && !isBroker.value) return false
  return !['AWAITING_PAYMENT', 'BROKER_REJECTED', 'REFUNDED', 'COMPLETED', 'CANCELLED', 'DISPUTE'].includes(
    props.deposit.status,
  )
})
const canAddEvidence = computed(() => {
  if (!isCustomer.value && !isBroker.value) return false
  const dispute = props.deposit.dispute
  if (!dispute || dispute.status === 'RESOLVED') return false
  return !dispute.evidence_deadline || new Date(dispute.evidence_deadline).getTime() > Date.now()
})
const canRate = computed(
  () => isCustomer.value && !props.deposit.has_rating && props.deposit.status.startsWith('VISITED'),
)

// Admin duyệt tài liệu mua nhà của đơn đang chờ
const canApprovePurchase = computed(() => isAdmin.value && props.deposit.can_approve_purchase)

const hasAnyAction = computed(
  () =>
    (isBroker.value && props.deposit.can_confirm) ||
    canGenerateOtp.value ||
    canCheckin.value ||
    canReport.value ||
    canOpenDispute.value ||
    canAddEvidence.value ||
    canRate.value ||
    canApprovePurchase.value,
)

function openModal(name: string) {
  modalVisible[name] = true
}

function closeModal(name: string) {
  modalVisible[name] = false
}

/** Chạy 1 thao tác, báo lỗi qua message và refresh dữ liệu khi thành công */
async function run(action: () => Promise<void>, successMessage: string) {
  processing.value = true
  try {
    await action()
    window.message?.success(successMessage)
    emit('changed')
    return true
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Thao tác thất bại'
    window.message?.error(message)
    return false
  } finally {
    processing.value = false
  }
}

async function confirmDeposit() {
  await run(() => depositService.confirmDeposit(props.deposit.id).then(() => undefined), 'Đã xác nhận lịch xem nhà')
}

async function rejectDeposit() {
  if (!rejectReason.value.trim()) {
    window.message?.warning('Vui lòng nhập lý do từ chối')
    return
  }
  const ok = await run(
    () => depositService.rejectDeposit(props.deposit.id, rejectReason.value).then(() => undefined),
    'Đã từ chối lịch, tiền cọc sẽ được hoàn cho khách',
  )
  if (ok) {
    closeModal('reject')
    rejectReason.value = ''
  }
}

async function generateOtp() {
  processing.value = true
  try {
    const result = await depositService.generateOtp(props.deposit.id)
    otpCode.value = result.otp
    otpExpiresAt.value = new Date(result.expires_at).toLocaleTimeString('vi-VN')
    openModal('otp')
    emit('changed')
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Không sinh được OTP'
    window.message?.error(message)
  } finally {
    processing.value = false
  }
}

async function submitCheckin() {
  if (!checkinOtp.value.trim()) {
    window.message?.warning('Vui lòng nhập mã OTP')
    return
  }
  const ok = await run(
    () => depositService.checkin(props.deposit.id, checkinOtp.value.trim()).then(() => undefined),
    'Check-in thành công, vui lòng báo cáo kết quả sau buổi xem',
  )
  if (ok) {
    closeModal('checkin')
    checkinOtp.value = ''
  }
}

async function submitReport() {
  if (!reportValue.value) {
    window.message?.warning('Vui lòng chọn kết quả')
    return
  }
  if (evidenceRequired.value && !evidenceUrls.value.length) {
    window.message?.warning('Vui lòng gửi kèm ít nhất 1 ảnh bằng chứng')
    return
  }
  if (purchaseProofRequired.value && !purchaseProofValue.value) {
    window.message?.warning('Khai đã mua nhà phải chọn loại tài liệu chứng minh')
    return
  }
  const ok = await run(
    () =>
      depositService
        .submitReport(props.deposit.id, reportValue.value, evidenceUrls.value, purchaseProofValue.value ?? '')
        .then(() => undefined),
    'Đã gửi báo cáo',
  )
  if (ok) {
    closeModal('report')
    reportValue.value = ''
    purchaseProofValue.value = null
    evidenceUrls.value = []
  }
}

async function submitDispute() {
  if (!disputeReason.value.trim()) {
    window.message?.warning('Vui lòng nhập lý do tranh chấp')
    return
  }
  const ok = await run(
    () => depositService.openDispute(props.deposit.id, disputeReason.value, evidenceUrls.value).then(() => undefined),
    'Đã mở tranh chấp, tiền cọc được tạm giữ',
  )
  if (ok) {
    closeModal('dispute')
    disputeReason.value = ''
    evidenceUrls.value = []
  }
}

async function submitEvidence() {
  if (!props.deposit.dispute) return
  if (!evidenceUrls.value.length) {
    window.message?.warning('Vui lòng chọn ít nhất 1 bằng chứng')
    return
  }
  const ok = await run(
    () => depositService.addEvidence(props.deposit.dispute!.id, evidenceUrls.value).then(() => undefined),
    'Đã gửi bằng chứng cho admin',
  )
  if (ok) {
    closeModal('evidence')
    evidenceUrls.value = []
  }
}

async function submitRating() {
  const ok = await run(
    () => depositService.rateBroker(props.deposit.id, ratingValue.value, ratingComment.value),
    'Cảm ơn bạn đã đánh giá môi giới',
  )
  if (ok) {
    closeModal('rating')
    ratingComment.value = ''
  }
}

/** Admin duyệt/từ chối tài liệu mua nhà — duyệt thì trừ luôn 1 căn của dự án */
async function submitPurchaseDecision(approved: boolean) {
  if (!approved && !purchaseNote.value.trim()) {
    window.message?.warning('Vui lòng nhập lý do từ chối tài liệu')
    return
  }
  const ok = await run(
    () => depositService.decidePurchase(props.deposit.id, approved, purchaseNote.value),
    approved
      ? 'Đã duyệt tài liệu, hoàn 100% tiền cọc và trừ 1 căn của dự án'
      : 'Đã từ chối tài liệu, đơn chuyển sang tranh chấp',
  )
  if (ok) {
    closeModal(approved ? 'purchase' : 'purchaseReject')
    purchaseNote.value = ''
  }
}
</script>
