<template>
  <n-layout class="bg-transparent">
    <n-layout-content class="mx-auto max-w-[1200px] px-6 py-6 bg-transparent">
      <n-breadcrumb class="mb-5">
        <n-breadcrumb-item @click="navigateTo('/')">Trang chủ</n-breadcrumb-item>
        <n-breadcrumb-item @click="navigateTo('/du-an')">Dự án</n-breadcrumb-item>
        <n-breadcrumb-item>{{ project?.name || 'Chi tiết dự án' }}</n-breadcrumb-item>
      </n-breadcrumb>

      <!-- Khi lỗi/không tìm thấy -->
      <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
        <SkeletonCard :count="6" type="project" />
      </div>

      <!-- Layout chính chia cột - Chỉ hiển thị khi project đã load thành công (khác null) -->
      <n-grid v-else-if="project" :cols="4" :x-gap="24" :y-gap="24" item-responsive>
        <!-- Cột trái: Thông tin chính (Chiếm 3/4) -->
        <n-grid-item :span="3" class="min-w-0">
          <ProjectMainInfo :project="project"/>
          <ProjectListings :project="project" />
        </n-grid-item>

        <!-- Cột phải: Khung liên hệ tư vấn (Chiếm 1/4) -->
        <n-grid-item :span="1" class="relative">
          <ProjectContactSidebar />
        </n-grid-item>
      </n-grid>
    </n-layout-content>
  </n-layout>
</template>

<script setup lang="ts">
import { useProjectDetail } from '~/stores/detail/project_detail';

const props = defineProps<{
  id: number
}>()

const store = useProjectDetail();

const { loading, project } = storeToRefs(store);

const projectTitle = computed(() => project.value?.name || "Chi tiết dự án");
useHead({
  title: projectTitle,
});

onMounted(async () => {
  await store.incrementView(props.id);
  store.fetchDetail(props.id);
});
</script>
