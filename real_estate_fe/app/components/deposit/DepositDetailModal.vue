<template>
  <n-modal :show="show" :mask-closable="false" style="width: 720px; max-width: 94vw" @update:show="handleShowChange">
    <div class="bg-white rounded-xl flex flex-col max-h-[88vh]">
      <!-- Header -->
      <div class="flex items-center justify-between gap-3 px-6 py-4 border-b border-gray-100 flex-shrink-0">
        <div class="flex items-center gap-3">
          <span class="font-semibold text-gray-800">Đơn đặt cọc #{{ deposit?.id ?? depositId }}</span>
          <DepositStatusTag v-if="deposit" :status="deposit.status" size="medium" />
        </div>
        <n-button quaternary circle @click="close">
          <template #icon>
            <n-icon>
              <IconCloseOutline />
            </n-icon>
          </template>
        </n-button>
      </div>

      <!-- Body -->
      <div class="px-6 py-5 overflow-y-auto flex-1">
        <n-spin :show="loading">
          <DepositSummaryCard v-if="deposit" :deposit="deposit" />
          <n-empty v-else-if="!loading" description="Không tìm thấy đơn đặt cọc" />
        </n-spin>
      </div>

      <!-- Actions -->
      <div v-if="deposit" class="px-6 pb-5 pt-1 flex-shrink-0">
        <DepositActions :deposit="deposit" :role="role" @changed="reload" />
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Deposit, DepositActorRole } from '~/types/deposit'
import { useDepositService } from '~/services/deposit.service'
import IconCloseOutline from '~/icons/IconCloseOutline.vue'

const props = defineProps<{
  show: boolean
  depositId: number | null
  role: DepositActorRole
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  changed: []
}>()

const depositService = useDepositService()

const deposit = ref<Deposit | null>(null)
const loading = ref(false)

watch(
  () => [props.show, props.depositId],
  ([visible]) => {
    if (visible && props.depositId) load()
  },
  { immediate: true },
)

async function load() {
  if (!props.depositId) return
  loading.value = true
  try {
    deposit.value = await depositService.getDeposit(props.depositId)
  } catch {
    deposit.value = null
  } finally {
    loading.value = false
  }
}

/** Tải lại chi tiết sau khi thao tác và báo trang cha refresh danh sách */
async function reload() {
  await load()
  emit('changed')
}

function close() {
  emit('update:show', false)
}

function handleShowChange(value: boolean) {
  emit('update:show', value)
}
</script>
