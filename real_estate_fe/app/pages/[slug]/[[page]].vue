<template>
  <div class="mx-auto max-w-[1200px] px-6 py-6">
    <div class="mb-4 flex flex-col gap-3">
      <SearchBar @search="fetchDataRealEstate" />
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
const pageSize = ref(10);
const totalRecords = ref(0);

// Slug = segment đầu, page = segment thứ 2 (nếu có)
const slugSegments = computed(() => {
  const s = route.params.slug;
  return Array.isArray(s) ? s : [s];
});

const slug = computed(() => slugSegments.value[0] || "");

const totalPages = computed(() =>
  Math.ceil(totalRecords.value / pageSize.value),
);

const payload = computed(() => ({
  page: realEstateStore.currentPage,
  size: pageSize.value,
  filter: filterStore.filters,
  search: filterStore.searchKeyword,
}))

const fetchDataRealEstate = async () => {
  loading.value = true;

  try {
    const res = await $api.post<PaginatedResponse<RealEstateResponse>>(
      `/real-estate/${slug.value}/${realEstateStore.currentPage}`,
      payload.value,
    );

    realEstates.value = res.data || [];
    totalRecords.value = res.total || 0;
  } catch (err) {
    const msg =
      err instanceof Error ? err.message : "Có lỗi xảy ra khi tải dữ liệu";
    window.message?.error(msg);
    console.error("Error fetching real estates:", err);
  } finally {
    loading.value = false;
  }
}

function goToPage(page: number) {
  realEstateStore.currentPage = page
  if (realEstateStore.currentPage < 1 || realEstateStore.currentPage > totalPages.value) return;
  if (realEstateStore.currentPage === 1) {
    navigateTo(`/${slug.value}`);
  } else {
    navigateTo(`/${slug.value}/${realEstateStore.currentPage}`);
  }
}

watch(
  () => [route.params.slug, route.params.page],
  async ([newSlug, newPage]) => {
    // Reset state
    realEstates.value = [];
    totalRecords.value = 0;

    // Đồng bộ page từ route
    realEstateStore.currentPage = newPage ? Number(newPage) : 1;

    await fetchDataRealEstate();
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

</script>
