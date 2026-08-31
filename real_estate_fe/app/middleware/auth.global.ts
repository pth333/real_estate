export default defineNuxtRouteMiddleware((to, _from) => {
  const authStore = useAuthStore()

  // Login và register không cần auth
  if (to.path === '/login' || to.path === '/register' || to.path === '/dang-nhap' || to.path === '/dang-ky') {
    // Nếu đã login, redirect về dashboard
    if (authStore.token) {
      return navigateTo('/')
    }
    return
  }

  // Các route khác cần auth
  if (!authStore.token) {
    return navigateTo('/dang-nhap')
  }
})
