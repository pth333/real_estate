<template>
  <div class="flex flex-1 ">
    <n-input v-model:value="filterStore.searchKeyword" placeholder="Tìm kiếm bất động sản theo địa chỉ, quận, loại..."
      clearable size="large" @keyup.enter="handleSearch">
      <template #suffix>
        <n-button type="primary" size="small" @click="handleSearch">
          <IconSearch class="h-4 w-4" /> Tìm kiếm
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
