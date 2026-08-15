import { onBeforeUnmount, onMounted } from "vue"
import { useRuntimeConfig } from "#app"
import { useSession } from "~/composables/useSession"
import { useAuthStore } from "~/stores/auth"

export const useTracking = () => {
  const config = useRuntimeConfig()
  const { getSessionId } = useSession()
  const authStore = useAuthStore()

  const apiBaseUrl = config.public.apiBaseUrl || "http://localhost:8000/api/2026"

  /**
   * Theo dõi thời gian xem một tin BĐS cụ thể
   * @param realEstateId ID của tin BĐS đang xem
   */
  const trackView = (realEstateId: number) => {
    if (!process.client) return

    let startTime = Date.now()
    let accumulatedTime = 0
    let isTracking = true

    // Hàm gửi dữ liệu lên Backend
    const sendTrackingData = () => {
      if (!isTracking) return
      isTracking = false // Đảm bảo chỉ gửi 1 lần khi rời trang

      // Tính toán tổng thời gian đã xem thực tế
      let durationMs = accumulatedTime
      if (startTime > 0) {
        durationMs += Date.now() - startTime
      }

      const durationSeconds = Math.round(durationMs / 1000)

      const sessionId = getSessionId()
      const payload = {
        real_estate_id: realEstateId,
        duration_seconds: durationSeconds,
        session_id: sessionId,
      }

      const url = `${apiBaseUrl}/tracking/view`

      // 1. Tối ưu: Nếu trình duyệt hỗ trợ sendBeacon (Rất tin cậy khi đóng tab / tắt trình duyệt)
      if (navigator.sendBeacon) {
        const blob = new Blob([JSON.stringify(payload)], {
          type: "application/json; charset=UTF-8",
        })
        navigator.sendBeacon(url, blob)
      } else {
        // 2. Fallback sang fetch thông thường
        const headers: Record<string, string> = {
          "Content-Type": "application/json",
        }
        if (authStore.token) {
          headers["Authorization"] = `Bearer ${authStore.token}`
        }
        fetch(url, {
          method: "POST",
          headers,
          body: JSON.stringify(payload),
          keepalive: true, // Giúp request tiếp tục chạy dù trang bị huỷ
        }).catch((err) => console.error("Failed to send fallback tracking", err))
      }
    }

    // Xử lý khi tab bị ẩn / hiện (Visibility Change)
    const handleVisibilityChange = () => {
      if (document.visibilityState === "hidden") {
        // Tab bị ẩn -> Pause timer, tích luỹ thời gian đã xem
        if (startTime > 0) {
          accumulatedTime += Date.now() - startTime
          startTime = 0
        }
      } else if (document.visibilityState === "visible") {
        // Tab hiện lại -> Resume timer
        startTime = Date.now()
      }
    }

    // Gán sự kiện lắng nghe
    document.addEventListener("visibilitychange", handleVisibilityChange)

    // Gửi tracking khi Huỷ Component (Rời trang)
    onBeforeUnmount(() => {
      sendTrackingData()
      document.removeEventListener("visibilitychange", handleVisibilityChange)
    })

    // Dự phòng: Gửi tracking khi người dùng tắt hẳn tab/trình duyệt
    const handlePageHide = () => {
      sendTrackingData()
    }
    window.addEventListener("pagehide", handlePageHide)

    onBeforeUnmount(() => {
      window.removeEventListener("pagehide", handlePageHide)
    })
  }

  /**
   * Theo dõi hành vi tìm kiếm của người dùng
   * @param query Từ khóa tìm kiếm
   * @param filters Bộ lọc tìm kiếm đi kèm
   */
  const trackSearch = async (query: string, filters: any = {}) => {
    if (!process.client) return

    const sessionId = getSessionId()
    const url = `${apiBaseUrl}/tracking/search`

    const payload = {
      query: query,
      session_id: sessionId,
      filters: {
        location: filters.location ? {
          city: filters.location.city || null,
          district: filters.location.district || null,
          ward: filters.location.ward || null,
        } : null,
        price_range: filters.price_range ? {
          min_price: filters.price_range.min_price || null,
          max_price: filters.price_range.max_price || null,
        } : null,
      }
    }

    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    }
    if (authStore.token) {
      headers["Authorization"] = `Bearer ${authStore.token}`
    }

    try {
      await fetch(url, {
        method: "POST",
        headers,
        body: JSON.stringify(payload),
      })
    } catch (err) {
      console.error("Failed to track search", err)
    }
  }

  return {
    trackView,
    trackSearch,
  }
}
