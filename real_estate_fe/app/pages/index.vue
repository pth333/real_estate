<template>
  <div class="flex flex-col gap-20">
    <section class="hero-bg py-40 px-2 relative overflow-hidden">
      <!-- Background image -->
      <div class="absolute inset-0 bg-cover bg-center bg-no-repeat"
        :style="{ backgroundImage: `url('https://pub-5eb4e976c2fe4062ba3cdabce48568cc.r2.dev/uploads/931c3b1968a2e54e32e00dace86c38de.jpg')` }">
        <!-- Overlay -->
        <div class="absolute inset-0 bg-white/70 backdrop-blur-[2px]" />
      </div>

      <!-- Content -->
      <div class="relative max-w-3xl mx-auto text-center">
        <h1 class="text-2xl sm:text-3xl font-bold text-primary leading-snug mb-2">
          Giải pháp giao dịch bất động sản từ trực tuyến<br class="hidden sm:block" /> đến trực tiếp của Phan Hieu Group
        </h1>
        <p class="text-sm text-primary/70 mb-6">Tìm kiếm Bất động sản theo nhu cầu của Quý khách</p>

        <!-- Search box -->
        <div class="bg-white  shadow-lg px-4 py-4 max-w-2xl mx-auto">
          <n-input-group>
            <n-input v-model:value="searchQuery" placeholder="Tìm kiếm theo khu vực hoặc dự án" size="large"
              style="text-align: left;" />
            <n-button type="primary" size="large" @click="handleSearch">
              <template #icon>
                <n-icon>
                  <IconSearch />
                </n-icon>
              </template>
            </n-button>
          </n-input-group>
        </div>
      </div>
    </section>

    <!-- Bất động sản theo khu vực -->
    <RealEstateByArea />
  </div>
</template>

<script setup lang="ts">
const searchQuery = ref('')

const trackSearch = async (query: string) => {
  const { $api } = useNuxtApp()
  try {
    await $api.post('/tracking/search', {
      query,
      user_id: window.menu?.user_id || '',
    })
  } catch (err) {
    console.error('Tracking search error:', err)
  }
}

const handleSearch = async () => {
  if (!searchQuery.value.trim()) return
  await trackSearch(searchQuery.value)
  navigateTo(`?search=${encodeURIComponent(searchQuery.value)}`)
}
</script>