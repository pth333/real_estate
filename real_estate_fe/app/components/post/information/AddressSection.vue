<template>
    <n-card title="Địa chỉ" size="small" :segmented="{ content: true }">
        <template #header-extra>
            <div class="cursor-pointer" @click="collapsed = !collapsed">
                <n-icon class="text-gray-500 transition-transform duration-200" :class="{ 'rotate-180': !collapsed }">
                    <IconChevronDownOutline />
                </n-icon>
            </div>
        </template>
        <div v-show="!collapsed" class="flex flex-col gap-4">
            <n-form-item label="Khu vực"
                :feedback="postStore.errorsAddress.province || postStore.errorsAddress.detail_address"
                :validation-status="hasAddressError ? 'error' : undefined">
                <n-input readonly :value="locationLabel" placeholder="Chọn khu vực..." class="cursor-pointer"
                    @click="openModal">
                    <template #suffix>
                        <IconChevronRight class="h-4 w-4" />
                    </template>
                </n-input>
            </n-form-item>

            <!-- Bản đồ thu nhỏ hiển thị ở ngoài sau khi đã lưu toạ độ -->
            <div v-if="postStore.form.latitude && postStore.form.longitude" class="border border-gray-200 rounded-lg overflow-hidden bg-gray-50">
                <div class="bg-gray-100 px-3 py-2 text-xs font-semibold text-gray-700 flex justify-between items-center border-b border-gray-200">
                    <span>Vị trí trên bản đồ</span>
                    <span class="text-gray-500 font-mono text-[10px]">
                        Lat: {{ postStore.form.latitude.toFixed(6) }}, Lng: {{ postStore.form.longitude.toFixed(6) }}
                    </span>
                </div>
                <div id="static-address-map" class="h-48 w-full z-10"></div>
            </div>
        </div>
    </n-card>

    <n-modal v-model:show="showLocationModal" title="Nhập địa chỉ" style="width: 750px; max-height: 850px;" preset="card" :mask-closable="true" @after-enter="onModalOpen">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Form nhập địa chỉ -->
            <div class="flex flex-col gap-1">
                <n-form-item label="Tỉnh / Thành phố" path="province" :feedback="postStore.errorsAddress.province"
                    :validation-status="postStore.errorsAddress.province ? 'error' : undefined">
                    <n-select v-model:value="postStore.form.province" placeholder="Chọn tỉnh/thành phố" clearable filterable
                        :options="provinceOption" @update:value="onProvinceChange" />
                </n-form-item>

                <n-form-item label="Phường / Xã" path="ward" :feedback="postStore.errorsAddress.ward"
                    :validation-status="postStore.errorsAddress.ward ? 'error' : undefined">
                    <n-select v-model:value="postStore.form.ward" placeholder="Chọn phường/xã" clearable filterable
                        :options="wardOptions" :loading="loadingWard" @update:value="onWardSelectChange"
                        :disabled="isDisabledWard" />
                </n-form-item>

                <n-form-item label="Địa chỉ chi tiết" path="detail" :feedback="postStore.errorsAddress.detail_address"
                    :validation-status="postStore.errorsAddress.detail_address ? 'error' : undefined">
                    <n-input v-model:value="postStore.form.detail_address" placeholder="Nhập số nhà, khu phố, ngõ hẻm..."
                        clearable @update:value="onDetailAddressInput" />
                </n-form-item>

                <n-form-item label="Dự án" path="project_id">
                    <n-select v-model:value="postStore.form.project_id" placeholder="Chọn dự án" clearable filterable
                        :options="projectOptions" :disabled="isDisabledPJ"/>
                </n-form-item>
            </div>

            <!-- Khu vực Bản đồ tương tác -->
            <div class="flex flex-col border border-gray-200 rounded-lg overflow-hidden h-[380px] md:h-auto bg-gray-50">
                <div class="bg-gray-100 px-3 py-2 text-xs font-semibold text-gray-700 flex justify-between items-center border-b border-gray-200">
                    <span>Chọn vị trí chính xác</span>
                    <n-button size="tiny" secondary type="info" :disabled="!isGeocodeBtnEnabled" @click="triggerGeocoding">
                        Định vị địa chỉ
                    </n-button>
                </div>
                <div class="relative flex-1">
                    <div id="interactive-address-map" class="absolute inset-0 z-10"></div>
                    <div v-if="loadingGeocode" class="absolute inset-0 bg-white/60 z-20 flex items-center justify-center">
                        <n-spin size="medium" />
                    </div>
                </div>
            </div>
        </div>

        <template #action>
            <div class="flex justify-end">
                <n-button type="primary" @click="applyAddress">Áp dụng</n-button>
            </div>
        </template>
    </n-modal>
</template>
<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import type { CityOption, WardOption, ProjectOption } from '~/types/real_estate'
import { useCreatePost } from '~/stores/create-post'
import { useGoogleMaps } from '~/composables/useGoogleMaps'

const { $api } = useNuxtApp()
const postStore = useCreatePost()
const { getMap, getMarker, geocodeAddress, toCoordinate, mapId } = useGoogleMaps()

const collapsed = ref(false)
const showLocationModal = ref(false)
const isDisabledWard = computed(() => !postStore.form.province)
const isDisabledPJ = computed(() => !postStore.form.province || !postStore.form.ward)

// Trạng thái load bản đồ và geocoding
const loadingGeocode = ref(false)
let interactiveMapInstance: any = null
let interactiveMarkerInstance: any = null
let staticMapInstance: any = null
let staticMarkerInstance: any = null

const openModal = () => {
    showLocationModal.value = true
}

const onModalOpen = async () => {
    try {
        const { Map } = await getMap();
        const { AdvancedMarkerElement } = await getMarker();
        if (!Map || !AdvancedMarkerElement) return;
        initInteractiveMap(Map, AdvancedMarkerElement);
    } catch (e) {
        console.error('Không thể tải Google Maps:', e);
    }
}

// Khởi tạo bản đồ tương tác trong Modal (Google Maps, AdvancedMarkerElement)
const initInteractiveMap = (Map: any, AdvancedMarkerElement: any) => {
    if (interactiveMapInstance) {
        interactiveMapInstance = null;
        interactiveMarkerInstance = null;
    }

    // Tọa độ mặc định: tâm Việt Nam (Đà Nẵng) hoặc Hồ Chí Minh nếu chưa có tọa độ
    const initLat = postStore.form.latitude || 10.7769;
    const initLng = postStore.form.longitude || 106.7009;
    const zoom = postStore.form.latitude ? 16 : 10;

    interactiveMapInstance = new Map(document.getElementById('interactive-address-map'), {
        center: { lat: initLat, lng: initLng },
        zoom,
        // mapId bắt buộc cho AdvancedMarkerElement
        ...(mapId.value ? { mapId: mapId.value } : {}),
    });

    interactiveMarkerInstance = new AdvancedMarkerElement({
        position: { lat: initLat, lng: initLng },
        map: interactiveMapInstance,
        gmpDraggable: true,
    });

    // Lắng nghe sự kiện kéo thả marker
    interactiveMarkerInstance.addListener('gmp-dragend', () => {
        const position = toCoordinate(interactiveMarkerInstance.position);
        postStore.form.latitude = position.lat;
        postStore.form.longitude = position.lng;
    });

    // Lắng nghe click lên bản đồ để di chuyển marker
    interactiveMapInstance.addListener('click', (e: any) => {
        interactiveMarkerInstance.position = e.latLng;
        postStore.form.latitude = e.latLng.lat();
        postStore.form.longitude = e.latLng.lng();
    });
};

// Khởi tạo bản đồ tĩnh bên ngoài khi có tọa độ và thu nhỏ form
const initStaticMap = async () => {
    if (!postStore.form.latitude || !postStore.form.longitude) return;

    try {
        const { Map } = await getMap();
        const { AdvancedMarkerElement } = await getMarker();
        if (!Map || !AdvancedMarkerElement) return;

        // Đợi DOM render xong thẻ chứa bản đồ tĩnh
        await nextTick();

        if (staticMapInstance) {
            staticMapInstance = null;
            staticMarkerInstance = null;
        }

        const mapEl = document.getElementById('static-address-map');
        if (!mapEl) return;

        // Bản đồ thu nhỏ, chỉ xem, không tương tác
        staticMapInstance = new Map(mapEl, {
            center: { lat: postStore.form.latitude, lng: postStore.form.longitude },
            zoom: 16,
            disableDefaultUI: true,
            draggable: false,
            scrollwheel: false,
            zoomControl: false,
            // mapId bắt buộc cho AdvancedMarkerElement
            ...(mapId.value ? { mapId: mapId.value } : {}),
        });

        staticMarkerInstance = new AdvancedMarkerElement({
            position: { lat: postStore.form.latitude, lng: postStore.form.longitude },
            map: staticMapInstance,
        });

    } catch (e) {
        console.error('Lỗi khi vẽ bản đồ tĩnh bên ngoài:', e);
    }
};

// Áp dụng địa chỉ: nếu thiếu field bắt buộc thì hiện lỗi ngay trong modal, không tắt
const applyAddress = () => {
    const valid = postStore.validateAddress()
    if (!valid) return

    showLocationModal.value = false;
    initStaticMap();
}

const loadingWard = ref(false)
const provinceOption = ref<SelectOption[]>([])
const wardOptions = ref<SelectOption[]>([])
const projectOptions = ref<SelectOption[]>([])

const hasAddressError = computed(() =>
    Object.values(postStore.errorsAddress).some((message) => message !== '')
)

const clearError = (field: keyof typeof postStore.errorsAddress) => {
    postStore.errorsAddress[field] = ''
}

const locationLabel = computed(() => {
    const city = provinceOption.value.find(item => item.value === postStore.form.province)
    const ward = wardOptions.value.find(item => item.value === postStore.form.ward)

    const parts: string[] = [];
    if (city) parts.push(city.label as string);
    if (ward) parts.push(ward.label as string);
    if (postStore.form.detail_address) parts.push(postStore.form.detail_address);
    if (parts.length === 0) return 'Chọn khu vực';
    return parts.join(', ');
});

const fetchListProvice = async () => {
    try {
        const res = await $api.get("/real-estate/list/city") as { data: CityOption[] }
        provinceOption.value = res.data.map((item: CityOption) => ({
            label: item.name,
            value: item.code
        }))
    } catch (error) {
        console.error("Lỗi khi tải danh sách tỉnh/thành phố:", error)
        provinceOption.value = []
    }
}

// Nếu form đang giữ NAME tỉnh (từ res.city khi edit tin cũ) → đổi sang CODE
// để khớp option select + fetch phường/dự án theo mã.
const normalizeProvinceToCode = () => {
    const current = postStore.form.province
    if (!current) return
    const matched = provinceOption.value.find(o => o.value === current)
    if (matched) return // đã là code
    const byLabel = provinceOption.value.find(o => o.label === current)
    if (byLabel) {
        postStore.form.province = byLabel.value as string
    }
}

const onProvinceChange = (provinceCode: string | null) => {
    onWardChange(provinceCode);
    // Tự động geocode lại nếu chọn tỉnh mới
    triggerGeocoding();
}

const onWardSelectChange = (wardCode: string | null) => {
    clearError('ward');
    // Tự động geocode lại khi chọn phường/xã
    triggerGeocoding();
}

const onDetailAddressInput = (val: string) => {
    clearError('detail_address');
}

// Hàm debounce/tránh gọi Geocoding API dồn dập
let geocodeTimeout: any = null;
const onDetailAddressDebounced = () => {
    if (geocodeTimeout) clearTimeout(geocodeTimeout);
    geocodeTimeout = setTimeout(() => {
        triggerGeocoding();
    }, 1500);
}

const onWardChange = async (provinceCode: string | null) => {
    // Giữ name phường cũ (từ res.district khi edit) để normalize sang code sau khi load
    const previousWard = postStore.form.ward
    postStore.form.ward = null
    wardOptions.value = []
    clearError('province')

    if (!provinceCode) return

    try {
        loadingWard.value = true
        const res = await $api.get<{ data: WardOption[] }>(`/real-estate/list/ward`, {
            params: { code: provinceCode }
        })
        wardOptions.value = res.data.map(item => ({
            label: item.name,
            value: item.code
        }))
        // Nếu trước đó form giữ NAME phường (edit tin cũ) → đổi sang CODE
        if (previousWard) {
            const byLabel = wardOptions.value.find(o => o.label === previousWard)
            if (byLabel) {
                postStore.form.ward = byLabel.value as string
            }
        }
    } catch (error) {
        console.error("Lỗi khi tải danh sách phường/xã:", error)
        wardOptions.value = []
    } finally {
        loadingWard.value = false
    }
}

const fetchListProject = async () => {
    try {
        const params: Record<string, string> = {}
        if (postStore.form.province) {
            params.province = postStore.form.province
        }
        if (postStore.form.ward) {
            params.ward = postStore.form.ward
        }

        const res = await $api.get<{ data: ProjectOption[] }>("/real-estate/list/project", {
            params
        })
        projectOptions.value = res.data.map((item: ProjectOption) => ({
            label: item.name,
            value: item.id
        }))
    } catch (error) {
        console.error("Lỗi khi tải danh sách dự án:", error)
        projectOptions.value = []
    }
}

// Geocoding địa chỉ chuỗi thông qua Google Maps Geocoder
const isGeocodeBtnEnabled = computed(() => {
    return !!(postStore.form.province || postStore.form.ward || postStore.form.detail_address);
});

const triggerGeocoding = async () => {
    if (!isGeocodeBtnEnabled.value) return;

    const city = (provinceOption.value.find(item => item.value === postStore.form.province)?.label as string) || '';
    const ward = (wardOptions.value.find(item => item.value === postStore.form.ward)?.label as string) || '';
    const detail = postStore.form.detail_address || '';

    // Xây dựng địa chỉ tìm kiếm theo cấu trúc chuẩn
    const queryParts: string[] = [];
    if (detail) queryParts.push(detail);
    if (ward) queryParts.push(ward as any);
    if (city) queryParts.push(city as any);
    queryParts.push('Việt Nam');

    const queryString = queryParts.join(', ');

    try {
        loadingGeocode.value = true;
        const coord = await geocodeAddress(queryString);

        if (coord) {
            postStore.form.latitude = coord.lat;
            postStore.form.longitude = coord.lng;

            // Cập nhật lên bản đồ tương tác
            if (interactiveMapInstance && interactiveMarkerInstance) {
                interactiveMarkerInstance.position = { lat: coord.lat, lng: coord.lng };
                interactiveMapInstance.setCenter({ lat: coord.lat, lng: coord.lng });
                interactiveMapInstance.setZoom(16);
            }
        }
    } catch (err) {
        console.error('Lỗi khi geocoding địa chỉ:', err);
    } finally {
        loadingGeocode.value = false;
    }
};

// Cập nhật vị trí marker khi gõ tay toạ độ trong modal
const onCoordsManualChange = () => {
    const lat = postStore.form.latitude;
    const lng = postStore.form.longitude;
    if (lat && lng && interactiveMapInstance && interactiveMarkerInstance) {
        interactiveMarkerInstance.position = { lat, lng };
        interactiveMapInstance.setCenter({ lat, lng });
        interactiveMapInstance.setZoom(16);
    }
}

watch(
    () => [postStore.form.province, postStore.form.ward],
    async ([newProvince, newWard], [oldProvince, oldWard]) => {
        if (newProvince !== oldProvince || newWard !== oldWard) {
            postStore.form.project_id = null
        }
        await fetchListProject()
    }
)

watch(
    () => postStore.form.detail_address,
    () => {
        onDetailAddressDebounced();
    }
)

// Khởi tạo bản đồ tĩnh bên ngoài nếu đã có dữ liệu toạ độ sẵn (VD: edit tin đăng cũ, hoặc load nháp)
watch(
    () => [postStore.form.latitude, postStore.form.longitude, collapsed.value],
    async ([lat, lng, isCollapsed]) => {
        if (lat && lng && !isCollapsed) {
            await nextTick();
            initStaticMap();
        }
    }
)

onMounted(async () => {
    await fetchListProvice()
    // Edit tin cũ: form.province đang là NAME → đổi sang CODE trước khi load phường
    normalizeProvinceToCode()
    if (postStore.form.province) {
        await onWardChange(postStore.form.province)
    }
    if (postStore.form.latitude && postStore.form.longitude) {
        initStaticMap();
    }
})
</script>

<style scoped>
/* Ensure map is rendered correctly on high DPI screen */
#interactive-address-map, #static-address-map {
    outline: none;
}
</style>
