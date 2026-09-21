/**
 * Thông tin thương hiệu dùng chung — một nguồn duy nhất cho tên và logo,
 * tránh mỗi nơi ghi một kiểu (header, trang đăng nhập/đăng ký, title...).
 */

export const BRAND_NAME = 'NhàViệt'

/** Nửa sau của tên thương hiệu được tô màu nhấn (Nhà + Việt) */
export const BRAND_NAME_ACCENT = 'Việt'

export const BRAND_TAGLINE = 'Nền tảng bất động sản'

export const LOGO_PATH = '/logo.svg'

/** Đổi màu logo.svg sang tông emerald của thương hiệu */
export const LOGO_FILTER =
  'invert(48%) sepia(79%) saturate(476%) hue-rotate(114deg) brightness(95%) contrast(95%)'

/** Tiêu đề trang theo thương hiệu: "Đăng nhập | NhàViệt" */
export function brandTitle(page: string): string {
  return `${page} | ${BRAND_NAME}`
}
