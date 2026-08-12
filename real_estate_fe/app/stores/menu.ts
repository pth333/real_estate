import { defineStore } from "pinia";
import { Menu } from "~/types/window";
import type { MenuSettings } from "~/types/window";

export const useMenuStore = defineStore("menu", () => {
  const menu = ref<MenuSettings>();
  const user_id = ref<number | undefined>();

  const fetchMenuItems = async () => {
    try {
      const { $api } = useNuxtApp();
      const res = await $api.get<MenuSettings>("/category");
      const data = (res as any).data;
      if (data) {
        menu.value = data;
        // Khởi tạo Menu trong store rồi gán vào window — nơi khác đọc qua global
        window.menu = new Menu(data);
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
