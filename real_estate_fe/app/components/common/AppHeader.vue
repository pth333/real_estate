<template>
  <header class="sticky top-0 z-50 bg-white shadow-sm">
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
          <n-button type="primary" tertiary size="small" class="text-sm text-gray-600 hover:text-emerald-600 transition"
            @click="$router.push('/login')">
            Đăng nhập
          </n-button>
          <n-button type="primary" size="small"
            class=" bg-emerald-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-700"
            @click="$router.push('/register')">
            Đăng ký
          </n-button>
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
