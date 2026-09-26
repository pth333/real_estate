<template>
  <div class="flex flex-col gap-8 md:gap-12 lg:gap-16">
    <section class="hero-bg relative overflow-hidden px-4 py-16 md:px-6 md:py-24 lg:py-40">
      <!-- Background image -->
      <div class="absolute inset-0 bg-cover bg-center bg-no-repeat" :style="backgroundStyle">
        <!-- Overlay -->
        <div class="absolute inset-0 bg-white/70 backdrop-blur-[2px]" />
      </div>

      <!-- Content -->
      <div class="relative mx-auto max-w-3xl text-center">
        <h1 class="mb-2 text-xl font-bold leading-snug text-primary md:text-2xl lg:text-3xl">
          Giải pháp giao dịch bất động sản từ trực tuyến<br class="hidden md:block" /> đến trực tiếp của Phan Hieu Group
        </h1>
        <p class="mb-5 text-sm text-primary/70 md:mb-6">Tìm kiếm Bất động sản theo nhu cầu của Quý khách</p>

        <!-- Search box -->
        <div class="mx-auto max-w-2xl bg-white px-3 py-3 shadow-lg md:px-4 md:py-4">
          <n-input-group>
            <n-input v-model:value="searchQuery" placeholder="Tìm kiếm theo khu vực hoặc dự án" size="large"
              style="text-align: left;" />
            <n-button type="primary" size="large" @click="handleSearchRealEstate">
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

    <RealEstateByArea />
    <RealEstateForYou />
    <FeaturedProjects />
  </div>
</template>

<script setup lang="ts">
import { useMenuStore } from "~/stores/menu"
import { useHeroBg } from "~/composables/useHeroBg"

useHead({
  title: "Trang chủ",
})

const menuStore = useMenuStore()
const { backgroundStyle } = useHeroBg()
const searchQuery = ref('')


const handleSearchRealEstate = async () => {
  const keyword = searchQuery.value.trim()
  if (!keyword) return

  // Lấy slug danh mục đầu tiên từ menu để điều hướng
  const firstCategorySlug = menuStore.menu?.categories?.[0]?.Slug || "nha-dat-ban"

  // Chuyển hướng người dùng sang trang danh sách BĐS kèm từ khóa tìm kiếm
  navigateTo(`/${firstCategorySlug}?search=${encodeURIComponent(keyword)}`)
}
</script>