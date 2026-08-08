<template>
  <!-- Route /{category} — nếu slug kết thúc "-rs{id}" → trang chi tiết, ngược lại trang list -->
  <RealEstateDetail v-if="isDetailSlug" :key="currentSlug" />
  <RealEstateList v-else />
</template>

<script setup lang="ts">
const route = useRoute();

// Slug 1 đoạn kết thúc "-rs{id}" là trang chi tiết (VD "nha-pho-...-rs123").
const currentSlug = computed<string>(() => {
  const v = route.params.category;
  return Array.isArray(v) ? v[0] ?? '' : v ?? '';
});


const isDetailSlug = computed<boolean>(() => /-rs\d+$/.test(currentSlug.value));
</script>