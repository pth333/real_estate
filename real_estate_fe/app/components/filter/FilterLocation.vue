<template>
  <div class="flex flex-col gap-4">
    <!-- Tỉnh/TP -->
    <div>
      <label class="mb-1 block text-xs font-medium text-gray-500">Tỉnh / Thành phố</label>
      <n-select v-model:value="filterStore.filterLocation.city" :options="filterStore.cityOptions"
        :loading="loadingCity" placeholder="Chọn tỉnh/thành phố" clearable filterable
        @update:value="onDistrictChange" />
    </div>

    <!-- Quận/Huyện -->
    <div>
      <label class="mb-1 block text-xs font-medium text-gray-500">Phường / Xã</label>
      <n-select v-model:value="filterStore.filterLocation.ward" :options="filterStore.wardOptions"
        :loading="loadingWard" placeholder="Chọn phường/xã" clearable filterable
        :disabled="!filterStore.filterLocation.city" />
    </div>

    <!-- Nút Áp dụng -->
    <div class="mt-2 flex justify-end">
      <n-button type="primary" @click="onApply">
        Áp dụng
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { CityOption, WardOption } from '~/types/real_estate'
import { useFilterStore } from '~/stores/filter';

const { $api } = useNuxtApp();
const filterStore = useFilterStore()

const loadingCity = ref(false)
const loadingWard = ref(false)

const fetchListCity = async () => {
  try {
    loadingCity.value = true
    const res = await $api.get<{ data: CityOption[] }>("/real-estate/list/city")
    filterStore.cityOptions = res.data.map((item: CityOption) => ({
      label: item.name,
      value: item.code
    }))
  } finally {
    loadingCity.value = false
  }
}


const onDistrictChange = async (districtCode: string | null) => {
  // Reset ward
  filterStore.filterLocation.ward = undefined
  filterStore.wardOptions = []

  if (!districtCode) return

  try {
    loadingWard.value = true
    const res = await $api.get<{ data: WardOption[] }>(`/real-estate/list/ward`, {
      params: { code: districtCode }
    })
    filterStore.wardOptions = res.data.map(item => ({
      label: item.name,
      value: item.code
    }))
  } finally {
    loadingWard.value = false
  }
}

const onApply = () => {
  filterStore.filters.location = { ...filterStore.filterLocation };
  filterStore.screen = 'main'
}

onMounted(() => {
  fetchListCity()
})

</script>
