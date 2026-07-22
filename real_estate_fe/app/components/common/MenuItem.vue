<template>
  <li class="relative group">
    <button
      @click="go(item.Slug)"
      class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-emerald-50 hover:text-emerald-700 whitespace-nowrap"
    >
      {{ item.Name }}
      <svg
        v-if="item.children?.length"
        class="inline w-3 h-3 ml-1"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M10 6l6 6-6 6"
        />
      </svg>
    </button>
    <ul
      v-if="item.children?.length"
      class="absolute left-full top-0 z-50 hidden group-hover:block bg-white shadow-lg border rounded-lg py-1 min-w-[200px]"
    >
      <MenuItem v-for="child in item.children" :key="child.ID" :item="child" />
    </ul>
  </li>
</template>

<script setup lang="ts">
import type { Category } from '~/types/menu';

defineProps<{ item: Category }>();

function go(slug: string) {
  navigateTo(`/${slug}`);
}
</script>
