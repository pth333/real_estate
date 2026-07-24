/**
 * Format thời gian tương đối: "Vừa xong", "5 phút trước", "Đăng hôm qua"
 */
export function useRelativeTime() {
  /**
   * Tính "cách đây bao lâu" từ date string
   */
  function fromNow(dateStr: string): string {
    if (!dateStr) return ''
    const d = new Date(dateStr)
    const now = new Date()
    const diffMs = now.getTime() - d.getTime()
    const diffMin = Math.floor(diffMs / 60000)

    if (diffMin < 1) return 'Vừa xong'
    if (diffMin < 60) return `${diffMin} phút trước`

    const diffHour = Math.floor(diffMin / 60)
    if (diffHour < 24) return `${diffHour} giờ trước`

    const diffDays = Math.floor(diffHour / 24)
    if (diffDays === 0) return 'Đăng hôm nay'
    if (diffDays === 1) return 'Đăng hôm qua'
    if (diffDays < 7) return `Đăng ${diffDays} ngày trước`

    return d.toLocaleDateString('vi-VN')
  }

  return { fromNow }
}
