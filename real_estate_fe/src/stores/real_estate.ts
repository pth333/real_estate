import { ref } from 'vue'
import { defineStore } from 'pinia'
import type {
  DashboardSummary,
  Filter,
  RealEstateModel,
} from '@/types/real_estate'
import realEstateApi from '@/api/real_estate.api'

export const useRealEstateStore = defineStore('realEstate', () => {
  const items = ref<RealEstateModel[]>([])
  const total = ref(0)
  const loading = ref(false)
  const summary = ref<DashboardSummary | null>(null)
  const summaryLoading = ref(false)
  const page = ref(1)
  const pageSize = ref(50)
  const filter = ref<Filter>({})

  // Chuyển price từ VND sang tỷ
  function formattedPrice(priceVND: number): string {
    return (priceVND / 1_000_000_000).toFixed(2) + ' tỷ'
  }

  const payload = computed(() => ({
    page: page.value,
    size: pageSize.value,
    filter: filter.value,
  }))

  async function fetchList() {
    loading.value = true
    try {
      const res = await realEstateApi.getList(payload.value)
      items.value = res.data.data
      total.value = res.data.total
    } catch (e) {
      console.error('Lỗi tải danh sách BĐS:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchSummary(from?: string, to?: string) {
    summaryLoading.value = true
    try {
      const res = await realEstateApi.getSummary(from, to)
      summary.value = res.data
    } catch (e) {
      console.error('Lỗi tải summary:', e)
    } finally {
      summaryLoading.value = false
    }
  }

  function setFilter(newFilter: Filter) {
    filter.value = newFilter
    page.value = 1
    fetchList()
  }

  function setPage(newPage: number) {
    page.value = newPage
    fetchList()
  }

  function setPageSize(newSize: number) {
    pageSize.value = newSize
    page.value = 1
    fetchList()
  }

  return {
    items,
    total,
    loading,
    summary,
    summaryLoading,
    page,
    pageSize,
    filter,
    formattedPrice,
    fetchList,
    fetchSummary,
    setFilter,
    setPage,
    setPageSize,
  }
})
