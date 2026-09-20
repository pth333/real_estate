/**
 * Composable thao tác bất động sản yêu thích (favorite).
 * Gọi API backend POST /real-estate/favorite/:id (toggle thêm/bỏ).
 *
 * Trang danh sách/chi tiết BĐS là trang công khai nên khách vãng lai vẫn bấm được
 * nút yêu thích. Vì yêu thích là dữ liệu cá nhân, khách sẽ được nhắc đăng nhập và
 * quay lại đúng trang đang xem.
 */
import { useFavoriteService } from '~/services/favorite.service'
import { useAuthStore } from '~/stores/auth'

export function useFavorite() {
  const favoriteService = useFavoriteService()
  const authStore = useAuthStore()
  const route = useRoute()
  const toggling = ref(false)

  /** Nhắc đăng nhập rồi quay lại đúng trang hiện tại */
  function requireLogin(): null {
    window.message?.info('Vui lòng đăng nhập để lưu bất động sản yêu thích')
    navigateTo({ path: '/dang-nhap', query: { redirect: route.fullPath } })
    return null
  }

  /**
   * Toggle yêu thích 1 tin đăng. Trả về trạng thái mới (true = đã yêu thích,
   * false = đã bỏ), hoặc null nếu chưa đăng nhập / lỗi.
   */
  const toggle = async (id: number): Promise<boolean | null> => {
    if (!authStore.isAuthenticated) return requireLogin()

    try {
      toggling.value = true
      return await favoriteService.toggle(id)
    } catch (e) {
      return null
    } finally {
      toggling.value = false
    }
  }

  return { toggle, toggling }
}
