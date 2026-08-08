import type { SelectOption } from "naive-ui";
import { defineStore } from "pinia";
import type {
  Filter,
  FilterAreaRange,
  FilterLocation,
  FilterPriceRange,
} from "~/types/real_estate";
export const useFilterStore = defineStore("filter", () => {
  const showModalAll = ref(false);
  const showPopover = ref(false);
  const showAreaPopover = ref(false);

  const filterLocation = ref<FilterLocation>({});

  const searchKeyword = ref<string | null>(null);

  const filters = ref<Filter>({});
  const filterPriceRange = ref<FilterPriceRange>({});
  const filterAreaRange = ref<FilterAreaRange>({});
  const showLocationModal = ref(false);

  const cityOptions = ref<SelectOption[]>([]);
  const districtOptions = ref<SelectOption[]>([]);

  const wardOptions = ref<SelectOption[]>([]);

  const screen = ref<"main" | "location">("main");

  // Lấy tên tỉnh/thành từ code (VD 79 → "Hà Nội") để ghép vào slug
  // cityOptions là SelectOption[]: value = code, label = name
  const slugCity = computed(() => {
    const found = cityOptions.value.find(
      (item) => item.value === filterLocation.value.city,
    );
    return found?.label ?? "";
  });

  return {
    screen,
    filters,
    searchKeyword,
    showModalAll,
    showPopover,
    showAreaPopover,
    filterPriceRange,
    filterAreaRange,
    showLocationModal,
    filterLocation,
    slugCity,
    cityOptions,
    districtOptions,
    wardOptions,
  };
});
