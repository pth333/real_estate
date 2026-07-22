import { defineStore } from "pinia";
import type { DashboardSummary, Filter, PaginatedResponse, RealEstateResponse } from "@/types/real_estate";

export const useRealEstateStore = defineStore("realEstate", () => {
  const items = ref<RealEstateResponse[]>([]);
  const total = ref(0);
  const loading = ref(false);
  const summary = ref<DashboardSummary | null>(null);
  const summaryLoading = ref(false);
  const page = ref(1);
  const pageSize = ref(50);
  const filter = ref<Filter>({});

  function formattedPrice(priceVND: number): string {
    return (priceVND / 1_000_000_000).toFixed(2) + " tỷ";
  }

  const payload = computed(() => ({
    page: page.value,
    size: pageSize.value,
    filter: filter.value,
  }));

  async function fetchList() {
    loading.value = true;
    try {
      const { $api } = useNuxtApp();
      const res = await $api.post<PaginatedResponse<RealEstateResponse>>(
        "/real-estate/list",
        payload.value,
      );
      items.value = res.data;
      total.value = res.total;
    } catch (e) {
      console.error("Lỗi tải danh sách BĐS:", e);
    } finally {
      loading.value = false;
    }
  }

  async function fetchSummary(from?: string, to?: string) {
    summaryLoading.value = true;
    try {
      const { $api } = useNuxtApp();
      const res = await $api.get<{ data: DashboardSummary }>(
        "/dashboard/summary",
        { params: { from, to } },
      );
      summary.value = (res as any).data;
    } catch (e) {
      console.error("Lỗi tải summary:", e);
    } finally {
      summaryLoading.value = false;
    }
  }

  function setFilter(newFilter: Filter) {
    filter.value = newFilter;
    page.value = 1;
    fetchList();
  }

  function setPage(newPage: number) {
    page.value = newPage;
    fetchList();
  }

  function setPageSize(newSize: number) {
    pageSize.value = newSize;
    page.value = 1;
    fetchList();
  }

  return {
    items, total, loading, summary, summaryLoading, page, pageSize, filter,
    formattedPrice, fetchList, fetchSummary, setFilter, setPage, setPageSize,
  };
});
