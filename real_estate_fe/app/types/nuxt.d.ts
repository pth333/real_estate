// types/nuxt.d.ts
import type { api } from '~/plugins/api.client'

declare module '#app' {
  interface NuxtApp {
    $api: typeof api
  }

  /**
   * Meta phân quyền ở tầng route — khai báo trong definePageMeta của page:
   *   definePageMeta({ requiresAuth: true })
   *   definePageMeta({ requiresRole: ['ADMIN'] })
   *   definePageMeta({ requiresPermission: ['broker.deposit.list'] })
   *
   * MẶC ĐỊNH MỌI TRANG LÀ CÔNG KHAI: khách vãng lai (chưa đăng ký) vẫn xem web bình
   * thường. Chỉ những trang khai báo 1 trong 3 meta trên mới bị middleware chặn.
   * requiresRole / requiresPermission cũng ngầm hiểu là phải đăng nhập.
   */
  interface PageMeta {
    /** Bắt buộc đăng nhập (dùng cho trang cá nhân chưa cần phân quyền chi tiết) */
    requiresAuth?: boolean
    requiresRole?: string[]
    requiresPermission?: string[]
  }
}

declare module 'vue' {
  interface ComponentCustomProperties {
    $api: typeof api
  }
}

export {}