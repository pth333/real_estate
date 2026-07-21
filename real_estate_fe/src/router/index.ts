import { createRouter, createWebHistory } from 'vue-router'

import DashboardV1 from '@/components/pages/DashboardV1.vue'
import LoginV1 from '@/components/pages/LoginV1.vue'
import RegisterV1 from '@/components/pages/RegisterV1.vue'
import ListRealEstateByCategory from '@/components/pages/ListRealEstateByCategory.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),

  routes: [
    {
      path: '/',
      name: 'Dashboard',
      component: DashboardV1,
      meta: { requiresAuth: true },
    },
    {
      path: '/login',
      name: 'Login',
      component: LoginV1,
    },
    {
      path: '/register',
      name: 'Register',
      component: RegisterV1,
    },
    {
      path: '/:slug/:page(\\d+)?',
      name: 'ListRealEstateByCategory',
      component: ListRealEstateByCategory,
      meta: { requiresAuth: true },
    },
  ],
})

/** Guard: bảo vệ route cần auth + chuyển hướng nếu đã login */
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')

  // Chưa đăng nhập → redirect về login
  if (to.meta.requiresAuth && !token) {
    next({ name: 'Login' })
    return
  }

  // Đã đăng nhập → không cho vào login/register, đẩy về dashboard
  if (token && (to.name === 'Login' || to.name === 'Register')) {
    next({ name: 'Dashboard' })
    return
  }

  next()
})

export default router
