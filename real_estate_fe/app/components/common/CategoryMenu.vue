<template>
  <nav class="flex flex-1">
    <ul class="flex items-center gap-1">
      <li
        v-for="cat in menuStore.categories"
        :key="cat.ID"
        class="relative group"
      >
        <button
          @click="handleClick(cat.Slug)"
          class="px-3 py-2 text-sm font-medium text-gray-700 hover:text-emerald-600 transition-colors"
        >
          {{ cat.Name }}
          <svg
            v-if="cat.children?.length"
            class="inline w-3 h-3 ml-0.5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M19 9l-7 7-7-7"
            />
          </svg>
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

const menuStore = useMenuStore();
onMounted(() => {
  menuStore.fetchMenuItems();
});

const handleClick = (slug: string) => {
  navigateTo(`/${slug}`);
};
</script>
