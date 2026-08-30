/**
 * Composable thao tác Geoapify — thay thế Google Maps.
 * - Bản đồ: Leaflet (CDN) + Geoapify tile layer (không cần Maps JS API).
 * - Geocoding API (REST): đổi chuỗi địa chỉ → toạ độ.
 * - Places API (REST): lấy tiện ích quanh vị trí.
 * - Static Map API (REST): ảnh bản đồ tĩnh.
 * Key đọc từ runtimeConfig.public.geoapifyApiKey (NUXT_PUBLIC_GEOAPIFY_API_KEY).
 */

export interface Coordinate {
  lat: number;
  lng: number;
}

export interface NearbyPlace {
  placeId: string;
  name: string;
  address: string;
  lat: number;
  lng: number;
}

// Singleton Leaflet (load từ CDN một lần)
let leafletLib: any = null;
let leafletPromise: Promise<any> | null = null;

// Trích lat/lng từ LatLng hoặc LatLngLiteral (Leaflet LatLng có .lat/.lng là hàm)
function toCoordinate(value: any): Coordinate {
  const lat = typeof value?.lat === "function" ? value.lat() : value?.lat;
  const lng = typeof value?.lng === "function" ? value.lng() : value?.lng;
  return { lat, lng };
}

export function useGeoapify() {
  const runtimeConfig = useRuntimeConfig();
  const apiKey = computed(() => runtimeConfig.public.geoapifyApiKey || "");
  const error = ref("");

  /**
   * Load Leaflet từ CDN (một lần duy nhất). Không phụ thuộc Maps JS API của Geoapify.
   */
  const loadLeaflet = (): Promise<any> => {
    if (leafletLib) return Promise.resolve(leafletLib);
    if ((window as any).L) {
      leafletLib = (window as any).L;
      return Promise.resolve(leafletLib);
    }
    if (leafletPromise) return leafletPromise;

    leafletPromise = new Promise((resolve, reject) => {
      const link = document.createElement("link");
      link.rel = "stylesheet";
      link.href = "https://unpkg.com/leaflet@1.9.4/dist/leaflet.css";
      document.head.appendChild(link);

      const script = document.createElement("script");
      script.src = "https://unpkg.com/leaflet@1.9.4/dist/leaflet.js";
      script.onload = () => {
        leafletLib = (window as any).L;
        resolve(leafletLib);
      };
      script.onerror = () => {
        leafletPromise = null;
        reject(new Error("Không thể tải Leaflet"));
      };
      document.body.appendChild(script);
    });
    return leafletPromise;
  };

  /**
   * Tạo bản đồ Leaflet với tile Geoapify. interactive=false → bản đồ tĩnh (chỉ xem).
   */
  const createMap = async (
    el: HTMLElement,
    opts: { lat: number; lng: number; zoom?: number; interactive?: boolean },
  ): Promise<any> => {
    const L = await loadLeaflet();
    if (!L) return null;

    const map = L.map(el).setView([opts.lat, opts.lng], opts.zoom ?? 14);

    // Tile OpenStreetMap (public, không cần key) — đảm bảo bản đồ luôn hiển thị
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors',
      maxZoom: 19,
    }).addTo(map);

    if (opts.interactive === false) {
      map.dragging.disable();
      map.scrollWheelZoom.disable();
      map.doubleClickZoom.disable();
      map.touchZoom.disable();
      map.keyboard.disable();
      map.zoomControl.remove();
    }

    // Container vừa được hiển thị (modal) → báo Leaflet tính lại kích thước để vẽ đúng
    requestAnimationFrame(() => map.invalidateSize());
    return map;
  };

  /** Trả Leaflet (L) để tạo marker (Geoapify dùng Leaflet). */
  const getLeaflet = (): any => leafletLib;

  /**
   * URL ảnh static map (Geoapify Static Map API) — không cần tải Maps JS.
   * @param marker true → vẽ marker đỏ tại tọa độ.
   */
  const staticMapUrl = (
    lat: number,
    lng: number,
    opts?: {
      width?: number;
      height?: number;
      zoom?: number;
      style?: string;
      marker?: boolean;
    },
  ): string => {
    if (!apiKey.value) return "";
    const {
      width = 600,
      height = 400,
      zoom = 16,
      style = "osm-bright",
      marker = true,
    } = opts || {};
    let url =
      `https://maps.geoapify.com/v1/staticmap?style=${style}&width=${width}&height=${height}` +
      `&center=lonlat:${lng},${lat}&zoom=${zoom}&apiKey=${apiKey.value}`;
    if (marker) {
      url += `&marker=lonlat:${lng},${lat};type:material;color:%23ef4444;size:medium;text:P`;
    }
    return url;
  };

  /**
   * Geocode chuỗi địa chỉ → toạ độ (Geoapify Geocoding REST).
   */
  const geocode = async (text: string): Promise<Coordinate | null> => {
    if (!apiKey.value || !text) return null;
    try {
      const url = `https://api.geoapify.com/v1/geocode/search?text=${encodeURIComponent(text)}&lang=vi&apiKey=${apiKey.value}`;
      const res = await fetch(url);
      const data = await res.json();
      console.log(data)
      if (data.features?.length) {
        const [lng, lat] = data.features[0].geometry.coordinates;
        return { lat, lng };
      }
      return null;
    } catch {
      return null;
    }
  };

  /**
   * Lấy danh sách tiện ích quanh vị trí (Geoapify Places REST).
   * @param categories danh mục Geoapify, VD ["education.school", "commercial.supermarket"]
   */
  const nearbyPlaces = async (
    coord: Coordinate,
    categories: string[],
    radius = 2000,
  ): Promise<NearbyPlace[]> => {
    if (!apiKey.value || categories.length === 0) return [];
    try {
      const url =
        `https://api.geoapify.com/v2/places?categories=${encodeURIComponent(categories.join(","))}` +
        `&filter=circle:${coord.lng},${coord.lat},${radius}&limit=20&apiKey=${apiKey.value}`;
      const res = await fetch(url);
      const data = await res.json();
      return (data.features || []).map((f: any) => ({
        placeId: String(f.properties.place_id || ""),
        name: f.properties.name || f.properties.address_line1 || "Không có tên",
        address: f.properties.address_line2 || f.properties.address_line1 || "",
        lat: f.geometry.coordinates[1],
        lng: f.geometry.coordinates[0],
      }));
    } catch {
      return [];
    }
  };

  /**
   * Khoảng cách giữa 2 toạ độ theo Haversine (mét).
   */
  const distanceMeters = (a: Coordinate, b: Coordinate): number => {
    const R = 6371000;
    const toRad = (deg: number) => (deg * Math.PI) / 180;
    const dLat = toRad(b.lat - a.lat);
    const dLng = toRad(b.lng - a.lng);
    const lat1 = toRad(a.lat);
    const lat2 = toRad(b.lat);
    const h =
      Math.sin(dLat / 2) ** 2 +
      Math.cos(lat1) * Math.cos(lat2) * Math.sin(dLng / 2) ** 2;
    return 2 * R * Math.asin(Math.sqrt(h));
  };

  /** Format mét → "1,2 km" hoặc "385,74 m". */
  const formatDistance = (meters: number): string => {
    if (meters >= 1000)
      return `${(meters / 1000).toFixed(1).replace(".", ",")} km`;
    return `${meters.toFixed(0).replace(".", ",")} m`;
  };

  return {
    apiKey,
    error,
    loadLeaflet,
    createMap,
    getLeaflet,
    staticMapUrl,
    geocode,
    nearbyPlaces,
    distanceMeters,
    formatDistance,
    toCoordinate,
  };
}
