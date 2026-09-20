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
          <span class="font-semibold text-gray-800">Đặt cọc giữ lịch xem nhà</span>
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

            <div class="grid grid-cols-2 gap-3">
              <n-form-item label="Giờ bắt đầu">
                <n-time-picker v-model:formatted-value="form.viewingStart" format="HH:mm" value-format="HH:mm"
                  class="w-full" placeholder="VD 09:00" />
              </n-form-item>
              <n-form-item label="Giờ kết thúc">
                <n-time-picker v-model:formatted-value="form.viewingEnd" format="HH:mm" value-format="HH:mm"
                  class="w-full" placeholder="VD 10:00" />
              </n-form-item>
            </div>

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

        <!-- Mức cọc do hệ thống tính theo giá BĐS — khách không sửa được -->
        <n-spin :show="loadingOptions">
          <div class="rounded-lg border border-gray-200 p-3 flex flex-col gap-2">
            <div class="flex items-center justify-between">
              <span class="text-xs font-semibold uppercase tracking-wide text-gray-400">Mức cọc hệ thống áp dụng</span>
              <n-tag v-if="options" size="tiny" :bordered="false" type="default">
                {{ options.policy_label }}
              </n-tag>
            </div>

            <div v-if="options" class="flex items-center justify-between">
              <span class="text-sm text-gray-600">Số tiền cọc</span>
              <span class="text-lg font-bold text-emerald-600">{{ formatVnd(options.amount) }}</span>
            </div>
            <div v-if="options" class="flex items-center justify-between">
              <span class="text-sm text-gray-600">Phí môi giới (nếu không mua)</span>
              <span class="text-sm font-semibold text-gray-800">{{ formatVnd(options.broker_fee) }}</span>
            </div>

            <n-alert v-if="options?.is_fallback" type="warning" :bordered="false" class="mt-1">
              Bất động sản chưa có giá nên đang áp dụng mức cọc mặc định của hệ thống.
            </n-alert>
          </div>
        </n-spin>

        <!-- Giải thích quy tắc hoàn tiền theo đúng số tiền đang áp dụng -->
        <div class="rounded-lg bg-emerald-50/60 border border-emerald-100 p-3 text-xs text-emerald-800 flex flex-col gap-1">
          <span>• Môi giới từ chối hoặc không phản hồi trong 24 giờ: hoàn 100% tiền cọc.</span>
          <span>• Khách đến xem và mua nhà: hoàn 100% tiền cọc.</span>
          <span>• Khách đến nhưng không mua: hoàn {{ formatVnd(refundIfNotBuy) }} (cọc trừ phí môi giới).</span>
          <span>• Khách không đến: mất toàn bộ tiền cọc.</span>
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
import { computed, reactive, ref, watch } from 'vue'
import { useDepositService } from '~/services/deposit.service'
import { formatVnd } from '~/utils/deposit'
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

const paymentMethods: { label: string; value: PaymentMethod }[] = [
  { label: 'VNPay', value: 'VNPAY' },
  { label: 'Momo', value: 'MOMO' },
  { label: 'ZaloPay', value: 'ZALOPAY' },
]

// Mức cọc + phí môi giới do backend trả về theo giá BĐS
const options = ref<BookingOptions | null>(null)
const loadingOptions = ref(false)
const submitting = ref(false)

const form = reactive({
  viewingDate: '' as string | null,
  viewingStart: '' as string | null,
  viewingEnd: '' as string | null,
  paymentMethod: 'VNPAY' as PaymentMethod,
})

const refundIfNotBuy = computed(() => {
  if (!options.value) return 0
  return Math.max(options.value.amount - options.value.broker_fee, 0)
})

// Không cho chọn ngày trong quá khứ
function disablePastDate(timestamp: number): boolean {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return timestamp < today.getTime()
}

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
  form.viewingEnd = null
  form.paymentMethod = 'VNPAY'
}

function close() {
  emit('update:show', false)
}

function handleShowChange(value: boolean) {
  emit('update:show', value)
}

function validate(): string | null {
  if (!options.value) return 'Chưa lấy được mức cọc của bất động sản này'
  if (!form.viewingDate) return 'Vui lòng chọn ngày xem nhà'
  if (!form.viewingStart || !form.viewingEnd) return 'Vui lòng chọn khung giờ xem nhà'
  if (form.viewingStart >= form.viewingEnd) return 'Giờ kết thúc phải sau giờ bắt đầu'
  return null
}

async function submit() {
  const error = validate()
  if (error) {
    window.message?.warning(error)
    return
  }

  submitting.value = true
  try {
    const result = await depositService.createDeposit({
      real_estate_id: props.realEstateId,
      viewing_date: form.viewingDate as string,
      viewing_start: form.viewingStart as string,
      viewing_end: form.viewingEnd as string,
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
