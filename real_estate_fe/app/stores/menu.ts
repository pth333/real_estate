import { defineStore } from "pinia";
import type { Category, MenuResponse } from "~/types/menu";

export const useMenuStore = defineStore("menu", () => {
  const menu = ref<Category[]>([]);
  const user_id = ref<number | undefined>();

  const fetchMenuItems = async () => {
    try {
      const { $api } = useNuxtApp();
      const res = await $api.get<MenuResponse>("/category");
      const data = (res as any).data as { user_id?: number; categories?: Category[] } | undefined;
      menu.value = data?.categories || (data as unknown as Category[]) || [];

      // Lưu user_id từ response /category
      if (typeof data?.user_id === "number") {
        user_id.value = data.user_id;
      }
    } catch (error) {
      console.error("Error fetching menu items:", error);
    }
  };

  return {
    menu,
    user_id,
    fetchMenuItems,
  };
});
