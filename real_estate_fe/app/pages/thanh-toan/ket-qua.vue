<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4 py-10">
    <div class="w-full max-w-[520px] bg-white rounded-xl border border-gray-200 p-8 flex flex-col items-center gap-4">
      <div class="flex h-14 w-14 items-center justify-center rounded-full"
        :class="isSuccess ? 'bg-emerald-50 text-emerald-600' : 'bg-red-50 text-red-500'">
        <n-icon size="30">
          <IconCheck v-if="isSuccess" />
          <IconXCircle v-else />
        </n-icon>
      </div>

      <h1 class="text-lg font-bold text-gray-900 text-center">
        {{ isLoading ? 'Đang xác nhận thanh toán...' : title }}
      </h1>

      <p class="text-sm text-gray-500 text-center">
        {{ message }}
      </p>

      <n-descriptions v-if="result" :column="1" size="small" bordered label-placement="left"
        class="w-full mt-2">
        <n-descriptions-item label="Mã đơn">
          #{{ result.deposit_id }}
        </n-descriptions-item>
        <n-descriptions-item label="Trạng thái">
          <DepositStatusTag :status="result.status" />
        </n-descriptions-item>
      </n-descriptions>

      <div class="flex gap-2 mt-2 w-full">
        <n-button class="flex-1" @click="goHome">Về trang chủ</n-button>
        <n-button class="flex-1" type="primary" @click="goToMyDeposits">Xem đơn đặt cọc</n-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import type { PaymentCallbackResult } from '~/types/deposit'
import { useDepositService } from '~/services/deposit.service'
import DepositStatusTag from '~/components/deposit/DepositStatusTag.vue'
import IconCheck from '~/icons/IconCheck.vue'
import IconXCircle from '~/icons/IconXCircle.vue'

definePageMeta({
  layout: 'empty',
})

const route = useRoute()
const depositService = useDepositService()

const isLoading = ref(true)
const result = ref<PaymentCallbackResult | null>(null)

const isSuccess = computed(() => result.value?.success === true)

const title = computed(() => (isSuccess.value ? 'Thanh toán thành công' : 'Thanh toán chưa hoàn tất'))

const message = computed(() => {
  if (result.value) return result.value.message
  if (isLoading.value) return 'Vui lòng chờ trong giây lát'
  return 'Không đọc được kết quả giao dịch từ cổng thanh toán'
})

onMounted(async () => {
  // Gom toàn bộ query cổng thanh toán trả về để backend xác thực chữ ký
  const params: Record<string, string> = {}
  Object.entries(route.query).forEach(([key, value]) => {
    if (typeof value === 'string') params[key] = value
  })

  if (!params.vnp_TxnRef && !params.payment_ref) {
    isLoading.value = false
    return
  }

  try {
    result.value = await depositService.confirmPayment(params)
  } catch {
    // Xác thực chữ ký thất bại hoặc giao dịch không tồn tại
    result.value = null
  } finally {
    isLoading.value = false
  }
})

function goHome() {
  navigateTo('/')
}

function goToMyDeposits() {
  navigateTo('/tai-khoan/dat-coc')
}
</script>
