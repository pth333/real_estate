export function useImageGrid(images: Ref<string[] | undefined>) {
  const DEFAULT_IMAGE = ''

  const mainImage = computed(() => images.value?.[0] || DEFAULT_IMAGE)

  // Cột phải hiện 2 ảnh phụ: ảnh 2 (hàng trên), ảnh 3 (hàng dưới trái); +N ở hàng dưới phải
  const thumbnails = computed(() => {
    return images.value?.slice(1, 3) || []
  })

  // Ảnh nền cho ô +N (ảnh tiếp theo sau ảnh 3)
  const overlayImage = computed(() => images.value?.[3] || DEFAULT_IMAGE)

  const remainingImagesCount = computed(() => {
    const total = images.value?.length || 0
    return Math.max(0, total - 3)
  })

  function handleImageError(event: Event) {
    const img = event.target as HTMLImageElement
    img.src = DEFAULT_IMAGE
  }

  return { mainImage, thumbnails, overlayImage, remainingImagesCount, handleImageError }
}