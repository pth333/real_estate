import { defineStore } from "pinia";

export interface ProjectDetail {
  id: number;
  name: string;
  slug: string;
  status: string;
  full_address: string;
  total_area_ha?: number;
  total_units?: number;
  price_min?: number;
  price_max?: number;
  view_count?: number;
  thumbnail?: string;
  description?: string;
}

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
        project.value = {
          ...res.data,
          thumbnail:
            res.data.thumbnail ||
            "https://placehold.co/600x400/e2e8f0/94a3b8?text=Project",
        };
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
