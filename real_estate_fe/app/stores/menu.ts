import { defineStore } from "pinia";

export interface Category {
  ID: number;
  Name: string;
  Slug: string;
  children?: Category[];
}

export interface FlatCategory {
  id: number;
  name: string;
  depth: number;
}

export const useMenuStore = defineStore("menu", () => {
  const categories = ref<Category[]>([]);

  const fetchMenuItems = async () => {
    try {
      const res = await $fetch("/api/category");
      categories.value = res.data?.data || [];
    } catch (error) {
      console.error("Error fetching menu items:", error);
    }
  };

  return {
    categories,
    fetchMenuItems,
  };
});
