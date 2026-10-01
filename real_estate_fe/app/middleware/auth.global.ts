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

  // ── Nạp phiên TRƯỚC khi kiểm tra ──
  // init() tự lo thứ tự đúng: nếu access token hết hạn/mất nhưng còn cookie refresh_token
  // thì gọi /auth/refresh trước, rồi mới lấy user info (roles/permissions).
  // Có chống gọi trùng nên GlobalInit gọi cùng lúc cũng không phát sinh request thứ 2,
  // và khách vãng lai (không có phiên) thì hàm trả về ngay, không gọi API nào.
  await authStore.init()

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

  // ── Trang công khai: khách vãng lai vào thẳng ──
  if (!needsLogin) return

  // ── Trang cần đăng nhập: init() đã cố refresh xong, giờ không có token nghĩa là hết phiên ──
  if (!authStore.token) {
    // Nhớ trang đang muốn vào để đăng nhập xong quay lại đúng chỗ
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
