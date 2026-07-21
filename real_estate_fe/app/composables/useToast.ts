export interface ToastMessage {
  id: number;
  type: "success" | "error" | "info";
  message: string;
}

// Singleton state dùng chung toàn app
const toasts = ref<ToastMessage[]>([]);
let nextId = 0;

export function useToast() {
  function showToast(
    type: ToastMessage["type"],
    message: string,
    duration = 4000,
  ) {
    const id = ++nextId;
    toasts.value.push({ id, type, message });

    // Tự động xoá sau duration ms
    setTimeout(() => dismissToast(id), duration);
  }

  function dismissToast(id: number) {
    const idx = toasts.value.findIndex((t: ToastMessage) => t.id === id);
    if (idx !== -1) toasts.value.splice(idx, 1);
  }

  return { toasts, showToast, dismissToast };
}
