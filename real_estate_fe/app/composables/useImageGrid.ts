/**
 * Xử lý image grid: main image, thumbnails, fallback
 */
export function useImageGrid(images: Ref<string[] | undefined>) {
  const DEFAULT_IMAGE = ''

  const mainImage = computed(() => images.value?.[0] || DEFAULT_IMAGE)

  const thumbnails = computed(() => {
    return images.value?.slice(1, 4) || []
  })

  const remainingImagesCount = computed(() => {
    const total = images.value?.length || 0
    return Math.max(0, total - 4)
  })

  function handleImageError(event: Event) {
    const img = event.target as HTMLImageElement
    img.src = DEFAULT_IMAGE
  }

  return { mainImage, thumbnails, remainingImagesCount, handleImageError }
}
