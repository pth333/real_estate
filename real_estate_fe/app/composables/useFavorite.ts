/**
 * Composable thao tác bất động sản yêu thích (favorite).
 * Gọi API backend POST /real-estate/favorite/:id (toggle thêm/bỏ).
 */

export function useFavorite() {
  const { $api } = useNuxtApp()
  const toggling = ref(false)

  /**
   * Toggle yêu thích 1 tin đăng. Trả về trạng thái mới (true = đã yêu thích,
   * false = đã bỏ), hoặc null nếu lỗi (vd chưa đăng nhập).
   */
  const toggle = async (id: number): Promise<boolean | null> => {
    try {
      toggling.value = true
      const res = await $api.post<{ data: { is_favorite: boolean } }>(
        `/real-estate/favorite/${id}`,
      )
      return res?.data?.is_favorite ?? null
    } catch (e) {
      return null
    } finally {
      toggling.value = false
    }
  }

  return { toggle, toggling }
}
