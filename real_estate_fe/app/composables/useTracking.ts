import { useRuntimeConfig } from "#app";
import { useSession } from "~/composables/useSession";
import { useAuthStore } from "~/stores/auth";

export const useTracking = () => {
  const config = useRuntimeConfig();
  const { getSessionId } = useSession();
  const authStore = useAuthStore();

  const apiBaseUrl =
    config.public.apiBaseUrl || "http://localhost:8000/api/2026";

  // State lưu trữ ID và thời gian bắt đầu của tin BĐS đang xem
  const activeRealEstateId = ref<number | null>(null);
  const startTime = ref<number>(0);
  const accumulatedTime = ref<number>(0);
  const isTracking = ref<boolean>(true);

  // Hàm gửi dữ liệu lên Backend
  const sendTrackingData = () => {
    if (!isTracking.value || !activeRealEstateId.value) return;
    isTracking.value = false; // Đảm bảo chỉ gửi 1 lần khi rời trang

    // Tính toán tổng thời gian đã xem thực tế
    let durationMs = accumulatedTime.value;
    if (startTime.value > 0) {
      durationMs += Date.now() - startTime.value;
    }

    const durationSeconds = Math.round(durationMs / 1000);

    const sessionId = getSessionId();
    const payload = {
      real_estate_id: activeRealEstateId.value,
      duration_seconds: durationSeconds,
      session_id: sessionId,
      user_id: window.menu?.settings?.user_id
    };

    const url = `${apiBaseUrl}/tracking/view`;

    // 1. Tối ưu: Nếu trình duyệt hỗ trợ sendBeacon (Rất tin cậy khi đóng tab / tắt trình duyệt)
    if (process.client && navigator.sendBeacon) {
      const blob = new Blob([JSON.stringify(payload)], {
        type: "application/json; charset=UTF-8",
      });
      navigator.sendBeacon(url, blob);
    } else {
      // 2. Fallback sang fetch thông thường
      const headers: Record<string, string> = {
        "Content-Type": "application/json",
      };
      if (authStore.token) {
        headers["Authorization"] = `Bearer ${authStore.token}`;
      }
      fetch(url, {
        method: "POST",
        headers,
        body: JSON.stringify(payload),
        keepalive: true, // Giúp request tiếp tục chạy dù trang bị huỷ
      }).catch((err) =>
        console.error("Failed to send fallback tracking", err),
      );
    }
  };

  /**
   * Theo dõi thời gian xem một tin BĐS cụ thể
   * @param realEstateId ID của tin BĐS đang xem
   */
  const trackView = (realEstateId: number) => {
    if (!process.client) return;

    // Nếu chuyển sang xem một tin đăng khác khi chưa rời trang (ví dụ click BĐS lân cận)
    // thì gửi luôn tracking của tin cũ trước
    if (activeRealEstateId.value && activeRealEstateId.value !== realEstateId) {
      sendTrackingData();
    }

    // Reset và bắt đầu theo dõi tin mới
    activeRealEstateId.value = realEstateId;
    startTime.value = Date.now();
    accumulatedTime.value = 0;
    isTracking.value = true;
  };

  // Xử lý khi tab bị ẩn / hiện (Visibility Change)
  const handleVisibilityChange = () => {
    if (document.visibilityState === "hidden") {
      // Tab bị ẩn -> Pause timer, tích luỹ thời gian đã xem
      if (startTime.value > 0) {
        accumulatedTime.value += Date.now() - startTime.value;
        startTime.value = 0;
      }
    } else if (document.visibilityState === "visible") {
      // Tab hiện lại -> Resume timer
      if (activeRealEstateId.value) {
        startTime.value = Date.now();
      }
    }
  };

  const handlePageHide = () => {
    sendTrackingData();
  };

  if (process.client) {
    // Đăng ký sự kiện lắng nghe tab ẩn/hiện và đóng tab
    document.addEventListener("visibilitychange", handleVisibilityChange);
    window.addEventListener("pagehide", handlePageHide);

    // Gửi tracking khi Huỷ Component (Rời trang)
    onBeforeUnmount(() => {
      sendTrackingData();
      document.removeEventListener("visibilitychange", handleVisibilityChange);
      window.removeEventListener("pagehide", handlePageHide);
    });
  }

  return {
    trackView,
  };
};
