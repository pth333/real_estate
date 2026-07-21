<template>

  <div class="list-real-estate-page">
     <!-- Header với tiêu đề và filters -->
    <div class="page-header">

      <h1 class="page-title">{{ pageTitle }}</h1>
       <!-- Filters -->
      <div class="filters">

        <div class="filter-group">
           <label>Giá tối thiểu:</label> <input
            v-model.number="filters.min_price"
            type="number"
            placeholder="VD: 1000000000"
            @change="handleFilterChange"
          />
        </div>

        <div class="filter-group">
           <label>Giá tối đa:</label> <input
            v-model.number="filters.max_price"
            type="number"
            placeholder="VD: 5000000000"
            @change="handleFilterChange"
          />
        </div>

        <div class="filter-group">
           <label>Quận/Huyện:</label> <input
            v-model="filters.district"
            type="text"
            placeholder="VD: Quận 1"
            @change="handleFilterChange"
          />
        </div>
         <button class="btn-reset" @click="resetFilters">Reset</button>
      </div>

    </div>
     <!-- Loading state -->
    <div v-if="loading" class="loading">

      <div class="spinner"></div>

      <p>Đang tải dữ liệu...</p>

    </div>
     <!-- Empty state -->
    <div v-else-if="realEstates.length === 0" class="empty">

      <p>Không tìm thấy bất động sản nào</p>

    </div>
     <!-- Grid hiển thị cards -->
    <div v-else class="cards-grid">
       <RealEstateCard
        v-for="estate in realEstates"
        :key="estate.ID"
        :estate="estate"
        @call="handleCall"
        @toggle-favorite="handleToggleFavorite"
      />
    </div>
     <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination">
       <button :disabled="currentPage === 1" @click="goToPage(currentPage - 1)"> ← Trước </button>
      <div class="page-numbers">
         <button
          v-for="page in visiblePages"
          :key="page"
          :class="{ active: page === currentPage }"
          @click="goToPage(page)"
        >
           {{ page }} </button
        >
      </div>
       <button :disabled="currentPage === totalPages" @click="goToPage(currentPage + 1)">
         Sau → </button
      >
    </div>

  </div>

</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { general } from '@/api/general.api'
import type { RealEstateModel, Filter } from '@/types/real_estate'
import RealEstateCard from '@/components/real-estate/RealEstateCard.vue'
import { useToast } from '@/composables/useToast'

const route = useRoute()

const { showToast } = useToast()

// State
const realEstates = ref<RealEstateModel[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const currentPage = ref(1)
const pageSize = ref(10)
const totalRecords = ref(0)

// Keyset pagination: lưu cursor_id cho mỗi trang đã xem (0 = trang đầu, không gửi cursor)
const cursorMap = ref<Record<number, number>>({ 1: 0 })

// Filters
const filters = ref<Filter>({
  min_price: undefined,
  max_price: undefined,
  district: undefined,
})

// Computed
const pageTitle = computed(() => {
  const slug = route.params.slug as string
  return slug ? `Danh sách BĐS - ${slug}` : 'Danh sách Bất Động Sản'
})

const totalPages = computed(() => {
  return Math.ceil(totalRecords.value / pageSize.value)
})

// Hiển thị tối đa 5 trang trong pagination
const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value
  const pages: number[] = []

  if (total <= 5) {
    // Nếu tổng số trang <= 5, hiển thị tất cả
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    // Hiển thị 5 trang xung quanh trang hiện tại
    let start = Math.max(1, current - 2)
    const end = Math.min(total, start + 4)

    // Điều chỉnh start nếu end đã chạm giới hạn
    if (end - start < 4) {
      start = Math.max(1, end - 4)
    }

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }
  }

  return pages
})

// Methods
const fetchData = async () => {
  loading.value = true
  error.value = null
  console.log('=== fetchData === page:', currentPage.value, 'cursorMap:', cursorMap.value)

  try {
    const cursorId = cursorMap.value[currentPage.value] ?? 0
    const payload = {
      cursor_id: cursorId > 0 ? cursorId : undefined,
      size: pageSize.value,
      filter: filters.value,
    }
    const res = await general.GetListByCategory(
      payload,
      route.params.slug as string,
      currentPage.value,
    )

    realEstates.value = res.data?.data || []
    totalRecords.value = res.data?.total || 0

    // Lưu id cuối để dùng làm cursor cho trang kế tiếp
    const items = res.data?.data || []
    if (items.length > 0) {
      cursorMap.value[currentPage.value + 1] = items[items.length - 1].ID
    }
  } catch (err) {
    const msg = err instanceof Error ? err.message : 'Có lỗi xảy ra khi tải dữ liệu'
    error.value = msg
    showToast('error', msg)
    console.error('Error fetching real estates:', err)
  } finally {
    loading.value = false
  }
}

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchData()
    // Scroll lên đầu trang
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

const handleFilterChange = () => {
  // Reset về trang 1 khi filter thay đổi
  cursorMap.value = { 1: 0 }
  currentPage.value = 1
  fetchData()
}

const resetFilters = () => {
  filters.value = {
    min_price: undefined,
    max_price: undefined,
    district: undefined,
  }
  cursorMap.value = { 1: 0 }
  currentPage.value = 1
  fetchData()
}

const handleCall = (phone: string) => {
  // Mở app gọi điện
  window.open(`tel:${phone}`, '_self')
}

const handleToggleFavorite = (id: number) => {
  // TODO: Gọi API để toggle favorite
  console.log('Toggle favorite for ID:', id)

  // Cập nhật UI optimistically
  const estate = realEstates.value.find((e) => e.ID === id)
  if (estate) {
    estate.IsFavorite = !estate.IsFavorite
  }
}

// Lifecycle
onMounted(() => {
  fetchData()
})

// Watch route params để reload data khi slug thay đổi
watch(
  () => route.params.slug,
  () => {
    cursorMap.value = { 1: 0 }
    currentPage.value = 1
    fetchData()
  },
)
</script>

<style scoped>
.list-real-estate-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 28px;
  font-weight: bold;
  color: #333;
  margin-bottom: 16px;
  text-transform: uppercase;
}

.filters {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  padding: 16px;
  background: #f9f9f9;
  border-radius: 8px;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.filter-group label {
  font-size: 13px;
  font-weight: 600;
  color: #666;
}

.filter-group input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  min-width: 180px;
}

.btn-reset {
  align-self: flex-end;
  padding: 8px 16px;
  background: #666;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
}

.btn-reset:hover {
  background: #555;
}

.loading,
.error,
.empty {
  text-align: center;
  padding: 64px 24px;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #f0f0f0;
  border-top-color: #4a90e2;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error p {
  color: #e63946;
  font-size: 16px;
  margin-bottom: 16px;
}

.error button {
  padding: 8px 24px;
  background: #4a90e2;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
}

.empty p {
  font-size: 16px;
  color: #999;
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-top: 32px;
}

.pagination button {
  padding: 8px 16px;
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s ease;
}

.pagination button:hover:not(:disabled) {
  background: #f0f0f0;
  border-color: #4a90e2;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-numbers {
  display: flex;
  gap: 4px;
}

.page-numbers button.active {
  background: #4a90e2;
  color: #fff;
  border-color: #4a90e2;
}

@media (max-width: 768px) {
  .cards-grid {
    grid-template-columns: 1fr;
  }

  .filters {
    flex-direction: column;
  }

  .filter-group input {
    min-width: 100%;
  }
}
</style>

