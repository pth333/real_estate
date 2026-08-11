<template>
  <nav class="flex flex-1">
    <ul class="flex items-center gap-1">
      <li
        v-for="cat in menuStore.menu"
        :key="cat.ID"
        class="relative group"
      >
        <button
          @click="handleClick(cat.Slug)"
          class="px-3 py-2 text-sm font-medium text-gray-700 hover:text-emerald-600 hover:underline underline-offset-4 transition-colors"
        >
          {{ cat.Name }}
        </button>
        <ul
          v-if="cat.children?.length"
          class="absolute left-0 top-full z-50 hidden group-hover:block bg-white shadow-lg py-1 min-w-[200px]"
        >
          <CategoryMenuItem
            v-for="child in cat.children"
            :key="child.ID"
            :item="child"
          />
        </ul>
      </li>
    </ul>
  </nav>
</template>

<script setup lang="ts">
import { useMenuStore } from "~/stores/menu";
import { useRealEstateStore } from "~/stores/real_estate";
import { Menu } from "~/types/window";

const menuStore = useMenuStore();
const realEstateStore = useRealEstateStore();
window.menu = new Menu(menuStore.user_id);

onMounted(() => {
  menuStore.fetchMenuItems();
});

const handleClick = (slug: string) => {
  // Lưu pure category slug để các filter build URL SEO đúng
  realEstateStore.categorySlug = slug;
  navigateTo(`/${slug}`);
};
</script>
