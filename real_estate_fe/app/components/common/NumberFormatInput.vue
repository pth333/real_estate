<template>
  <n-input
    :value="displayed"
    :placeholder="placeholder"
    :status="status"
    class="w-full"
    @update:value="onInput"
  >
    <template v-if="suffix" #suffix>{{ suffix }}</template>
  </n-input>
</template>

<script setup lang="ts">
/**
 * Ô nhập số có tự động format dấu chấm hàng nghìn (vi-VN).
 * - Người dùng gõ số → ô hiển thị ngay dạng đã format, đồng thời emit số thực ra ngoài qua v-model.
 * - Bên ngoài đổi giá trị (reset form...) → ô tự đồng bộ lại hiển thị.
 *
 * VD:
 *   <NumberFormatInput v-model="postStore.form.price" placeholder="Nhập mức giá" suffix="m²" />
 */

type InputStatus = 'error' | 'warning' | 'success' | 'default'

const props = withDefaults(defineProps<{
  modelValue: number | null
  placeholder?: string
  suffix?: string
  status?: InputStatus
}>(), {
  placeholder: '',
  suffix: '',
  status: 'default',
})

const emit = defineEmits<{
  'update:modelValue': [value: number | null]
}>()

// Giá trị đang hiển thị trên ô (đã format dấu chấm)
const displayed = ref<string>(formatPriceNumber(props.modelValue))

// Người dùng gõ: lọc số, format lại ô hiển thị, emit số thực về store
function onInput(raw: string) {
  const parsed = parsePriceNumber(raw)
  displayed.value = formatPriceNumber(parsed)
  emit('update:modelValue', parsed)
}

// Bên ngoài thay đổi (reset form...) → đồng bộ lại ô hiển thị
watch(() => props.modelValue, (val) => {
  const current = parsePriceNumber(displayed.value)
  if (current !== val) {
    displayed.value = formatPriceNumber(val)
  }
})
</script>
