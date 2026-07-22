
import { defineStore } from "pinia";
import type { FilterRealEstate } from "~/types/real_estate";
export const useFilterStore = defineStore("filter", () => {
    const filters = ref<FilterRealEstate>({
        min_price: undefined,
        max_price: undefined,
        district: undefined,
    })

    return {
        filters,
    }
})