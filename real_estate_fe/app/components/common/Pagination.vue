<template>
  <div v-if="totalPages > 1" class="flex items-center justify-center gap-2">
    <button
      :disabled="currentPage <= 1"
      class="cursor-pointer rounded border border-gray-300 bg-white px-4 py-2 text-sm font-semibold text-gray-500 transition hover:border-blue-500 hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
      @click="handlePage(currentPage - 1)"
    >
      ← Trước
    </button>

    <div class="flex gap-1">
      <button
        v-for="page in visiblePages"
        :key="page"
        class="cursor-pointer rounded border px-4 py-2 text-sm font-semibold transition"
        :class="
          page === currentPage
            ? 'border-blue-500 bg-blue-500 text-white'
            : 'border-gray-300 bg-white text-gray-500 hover:border-blue-500 hover:bg-gray-50'
        "
        @click="handlePage(page)"
      >
        {{ page }}
      </button>
    </div>

    <button
      :disabled="currentPage >= totalPages"
      class="cursor-pointer rounded border border-gray-300 bg-white px-4 py-2 text-sm font-semibold text-gray-500 transition hover:border-blue-500 hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
      @click="handlePage(currentPage + 1)"
    >
      Sau →
    </button>
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

const visiblePages = computed(() => {
  const total = props.totalPages
  const current = props.currentPage
  const pages: number[] = []
  if (total <= 5) {
    for (let i = 1; i <= total; i++) pages.push(i)
  } else {
    let start = Math.max(1, current - 2)
    const end = Math.min(total, start + 4)
    if (end - start < 4) start = Math.max(1, end - 4)
    for (let i = start; i <= end; i++) pages.push(i)
  }
  return pages
})

function handlePage(page: number) {
  if (page >= 1 && page <= props.totalPages) {
    emit("pageChange", page)
  }
}
</script>
