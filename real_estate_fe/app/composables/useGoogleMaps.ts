/**
 * Composable dùng chung cho Google Maps JavaScript API.
 *
 * Load theo chuẩn khuyến nghị của Google:
 * - Script được load với `loading=async` (không dùng callback).
 * - Các thư viện lấy qua `google.maps.importLibrary()`.
 * - Dùng AdvancedMarkerElement / Place thay cho Marker / PlacesService đã deprecated.
 */

// Singleton: lưu google lib và mọi lần gọi đều dùng chung 1 script
let googleLib: any = null
let loadPromise: Promise<any> | null = null

export interface Coordinate {
  lat: number
  lng: number
}

// Trích lat/lng từ một đối tượng có thể là LatLngLiteral hoặc google.maps.LatLng
function toCoordinate(value: any): Coordinate {
  const lat = typeof value?.lat === 'function' ? value.lat() : value?.lat
  const lng = typeof value?.lng === 'function' ? value.lng() : value?.lng
  return { lat, lng }
}

export function useGoogleMaps() {
  const runtimeConfig = useRuntimeConfig()
  const apiKey = computed(() => runtimeConfig.public.googleMapsApiKey || '')
  // Map ID bắt buộc để dùng AdvancedMarkerElement (tạo trong Google Cloud Console)
  const mapId = computed(() => runtimeConfig.public.googleMapsMapId || '')
  const isLoaded = ref(false)
  const error = ref('')

  /**
   * Load Google Maps (lazily, một lần duy nhất) bằng `loading=async`.
   */
  const loadGoogleMaps = async (): Promise<any> => {
    if (googleLib) return googleLib
    if (!apiKey.value) {
      error.value = 'Thiếu Google Maps API key. Hãy điền NUXT_PUBLIC_GOOGLE_MAPS_API_KEY vào .env.'
      return null
    }
    if (loadPromise) return loadPromise

    loadPromise = new Promise((resolve, reject) => {
      const script = document.createElement('script')
      // loading=async: không dùng callback, các thư viện nạp qua importLibrary
      script.src = `https://maps.googleapis.com/maps/api/js?key=${apiKey.value}&loading=async`
      script.async = true
      script.defer = true
      script.onload = () => {
        const g = (window as any).google
        googleLib = g
        // importLibrary được gắn bất đồng bộ sau khi script load → chờ cho tới khi khả dụng
        const waitForImport = (attempt = 0) => {
          if (g?.maps?.importLibrary) {
            isLoaded.value = true
            resolve(g)
            return
          }
          if (attempt > 50) {
            error.value = 'Google Maps importLibrary không khả dụng'
            reject(new Error(error.value))
            return
          }
          setTimeout(() => waitForImport(attempt + 1), 100)
        }
        waitForImport()
      }
      script.onerror = () => {
        error.value = 'Không thể tải Google Maps'
        reject(new Error(error.value))
      }
      document.head.appendChild(script)
    })
    return loadPromise
  }

  /**
   * Lấy một thư viện Google Maps trả về promise (vd 'maps', 'marker', 'places', 'geometry').
   */
  const importLibrary = async (name: string): Promise<any> => {
    const g = await loadGoogleMaps()
    if (!g) return null
    return g.maps.importLibrary(name)
  }

  /** Lấy class Map. */
  const getMap = async () => (await importLibrary('maps')) || null

  /** Lấy class AdvancedMarkerElement (thay Marker deprecated). */
  const getMarker = async () => (await importLibrary('marker')) || null

  /** Đổi chuỗi địa chỉ → toạ độ (Geocoder). Trả null nếu không tìm thấy. */
  const geocodeAddress = async (address: string): Promise<Coordinate | null> => {
    const lib = await importLibrary('geocoding')
    const g = googleLib
    if (!lib?.Geocoder || !g) return null
    return new Promise((resolve) => {
      const geocoder = new lib.Geocoder()
      geocoder.geocode({ address }, (results: any, status: string) => {
        if (status === 'OK' && results?.[0]?.geometry?.location) {
          resolve({ lat: results[0].geometry.location.lat(), lng: results[0].geometry.location.lng() })
        } else {
          resolve(null)
        }
      })
    })
  }

  /**
   * Khoảng cách giữa 2 toạ độ theo công thức Haversine (trả về mét).
   * Tự tính đồng bộ để không phụ thuộc geometry library khi loading=async.
   */
  const distanceMeters = (a: Coordinate, b: Coordinate): number => {
    const R = 6371000
    const toRad = (deg: number) => (deg * Math.PI) / 180
    const dLat = toRad(b.lat - a.lat)
    const dLng = toRad(b.lng - a.lng)
    const lat1 = toRad(a.lat)
    const lat2 = toRad(b.lat)
    const h = Math.sin(dLat / 2) ** 2 + Math.cos(lat1) * Math.cos(lat2) * Math.sin(dLng / 2) ** 2
    return 2 * R * Math.asin(Math.sqrt(h))
  }

  /** Format khoảng cách mét → "1,2 km" hoặc "385,74 m". */
  const formatDistance = (meters: number): string => {
    if (meters >= 1000) return `${(meters / 1000).toFixed(1).replace('.', ',')} km`
    return `${meters.toFixed(0).replace('.', ',')} m`
  }

  return {
    apiKey,
    mapId,
    isLoaded,
    error,
    loadGoogleMaps,
    importLibrary,
    getMap,
    getMarker,
    geocodeAddress,
    distanceMeters,
    formatDistance,
    toCoordinate,
  }
}
