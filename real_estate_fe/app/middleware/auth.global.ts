/**
 * Middleware chạy cho MỌI điều hướng (client-side).
 *
 * MẶC ĐỊNH MỌI TRANG LÀ CÔNG KHAI — khách vãng lai (chưa đăng ký) vẫn xem web bình
 * thường: trang chủ, danh sách BĐS, chi tiết BĐS, dự án, danh mục...
 * Chỉ trang khai báo meta mới bị chặn:
 *   - requiresAuth: true                → phải đăng nhập
 *   - requiresRole: ['ADMIN']           → phải có 1 trong các role
 *   - requiresPermission: ['...']       → phải có 1 trong các permission
 *
 * Chặn ở đây là để gõ tay URL cũng không vào được page (UX), KHÔNG phải ranh giới
 * bảo mật — backend vẫn kiểm tra lại quyền cho từng API.
 */
export default defineNuxtRouteMiddleware(async (to) => {
  const authStore = useAuthStore()

  // ── Trang đăng nhập / đăng ký ──
  const guestPaths = ['/login', '/register', '/dang-nhap', '/dang-ky']
  if (guestPaths.includes(to.path)) {
    // Đã đăng nhập rồi thì không cần vào trang đăng nhập nữa
    if (authStore.token) return navigateTo('/')
    return
  }

  const requiredRoles = to.meta.requiresRole ?? []
  const requiredPermissions = to.meta.requiresPermission ?? []
  const needsLogin =
    to.meta.requiresAuth === true ||
    requiredRoles.length > 0 ||
    requiredPermissions.length > 0

  // ── Trang công khai: khách vãng lai vào thẳng, không gọi API quyền ──
  if (!needsLogin) return

  // ── Trang cần đăng nhập ──
  if (!authStore.token) {
    // Nhớ trang đang muốn vào để đăng nhập xong quay lại đúng chỗ
    return navigateTo({ path: '/dang-nhap', query: { redirect: to.fullPath } })
  }

  // Nạp quyền mới nhất (có chống gọi trùng) trước khi kiểm tra meta,
  // vì user đọc từ cookie có thể đã cũ sau khi admin đổi role.
  await authStore.ensureAccessLoaded()
  if (!authStore.token) {
    return navigateTo({ path: '/dang-nhap', query: { redirect: to.fullPath } })
  }

  // ── Chặn theo role khai báo trong definePageMeta ──
  if (requiredRoles.length && !authStore.hasAnyRole(requiredRoles)) {
    return navigateTo('/forbidden')
  }

  // ── Chặn theo permission chi tiết ──
  if (requiredPermissions.length && !authStore.hasAnyPermission(requiredPermissions)) {
    return navigateTo('/forbidden')
  }
})
