<template>
    <div class="bg-white rounded-lg p-4">
        <h2 class="mb-3 border-b border-gray-200 pb-2 text-base font-bold text-gray-800">
            Vị trí dự án {{ projectName }}
        </h2>

        <!-- Bản đồ -->
        <div class="relative overflow-hidden rounded-lg border border-gray-200">
            <div id="nearby-amenities-map" class="h-[260px] w-full"></div>
            <p v-if="!hasCoords"
                class="absolute inset-0 flex items-center justify-center bg-gray-100 text-sm text-gray-400">
                Không có vị trí trên bản đồ
            </p>
            <p v-else-if="mapsError"
                class="absolute inset-0 flex items-center justify-center bg-gray-100 p-4 text-center text-sm text-gray-500">
                {{ mapsError }}
            </p>
        </div>

        <!-- Tab danh mục tiện ích -->
        <div class="mt-2 flex items-center gap-5 overflow-x-auto border-b border-gray-100 py-2">
            <button v-for="cat in categories" :key="cat.key"
                class="flex shrink-0 items-center gap-1.5 text-sm font-medium transition-colors"
                :class="activeCategory === cat.key ? 'text-emerald-600' : 'text-gray-500 hover:text-gray-700'"
                @click="onSelectCategory(cat.key)">
                <component :is="cat.icon" class="h-4 w-4 shrink-0" />
                <span>{{ cat.label }}</span>
            </button>
            <button class="flex shrink-0 items-center gap-1.5 text-sm font-medium text-gray-500 hover:text-gray-700"
                @click="locateMe">
                <IconMapPin class="h-4 w-4 shrink-0" />
                <span>Vị trí của bạn</span>
            </button>
        </div>

        <!-- Danh sách tiện ích -->
        <div class="mt-3">
            <p class="text-sm text-gray-500">
                Có {{ totalCount }} {{ activeCategoryHeading }} trong vòng 2 km
            </p>

            <!-- Loading -->
            <div v-if="loadingPlaces" class="flex items-center justify-center py-8">
                <n-spin size="small" />
            </div>

            <!-- Không có dữ liệu -->
            <p v-else-if="activeItems.length === 0" class="py-8 text-center text-sm text-gray-400">
                Không tìm thấy {{ activeCategoryHeading }} nào trong khu vực
            </p>

            <!-- Danh sách -->
            <ul v-else class="mt-1">
                <li v-for="item in activeItems" :key="item.placeId"
                    class="flex items-center gap-3 border-b border-gray-100 py-3">
                    <component :is="activeCategoryIcon" class="h-4 w-4 shrink-0 text-emerald-500" />
                    <div class="min-w-0 flex-1">
                        <p class="truncate text-sm font-medium text-gray-800">{{ item.name }}</p>
                        <p class="truncate text-xs text-gray-400">{{ item.address }}</p>
                    </div>
                    <div class="shrink-0 text-right">
                        <p class="text-sm font-medium text-gray-700">{{ item.distance }}</p>
                        <p class="text-xs text-gray-400">{{ item.duration }}</p>
                    </div>
                </li>
            </ul>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useRealEstateDetail } from '~/stores/detail/real_estate_detail'
import { useGoogleMaps, type Coordinate } from '~/composables/useGoogleMaps'
import IconSchool from '~/icons/IconSchool.vue'
import IconShopping from '~/icons/IconShopping.vue'
import IconPark from '~/icons/IconPark.vue'
import IconHospital from '~/icons/IconHospital.vue'
import IconRestaurant from '~/icons/IconRestaurant.vue'

const realEstateDetailStore = useRealEstateDetail()
const { loadGoogleMaps, getMap, getMarker, importLibrary, distanceMeters, formatDistance, error: googleMapsError, toCoordinate, mapId } = useGoogleMaps()

const listing = computed(() => realEstateDetailStore.listing)
const projectName = computed(() => listing.value?.title || 'Bất động sản')
const baseLat = computed(() => listing.value?.latitude ?? null)
const baseLng = computed(() => listing.value?.longitude ?? null)
const listingCoords = computed<Coordinate | null>(() =>
    baseLat.value !== null && baseLng.value !== null ? { lat: baseLat.value, lng: baseLng.value } : null,
)
const hasCoords = computed(() => listingCoords.value !== null)

interface NearbyItem {
    placeId: string
    name: string
    address: string
    distance: string
    duration: string
    lat: number
    lng: number
}

interface NearbyCategory {
    key: string
    label: string
    heading: string
    // Loại Places API (Google Nearby Search)
    placeType: 'school' | 'supermarket' | 'park' | 'hospital' | 'restaurant'
    icon: any
}

const categories: NearbyCategory[] = [
    { key: 'school', label: 'Trường học', heading: 'trường học', placeType: 'school', icon: IconSchool },
    { key: 'supermarket', label: 'Siêu thị', heading: 'siêu thị', placeType: 'supermarket', icon: IconShopping },
    { key: 'park', label: 'Công viên', heading: 'công viên', placeType: 'park', icon: IconPark },
    { key: 'hospital', label: 'Bệnh viện', heading: 'bệnh viện', placeType: 'hospital', icon: IconHospital },
    { key: 'restaurant', label: 'Nhà hàng', heading: 'nhà hàng', placeType: 'restaurant', icon: IconRestaurant },
]

const activeCategory = ref<NearbyCategory['key']>('school')
const activeCategoryData = computed<NearbyCategory>(() => categories.find(c => c.key === activeCategory.value) ?? categories[0]!)
const activeItems = ref<NearbyItem[]>([])
const loadingPlaces = ref(false)
const mapsError = computed(() => googleMapsError.value)

const activeCategoryHeading = computed(() => activeCategoryData.value.heading)
const activeCategoryIcon = computed(() => activeCategoryData.value.icon)
const totalCount = computed(() => activeItems.value.length)

// ── Google Maps ──
let mapInstance: any = null
let projectMarkerInstance: any = null
let userMarkerInstance: any = null
const placeMarkerInstances: any[] = []

const onSelectCategory = (key: string) => {
    activeCategory.value = key as NearbyCategory['key']
    fetchPlaces()
}

// Ước lượng thời gian đi bộ từ khoảng cách (m). ~80m/phút, tối thiểu 1 phút.
const estimateWalkDuration = (distance: number): string => {
    const minutes = Math.max(1, Math.round(distance / 80))
    return `${minutes} phút`
}

// Vẽ marker các địa điểm (xoá marker cũ trước) bằng AdvancedMarkerElement
const renderPlaceMarkers = (AdvancedMarkerElement: any, g: any) => {
    placeMarkerInstances.forEach(m => { m.map = null })
    placeMarkerInstances.length = 0
    if (!mapInstance || !listingCoords.value) return

    activeItems.value.forEach(item => {
        const marker = new AdvancedMarkerElement({
            position: { lat: item.lat, lng: item.lng },
            map: mapInstance,
        })
        marker.addListener('gmp-click', () => {
            const info = new g.maps.InfoWindow({ content: `<strong>${item.name}</strong>` })
            info.open({ map: mapInstance, anchor: marker })
        })
        placeMarkerInstances.push(marker)
    })
}

// Lấy danh sách tiện ích quanh vị trí bằng Google Place (searchNearby) - API mới
const fetchPlaces = async () => {
    if (!listingCoords.value) return
    const { AdvancedMarkerElement } = await getMarker()
    // google global (googleLib) để dùng InfoWindow (core)
    const g = await loadGoogleMaps()
    if (!AdvancedMarkerElement || !g || !mapInstance) return

    const { Place } = await importLibrary('places')
    if (!Place) return

    loadingPlaces.value = true
    try {
        Place.searchNearby(
            {
                // location phải là LatLng instance (LatLngLiteral đôi khi gây "unknown property location")
                location: new g.maps.LatLng(listingCoords.value.lat, listingCoords.value.lng),
                radius: 2000,
                includedTypes: [activeCategoryData.value.placeType],
                maxResultCount: 20,
            },
            (results: any, status: string) => {
                if (status === 'OK' && results?.length) {
                    activeItems.value = results.map((r: any) => {
                        const placeCoord = toCoordinate(r.location)
                        const distance = distanceMeters(listingCoords.value!, placeCoord)
                        return {
                            placeId: r.placeId,
                            name: r.displayName?.text || 'Không có tên',
                            address: r.formattedAddress || r.displayName?.text || '',
                            distance: formatDistance(distance),
                            duration: estimateWalkDuration(distance),
                            lat: placeCoord.lat,
                            lng: placeCoord.lng,
                        }
                    })
                } else if (status === 'ZERO_RESULTS') {
                    activeItems.value = []
                } else {
                    // OVER_QUERY_LIMIT, REQUEST_DENIED... → giữ danh sách cũ
                    console.warn('Place searchNearby status:', status)
                }
                loadingPlaces.value = false
                renderPlaceMarkers(AdvancedMarkerElement, g)
            },
        )
    } catch (e) {
        console.error('Lỗi khi tải Places Nearby:', e)
        loadingPlaces.value = false
    }
}

const initMap = async () => {
    if (!listingCoords.value) return
    const { Map } = await getMap()
    const { AdvancedMarkerElement } = await getMarker()
    if (!Map || !AdvancedMarkerElement) return
    if (mapInstance) {
        mapInstance = null
        projectMarkerInstance = null
    }
    mapInstance = new Map(document.getElementById('nearby-amenities-map'), {
        center: listingCoords.value,
        zoom: 15,
        // mapId bắt buộc cho AdvancedMarkerElement
        ...(mapId.value ? { mapId: mapId.value } : {}),
    })
    projectMarkerInstance = new AdvancedMarkerElement({
        position: listingCoords.value,
        map: mapInstance,
    })
    fetchPlaces()
}

// Định vị vị trí người dùng trên bản đồ
const locateMe = () => {
    if (!navigator.geolocation) return
    navigator.geolocation.getCurrentPosition(
        async (pos) => {
            const { AdvancedMarkerElement } = await getMarker()
            if (!AdvancedMarkerElement || !mapInstance) return
            const coord = { lat: pos.coords.latitude, lng: pos.coords.longitude }
            if (userMarkerInstance) userMarkerInstance.map = null
            userMarkerInstance = new AdvancedMarkerElement({ position: coord, map: mapInstance })
            mapInstance.setCenter(coord)
            mapInstance.setZoom(16)
        },
        () => window.message?.warning('Không thể lấy vị trí của bạn'),
    )
}

watch(() => listing.value?.id, async () => {
    if (!import.meta.client) return
    await nextTick()
    initMap()
}, { immediate: true })

onUnmounted(() => {
    mapInstance = null
    projectMarkerInstance = null
    userMarkerInstance = null
    placeMarkerInstances.length = 0
})
</script>

<style scoped>
#nearby-amenities-map {
    outline: none;
}
</style>
