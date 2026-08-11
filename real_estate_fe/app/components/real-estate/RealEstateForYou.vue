<template>
  <section class="py-8">
    <div class="container mx-auto px-24">
    <!-- Header -->
    <div class="flex items-center justify-between mb-5">
      <h2 class="text-xl font-bold text-gray-900">Bất động sản dành cho bạn</h2>
      <div class="flex items-center gap-3 text-sm">
        <a href="#" class="text-gray-600 hover:text-red-500 transition-colors">Tin nhà đất bán mới nhất</a>
        <span class="text-gray-300">|</span>
        <a href="#" class="text-gray-600 hover:text-red-500 transition-colors">Tin nhà đất cho thuê mới nhất</a>
      </div>
    </div>

    <!-- Grid -->
    <div class="grid grid-cols-4 gap-4">
      <div
        v-for="item in visibleItems"
        :key="item.id"
        class="bg-white overflow-hidden border border-gray-100 hover:shadow-md transition-shadow cursor-pointer"
      >
        <!-- Ảnh -->
        <div class="relative aspect-[4/3] overflow-hidden bg-gray-100">
          <img
            :src="item.thumbnail"
            :alt="item.title"
            class="w-full h-full object-cover hover:scale-105 transition-transform duration-300"
          />
        </div>

        <!-- Nội dung -->
        <div class="p-3 flex flex-col gap-1.5">
          <!-- Badge xác thực + tiêu đề -->
          <div class="text-sm font-medium text-gray-800 leading-snug line-clamp-2">
            <span v-if="item.verified" class="inline-flex items-center gap-1 text-green-600 font-semibold mr-1">
              <IconShieldCheck class="h-3.5 w-3.5" />
              XÁC THỰC
            </span>
            {{ item.title }}
          </div>

          <!-- Giá + diện tích -->
          <div class="flex items-center gap-2 text-sm">
            <span class="text-red-500 font-semibold">{{ item.price }}</span>
            <span class="text-gray-300">·</span>
            <span class="text-gray-600">{{ item.area }}</span>
          </div>

          <!-- Địa chỉ -->
          <div class="flex items-center gap-1 text-xs text-gray-500">
            <IconMapPin class="h-3 w-3 shrink-0" />
            <span class="truncate">{{ item.location }}</span>
          </div>

          <!-- Footer: ngày đăng + yêu thích -->
          <div class="flex items-center justify-between mt-1">
            <span class="text-xs text-gray-400">{{ item.postedAt }}</span>
            <button
              class="p-1.5 border border-gray-200 hover:border-red-300 hover:text-red-500 transition-colors"
              @click.stop="toggleFavorite(item.id)"
            >
              <IconHeart
                class="h-4 w-4"
                :class="item.isFavorite ? 'fill-red-500 text-red-500' : 'text-gray-400'"
              />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Mở rộng -->
    <div class="flex justify-center mt-6">
      <button
        class="flex items-center gap-2 px-6 py-2.5 border border-gray-300 text-sm text-gray-700 hover:border-gray-400 transition-colors"
        @click="expanded = !expanded"
      >
        {{ expanded ? 'Thu gọn' : 'Mở rộng' }}
        <IconChevronDown class="h-4 w-4 transition-transform" :class="{ 'rotate-180': expanded }" />
      </button>
    </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const expanded = ref(false)

// Mock data — thay bằng API sau
const items = ref([
  {
    id: 1,
    title: 'Sun Festo khoảng nóng Ocean - 2PN 68m2 giá full 3,9tỷ chiết kh...',
    price: '3,9 tỷ',
    area: '68 m²',
    location: 'Hạ Long, Quảng Ninh',
    postedAt: 'Đăng hôm nay',
    verified: false,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
  {
    id: 2,
    title: 'Cho thuê căn hộ 2 phòng ngủ Estella Heights - Giá...',
    price: '35 triệu/tháng',
    area: '103 m²',
    location: 'Quận 2, Hồ Chí Minh',
    postedAt: 'Đăng hôm nay',
    verified: true,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
  {
    id: 3,
    title: 'Nhà hẻm Chu Văn An | 40m2 (4.6 x 6.5m) | 2 tầng, 2PN/2WC | 5.5 t...',
    price: '5,5 tỷ',
    area: '40 m²',
    location: 'Bình Thạnh, Hồ Chí Minh',
    postedAt: 'Đăng hôm nay',
    verified: false,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
  {
    id: 4,
    title: 'Cho thuê căn hộ 3 phòng ngủ Estella Heights - thá...',
    price: '78 triệu/tháng',
    area: '142 m²',
    location: 'Quận 2, Hồ Chí Minh',
    postedAt: 'Đăng hôm nay',
    verified: true,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
  {
    id: 5,
    title: 'Cho thuê căn hộ 1 phòng ngủ Estella Heights - Giá...',
    price: '25 triệu/tháng',
    area: '60 m²',
    location: 'Quận 2, Hồ Chí Minh',
    postedAt: 'Đăng hôm nay',
    verified: true,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
  {
    id: 6,
    title: 'Căn hộ cho thuê 52m2 Quận 9, gần đại học FPT, ĐH HUTECH,...',
    price: '4,4 triệu/tháng',
    area: '52 m²',
    location: 'Quận 9, Hồ Chí Minh',
    postedAt: 'Đăng hôm nay',
    verified: false,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
  {
    id: 7,
    title: 'Bán Nhà KDC Bình Hưng, Bình Chánh, ngay bến xe Quận 8 vào...',
    price: '9,8 tỷ',
    area: '96 m²',
    location: 'Bình Chánh, Hồ Chí Minh',
    postedAt: 'Đăng hôm nay',
    verified: false,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
  {
    id: 8,
    title: 'Căn hộ 1 phòng ngủ mini xịn sò - đủ nội thất - thang...',
    price: '4,8 triệu/tháng',
    area: '30 m²',
    location: 'Bình Tân, Hồ Chí Minh',
    postedAt: 'Đăng hôm nay',
    verified: true,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
  // Row 3 — ẩn khi chưa mở rộng
  {
    id: 9,
    title: 'Bán căn hộ cao cấp view sông thoáng mát, full nội thất...',
    price: '6,2 tỷ',
    area: '88 m²',
    location: 'Quận 7, Hồ Chí Minh',
    postedAt: 'Đăng hôm nay',
    verified: true,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
  {
    id: 10,
    title: 'Cho thuê nhà nguyên căn hẻm xe hơi, 4PN, sân thượng...',
    price: '18 triệu/tháng',
    area: '72 m²',
    location: 'Gò Vấp, Hồ Chí Minh',
    postedAt: 'Đăng hôm nay',
    verified: false,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
  {
    id: 11,
    title: 'Đất nền dự án KDC Nam Sài Gòn, sổ hồng riêng...',
    price: '2,1 tỷ',
    area: '120 m²',
    location: 'Nhà Bè, Hồ Chí Minh',
    postedAt: 'Đăng hôm nay',
    verified: false,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
  {
    id: 12,
    title: 'Căn hộ Studio full nội thất cao cấp, ban công thoáng...',
    price: '3,2 triệu/tháng',
    area: '28 m²',
    location: 'Thủ Đức, Hồ Chí Minh',
    postedAt: 'Đăng hôm nay',
    verified: true,
    isFavorite: false,
    thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=BDS',
  },
])

// Hiển thị 8 khi thu gọn, tất cả khi mở rộng
const visibleItems = computed(() => expanded.value ? items.value : items.value.slice(0, 8))

function toggleFavorite(id: number) {
  const item = items.value.find(i => i.id === id)
  if (item) item.isFavorite = !item.isFavorite
}
</script>