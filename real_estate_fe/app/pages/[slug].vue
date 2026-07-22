<template>
  <div class="mx-auto max-w-[1200px] px-6 py-6">
    <FilterManager @apply-filters="fetchData" />

    <!-- Loading: skeleton cards -->
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
    <Pagination :current-page="currentPage" :total-pages="totalPages" @page-change="goToPage" />
  </div>
</template>

<style scoped></style>

<script setup lang="ts">
import type { RealEstateResponse, PaginatedResponse } from "~/types/real_estate";
import { useToast } from "~/composables/useToast";
import { useFilterStore } from "~/stores/filter";

const route = useRoute();
const { $api } = useNuxtApp();
const { showToast } = useToast();
const filterStore = useFilterStore();

const realEstates = ref<RealEstateResponse[]>([]);
const loading = ref(false);
const currentPage = ref(Number(route.query.page) || 1);
const pageSize = ref(10);
const totalRecords = ref(0);

const totalPages = computed(() =>
  Math.ceil(totalRecords.value / pageSize.value),
);

async function fetchData() {
  loading.value = true;

  try {
    const payload = {
      page: currentPage.value,
      size: pageSize.value,
      filter: filterStore.filters,
    };
    const res = await $api.post<PaginatedResponse<RealEstateResponse>>(
      `/real-estate/${route.params.slug as string}/${currentPage.value}`,
      payload,
    );

    realEstates.value = res.data || [];
    totalRecords.value = res.total || 0;
  } catch (err) {
    const msg =
      err instanceof Error ? err.message : "Có lỗi xảy ra khi tải dữ liệu";
    showToast("error", msg);
    console.error("Error fetching real estates:", err);
  } finally {
    loading.value = false;
  }
}

function goToPage(page: number) {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  fetchData();
  window.scrollTo({ top: 0, behavior: "smooth" });
}

function handleCall(phone: string) {
  window.open(`tel:${phone}`, "_self");
}

function handleToggleFavorite(id: number) {
  const estate = realEstates.value.find((e: RealEstateResponse) => e.id === id);
  if (estate) {
    estate.is_favorite = !estate.is_favorite;
  }
}

onMounted(() => {
  fetchData();
});
</script>
