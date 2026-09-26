<template>
  <!-- Sticky bar nằm ngoài. top-16 = chiều cao AppHeader (h-16) đang sticky ở top-0,
       nếu để top-0 thì thanh search luôn bị header che khi cuộn. -->
  <div class="sticky top-16 z-10 border-b border-gray-100 bg-white lg:border-0">
    <div class="mx-auto max-w-[1200px] px-4 py-2 md:px-6 md:py-3">
      <div class="flex flex-col gap-2 md:gap-3">
        <SearchBar @search="handleSearch" />
        <FilterManager />
      </div>
    </div>
  </div>

  <!-- Content bình thường -->
  <div class="mx-auto max-w-[1200px] px-4 py-4 md:px-6 md:py-6">
    <!-- Loading -->
    <SkeletonCard v-if="loading" :count="pageSize" />

    <!-- Empty -->
    <div v-else-if="realEstates.length === 0" class="px-4 py-12 text-center md:py-16">
      <p class="text-base text-gray-400">Không tìm thấy bất động sản nào</p>
    </div>

    <!-- Data grid -->
    <div v-else class="mb-6 flex flex-col gap-3 md:mb-8 md:gap-4">
      <RealEstateCard :realEstates="realEstates" />
    </div>

    <!-- Pagination -->
    <Pagination :current-page="realEstateStore.currentPage" :total-pages="totalPages" @page-change="goToPage" />
  </div>
</template>

<style scoped></style>

<script setup lang="ts">
import type { RealEstateResponse } from "~/types/real_estate";
import { useFilterStore } from "~/stores/filter";
import { useRealEstateStore } from "~/stores/real_estate";
import { useRealEstateService } from "~/services/real-estate.service";

const route = useRoute();
const realEstateService = useRealEstateService();
const filterStore = useFilterStore();
const realEstateStore = useRealEstateStore()
const realEstates = ref<RealEstateResponse[]>([]);
const loading = ref(false);
const pageSize = ref(12);
const totalRecords = ref(0);

const props = defineProps<{
  categorySlug: string
}>()

const query = computed<string>(() => {
  const v = route.query.search as string
  return v ? v : ''
})

const filterSegments = computed<string[]>(() => {
  const v = route.params.filters;
  if (Array.isArray(v)) return v;
  return v ? [v] : [];
});

const apiPath = computed<string>(() => {
  const parts: string[] = [props.categorySlug];
  if (filterSegments.value.length > 0) parts.push(filterSegments.value.join("/"));
  return parts.join("/");
});

const totalPages = computed(() =>
  Math.ceil(totalRecords.value / pageSize.value),
);

const fetchDataRealEstate = async () => {
  loading.value = true;
  try {
    const res = await realEstateService.getListByCategoryPath(apiPath.value, {
      page: realEstateStore.currentPage,
      size: pageSize.value,
      search: query.value,
      ...buildListParams(filterStore.filters),
    });
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
  const parts: string[] = [props.categorySlug];
  if (filterSegments.value.length > 0) parts.push(filterSegments.value.join("/"));
  // page > 1 → query string (cấu trúc [...filters] không có segment page)
  let url = `/${parts.join("/")}`;
  if (page > 1) url += `?page=${page}`;
  navigateTo(url);
}

watch(
  () => [route.params.category, route.params.filters, route.query],
  ([_cat, _filt, newQuery]) => {
    realEstates.value = [];
    totalRecords.value = 0;

    const query = newQuery as Record<string, string>;
    const pageNumber = query?.page ? Number(query.page) : 1;
    realEstateStore.currentPage = Number.isNaN(pageNumber) ? 1 : pageNumber;

    fetchDataRealEstate();
  },
  { immediate: true },
);



const handleSearch = async () => {
  // Server-driven: đưa keyword vào query để server chạy FULLTEXT
  const q = filterStore.searchKeyword || "";

  const parts: string[] = [props.categorySlug];
  if (filterSegments.value.length > 0) parts.push(filterSegments.value.join("/"));
  let url = `/${parts.join("/")}`;
  const params: string[] = [];
  if (realEstateStore.currentPage > 1) params.push(`page=${realEstateStore.currentPage}`);
  if (q) params.push(`search=${encodeURIComponent(q)}`);
  if (params.length) url += `?${params.join("&")}`;
  navigateTo(url);
}
</script>