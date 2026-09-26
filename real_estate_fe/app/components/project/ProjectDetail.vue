<template>
  <n-layout class="bg-transparent">
    <!-- pb lớn ở mobile/tablet để thanh CTA cố định không che nội dung -->
    <n-layout-content class="mx-auto max-w-[1200px] bg-transparent px-4 py-4 pb-24 md:px-6 md:py-6 lg:pb-6">
      <n-breadcrumb class="mb-4 md:mb-5">
        <n-breadcrumb-item @click="navigateTo('/')">Trang chủ</n-breadcrumb-item>
        <n-breadcrumb-item @click="navigateTo('/du-an')">Dự án</n-breadcrumb-item>
        <n-breadcrumb-item>{{ project?.name || 'Chi tiết dự án' }}</n-breadcrumb-item>
      </n-breadcrumb>

      <!-- Khi lỗi/không tìm thấy -->
      <div v-if="loading" class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        <SkeletonCard :count="6" type="project" />
      </div>

      <!-- Layout chính: 1 cột mobile/tablet, 2 cột desktop (nội dung + liên hệ) -->
      <n-grid v-else-if="project" :cols="'1 md:2 lg:4'" :x-gap="'0 md:16 lg:24'" :y-gap="'16 lg:24'"
        responsive="screen">
        <!-- Cột trái: Thông tin chính (Chiếm 3/4 ở desktop) -->
        <n-grid-item :span="'1 md:2 lg:3'" class="min-w-0">
          <ProjectMainInfo :project="project"/>
          <ProjectListings :project="project" />
        </n-grid-item>

        <!-- Cột phải: Khung liên hệ tư vấn (Chiếm 1/4 ở desktop, nằm dưới ở mobile/tablet) -->
        <n-grid-item :span="'1 md:2 lg:1'" class="relative">
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
