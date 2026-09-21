import type { Component } from 'vue'

/**
 * Một mục menu trong khu vực quản lý (dùng chung cho trang môi giới và admin).
 * `permission` khớp với meta `requiresPermission` của trang tương ứng để menu và
 * quyền vào trang luôn nhất quán.
 */
export interface DashboardMenuItem {
  /** Khoá định danh mục menu */
  key: string
  /** Nhãn hiển thị trên menu */
  label: string
  /** Tiêu đề trang (hiện ở topbar và breadcrumb) */
  pageTitle: string
  /** Đường dẫn điều hướng khi bấm */
  path: string
  /** Icon của mục */
  icon: Component
  /** Quyền tối thiểu để thấy mục (bỏ trống = ai đã đăng nhập cũng thấy) */
  permission?: string
  /** Các đoạn đường dẫn dùng để nhận biết mục đang mở */
  match: string[]
  /**
   * Mục chỉ dùng để nhận biết trang (tiêu đề + highlight menu cha), KHÔNG hiện trên sidebar.
   * Dùng cho trang con như "Tạo dự án mới".
   */
  hidden?: boolean
}

/** Nhóm menu (một khu vực có thể có nhiều nhóm) */
export interface DashboardMenuGroup {
  label: string
  items: DashboardMenuItem[]
}
