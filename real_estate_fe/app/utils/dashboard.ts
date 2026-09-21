import type { DashboardMenuItem } from '~/types/dashboard'

/**
 * Tìm mục menu đang mở theo đường dẫn hiện tại.
 *
 * Chọn mục có đoạn khớp DÀI NHẤT thay vì mục đầu tiên khớp, để xử lý đúng trường hợp
 * một mục dùng đoạn ngắn bao trùm mục khác: mục "Tổng quan escrow" khớp `/admin`,
 * nhưng `/admin/users` phải thuộc mục "Người dùng & phân quyền".
 *
 * Trả về undefined nếu không khớp mục nào (nơi gọi tự quyết định mặc định).
 */
export function findActiveItem(
  items: DashboardMenuItem[],
  path: string,
): DashboardMenuItem | undefined {
  let best: DashboardMenuItem | undefined
  let bestScore = 0

  for (const item of items) {
    let score = 0
    for (const pattern of item.match) {
      if (path.includes(pattern) && pattern.length > score) score = pattern.length
    }
    if (score > bestScore) {
      best = item
      bestScore = score
    }
  }

  return best
}
