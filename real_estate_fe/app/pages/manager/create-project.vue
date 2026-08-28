<template>
  <div class="w-full flex-1 min-h-0 bg-white p-6 rounded-lg border border-gray-200 flex flex-col gap-4">
    <div class="flex-shrink-0">
      <h2 class="text-lg font-bold text-gray-800">{{ isEdit ? 'Chỉnh sửa dự án' : 'Tạo dự án mới' }}</h2>
    </div>

    <!-- Loading khi tải chi tiết dự án (chế độ chỉnh sửa) -->
    <div v-if="loadingDetail" class="flex flex-1 items-center justify-center py-10">
      <n-spin size="small" />
    </div>

    <!-- Hai tab: Nội dung / Hình ảnh dự án -->
    <n-tabs v-else type="line" class="flex-1 min-h-0">
      <n-tab-pane name="info" tab="Nội dung">
        <ProjectInfoForm
          ref="infoFormRef"
          v-model:form="form"
          :province-options="provinceOptions"
          :ward-options="wardOptions"
          :loading-ward="loadingWard"
          :status-options="statusOptions"
          @province-change="onProvinceChange"
        />
      </n-tab-pane>

      <n-tab-pane name="images" tab="Hình ảnh dự án">
        <ProjectImageUpload v-model:ids="form.image_ids" :existing="existingImages" />
      </n-tab-pane>
    </n-tabs>

    <!-- Actions sticky luôn hiển thị ở đáy -->
    <div
      class="sticky bottom-0 flex flex-shrink-0 justify-end gap-2 pt-3 border-t border-gray-100 bg-white -mx-6 -mb-6 px-6 pb-6 mt-1">
      <n-button @click="navigateTo('/nguoi-ban/quan-ly-du-an')">Hủy</n-button>
      <n-button type="error" :loading="submitting" @click="handleSubmit">{{ isEdit ? 'Cập nhật' : 'Tạo dự án' }}</n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import type { CityOption, WardOption } from '~/types/real_estate'
import type { ProjectFormData } from '~/types/manager'
import { useManagerStore } from '~/stores/manager'
import ProjectInfoForm from '~/components/manager/ProjectInfoForm.vue'

definePageMeta({
  alias: [
    '/nguoi-ban/tao-du-an'
  ],
})

const { $api } = useNuxtApp()
const managerStore = useManagerStore()
const route = useRoute()

// Chế độ chỉnh sửa: ?id=<projectId>
const projectId = computed<number | null>(() => {
  const raw = route.query.id
  const id = raw ? Number(raw) : NaN
  return Number.isFinite(id) && id > 0 ? id : null
})
const isEdit = computed(() => projectId.value !== null)

const form = ref<ProjectFormData>({
  name: '',
  alternative_name: '',
  province: null,
  ward: null,
  full_address: '',
  status: 'active',
  category_id: null,
  total_area_ha: null,
  total_units: null,
  price_min: null,
  price_max: null,
  construction_start_date: null,
  handover_date: null,
  image_ids: [],
})

// Ảnh dự án đã lưu (hiển thị khi chỉnh sửa)
const existingImages = ref<{ id: number; url: string }[]>([])

const infoFormRef = ref<InstanceType<typeof ProjectInfoForm> | null>(null)
const submitting = ref(false)
const loadingDetail = ref(false)

const statusOptions: SelectOption[] = [
  { label: 'Đang mở bán', value: 'active' },
  { label: 'Sắp mở bán', value: 'upcoming' },
  { label: 'Đã bàn giao', value: 'handed_over' },
  { label: 'Tạm dừng', value: 'paused' },
]

// ── Vị trí (tỉnh/phường) ──
const provinceOptions = ref<SelectOption[]>([])
const wardOptions = ref<SelectOption[]>([])
const loadingWard = ref(false)

const fetchProvinces = async () => {
  try {
    const res = await $api.get<{ data: CityOption[] }>('/real-estate/list/city')
    provinceOptions.value = res.data.map((item: CityOption) => ({
      label: item.name,
      value: item.code,
    }))
  } catch {
    provinceOptions.value = []
  }
}

const onProvinceChange = async (provinceCode: string | null) => {
  form.value.ward = null
  wardOptions.value = []
  await loadWards(provinceCode)
}

// Load danh sách phường/xã theo mã tỉnh (không reset ward — dùng khi init form edit)
const loadWards = async (provinceCode: string | null) => {
  if (!provinceCode) {
    wardOptions.value = []
    return
  }
  loadingWard.value = true
  try {
    const res = await $api.get<{ data: WardOption[] }>('/real-estate/list/ward', {
      params: { code: provinceCode },
    })
    wardOptions.value = res.data.map((item: WardOption) => ({
      label: item.name,
      value: item.code,
    }))
  } catch {
    wardOptions.value = []
  } finally {
    loadingWard.value = false
  }
}

// Tải chi tiết dự án để điền form khi ở chế độ chỉnh sửa
const loadProjectDetail = async (id: number) => {
  loadingDetail.value = true
  try {
    const res = await $api.get<{ data: any }>(`/manager/projects/${id}`)
    const p = res?.data
    if (!p) return
    form.value = {
      name: p.name || '',
      alternative_name: p.alternative_name || '',
      province: p.province || null,
      ward: p.ward || null,
      full_address: p.full_address || '',
      status: p.status || 'active',
      category_id: p.category_id ?? null,
      total_area_ha: p.total_area_ha ?? null,
      total_units: p.total_units ?? null,
      price_min: p.price_min ?? null,
      price_max: p.price_max ?? null,
      construction_start_date: p.construction_start_date || null,
      handover_date: p.handover_date || null,
      image_ids: (p.images || []).map((img: any) => img.id),
    }
    existingImages.value = (p.images || []).map((img: any) => ({ id: img.id, url: img.url }))
    // Nạp phường/xã theo tỉnh đã lưu rồi mới set ward (để select hiển thị đúng)
    await loadWards(form.value.province)
    form.value.ward = p.ward || null
  } catch (error: any) {
    window.message?.error('Lỗi khi tải thông tin dự án: ' + (error?.message || 'Lỗi máy chủ'))
  } finally {
    loadingDetail.value = false
  }
}

onMounted(async () => {
  await fetchProvinces()
  if (isEdit.value && projectId.value) {
    await loadProjectDetail(projectId.value)
  }
})

// Ép ngày về dạng string "yyyy-MM-dd" (n-date-picker đôi khi trả timestamp number)
const normalizeDate = (value: string | number | null): string | null => {
  if (!value) return null
  if (typeof value === 'number') {
    const d = new Date(value)
    if (Number.isNaN(d.getTime())) return null
    const mm = String(d.getMonth() + 1).padStart(2, '0')
    const dd = String(d.getDate()).padStart(2, '0')
    return `${d.getFullYear()}-${mm}-${dd}`
  }
  return value
}

const handleSubmit = async () => {
  try {
    await infoFormRef.value?.validate()
  } catch {
    window.message?.warning('Vui lòng điền đầy đủ thông tin bắt buộc')
    return
  }

  submitting.value = true
  try {
    const payload = {
      ...form.value,
      construction_start_date: normalizeDate(form.value.construction_start_date),
      handover_date: normalizeDate(form.value.handover_date),
    }
    let res: any
    if (isEdit.value && projectId.value) {
      res = await $api.put<{ success: boolean }>(`/manager/update-project/${projectId.value}`, payload)
    } else {
      res = await $api.post<{ success: boolean; data: { id: number } }>('/manager/create-project', payload)
    }
    if (res?.success) {
      window.message?.success(isEdit.value ? 'Cập nhật dự án thành công' : 'Tạo dự án thành công')
      // Đánh dấu cache dự án cũ → trang quản lý dự án sẽ fetch lại
      managerStore.invalidateProjects()
      navigateTo('/nguoi-ban/quan-ly-du-an')
    }
  } catch (error: any) {
    window.message?.error((isEdit.value ? 'Lỗi khi cập nhật dự án: ' : 'Lỗi khi tạo dự án: ') + (error?.message || 'Lỗi máy chủ'))
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
/* Tab content cuộn nội bộ, không cuốn cả trang */
:deep(.n-tabs) {
  display: flex;
  flex-direction: column;
  height: 100%;
}
:deep(.n-tabs .n-tab-pane) {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}
</style>
