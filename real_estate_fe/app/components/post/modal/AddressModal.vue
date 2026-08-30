<template>
    <n-modal :show="show" title="Nhập địa chỉ" style="width: 750px; max-height: 850px;" preset="card"
        :mask-closable="true" @update:show="onShowChange" @after-enter="onModalOpen">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Form nhập địa chỉ -->
            <div class="flex flex-col gap-1">
                <n-form-item label="Tỉnh / Thành phố" path="province" :feedback="postStore.errorsAddress.province"
                    :validation-status="postStore.errorsAddress.province ? 'error' : undefined">
                    <n-select v-model:value="postStore.form.province" placeholder="Chọn tỉnh/thành phố" clearable
                        filterable :options="provinceOptions" @update:value="onProvinceChange" />
                </n-form-item>

                <n-form-item label="Phường / Xã" path="ward" :feedback="postStore.errorsAddress.ward"
                    :validation-status="postStore.errorsAddress.ward ? 'error' : undefined">
                    <n-select v-model:value="postStore.form.ward" placeholder="Chọn phường/xã" clearable filterable
                        :options="wardOptions" :loading="loadingWard" @update:value="onWardSelectChange"
                        :disabled="isDisabledWard" />
                </n-form-item>

                <n-form-item label="Địa chỉ chi tiết" path="detail" :feedback="postStore.errorsAddress.detail_address"
                    :validation-status="postStore.errorsAddress.detail_address ? 'error' : undefined">
                    <n-input v-model:value="postStore.form.detail_address"
                        placeholder="Nhập số nhà, khu phố, ngõ hẻm..." clearable @update:value="onDetailAddressInput" />
                </n-form-item>

                <n-form-item label="Dự án" path="project_id">
                    <n-select v-model:value="postStore.form.project_id" placeholder="Chọn dự án" clearable filterable
                        :options="projectOptions" :disabled="isDisabledPJ" />
                </n-form-item>
            </div>

            <!-- Bản đồ tương tác: kéo/click marker để chọn toạ độ chính xác -->
            <div class="flex flex-col border border-gray-200 rounded-lg overflow-hidden h-[380px] bg-gray-50">
                <div
                    class="bg-gray-100 px-3 py-2 text-xs font-semibold text-gray-700 flex justify-between items-center border-b border-gray-200">
                    <span>Chọn vị trí chính xác</span>
                </div>
                <div class="relative flex-1 min-h-0">
                    <div id="interactive-address-map" class="absolute inset-0 z-10"></div>
                    <!-- Loading khi đang định vị địa chỉ -->
                    <div v-if="loadingGeocode"
                        class="absolute inset-0 bg-white/60 z-20 flex items-center justify-center">
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
import { useGeoapify } from '~/composables/useGeoapify'

const props = defineProps<{
    show: boolean
}>()

const emit = defineEmits<{
    'update:show': [value: boolean]
    'update:locationLabel': [value: string]
}>()

const { $api } = useNuxtApp()
// Form dùng chung qua store; options do chính modal load
const postStore = useCreatePost()
const { createMap, getLeaflet, geocode } = useGeoapify()

const isDisabledWard = computed(() => !postStore.form.province)
const isDisabledPJ = computed(() => !postStore.form.province || !postStore.form.ward)

// ── Options tỉnh/phường — modal tự load ──
const provinceOptions = ref<SelectOption[]>([])
const wardOptions = ref<SelectOption[]>([])
const loadingWard = ref(false)

const fetchListProvice = async () => {
    try {
        const res = await $api.get("/real-estate/list/city") as { data: CityOption[] }
        provinceOptions.value = res.data.map((item: CityOption) => ({
            label: item.name,
            value: item.code
        }))
    } catch (error) {
        console.error("Lỗi khi tải danh sách tỉnh/thành phố:", error)
        provinceOptions.value = []
    }
}

// Load phường/xã theo tỉnh (gọi từ watch form.province — gồm cả khi edit API set province)
const loadWards = async (provinceCode: string | null) => {
    // Giữ name phường cũ (từ res.district khi edit) để normalize sang code sau khi load
    const previousWard = postStore.form.ward
    postStore.form.ward = null
    wardOptions.value = []
    postStore.errorsAddress.province = ''

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
        // Edit tin cũ: form.ward đang là NAME → đổi sang CODE để khớp option
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

// Khi province thay đổi (set từ API khi edit, hoặc user chọn) → load phường/xã.
// Edit: form.province đang là NAME (từ res.city) → đổi sang CODE trước khi fetch.
// Watch cả provinceOptions để chờ danh sách tỉnh load xong rồi mới normalize.
watch(
    [() => postStore.form.province, provinceOptions],
    ([province]) => {
        if (!province) return
        const matched = provinceOptions.value.find(o => o.value === province)
        if (matched) {
            // Đã là code → load phường trực tiếp
            loadWards(province)
            return
        }
        // Đang là NAME → đổi sang CODE (watch sẽ chạy lại với code)
        const byLabel = provinceOptions.value.find(o => o.label === province)
        if (byLabel && byLabel.value !== province) {
            postStore.form.province = byLabel.value as string
        }
    },
    { immediate: true },
)

onMounted(fetchListProvice)

// ── Bản đồ ──
const loadingGeocode = ref(false)
let interactiveMapInstance: any = null
let interactiveMarkerInstance: any = null

const onShowChange = (value: boolean) => {
    emit('update:show', value)
}

// Mở modal → khởi tạo bản đồ tại toạ độ hiện có (hoặc mặc định) + marker
const onModalOpen = async () => {
    try {
        await nextTick();
        const el = document.getElementById('interactive-address-map');
        if (!el) return;
        const map = await createMap(el, {
            lat: postStore.form.latitude || 10.7769,
            lng: postStore.form.longitude || 106.7009,
            zoom: postStore.form.latitude ? 16 : 10,
        });
        if (!map) return;
        initInteractiveMap(map);
    } catch (e) {
        console.error('Lỗi khởi tạo bản đồ:', e);
    }
}

const initInteractiveMap = (map: any) => {
    if (interactiveMapInstance) {
        interactiveMapInstance.remove();
        interactiveMapInstance = null;
        interactiveMarkerInstance = null;
    }

    interactiveMapInstance = map;
    const L = getLeaflet();
    if (!L) return;

    const initLat = postStore.form.latitude || 10.7769;
    const initLng = postStore.form.longitude || 106.7009;

    interactiveMarkerInstance = L.marker([initLat, initLng], { draggable: true }).addTo(map);

    // Kéo marker → cập nhật lat/lng
    interactiveMarkerInstance.on('dragend', () => {
        const position = interactiveMarkerInstance.getLatLng();
        postStore.form.latitude = position.lat;
        postStore.form.longitude = position.lng;
    });

    // Click map → di chuyển marker + cập nhật lat/lng
    map.on('click', (e: any) => {
        interactiveMarkerInstance.setLatLng(e.latlng);
        postStore.form.latitude = e.latlng.lat;
        postStore.form.longitude = e.latlng.lng;
    });
};

// Đổi địa chỉ (tỉnh/phường/detail) → tự động geocode → cập nhật lat/lng + marker
const triggerGeocoding = async () => {
    if (!(postStore.form.province || postStore.form.ward || postStore.form.detail_address)) return;

    const city = (provinceOptions.value.find(item => item.value === postStore.form.province)?.label as string) || '';
    const ward = (wardOptions.value.find(item => item.value === postStore.form.ward)?.label as string) || '';
    const detail = postStore.form.detail_address || '';

    const queryParts: string[] = [];
    if (detail) queryParts.push(detail);
    if (ward) queryParts.push(ward);
    if (city) queryParts.push(city);
    queryParts.push('Việt Nam');

    try {
        loadingGeocode.value = true;
        const coord = await geocode(queryParts.join(', '));
        if (coord) {
            postStore.form.latitude = coord.lat;
            postStore.form.longitude = coord.lng;
            if (interactiveMapInstance && interactiveMarkerInstance) {
                interactiveMarkerInstance.setLatLng([coord.lat, coord.lng]);
                interactiveMapInstance.setView([coord.lat, coord.lng], 16);
            }
        } else {
            window.message?.warning('Không định vị được địa chỉ. Kiểm tra lại địa chỉ hoặc key Geoapify.');
        }
    } catch (err) {
        console.error('Lỗi khi định vị địa chỉ:', err);
        window.message?.error('Lỗi khi định vị địa chỉ.');
    } finally {
        loadingGeocode.value = false;
    }
};

// Áp dụng: validate + đóng modal (lat/lng đã nằm trong form)
const applyAddress = () => {
    const valid = postStore.validateAddress()
    if (!valid) return
    emit('update:show', false);
}

// ── Dữ liệu dự án ──
const projectOptions = ref<SelectOption[]>([])

const clearError = (field: keyof typeof postStore.errorsAddress) => {
    postStore.errorsAddress[field] = ''
}

const locationLabel = computed(() => {
    const city = provinceOptions.value.find(item => item.value === postStore.form.province)
    const ward = wardOptions.value.find(item => item.value === postStore.form.ward)

    const parts: string[] = [];
    // Tìm được label theo code → dùng label; ngược lại (edit: form đang giữ NAME,
    // hoặc options chưa load) → dùng chính giá trị form
    if (city) parts.push(city.label as string)
    else if (postStore.form.province) parts.push(postStore.form.province)
    if (ward) parts.push(ward.label as string)
    else if (postStore.form.ward) parts.push(postStore.form.ward)
    if (postStore.form.detail_address) parts.push(postStore.form.detail_address);
    if (parts.length === 0) return 'Chọn khu vực';
    return parts.join(', ');
})

watch(() => locationLabel.value, (value) => {
    emit('update:locationLabel', value)
}, { immediate: true })

const onProvinceChange = (provinceCode: string | null) => {
    
    triggerGeocoding()
}

const onWardSelectChange = (wardCode: string | null) => {
    clearError('ward');
    // Tự động định vị khi chọn phường/xã
    triggerGeocoding();
}

const onDetailAddressInput = (val: string) => {
    clearError('detail_address');
    // Tự động định vị khi gõ địa chỉ (debounce 1.5s)
    if (geocodeTimeout) clearTimeout(geocodeTimeout);
    geocodeTimeout = setTimeout(() => {
        triggerGeocoding();
    }, 1500);
}

// Timer debounce cho việc gõ địa chỉ chi tiết
let geocodeTimeout: any = null;

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

// Đổi tỉnh/phường → reset dự án + tải lại danh sách dự án theo vị trí
watch(
    () => [postStore.form.province, postStore.form.ward],
    async ([newProvince, newWard], [oldProvince, oldWard]) => {
        if (newProvince !== oldProvince || newWard !== oldWard) {
            postStore.form.project_id = null
        }
        await fetchListProject()
    }
)
</script>

<style scoped>
#interactive-address-map {
    outline: none;
}
</style>
