<template>
  <n-form ref="formRef" :model="form" :rules="rules" label-placement="top" class="w-full">
    <!-- Thông tin dự án -->
    <n-form-item label="Tên dự án" path="name">
      <n-input v-model:value="form.name" placeholder="VD: Mandarin Garden 3" clearable />
    </n-form-item>

    <n-form-item label="Tên gọi khác" path="alternative_name">
      <n-input v-model:value="form.alternative_name" placeholder="VD: Mandarin Garden 3 Yên Sở" clearable />
    </n-form-item>

    <!-- Vị trí dự án — quan trọng để map với bất động sản khi tạo tin -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <n-form-item label="Tỉnh / Thành phố" path="province">
        <n-select v-model:value="form.province" placeholder="Chọn tỉnh/thành phố" clearable filterable
          :options="provinceOptions" @update:value="onProvinceChange" />
      </n-form-item>

      <n-form-item label="Phường / Xã" path="ward">
        <n-select v-model:value="form.ward" placeholder="Chọn phường/xã" clearable filterable :options="wardOptions"
          :loading="loadingWard" :disabled="!form.province" />
      </n-form-item>
    </div>

    <n-form-item label="Địa chỉ dự án" path="full_address">
      <n-input v-model:value="form.full_address" placeholder="Số nhà, đường, khu vực..." clearable />
    </n-form-item>

    <n-form-item label="Trạng thái" path="status">
      <n-select v-model:value="form.status" placeholder="Chọn trạng thái" :options="statusOptions" />
    </n-form-item>

    <n-form-item label="Danh mục dự án" path="category_id">
      <n-select v-model:value="form.category_id" placeholder="Chọn danh mục" clearable filterable
        :options="categoryOptions" />
    </n-form-item>

    <!-- Quy mô & giá -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <n-form-item label="Tổng diện tích (ha)">
        <n-input-number v-model:value="form.total_area_ha" placeholder="VD: 5.2" class="w-full" :min="0" />
      </n-form-item>
      <n-form-item label="Tổng số căn">
        <n-input-number v-model:value="form.total_units" placeholder="VD: 1200" class="w-full" :min="0" />
      </n-form-item>
      <n-form-item label="Giá từ (VNĐ)">
        <n-input-number v-model:value="form.price_min" placeholder="VD: 2000000000" class="w-full" :min="0" />
      </n-form-item>
      <n-form-item label="Giá đến (VNĐ)">
        <n-input-number v-model:value="form.price_max" placeholder="VD: 5000000000" class="w-full" :min="0" />
      </n-form-item>
      <n-form-item label="Ngày khởi công">
        <n-date-picker v-model:value="form.construction_start_date" type="date" value-format="yyyy-MM-dd" class="w-full"
          clearable />
      </n-form-item>
      <n-form-item label="Ngày bàn giao">
        <n-date-picker v-model:value="form.handover_date" type="date" value-format="yyyy-MM-dd" class="w-full"
          clearable />
      </n-form-item>
    </div>
  </n-form>
</template>

<script setup lang="ts">
import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import type { ProjectFormData } from '~/types/manager'

// v-model:form — object dữ liệu form (cùng reference với parent)
const form = defineModel<ProjectFormData>('form', { required: true })

const props = defineProps<{
  provinceOptions: SelectOption[]
  wardOptions: SelectOption[]
  loadingWard: boolean
  statusOptions: SelectOption[]
}>()

const emit = defineEmits<{
  'province-change': [provinceCode: string | null]
}>()

const formRef = ref<FormInst | null>(null)

// Danh mục dự án — đọc từ window.menu (đã load qua menu store),
// duyệt đệ quy cả children, lọc loại "project"
const categoryOptions = computed<SelectOption[]>(() => {
  if (!import.meta.client) return []
  const categories = window.menu?.settings?.categories || []
  const result: SelectOption[] = []

  const walk = (list: typeof categories) => {
    list.forEach((cat) => {
      if (cat.Type === 'project') {
        result.push({ label: cat.Name, value: cat.ID })
      }
      if (cat.children?.length) {
        walk(cat.children)
      }
    })
  }

  walk(categories)
  // Bỏ option đầu tiên (thường là danh mục cha tổng, không phải dự án cụ thể)
  result.shift()
  return result
})

const rules: FormRules = {
  name: { required: true, message: 'Vui lòng nhập tên dự án', trigger: 'blur' },
  province: { required: true, message: 'Vui lòng chọn tỉnh/thành phố', trigger: 'change' },
  ward: { required: true, message: 'Vui lòng chọn phường/xã', trigger: 'change' },
}

const onProvinceChange = (provinceCode: string | null) => {
  emit('province-change', provinceCode)
}

// Cho parent gọi validate khi submit
defineExpose({
  validate: () => formRef.value?.validate(),
})
</script>
