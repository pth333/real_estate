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

            <!-- Danh sách (giới hạn ~5 dòng, nhiều hơn thì cuộn trong khung) -->
            <ul v-else class="mt-1 max-h-[300px] overflow-y-auto pr-1">
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
import { useGeoapify, type Coordinate, type NearbyPlace } from '~/composables/useGeoapify'
import IconSchool from '~/icons/IconSchool.vue'
import IconShopping from '~/icons/IconShopping.vue'
import IconPark from '~/icons/IconPark.vue'
import IconHospital from '~/icons/IconHospital.vue'
import IconRestaurant from '~/icons/IconRestaurant.vue'

const realEstateDetailStore = useRealEstateDetail()
const { createMap, getLeaflet, nearbyPlaces, distanceMeters, formatDistance, error: geoapifyError } = useGeoapify()

const listing = computed(() => realEstateDetailStore.listing)
const projectName = computed(() => listing.value?.title || 'Bất động sản')
const baseLat = computed(() => listing.value?.latitude ?? null)
const baseLng = computed(() => listing.value?.longitude ?? null)
const listingCoords = computed<Coordinate | null>(() =>
    baseLat.value !== null && baseLng.value !== null ? { lat: baseLat.value, lng: baseLng.value } : null,
)
const hasCoords = computed(() => listingCoords.value !== null)

interface NearbyCategory {
    key: string
    label: string
    heading: string
    // Danh mục Geoapify Places
    geoapifyCategory: string
    icon: any
}

// Địa điểm hiển thị: dữ liệu Places + khoảng cách/thời gian đi bộ
interface NearbyItem extends NearbyPlace {
    distance: string
    duration: string
}

const categories: NearbyCategory[] = [
    { key: 'school', label: 'Trường học', heading: 'trường học', geoapifyCategory: 'education.school', icon: IconSchool },
    { key: 'supermarket', label: 'Siêu thị', heading: 'siêu thị', geoapifyCategory: 'commercial.supermarket', icon: IconShopping },
    { key: 'park', label: 'Công viên', heading: 'công viên', geoapifyCategory: 'leisure.park', icon: IconPark },
    { key: 'hospital', label: 'Bệnh viện', heading: 'bệnh viện', geoapifyCategory: 'healthcare.hospital', icon: IconHospital },
    { key: 'restaurant', label: 'Nhà hàng', heading: 'nhà hàng', geoapifyCategory: 'catering.restaurant', icon: IconRestaurant },
]

const activeCategory = ref<NearbyCategory['key']>('school')
const activeCategoryData = computed<NearbyCategory>(() => categories.find(c => c.key === activeCategory.value) ?? categories[0]!)
const activeItems = ref<NearbyItem[]>([])
const loadingPlaces = ref(false)
const mapsError = computed(() => geoapifyError.value)

const activeCategoryHeading = computed(() => activeCategoryData.value.heading)
const activeCategoryIcon = computed(() => activeCategoryData.value.icon)
const totalCount = computed(() => activeItems.value.length)

// ── Geoapify / Leaflet ──
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

// Vẽ marker các địa điểm (xoá marker cũ trước)
const renderPlaceMarkers = (L: any) => {
    placeMarkerInstances.forEach(m => m.remove())
    placeMarkerInstances.length = 0
    if (!mapInstance || !listingCoords.value) return

    activeItems.value.forEach(item => {
        const marker = L.marker([item.lat, item.lng]).addTo(mapInstance)
        marker.bindPopup(`<strong>${item.name}</strong>`)
        placeMarkerInstances.push(marker)
    })
}

// Lấy danh sách tiện ích quanh vị trí bằng Geoapify Places API
const fetchPlaces = async () => {
    if (!listingCoords.value) return

    loadingPlaces.value = true
    try {
        const places = await nearbyPlaces(listingCoords.value, [activeCategoryData.value.geoapifyCategory], 2000)
        // Tính khoảng cách + thời gian đi bộ
        activeItems.value = places.map((p) => {
            const distance = distanceMeters(listingCoords.value!, { lat: p.lat, lng: p.lng })
            return {
                ...p,
                distance: formatDistance(distance),
                duration: estimateWalkDuration(distance),
            }
        })
    } catch (e) {
        console.error('Lỗi khi tải Places Nearby:', e)
        activeItems.value = []
    } finally {
        loadingPlaces.value = false
        const L = getLeaflet()
        if (L) renderPlaceMarkers(L)
    }
}

const initMap = async () => {
    if (!listingCoords.value) return
    if (mapInstance) {
        mapInstance.remove()
        mapInstance = null
        projectMarkerInstance = null
    }

    const el = document.getElementById('nearby-amenities-map')
    if (!el) return

    mapInstance = await createMap(el, { lat: listingCoords.value.lat, lng: listingCoords.value.lng, zoom: 15 })
    if (!mapInstance) return

    const L = getLeaflet()
    if (!L) return
    projectMarkerInstance = L.marker([listingCoords.value.lat, listingCoords.value.lng]).addTo(mapInstance)
    fetchPlaces()
}

// Định vị vị trí người dùng trên bản đồ
const locateMe = () => {
    if (!navigator.geolocation) return
    navigator.geolocation.getCurrentPosition(
        (pos) => {
            const L = getLeaflet()
            if (!L || !mapInstance) return
            const coord = { lat: pos.coords.latitude, lng: pos.coords.longitude }
            if (userMarkerInstance) userMarkerInstance.remove()
            userMarkerInstance = L.marker([coord.lat, coord.lng]).addTo(mapInstance)
            mapInstance.setView([coord.lat, coord.lng], 16)
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
