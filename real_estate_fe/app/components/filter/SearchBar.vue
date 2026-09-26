<template>
  <div class="flex min-w-0 flex-1">
    <n-input v-model:value="filterStore.searchKeyword" placeholder="Tìm theo địa chỉ, quận, loại BĐS..."
      clearable size="large" @keyup.enter="handleSearch">
      <template #suffix>
        <!-- Mobile chỉ hiện icon cho gọn, từ tablet trở lên mới hiện chữ -->
        <n-button type="primary" size="small" @click="handleSearch">
          <IconSearch class="h-4 w-4" />
          <span class="ml-1 hidden md:inline">Tìm kiếm</span>
        </n-button>
      </template>
    </n-input>

  </div>
</template>

<script setup lang="ts">
import { useFilterStore } from '~/stores/filter';

const emit = defineEmits<{
  (e: 'search'): void;
}>();

const filterStore = useFilterStore();

function handleSearch() {
  emit('search');
}
</script>
