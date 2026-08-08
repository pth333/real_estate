<template>
  <!-- Route /{category} — nếu slug kết thúc "-rs{id}" → trang chi tiết (gọi API theo id), ngược lại trang list -->
  <RealEstateDetail v-if="detailId > 0" :key="detailId" :id="detailId" />
  <RealEstateList v-else />
</template>

<script setup lang="ts">
const route = useRoute();

// Slug 1 đoạn kết thúc "-rs{id}" là trang chi tiết (VD "nha-pho-...-rs123").
// Dựa vào đuôi "-rs{id}" để tách nhánh và lấy id truyền vào GetByDetail.
const detailId = computed<number>(() => {
  const v = route.params.category;
  const slug = Array.isArray(v) ? v[0] ?? '' : v ?? '';
  if (!slug) return 0;
  const match = slug.match(/-rs(\d+)$/);
  return match ? Number(match[1]) : 0;
});
</script>