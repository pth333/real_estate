import type { SelectOption } from "naive-ui";
import { defineStore } from "pinia";
import type {
  Filter,
  FilterLocation,
  FilterPriceRange,
} from "~/types/real_estate";
export const useFilterStore = defineStore("filter", () => {
  const showModalAll = ref(false);
  const showPopover = ref(false);

  const filterLocation = ref<FilterLocation>({});

  const searchKeyword = ref<string | null>(null);

  const filters = ref<Filter>({});
  const filterPriceRange = ref<FilterPriceRange>({});
  const showLocationModal = ref(false);

  // Mock cityOptions
  const cityOptions = ref<SelectOption[]>([]);

  // Mock districtOptions
  const districtOptions = ref<SelectOption[]>([]);

  // Mock wardOptions
  const wardOptions = ref<SelectOption[]>([]);

  const screen = ref<"main" | "location">("main");

  return {
    screen,
    filters,
    searchKeyword,
    showModalAll,
    showPopover,
    filterPriceRange,
    showLocationModal,
    filterLocation,

    cityOptions,
    districtOptions,
    wardOptions,
  };
});
