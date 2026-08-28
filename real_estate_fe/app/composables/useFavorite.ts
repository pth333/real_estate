/**
 * Composable thao tác bất động sản yêu thích (favorite).
 * Gọi API backend POST /real-estate/favorite/:id (toggle thêm/bỏ).
 * Khi THÊM mới yêu thích → hiện modal xác nhận (NDialog).
 */

import { useDialog } from 'naive-ui'

export function useFavorite() {
  const { $api } = useNuxtApp()
  const dialog = useDialog()
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

  /**
   * Toggle kèm modal xác nhận khi THÊM (chưa yêu thích).
   * Nếu đang yêu thích (bỏ) → thao tác luôn, không cần xác nhận.
   */
  const toggleWithConfirm = (id: number, isFavorite: boolean): Promise<boolean | null> => {
    if (isFavorite) {
      return toggle(id)
    }
    return new Promise((resolve) => {
      dialog.warning({
        title: 'Thêm vào danh mục yêu thích',
        content: 'Bạn có muốn thêm bất động sản này vào danh mục yêu thích?',
        positiveText: 'Thêm',
        negativeText: 'Hủy',
        onPositiveClick: () => toggle(id).then(resolve),
        onNegativeClick: () => resolve(null),
      })
    })
  }

  return { toggle, toggleWithConfirm, toggling }
}
