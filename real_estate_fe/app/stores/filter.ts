import type { SelectOption } from "naive-ui";
import { defineStore } from "pinia";
import type { Filter } from "~/types/real_estate";

// Store filter dùng chung cho list BĐS — giữ toàn bộ tiêu chí trong 1 object
// phẳng `filters` (khớp dto.Filter backend). Các state tách nhỏ
// (filterLocation/filterPriceRange/...) đã gộp về đây để URL là nguồn truth.
export const useFilterStore = defineStore("filter", () => {
  const showModalAll = ref(false);
  const showPopover = ref(false);
  const showAreaPopover = ref(false);

  const searchKeyword = ref<string | null>(null);

  // Filter phẳng duy nhất — mọi tiêu chí lọc nằm trong đây
  const filters = ref<Filter>({});

  // Options location (dùng chung: FilterLocation, build URL...)
  const cityOptions = ref<SelectOption[]>([]);
  const districtOptions = ref<SelectOption[]>([]);
  const wardOptions = ref<SelectOption[]>([]);

  const screen = ref<"main" | "location">("main");

  // Lấy tên tỉnh/thành từ code (VD 79 → "Hà Nội") để ghép vào URL
  // cityOptions là SelectOption[]: value = code, label = name
  const slugCity = computed(() => {
    const found = cityOptions.value.find(
      (item) => item.value === filters.value.city,
    );
    return found?.label ?? "";
  });

  // Đặt lại toàn bộ filter (dùng cho nút "Đặt lại" trong modal)
  function resetFilters() {
    filters.value = {};
  }

  return {
    screen,
    filters,
    searchKeyword,
    showModalAll,
    showPopover,
    showAreaPopover,
    slugCity,
    cityOptions,
    districtOptions,
    wardOptions,
    resetFilters,
  };
});
