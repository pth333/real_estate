<template>
    <RealEstateList :category-slug="categorySlug" />
</template>
<script setup lang="ts">
const route = useRoute()

const categorySlug = computed<string>(() => {
    const v = route.params.category
    return Array.isArray(v) ? v[0] ?? '' : v ?? ''
})
</script>