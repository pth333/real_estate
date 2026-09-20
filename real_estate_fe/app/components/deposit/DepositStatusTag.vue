<template>
  <n-tag :type="meta.type" :size="size" round :bordered="false">
    {{ meta.label }}
  </n-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { DepositStatus } from '~/types/deposit'
import { DEPOSIT_STATUS_META } from '~/types/deposit'

const props = withDefaults(
  defineProps<{
    status: DepositStatus
    size?: 'small' | 'medium' | 'large'
  }>(),
  { size: 'small' },
)

const meta = computed(
  () => DEPOSIT_STATUS_META[props.status] ?? { label: props.status, type: 'default' as const },
)
</script>
