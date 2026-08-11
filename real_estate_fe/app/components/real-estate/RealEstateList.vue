<template>
  <div class="mx-auto max-w-[1200px] px-6 py-6">
    <div class="mb-4 flex flex-col gap-3">
      <SearchBar @search="handleSearch" />
      <FilterManager />
    </div>

    <!-- Loading -->
    <SkeletonCard v-if="loading" :count="pageSize" />

    <!-- Empty -->
    <div v-else-if="realEstates.length === 0" class="px-6 py-16 text-center">
      <p class="text-base text-gray-400">Không tìm thấy bất động sản nào</p>
    </div>

    <!-- Data grid -->
    <div v-else class="mb-8 flex flex-col gap-4">
      <RealEstateCard v-for="estate in realEstates" :key="estate.id" :estate="estate" @call="handleCall"
        @toggle-favorite="handleToggleFavorite" />

    </div>

    <!-- Pagination -->
    <Pagination :current-page="realEstateStore.currentPage" :total-pages="totalPages" @page-change="goToPage" />
  </div>
</template>

<style scoped></style>

<script setup lang="ts">
import type { RealEstateResponse, PaginatedResponse } from "~/types/real_estate";
import { useFilterStore } from "~/stores/filter";
import { useRealEstateStore } from "~/stores/real_estate";

const route = useRoute();
const { $api } = useNuxtApp();
const filterStore = useFilterStore();
const realEstateStore = useRealEstateStore()
const realEstates = ref<RealEstateResponse[]>([]);
const loading = ref(false);
const pageSize = ref(12);
const totalRecords = ref(0);

const categorySlug = computed<string>(() => {
  const v = route.params.category;
  return Array.isArray(v) ? v[0] ?? "" : v ?? "";
});

const filterSegments = computed<string[]>(() => {
  const v = route.params.filters;
  if (Array.isArray(v)) return v;
  return v ? [v] : [];
});

const apiPath = computed<string>(() => {
  const parts: string[] = [categorySlug.value];
  if (filterSegments.value.length > 0) parts.push(filterSegments.value.join("/"));
  return parts.join("/");
});


const totalPages = computed(() =>
  Math.ceil(totalRecords.value / pageSize.value),
);

const fetchDataRealEstate = async () => {
  loading.value = true;
  try {
    const res = await $api.get<PaginatedResponse<RealEstateResponse>>(
      `/real-estate/${apiPath.value}`,
      {
        params: {
          page: realEstateStore.currentPage,
          size: pageSize.value,
          ...buildListParams(filterStore.filters)
        },
      },
    );
    realEstates.value = res.data || [];
    totalRecords.value = res.total || 0;
  } catch (err) {
    const msg =
      err instanceof Error ? err.message : "Có lỗi xảy ra khi tải dữ liệu";
    window.message?.error(msg);
  } finally {
    loading.value = false;
  }
}

function goToPage(page: number) {
  if (page < 1 || page > totalPages.value) return;
  realEstateStore.currentPage = page;
  const parts: string[] = [categorySlug.value];
  if (filterSegments.value.length > 0) parts.push(filterSegments.value.join("/"));
  // page > 1 → query string (cấu trúc [...filters] không có segment page)
  let url = `/${parts.join("/")}`;
  if (page > 1) url += `?page=${page}`;
  navigateTo(url);
}

watch(
  () => [route.params.category, route.params.filters, route.query],
  ([_cat, _filt, newPage]) => {
    realEstates.value = [];
    totalRecords.value = 0;

    const pageNumber = newPage ? Number(newPage) : 1;
    realEstateStore.currentPage = Number.isNaN(pageNumber) ? 1 : pageNumber;

    fetchDataRealEstate();
  },
  { immediate: true },
);

function handleCall(phone: string) {
  window.open(`tel:${phone}`, "_self");
}

function handleToggleFavorite(id: number) {
  const estate = realEstates.value.find((e: RealEstateResponse) => e.id === id);
  if (estate) {
    estate.is_favorite = !estate.is_favorite;
  }
}

const handleSearch = async () => {
  // Server-driven: đưa keyword vào query để server chạy FULLTEXT
  const q = filterStore.searchKeyword || "";
  const parts: string[] = [categorySlug.value];
  if (filterSegments.value.length > 0) parts.push(filterSegments.value.join("/"));
  let url = `/${parts.join("/")}`;
  const params: string[] = [];
  if (realEstateStore.currentPage > 1) params.push(`page=${realEstateStore.currentPage}`);
  if (q) params.push(`search=${encodeURIComponent(q)}`);
  if (params.length) url += `?${params.join("&")}`;
  navigateTo(url);
}
</script>