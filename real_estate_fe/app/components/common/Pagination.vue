<template>
  <div v-if="totalPages > 1" class="flex justify-center">
    <n-pagination
      :page="currentPage"
      :page-count="totalPages"
      :page-slot="5"
      @update:page="handlePage"
    />
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  currentPage: number
  totalPages: number
}>()

const emit = defineEmits<{
  pageChange: [page: number]
}>()

function handlePage(page: number) {
  if (page >= 1 && page <= props.totalPages) {
    emit("pageChange", page)
  }
}
</script>
