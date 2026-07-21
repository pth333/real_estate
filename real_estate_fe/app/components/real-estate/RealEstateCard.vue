<template>

  <div class="real-estate-card">
     <!-- Badge VIP -->
    <div v-if="estate.Badge" class="badge">{{ estate.Badge }}</div>
     <!-- Grid ảnh: 1 ảnh lớn + 3 ảnh nhỏ -->
    <div class="images-grid">

      <div class="main-image">
         <img :src="mainImage" :alt="estate.Title" @error="handleImageError" />
      </div>

      <div class="thumbnail-grid">

        <div v-for="(img, idx) in thumbnails" :key="idx" class="thumbnail">
           <img :src="img" :alt="`${estate.Title} ${idx + 2}`" @error="handleImageError" />
        </div>
         <!-- Hiển thị số ảnh còn lại nếu có -->
        <div v-if="remainingImagesCount > 0" class="thumbnail remaining-count">
           <span>{{ remainingImagesCount }}</span
          >
        </div>

      </div>

    </div>
     <!-- Thông tin chính -->
    <div class="content">

      <h3 class="title">{{ estate.Title }}</h3>

      <div class="stats">
         <span class="price">{{ formattedPrice }}</span
        > <span class="area">{{ formattedArea }}</span
        > <span class="price-per-m2">{{ formattedPricePerM2 }}</span
        > <!-- Phòng ngủ & toilet -->
        <div class="rooms">
           <span v-if="estate.Bedrooms" class="bedrooms"
            > <i class="icon-bed"></i> {{ estate.Bedrooms }} </span
          > <span v-if="estate.Bathrooms" class="bathrooms"
            > <i class="icon-bath"></i> {{ estate.Bathrooms }} </span
          >
        </div>
         <span class="location">{{ fullLocation }}</span
        >
      </div>
       <!-- Mô tả -->
      <p v-if="estate.Description" class="description">{{ truncatedDescription }}</p>
       <!-- Footer: người đăng & actions -->
      <div class="footer">

        <div class="agent-info">

          <div class="avatar">{{ agentInitial }}</div>

          <div class="agent-details">

            <div class="agent-name">{{ estate.AgentName || 'Người đăng' }}</div>

            <div class="post-time">{{ postTime }}</div>

          </div>

        </div>

        <div class="actions">
           <button v-if="estate.AgentPhone" class="btn-call" @click="handleCall">
             <i class="icon-phone"></i> {{ formattedPhone }} </button
          > <button
            class="btn-favorite"
            :class="{ active: isFavorite }"
            @click="handleToggleFavorite"
          >
             <i class="icon-heart"></i> </button
          >
        </div>

      </div>

    </div>

  </div>

</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { RealEstateModel } from '@/types/real_estate'

interface Props {
  estate: RealEstateModel
}

const props = defineProps<Props>()

const emit = defineEmits<{
  call: [phone: string]
  toggleFavorite: [id: number]
}>()

// Placeholder ảnh mặc định
const DEFAULT_IMAGE = ''

const isFavorite = ref(props.estate.IsFavorite || false)

// Ảnh chính (ảnh đầu tiên hoặc placeholder)
const mainImage = computed(() => {
  return props.estate.Images?.[0] || DEFAULT_IMAGE
})

// 3 ảnh thumbnail tiếp theo
const thumbnails = computed(() => {
  const images = props.estate.Images || []
  return images.slice(1, 4)
})

// Số ảnh còn lại (nếu có > 4 ảnh)
const remainingImagesCount = computed(() => {
  const total = props.estate.Images?.length || 0
  return Math.max(0, total - 4)
})

// Format giá: 7.7 tỷ, 850 triệu, etc.
const formattedPrice = computed(() => {
  const price = props.estate.PriceVND
  if (price >= 1_000_000_000) {
    return `${(price / 1_000_000_000).toFixed(1)} tỷ`
  }
  if (price >= 1_000_000) {
    return `${(price / 1_000_000).toFixed(0)} triệu`
  }
  return `${price.toLocaleString('vi-VN')} đ`
})

// Format diện tích: 41.8 m²
const formattedArea = computed(() => {
  return `${props.estate.Acreage.toFixed(1)} m²`
})

// Format đơn giá: 184.21 tr/m²
const formattedPricePerM2 = computed(() => {
  const pricePerM2 = props.estate.PricePerM2
  if (pricePerM2 >= 1_000_000) {
    return `${(pricePerM2 / 1_000_000).toFixed(2)} tr/m²`
  }
  return `${pricePerM2.toLocaleString('vi-VN')} đ/m²`
})

// Vị trí đầy đủ
const fullLocation = computed(() => {
  const parts = [props.estate.District, props.estate.City].filter(Boolean)
  return parts.join(', ')
})

// Mô tả rút gọn (tối đa 150 ký tự)
const truncatedDescription = computed(() => {
  const desc = props.estate.Description || ''
  return desc.length > 150 ? desc.substring(0, 150) + '...' : desc
})

// Chữ cái đầu tên người đăng
const agentInitial = computed(() => {
  const name = props.estate.AgentName || 'Q'
  return name.charAt(0).toUpperCase()
})

// Thời gian đăng
const postTime = computed(() => {
  const created = new Date(props.estate.CreatedAt)
  const now = new Date()
  const diffMs = now.getTime() - created.getTime()
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffDays = Math.floor(diffHours / 24)

  if (diffDays === 0) {
    return 'Đăng hôm nay'
  }
  if (diffDays === 1) {
    return 'Đăng hôm qua'
  }
  if (diffDays < 7) {
    return `Đăng ${diffDays} ngày trước`
  }
  return created.toLocaleDateString('vi-VN')
})

// Format số điện thoại
const formattedPhone = computed(() => {
  const phone = props.estate.AgentPhone || ''
  // Ẩn 3 số giữa: 0936 863 *** - Hiện số
  if (phone.length >= 10) {
    return `${phone.substring(0, 4)} ${phone.substring(4, 7)} *** - Hiện số`
  }
  return phone
})

// Xử lý lỗi load ảnh
const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = DEFAULT_IMAGE
}

// Xử lý gọi điện
const handleCall = () => {
  if (props.estate.AgentPhone) {
    emit('call', props.estate.AgentPhone)
  }
}

// Xử lý toggle yêu thích
const handleToggleFavorite = () => {
  isFavorite.value = !isFavorite.value
  emit('toggleFavorite', props.estate.ID)
}
</script>

<style scoped>
.real-estate-card {
  position: relative;
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: box-shadow 0.3s ease;
}

.real-estate-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

/* Badge VIP */
.badge {
  position: absolute;
  top: 12px;
  left: 12px;
  background: #e63946;
  color: #fff;
  padding: 4px 12px;
  font-size: 12px;
  font-weight: bold;
  border-radius: 4px;
  z-index: 10;
  text-transform: uppercase;
}

/* Grid ảnh */
.images-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 4px;
  height: 240px;
  background: #f0f0f0;
}

.main-image {
  grid-row: span 2;
}

.main-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumbnail-grid {
  display: grid;
  grid-template-rows: repeat(2, 1fr);
  gap: 4px;
}

.thumbnail {
  position: relative;
  overflow: hidden;
  background: #e0e0e0;
}

.thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumbnail.remaining-count {
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  font-size: 24px;
  font-weight: bold;
}

/* Nội dung */
.content {
  padding: 16px;
}

.title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 12px 0;
  text-transform: uppercase;
  line-height: 1.4;
}

.stats {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 12px;
  font-size: 14px;
  color: #666;
}

.price {
  color: #e63946;
  font-weight: bold;
  font-size: 18px;
}

.area,
.price-per-m2 {
  color: #333;
}

.rooms {
  display: flex;
  gap: 12px;
}

.bedrooms,
.bathrooms {
  display: flex;
  align-items: center;
  gap: 4px;
}

.icon-bed::before {
  content: '🛏️';
}

.icon-bath::before {
  content: '🚿';
}

.location {
  color: #666;
}

.description {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  margin: 12px 0;
}

/* Footer */
.footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #e0e0e0;
}

.agent-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #4a90e2;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 18px;
}

.agent-details {
  font-size: 13px;
}

.agent-name {
  font-weight: 600;
  color: #333;
}

.post-time {
  color: #999;
  font-size: 12px;
}

.actions {
  display: flex;
  gap: 8px;
}

.btn-call {
  background: #00bfa5;
  color: #fff;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: background 0.3s ease;
}

.btn-call:hover {
  background: #00a68a;
}

.icon-phone::before {
  content: '📞';
}

.btn-favorite {
  background: transparent;
  border: 2px solid #ddd;
  width: 40px;
  height: 40px;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.btn-favorite:hover {
  border-color: #e63946;
}

.btn-favorite.active {
  background: #e63946;
  border-color: #e63946;
}

.icon-heart::before {
  content: '🤍';
  font-size: 18px;
}

.btn-favorite.active .icon-heart::before {
  content: '❤️';
}
</style>

