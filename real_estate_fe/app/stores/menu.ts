import { defineStore } from "pinia";
import type { Category, MenuResponse } from "@/types/menu";

export const useMenuStore = defineStore("menu", () => {
  const categories = ref<Category[]>([]);

  const fetchMenuItems = async () => {
    try {
      const { $api } = useNuxtApp();
      const res = await $api.get<MenuResponse>("/category");
      categories.value = (res as any).data || [];
    } catch (error) {
      console.error("Error fetching menu items:", error);
    }
  };

  return {
    categories,
    fetchMenuItems,
  };
});
