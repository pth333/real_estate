/**
 * Composable quản lý phân trang dùng chung
 * - page, pageSize, total, totalPages
 * - goToPage, setPageSize
 * - triggers fetchFn mỗi khi page/pageSize thay đổi
 */
export function usePagination(fetchFn: () => Promise<void>) {
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  const totalPages = computed(() =>
    Math.ceil(total.value / pageSize.value),
  )

  /** Chuyển trang, gọi fetch */
  function goToPage(p: number) {
    if (p < 1 || p > totalPages.value) return
    page.value = p
    fetchFn()
  }

  /** Đổi số dòng/trang, reset về trang 1 */
  function setPageSize(newSize: number) {
    pageSize.value = newSize
    page.value = 1
    fetchFn()
  }

  /** Reset về đầu */
  function reset() {
    page.value = 1
    total.value = 0
  }

  return {
    page, pageSize, total, totalPages,
    goToPage, setPageSize, reset,
  }
}
