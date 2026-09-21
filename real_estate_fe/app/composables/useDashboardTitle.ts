/**
 * Cho phép TRANG CON ghi đè tiêu đề hiển thị ở khu vực quản lý (topbar + breadcrumb).
 *
 * Cần thiết khi tiêu đề phụ thuộc dữ liệu của trang, ví dụ cùng 1 trang nhưng
 * "Tạo dự án mới" hay "Chỉnh sửa dự án" tuỳ theo có tham số id hay không.
 * Khu vực quản lý tự xoá ghi đè khi đổi trang.
 */
export function useDashboardTitle() {
  const titleOverride = useState<string | null>('dashboard:title-override', () => null)

  function setTitle(title: string) {
    titleOverride.value = title
  }

  function clearTitle() {
    titleOverride.value = null
  }

  return { titleOverride, setTitle, clearTitle }
}
