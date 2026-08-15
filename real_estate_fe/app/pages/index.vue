<template>
  <div class="flex flex-col gap-16">
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

const menuStore = useMenuStore()
const searchQuery = ref('')

/**
 * Xử lý tìm kiếm bất động sản từ trang chủ
 * Điều hướng người dùng sang trang danh mục mặc định kèm từ khóa tìm kiếm
 */
const handleSearchRealEstate = async () => {
  const keyword = searchQuery.value.trim()
  if (!keyword) return

  // Lấy slug danh mục đầu tiên từ menu để điều hướng
  const firstCategorySlug = menuStore.menu?.categories?.[0]?.Slug || "nha-dat-ban"

  // Chuyển hướng người dùng sang trang danh sách BĐS kèm từ khóa tìm kiếm
  navigateTo(`/${firstCategorySlug}?search=${encodeURIComponent(keyword)}`)
}
</script>