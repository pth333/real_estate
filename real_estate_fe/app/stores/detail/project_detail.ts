import { defineStore } from "pinia";
import type { ProjectDetail } from "~/types/project";

export const useProjectDetail = defineStore("project_detail", () => {
  const loading = ref(false);
  const project = ref<ProjectDetail>();

  async function fetchDetail(id: number) {
    const { $api } = useNuxtApp();
    loading.value = true;
    try {
      const res = await $api.get<{ data: ProjectDetail }>(
        `/real-estate/project/detail/${id}`,
      );
      if (res.data) {
        project.value = res.data;
      }
    } catch (err) {
      console.error("Lỗi khi tải thông tin dự án:", err);
    } finally {
      loading.value = false;
    }
  }

  async function incrementView(id: number) {
    const { $api } = useNuxtApp();
    try {
      await $api.post(`/real-estate/project/view/${id}`);
    } catch (err) {
      console.error("Lỗi khi tăng lượt xem dự án:", err);
    }
  }

  return {
    loading,
    project,
    fetchDetail,
    incrementView,
  };
});
