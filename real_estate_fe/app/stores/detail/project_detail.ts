import { defineStore } from "pinia";
import type { ProjectDetail } from "~/types/project";
import { useProjectService } from "~/services/project.service";

export const useProjectDetail = defineStore("project_detail", () => {
  const loading = ref(false);
  const project = ref<ProjectDetail>();

  async function fetchDetail(id: number) {
    const projectService = useProjectService();
    loading.value = true;
    try {
      const data = await projectService.getDetail(id);
      if (data) {
        project.value = data;
      }
    } catch (err) {
      console.error("Lỗi khi tải thông tin dự án:", err);
    } finally {
      loading.value = false;
    }
  }

  async function incrementView(id: number) {
    const projectService = useProjectService();
    try {
      await projectService.incrementView(id);
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
