<template>
  <div v-if="hasAnyAction" class="rounded-lg border border-emerald-100 bg-emerald-50/50 p-4 flex flex-col gap-3">
    <span class="text-xs font-semibold uppercase tracking-wide text-emerald-700">Thao tác của bạn</span>

    <div class="flex flex-wrap items-center gap-2">
      <!-- Môi giới: 1 chạm = xác nhận đã tới, hệ thống tự ghi vị trí + sinh mã cho khách -->
      <n-button v-if="isBroker && deposit.can_confirm" key="confirm-schedule" type="primary" :loading="processing"
        @click="confirmDeposit">
        Xác nhận lịch
      </n-button>
      <n-button v-if="isBroker && deposit.can_confirm" key="reject" type="error" ghost @click="openModal('reject')">
        Từ chối
      </n-button>
      <n-button v-if="isBroker && canGenerateOtp" key="broker-arrived" type="warning" :loading="processing"
        @click="generateOtp">
        Tôi đã tới - Sinh mã OTP
      </n-button>

      <!-- Khách: 1 chạm = xác nhận đã tới, hệ thống tự ghi vị trí (không cần nhập gì) -->
      <n-button v-if="isCustomer && canCheckin" key="customer-otp" type="primary" @click="openModal('checkin')">
        Nhập mã OTP từ môi giới
      </n-button>
      <n-button v-if="isCustomer && canCheckin" key="customer-arrived" type="warning" :loading="processing"
        @click="confirmArrived">
        Tôi đã tới
      </n-button>
      <n-button v-if="isCustomer && canRate" key="rate" secondary @click="openModal('rating')">
        Đánh giá môi giới
      </n-button>

      <!-- Dùng chung -->
      <n-button v-if="canReport" key="report" type="info" ghost @click="openModal('report')">
        Báo cáo kết quả
      </n-button>
      <n-button v-if="canAddEvidence" key="evidence" warning ghost @click="openModal('evidence')">
        Gửi bằng chứng
      </n-button>
      <n-button v-if="canOpenDispute" key="dispute" type="error" ghost @click="openModal('dispute')">
        Mở tranh chấp
      </n-button>

      <!-- Admin duyệt tài liệu mua nhà -->
      <n-button v-if="canApprovePurchase" key="purchase-approve" type="primary" @click="openModal('purchase')">
        Duyệt tài liệu mua nhà
      </n-button>
      <n-button v-if="canApprovePurchase" key="purchase-reject" type="error" ghost @click="openModal('purchaseReject')">
        Từ chối tài liệu
      </n-button>
    </div>

    <!-- Cảnh báo bằng chứng vị trí: HIỂN THỊ CỨNG (theo dữ liệu đơn, không phải toast) -->
    <n-alert v-if="evidenceWarning || actionWarning" type="warning" :bordered="false" class="text-xs">
      {{ evidenceWarning || actionWarning }}
    </n-alert>

    <!-- Khách: nhập mã OTP là cách duy nhất để đơn chuyển sang "Đã check-in" -->
    <n-alert v-if="isCustomer && canCheckin" type="info" :bordered="false" class="text-xs">
      <strong>Nhập mã OTP từ môi giới</strong> để xác nhận buổi xem (chỉ khi đó đơn mới chuyển sang "Đã check-in").<br />
      Môi giới chưa đưa mã? Bấm <strong>Tôi đã tới</strong> — hệ thống tự ghi nhận vị trí của bạn làm bằng chứng
      bạn đã có mặt. Nếu môi giới không đến, bạn được hoàn 100% phí môi giới.
    </n-alert>

    <!-- Môi giới: vị trí là căn cứ chuyển phí, khách nhập mã OTP là căn cứ chắc nhất -->
    <n-alert v-if="isBroker && canGenerateOtp" type="info" :bordered="false" class="text-xs">
      Bấm <strong>Tôi đã tới - Sinh mã OTP</strong> rồi đọc mã cho khách nhập: khách nhập mã là căn cứ chắc chắn
      nhất cho buổi xem. Vị trí của bạn cũng được ghi nhận làm bằng chứng (nên bật định vị đúng chỗ).
    </n-alert>

    <!-- Từ chối lịch -->
    <n-modal v-model:show="modalVisible.reject" :mask-closable="false">
      <div class="w-[440px] bg-white rounded-xl p-6 flex flex-col gap-4">
        <span class="font-semibold text-gray-800">Từ chối lịch xem nhà</span>
        <p class="text-sm text-gray-500">
          Khách sẽ được hoàn 100% phí môi giới. Vui lòng cho khách biết lý do.
        </p>
        <n-input v-model:value="rejectReason" type="textarea" :rows="3" placeholder="VD: Bất động sản đã có khách mua" />
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
        <n-alert v-if="locationHint" type="info" :bordered="false" class="text-xs">
          {{ locationHint }}
        </n-alert>
        <n-alert type="warning" :bordered="false" class="text-xs">
          Môi giới chưa đưa mã? Đóng hộp này và bấm <strong>Tôi đã tới</strong> — hệ thống tự ghi nhận vị trí của bạn.
        </n-alert>
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
          phí môi giới sẽ được tạm giữ và chuyển cho admin xác minh dựa trên bằng chứng.
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
          Phí môi giới sẽ bị tạm giữ cho tới khi admin ra quyết định. Vui lòng mô tả rõ sự việc và gửi kèm bằng chứng.
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
          Duyệt sẽ <strong>hoàn 100% phí môi giới</strong> cho khách và
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
          Đơn sẽ chuyển sang <strong>tranh chấp</strong> để xử lý tiếp, phí môi giới vẫn bị tạm giữ.
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
import { NButton, NInput, NModal, NRadio, NRadioGroup, NRate } from 'naive-ui'
import type { Deposit, DepositActorRole } from '~/types/deposit'
import { PURCHASE_PROOF_LABEL, PURCHASE_PROOF_OPTIONS } from '~/types/deposit'
import { useDepositService } from '~/services/deposit.service'
import { reportOptions, getCurrentLocation } from '~/utils/deposit'

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
// Lỗi tức thời khi xin quyền vị trí / gọi API (cảnh báo bằng chứng thì tính từ dữ liệu đơn, xem evidenceWarning)
const actionWarning = ref('')
const reportValue = ref('')
const purchaseProofValue = ref<string | null>(null)
const disputeReason = ref('')
const ratingValue = ref(5)
const ratingComment = ref('')
const purchaseNote = ref('')

const options = computed(() => reportOptions(props.deposit, isBroker.value))

/**
 * Cảnh báo bằng chứng vị trí — tính TỪ DỮ LIỆU ĐƠN nên hiển thị cứng:
 * mở lại chi tiết đơn, refresh trang vẫn còn, chỉ mất khi tình trạng thực sự được xử lý
 * (bấm lại với định vị đúng chỗ ⇒ evidence = AT_ESTATE, hoặc đơn đã check-in).
 */
const evidenceWarning = computed(() => {
  if (props.deposit.status !== 'BROKER_CONFIRMED') return ''

  const mine = isBroker.value ? props.deposit.broker_checkin_evidence : props.deposit.customer_checkin_evidence
  const role = isBroker.value ? 'môi giới' : 'khách'

  if (mine === 'FAR') {
    return `Vị trí bạn gửi khá xa bất động sản nên chưa đủ căn cứ xác nhận ${role} đã tới — trừ khi khách nhập mã OTP (lúc đó hệ thống coi như 2 bên đã gặp nhau). Hãy bật định vị đúng chỗ rồi bấm lại.`
  }
  if (mine === 'NO_LOCATION') {
    return `Chưa ghi nhận được vị trí của bạn nên chưa đủ căn cứ xác nhận ${role} đã tới — trừ khi khách nhập mã OTP. Hãy bật định vị rồi bấm lại.`
  }
  return ''
})

// Hệ thống tự lấy vị trí khi khách bấm xác nhận — chỉ lưu toạ độ làm tròn ~100m, không lưu vị trí chính xác
const locationHint = computed(() =>
  'Hệ thống tự ghi nhận vị trí của bạn (chỉ lưu toạ độ làm tròn ~100m), bạn không phải nhập gì.',
)

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

// Cho phép sinh lại OTP trong suốt cửa sổ check-in (mã cũ hết hạn / khách tới muộn)
const canGenerateOtp = computed(() => isBroker.value && props.deposit.can_checkin)
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
// Cho phép đánh giá sau khi buổi xem đã diễn ra (đã check-in) hoặc đơn đã tất toán
const canRate = computed(
  () =>
    isCustomer.value &&
    !props.deposit.has_rating &&
    ['CHECKED_IN', 'COMPLETED', 'REFUNDED', 'VISITED_BOUGHT', 'VISITED_NOT_BUY'].includes(props.deposit.status),
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
    'Đã từ chối lịch, phí môi giới sẽ được hoàn cho khách',
  )
  if (ok) {
    closeModal('reject')
    rejectReason.value = ''
  }
}

async function generateOtp() {
  processing.value = true
  try {
    // Kèm vị trí: hệ thống dùng toạ độ làm căn cứ chuyển phí cho môi giới
    const location = await getCurrentLocation()
    const result = await depositService.generateOtp(props.deposit.id, location)
    otpCode.value = result.otp
    otpExpiresAt.value = new Date(result.expires_at).toLocaleTimeString('vi-VN')
    // Cảnh báo bằng chứng tính từ dữ liệu đơn (evidenceWarning) nên không cần set ở đây
    actionWarning.value = ''
    openModal('otp')
    emit('changed')
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Không sinh được OTP'
    window.message?.error(message)
  } finally {
    processing.value = false
  }
}

/** Khách xác nhận đã tới — 1 chạm, hệ thống tự lấy vị trí (không cần mã OTP của môi giới) */
async function confirmArrived() {
  processing.value = true
  try {
    const location = await getCurrentLocation()
    if (!location) {
      // Lỗi này không phải toast: khách cần thấy để biết phải bật quyền vị trí rồi bấm lại
      actionWarning.value = 'Cần quyền truy cập vị trí để xác nhận bạn đã tới. Hãy cho phép truy cập vị trí rồi bấm lại "Tôi đã tới".'
      return
    }
    await depositService.checkinLocation(props.deposit.id, location)
    actionWarning.value = ''
    window.message?.success('Đã xác nhận bạn có mặt tại buổi xem nhà')
    emit('changed')
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Không xác nhận được'
    actionWarning.value = message
  } finally {
    processing.value = false
  }
}

async function submitCheckin() {
  if (!checkinOtp.value.trim()) {
    window.message?.warning('Vui lòng nhập mã OTP')
    return
  }
  const ok = await run(async () => {
    // Vị trí là bằng chứng bổ trợ: không lấy được vẫn check-in được bằng OTP
    const location = await getCurrentLocation()
    await depositService.checkin(props.deposit.id, checkinOtp.value.trim(), location).then(() => undefined)
  }, 'Check-in thành công, hệ thống đã ghi nhận 2 bên gặp nhau')
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
    'Đã mở tranh chấp, phí môi giới được tạm giữ',
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
      ? 'Đã duyệt tài liệu, hoàn 100% phí môi giới và trừ 1 căn của dự án'
      : 'Đã từ chối tài liệu, đơn chuyển sang tranh chấp',
  )
  if (ok) {
    closeModal(approved ? 'purchase' : 'purchaseReject')
    purchaseNote.value = ''
  }
}
</script>
