<template>
  <n-modal :show="show" :mask-closable="false" @update:show="handleShowChange">
    <div class="w-[520px] max-w-[92vw] bg-white rounded-xl flex flex-col">
      <!-- Header -->
      <div class="flex items-center gap-3 px-6 pt-6 pb-4 border-b border-gray-100">
        <div class="flex h-9 w-9 items-center justify-center rounded-full bg-emerald-50 text-emerald-600">
          <n-icon size="20">
            <IconShieldCheck />
          </n-icon>
        </div>
        <div class="flex flex-col leading-tight">
          <span class="font-semibold text-gray-800">Đặt lịch xem nhà</span>
          <span class="text-xs text-gray-400">Tiền được platform giữ an toàn, không chuyển cho môi giới ngay</span>
        </div>
      </div>

      <!-- Body -->
      <div class="px-6 py-5 flex flex-col gap-4">
        <div v-if="estateTitle" class="rounded-lg bg-gray-50 px-3 py-2 text-sm text-gray-600">
          {{ estateTitle }}
        </div>

        <n-form label-placement="top" :show-feedback="false">
          <div class="flex flex-col gap-4">
            <n-form-item label="Ngày xem nhà">
              <n-date-picker v-model:formatted-value="form.viewingDate" type="date" value-format="yyyy-MM-dd"
                :is-date-disabled="disablePastDate" clearable class="w-full" placeholder="Chọn ngày" />
            </n-form-item>

            <n-form-item label="Khung giờ xem nhà">
              <div class="flex w-full flex-col gap-1">
                <n-select v-model:value="form.viewingStart" :options="slotOptions" class="w-full"
                  :disabled="!form.viewingDate || slotOptions.length === 0"
                  :placeholder="form.viewingDate ? 'Chọn khung giờ' : 'Chọn ngày xem nhà trước'" />
                <span v-if="form.viewingDate && slotOptions.length === 0" class="text-xs text-amber-600">
                  Hôm nay đã hết khung giờ xem nhà, vui lòng chọn ngày khác.
                </span>
              </div>
            </n-form-item>

            <n-form-item label="Họ tên người liên hệ">
              <n-input v-model:value="form.contactName" class="w-full" maxlength="100"
                placeholder="VD Nguyễn Văn A" />
            </n-form-item>

            <n-form-item label="Số điện thoại liên hệ">
              <n-input v-model:value="form.contactPhone" class="w-full" maxlength="20"
                placeholder="VD 0901234567" />
            </n-form-item>

            <n-form-item label="Phương thức thanh toán">
              <n-radio-group v-model:value="form.paymentMethod" name="payment-method">
                <n-space>
                  <n-radio-button v-for="method in paymentMethods" :key="method.value" :value="method.value">
                    {{ method.label }}
                  </n-radio-button>
                </n-space>
              </n-radio-group>
            </n-form-item>
          </div>
        </n-form>

        <!-- Phí môi giới do hệ thống tính theo giá BĐS — khách không sửa được -->
        <n-spin :show="loadingOptions">
          <div class="rounded-lg border border-gray-200 p-3 flex flex-col gap-2">
            <div class="flex items-center justify-between">
              <span class="text-xs font-semibold uppercase tracking-wide text-gray-400">Phí môi giới hệ thống áp dụng</span>
              <n-tag v-if="options" size="tiny" :bordered="false" type="default">
                {{ options.policy_label }}
              </n-tag>
            </div>

            <div v-if="options" class="flex items-center justify-between">
              <span class="text-sm text-gray-600">Phí môi giới (nếu không mua)</span>
              <span class="text-sm font-semibold text-gray-800">{{ formatVnd(options.broker_fee) }}</span>
            </div>

            <n-alert v-if="options?.is_fallback" type="warning" :bordered="false" class="mt-1">
              Bất động sản chưa có giá nên đang áp dụng phí môi giới mặc định của hệ thống.
            </n-alert>
          </div>
        </n-spin>

        <!-- Quy tắc xử lý phí môi giới sau buổi xem -->
        <div class="rounded-lg bg-emerald-50/60 border border-emerald-100 p-3 text-xs text-emerald-800 flex flex-col gap-1">
          <span>• Bạn đặt cọc mua bất động sản trong <strong>4 ngày</strong> sau buổi xem: hoàn 100% phí.</span>
          <span>• Quá 4 ngày mà bạn không đặt cọc mua: phí thuộc về môi giới.</span>
          <span>• Môi giới từ chối hoặc không phản hồi trong 24 giờ: hoàn 100% phí.</span>
          <span>• Bạn không đến buổi xem: môi giới nhận toàn bộ phí.</span>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex justify-end gap-2 px-6 pb-6">
        <n-button @click="close">Huỷ</n-button>
        <n-button type="primary" :loading="submitting" :disabled="!options" @click="submit">
          Thanh toán và giữ lịch
        </n-button>
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { useDepositService } from '~/services/deposit.service'
import { buildViewingSlots, formatVnd, isVietnamesePhone, normalizeVietnamesePhone } from '~/utils/deposit'
import type { ViewingSlot } from '~/utils/deposit'
import type { BookingOptions, PaymentMethod } from '~/types/deposit'
import IconShieldCheck from '~/icons/IconShieldCheck.vue'

const props = defineProps<{
  show: boolean
  realEstateId: number
  estateTitle?: string
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
}>()

const depositService = useDepositService()
const authStore = useAuthStore()

const paymentMethods: { label: string; value: PaymentMethod }[] = [
  { label: 'VNPay', value: 'VNPAY' },
  { label: 'Momo', value: 'MOMO' },
  { label: 'ZaloPay', value: 'ZALOPAY' },
]

// Phí môi giới do backend trả về theo giá BĐS
const options = ref<BookingOptions | null>(null)
const loadingOptions = ref(false)
const submitting = ref(false)

const form = reactive({
  viewingDate: '' as string | null,
  viewingStart: '' as string | null,
  // Liên hệ của buổi xem — mặc định lấy theo hồ sơ, khách sửa được
  contactName: '',
  contactPhone: '',
  paymentMethod: 'VNPAY' as PaymentMethod,
})

// Danh sách khung giờ cố định 1 tiếng; ngày hôm nay chỉ còn khung sau giờ hiện tại
const slotOptions = computed(() =>
  buildViewingSlots(form.viewingDate).map((slot) => ({ label: slot.label, value: slot.start })),
)

/** Khung giờ đang chọn — giờ kết thúc lấy từ khung, khách không sửa được */
function currentSlot(): ViewingSlot | null {
  if (!form.viewingDate || !form.viewingStart) return null
  return buildViewingSlots(form.viewingDate).find((slot) => slot.start === form.viewingStart) ?? null
}

// Không cho chọn ngày trong quá khứ
function disablePastDate(timestamp: number): boolean {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return timestamp < today.getTime()
}

// Đổi ngày thì khung giờ cũ không còn hợp lệ (nhất là khi chuyển sang ngày hôm nay)
watch(
  () => form.viewingDate,
  () => {
    form.viewingStart = null
  },
)

watch(
  () => props.show,
  (visible) => {
    if (!visible) return
    resetForm()
    fetchOptions()
  },
)

async function fetchOptions() {
  loadingOptions.value = true
  try {
    options.value = await depositService.getBookingOptions(props.realEstateId)
  } catch {
    options.value = null
  } finally {
    loadingOptions.value = false
  }
}

function resetForm() {
  form.viewingDate = null
  form.viewingStart = null
  // Prefill liên hệ từ hồ sơ user đang đăng nhập, khách vẫn sửa được
  form.contactName = authStore.user?.name ?? ''
  form.contactPhone = authStore.user?.phone ?? ''
  form.paymentMethod = 'VNPAY'
}

function close() {
  emit('update:show', false)
}

function handleShowChange(value: boolean) {
  emit('update:show', value)
}

function validate(): string | null {
  if (!options.value) return 'Chưa lấy được phí môi giới của bất động sản này'
  if (form.contactName.trim().length < 2) return 'Vui lòng nhập họ tên người liên hệ (2 - 100 ký tự)'
  if (form.contactName.trim().length > 100) return 'Họ tên người liên hệ tối đa 100 ký tự'
  if (!isVietnamesePhone(form.contactPhone)) return 'Số điện thoại liên hệ không hợp lệ (VD 0901234567)'
  if (!form.viewingDate) return 'Vui lòng chọn ngày xem nhà'
  if (!form.viewingStart) return 'Vui lòng chọn khung giờ xem nhà'
  // Khách có thể để form mở lâu nên kiểm tra lại theo giờ hiện tại (khung hôm nay có thể đã qua)
  if (!currentSlot()) return 'Khung giờ này đã qua, vui lòng chọn khung giờ khác'
  return null
}

async function submit() {
  const error = validate()
  if (error) {
    window.message?.warning(error)
    return
  }

  const slot = currentSlot()
  if (!slot) return

  submitting.value = true
  try {
    const result = await depositService.createDeposit({
      real_estate_id: props.realEstateId,
      viewing_date: form.viewingDate as string,
      viewing_start: slot.start,
      viewing_end: slot.end,
      contact_name: form.contactName.trim(),
      contact_phone: normalizeVietnamesePhone(form.contactPhone),
      payment_method: form.paymentMethod,
    })

    window.message?.success('Đã giữ khung giờ, đang chuyển tới cổng thanh toán')
    close()
    // Chuyển sang cổng thanh toán (cổng thật hoặc trang mô phỏng khi dùng mock)
    await navigateTo(result.payment_url, { external: true })
  } catch {
    // $api đã hiển thị message lỗi
  } finally {
    submitting.value = false
  }
}
</script>
