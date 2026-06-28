import { defineStore } from 'pinia'
import type {
  DashboardSummary,
  Filter,
  PaginatedResponse,
  RealEstateModel,
} from '@/types/real_estate'
import realEstateApi from '@/api/real_estate.api'

interface RealEstateState {
  items: RealEstateModel[]
  total: number
  loading: boolean
  summary: DashboardSummary | null
  summaryLoading: boolean
  page: number
  pageSize: number
  filter: Filter
}

export const useRealEstateStore = defineStore('realEstate', {
  state: (): RealEstateState => ({
    items: [],
    total: 0,
    loading: false,
    summary: null,
    summaryLoading: false,
    page: 1,
    pageSize: 50,
    filter: {},
  }),

  getters: {
    // Chuyển price từ VND sang tỷ
    formattedPrice:
      () =>
      (priceVND: number): string => {
        return (priceVND / 1_000_000_000).toFixed(2) + ' tỷ'
      },
  },

  actions: {
    async fetchList() {
      this.loading = true
      try {
        const res = await realEstateApi.getList({
          page: this.page,
          size: this.pageSize,
          filter: this.filter,
        })
        this.items = res.data.data
        this.total = res.data.total
      } catch (e) {
        console.error('Lỗi tải danh sách BĐS:', e)
      } finally {
        this.loading = false
      }
    },

    async fetchSummary(from?: string, to?: string) {
      this.summaryLoading = true
      try {
        const res = await realEstateApi.getSummary(from, to)
        this.summary = res.data
      } catch (e) {
        console.error('Lỗi tải summary:', e)
      } finally {
        this.summaryLoading = false
      }
    },

    setFilter(filter: Filter) {
      this.filter = filter
      this.page = 1
      this.fetchList()
    },

    setPage(page: number) {
      this.page = page
      this.fetchList()
    },
  },
})
