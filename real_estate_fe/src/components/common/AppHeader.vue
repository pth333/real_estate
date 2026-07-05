<template>
  <header class="bg-white shadow-sm border-b border-gray-200">
    <div class="container mx-auto flex justify-between items-center p-4">
      <h1 class="text-xl font-bold text-blue-600">RealEstate</h1>

      <nav class="flex items-center gap-6 text-gray-700">
        <!-- Notification (chỉ hiển thị khi đã login) -->
        <template v-if="auth.isAuthenticated">
          <span class="text-sm font-medium text-gray-600">{{ auth.userName || auth.userEmail }}</span>
          <NotificationBell />
          <button
            @click="handleLogout"
            class="text-sm text-red-500 hover:text-red-700 transition cursor-pointer"
          >
            Đăng xuất
          </button>
        </template>

        <!-- Chưa login: hiển thị nút Đăng nhập / Đăng ký -->
        <template v-else>
          <router-link
            :to="{ name: 'Login' }"
            class="text-sm text-gray-600 hover:text-blue-600 transition"
          >
            Đăng nhập
          </router-link>
          <router-link
            :to="{ name: 'Register' }"
            class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700"
          >
            Đăng ký
          </router-link>
        </template>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NotificationBell from '@/components/notification/NotificationBell.vue'

const router = useRouter()
const auth = useAuthStore()

function handleLogout() {
  auth.logout()
  router.push({ name: 'Login' })
}
</script>
