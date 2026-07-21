<template>
  <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
    <TransitionGroup name="toast">
      <div
        v-for="(toast, idx) in store.toasts"
        :key="idx"
        class="flex w-80 items-start gap-3 rounded-lg border border-green-200 bg-white p-4 shadow-lg"
      >
        <!-- Icon -->
        <div
          class="flex size-8 shrink-0 items-center justify-center rounded-full bg-green-100"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="size-5 text-green-600"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
            />
          </svg>
        </div>

        <!-- Nội dung -->
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-gray-800">BĐS mới!</p>
          <p class="mt-0.5 truncate text-xs text-gray-500">
            {{ toast.title }}
          </p>
          <p class="mt-0.5 text-xs text-gray-400">
            {{ formatPrice(toast.price_vnd) }} ·
            {{ toast.acreage ? toast.acreage + ' m²' : '' }}
          </p>
        </div>

        <!-- Nút đóng -->
        <button
          class="text-xs text-gray-400 hover:text-gray-600"
          @click="store.dismissToast(toast)"
        >
          ✕
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore } from '@/stores/notification'

const store = useNotificationStore()

function formatPrice(priceVND?: number): string {
  if (priceVND == null) return ''
  if (priceVND >= 1_000_000_000) return (priceVND / 1_000_000_000).toFixed(2) + ' tỷ'
  if (priceVND >= 1_000_000) return (priceVND / 1_000_000).toFixed(0) + ' tr'
  return priceVND.toLocaleString() + ' đ'
}
</script>

<style scoped>
.toast-enter-active {
  transition: all 0.3s ease-out;
}
.toast-leave-active {
  transition: all 0.2s ease-in;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
