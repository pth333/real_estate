<template>
  <Teleport to="body">
    <ClientOnly>
      <div class="toast-container">
        <TransitionGroup name="toast-pop">
          <div
            v-for="toast in toasts"
            :key="toast.id"
            :class="['toast-item', `toast-${toast.type}`]"
          >
            <span class="toast-icon">
              <svg
                v-if="toast.type === 'success'"
                class="icon"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M20 6L9 17l-5-5" />
              </svg>
              <svg
                v-else-if="toast.type === 'error'"
                class="icon"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <circle cx="12" cy="12" r="10" />
                <path d="M15 9l-6 6M9 9l6 6" />
              </svg>
              <svg
                v-else
                class="icon"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <circle cx="12" cy="12" r="10" />
                <path d="M12 16v-4M12 8h.01" />
              </svg>
            </span>
            <span class="toast-message">{{ toast.message }}</span>
            <button class="toast-close" @click="dismissToast(toast.id)">
              ✕
            </button>
          </div>
        </TransitionGroup>
      </div>
    </ClientOnly>
  </Teleport>
</template>

<script setup lang="ts">
import { useToast } from "~/app/composables/useToast";
const { toasts, dismissToast } = useToast();
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 16px;
  right: 16px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-width: 400px;
}
.toast-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  font-size: 14px;
  line-height: 1.4;
}
.toast-success {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  color: #166534;
}
.toast-error {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #991b1b;
}
.toast-info {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  color: #1e40af;
}
.toast-icon {
  flex-shrink: 0;
}
.icon {
  width: 20px;
  height: 20px;
}
.toast-message {
  flex: 1;
}
.toast-close {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 14px;
  color: inherit;
  opacity: 0.5;
  padding: 0;
}
.toast-close:hover {
  opacity: 1;
}
.toast-pop-enter-active {
  transition: all 0.3s ease-out;
}
.toast-pop-leave-active {
  transition: all 0.2s ease-in;
}
.toast-pop-enter-from {
  opacity: 0;
  transform: translateX(40px);
}
.toast-pop-leave-to {
  opacity: 0;
  transform: translateX(40px);
}
</style>
