<template>
  <header class="bg-white shadow-sm border-b border-gray-200">
    <div class="container mx-auto flex justify-between items-center p-2 gap-4">
      <!-- <h1 class="text-xl font-bold text-blue-600 whitespace-nowrap">
        RealEstate
      </h1> -->
      <CategoryMenu />
      <nav class="flex items-center gap-6 text-gray-700 whitespace-nowrap">
        <template v-if="auth.isAuthenticated">
          <NotificationBell />
          <n-button @click="handleLogout">
            Đăng xuất
          </n-button>

        </template>
        <template v-else>
          <NuxtLink to="/login" class="text-sm text-gray-600 hover:text-blue-600 transition">
            Đăng nhập
          </NuxtLink>
          <NuxtLink to="/register"
            class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700">
            Đăng ký
          </NuxtLink>
        </template>
      </nav>
      <n-button type="primary" @click="handlePushlishPost">Đăng tin</n-button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";

const auth = useAuthStore();

async function handleLogout() {
  await auth.logout();
  navigateTo("/login");
}

const handlePushlishPost = () => {
  navigateTo("/nguoi-ban/dang-tin");
}
</script>
