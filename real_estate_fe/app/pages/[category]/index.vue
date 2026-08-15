<template>
  <div>
    <RealEstateDetail v-if="detailId > 0" :key="detailId" :id="detailId" />
    <ProjectDetail v-if="projectId > 0" :key="projectId" :id="projectId" />


    <RealEstateList v-if="pageType === 'real_estate'" :category-slug="categorySlug" />
    <ProjectList v-if="pageType === 'project'" :category-slug="categorySlug" />
  </div>
</template>

<script setup lang="ts">
import { useMenuStore } from '~/stores/menu'
import type { Category } from '~/types/menu'
import RealEstateList from '~/components/real-estate/RealEstateList.vue'
import ProjectList from '~/components/project/ProjectList.vue'

const route = useRoute()
const menuStore = useMenuStore()

const categorySlug = computed<string>(() => {
  const v = route.params.category
  return Array.isArray(v) ? v[0] ?? '' : v ?? ''
})

const detailId = computed<number>(() => {
  const slug = categorySlug.value
  if (!slug) return 0
  const match = slug.match(/-rs(\d+)$/)
  return match ? Number(match[1]) : 0
})

const projectId = computed<number>(() => {
  const slug = categorySlug.value
  if (!slug) return 0
  const match = slug.match(/-pj(\d+)$/)
  return match ? Number(match[1]) : 0
})

const findCategory = (categories: Category[], slug: string): Category | null => {
  let bestMatch: Category | null = null

  const traverse = (list: Category[]) => {
    for (const cat of list) {
      // Điều kiện khớp: Khớp hoàn toàn OR khớp tiền tố và có dấu gạch ngang phân tách (ví dụ "nha-dat-ban-ha-noi" khớp "nha-dat-ban")
      const isMatch = cat.Slug === slug || slug.startsWith(cat.Slug + '-')

      if (isMatch) {
        // Ưu tiên danh mục có Slug dài hơn (độ khớp chính xác cao hơn, tránh khớp nhầm nha-dat thay vì nha-dat-ban)
        if (!bestMatch || cat.Slug.length > bestMatch.Slug.length) {
          bestMatch = cat
        }
      }

      if (cat.children && cat.children.length > 0) {
        traverse(cat.children)
      }
    }
  }

  traverse(categories)
  return bestMatch
}

// Phân giải pageType trực tiếp từ Menu Store
const pageType = computed<'real_estate' | 'project' | null>(() => {
  if (detailId.value > 0) return null

  const categories = menuStore.menu?.categories || []
  const category = findCategory(categories, categorySlug.value)
  if (!category) return null

  return category.Type === 'project' ? 'project' : 'real_estate'
})
</script>
