<template>
  <section class="hero-bg py-14 px-4 relative overflow-hidden">
    <!-- Background image -->
    <div
      class="absolute inset-0 bg-cover bg-center bg-no-repeat"
      :style="{ backgroundImage: `url('https://images.unsplash.com/photo-1545324418-cc1a3fa10c00?w=1200')` }"
    >
      <!-- Overlay -->
      <div class="absolute inset-0 bg-white/70 backdrop-blur-[2px]" />
    </div>

    <!-- Content -->
    <div class="relative max-w-3xl mx-auto text-center">
      <h1 class="text-2xl sm:text-3xl font-bold text-primary leading-snug mb-2">
        Giải pháp giao dịch bất động sản từ trực tuyến<br class="hidden sm:block" /> đến trực tiếp của Vinhomes
      </h1>
      <p class="text-sm text-primary/70 mb-6">Tìm kiếm Bất động sản theo nhu cầu của Quý khách</p>

      <!-- Search box -->
      <div class="bg-white rounded-xl shadow-lg px-4 py-4 max-w-2xl mx-auto">
        <n-input-group>
          <n-input
            v-model:value="searchQuery"
            placeholder="Tìm kiếm theo khu vực hoặc dự án"
            size="large"
          />
          <n-button type="primary" size="large">
            <template #icon>
              <n-icon><IconSearch /></n-icon>
            </template>
          </n-button>
        </n-input-group>

        <!-- trending tags -->
        <div class="flex items-center gap-2 mt-3 flex-wrap">
          <n-text depth="3" class="text-xs flex items-center gap-2">
            <n-icon size="16"><TrendingUpOutline /></n-icon>
            Được quan tâm
          </n-text>
          <n-tag
            v-for="tag in trendingTags"
            :key="tag"
            size="small"
            round
            :bordered="false"
            class="cursor-pointer"
            @click="searchQuery = tag"
          >
            {{ tag }}
          </n-tag>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const searchQuery = ref('')
const trendingTags = ref([])
</script>
