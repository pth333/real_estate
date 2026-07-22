<template>
  <Teleport to="body">
    <ClientOnly>
      <div class="fixed right-4 top-4 z-[9999] flex max-w-[400px] flex-col gap-2">
        <TransitionGroup name="toast-pop">
          <div
            v-for="toast in toasts"
            :key="toast.id"
            class="flex items-center gap-2.5 rounded-lg px-4 py-3 text-sm leading-relaxed shadow-lg"
            :class="{
              'border border-green-200 bg-green-50 text-green-800': toast.type === 'success',
              'border border-red-200 bg-red-50 text-red-800': toast.type === 'error',
              'border border-blue-200 bg-blue-50 text-blue-800': toast.type === 'info',
            }"
          >
            <span class="shrink-0">
              <svg v-if="toast.type === 'success'" class="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 6L9 17l-5-5" />
              </svg>
              <svg v-else-if="toast.type === 'error'" class="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10" />
                <path d="M15 9l-6 6M9 9l6 6" />
              </svg>
              <svg v-else class="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10" />
                <path d="M12 16v-4M12 8h.01" />
              </svg>
            </span>
            <span class="flex-1">{{ toast.message }}</span>
            <button class="cursor-pointer border-none bg-transparent p-0 text-sm opacity-50 hover:opacity-100" @click="dismissToast(toast.id)">
              ✕
            </button>
          </div>
        </TransitionGroup>
      </div>
    </ClientOnly>
  </Teleport>
</template>

<script setup lang="ts">
import { useToast } from "~/composables/useToast";
const { toasts, dismissToast } = useToast();
</script>

<style scoped>
.toast-pop-enter-active { transition: all 0.3s ease-out; }
.toast-pop-leave-active { transition: all 0.2s ease-in; }
.toast-pop-enter-from { opacity: 0; transform: translateX(40px); }
.toast-pop-leave-to { opacity: 0; transform: translateX(40px); }
</style>
