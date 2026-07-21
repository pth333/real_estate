<template>
  <div class="flex flex-col gap-6">
    <!-- Hero Section -->
    <div class="bg-foundation px-8 py-12">
      <div class="mx-auto max-w-7xl">
        <h1 class="mb-3 font-bold tracking-tight text-limestone text-4xl">
          Dashboard
        </h1>
        <p class="mb-8 text-patina/80">Tổng quan dữ liệu bất động sản</p>

        <!-- Search bar -->
        <div class="relative max-w-2xl">
          <input
            type="text"
            placeholder="Tìm kiếm bất động sản theo địa chỉ, quận, loại..."
            class="w-full border border-limestone/20 bg-foundation/50 px-6 py-4 text-sm text-limestone placeholder:text-limestone/40 focus:border-patina focus:outline-none"
          />
          <div class="absolute right-4 top-1/2 -translate-y-1/2">
            <svg
              class="h-5 w-5 text-patina"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
          </div>
          <div
            class="absolute -right-3 top-1/2 h-px w-6 -translate-y-1/2 bg-limestone/20"
          ></div>
        </div>

        <!-- Quick stats -->
        <div class="mt-8 grid grid-cols-3 gap-6">
          <div class="border-l-2 border-patina pl-4">
            <div
              class="font-mono text-xs uppercase tracking-wider text-limestone/50"
            >
              Tổng BĐS
            </div>
            <div class="mt-1 font-bold tabular-nums text-limestone text-3xl">
              {{ totalProperties }}
            </div>
          </div>
          <div class="border-l-2 border-oak pl-4">
            <div
              class="font-mono text-xs uppercase tracking-wider text-limestone/50"
            >
              Giá trung bình
            </div>
            <div class="mt-1 font-bold tabular-nums text-limestone text-3xl">
              {{ avgPrice }}
            </div>
          </div>
          <div class="border-l-2 border-limestone/30 pl-4">
            <div
              class="font-mono text-xs uppercase tracking-wider text-limestone/50"
            >
              Quận
            </div>
            <div class="mt-1 font-bold tabular-nums text-limestone text-3xl">
              {{ totalDistricts }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Content zones -->
    <div class="px-8">
      <div class="mx-auto max-w-7xl">
        <SummaryCards />
        <!-- Charts -->
        <div class="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-2">
          <ClientOnly><PriceChart /></ClientOnly>
          <ClientOnly><DistrictChart /></ClientOnly>
        </div>
        <!-- Bảng dữ liệu -->
        <div class="mt-6">
          <ClientOnly><RealEstateTable /></ClientOnly>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRealEstateStore } from "~/app/stores/real_estate";

const store = useRealEstateStore();

const totalProperties = computed(() => store.summary?.total_posts ?? 0);
const avgPrice = computed(() => {
  const avg = store.summary?.avg_price_m2 ?? 0;
  return avg > 0 ? `${(avg / 1000).toFixed(1)}K` : "—";
});
const totalDistricts = computed(() => {
  const districts = new Set(store.items.map((item: any) => item.District));
  return districts.size;
});

onMounted(() => {
  store.fetchList();
  store.fetchSummary();
});
</script>
