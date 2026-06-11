import { createRouter, createWebHistory } from 'vue-router'

import DashboardV1 from '@/components/pages/DashboardV1.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),

  routes: [
    {
      path: '/',
      name: 'Dashboard',
      component: DashboardV1,
    },
  ],
})

export default router
