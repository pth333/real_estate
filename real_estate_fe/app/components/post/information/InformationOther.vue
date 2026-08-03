<template>
    <n-card title="Thông tin khác" size="small" :segmented="{ content: true }">
        <template #header-extra>
            <div class="cursor-pointer" @click="collapsed = !collapsed">
                <n-icon class="text-gray-500 transition-transform duration-200" :class="{ 'rotate-180': !collapsed }">
                    <IconChevronDownOutline />
                </n-icon>
            </div>
        </template>

        <div v-show="!collapsed">
            <!-- Giấy tờ pháp lý + Nội thất -->
            <div class="flex gap-2 mb-4">
                <div class="flex flex-col flex-1">
                    <label class="text-sm text-gray-700 mb-1.5">Giấy tờ pháp lý</label>
                    <n-select v-model:value="postStore.form.legal_docs" placeholder="Chọn giấy tờ pháp lý"
                        :options="legalDocOptions" />
                </div>
                <div class="flex flex-col flex-1">
                    <label class="text-sm text-gray-700 mb-1.5">Nội thất</label>
                    <n-select v-model:value="postStore.form.interior" placeholder="Chọn nội thất"
                        :options="interiorOptions" />
                </div>
            </div>

            <!-- Số phòng tắm + Số phòng ngủ -->
            <div class="flex gap-2 mb-4">
                <div class="flex flex-col flex-1">
                    <label class="text-sm text-gray-700 mb-1.5">Số phòng tắm</label>
                    <n-input-number v-model:value="postStore.form.bathroom_count" placeholder="Số phòng tắm" :min="0"
                        class="w-full" />
                </div>
                <div class="flex flex-col flex-1">
                    <label class="text-sm text-gray-700 mb-1.5">Số phòng ngủ</label>
                    <n-input-number v-model:value="postStore.form.bedroom_count" placeholder="Số phòng ngủ" :min="0"
                        class="w-full" />
                </div>
            </div>

            <!-- Hướng nhà + Hướng ban công -->
            <div class="flex gap-2 mb-4">
                <div class="flex flex-col flex-1">
                    <label class="text-sm text-gray-700 mb-1.5">Hướng nhà</label>
                    <n-select v-model:value="postStore.form.house_direction" placeholder="Chọn hướng nhà"
                        :options="directionOptions" />
                </div>
                <div class="flex flex-col flex-1">
                    <label class="text-sm text-gray-700 mb-1.5">Hướng ban công</label>
                    <n-select v-model:value="postStore.form.balcony_direction" placeholder="Chọn hướng ban công"
                        :options="directionOptions" />
                </div>
            </div>

            <!-- Thời gian vào ở -->
            <div class="mb-4">
                <label class="text-sm text-gray-700 block mb-1.5">Thời gian vào ở</label>
                <n-select v-model:value="postStore.form.move_in_time" placeholder="Chọn thời gian vào ở"
                    :options="moveInTimeOptions" />
            </div>

            <!-- Mức giá điện / nước / internet -->
            <div class="flex gap-2 mb-4">
                <div class="flex flex-col flex-1">
                    <label class="text-sm text-gray-700 mb-1.5">Mức giá điện</label>
                    <NumberFormatInput v-model="postStore.form.price_electricity" placeholder="VD: 3.500"
                        suffix="đ/kWh" />
                </div>
                <div class="flex flex-col flex-1">
                    <label class="text-sm text-gray-700 mb-1.5">Mức giá nước</label>
                    <NumberFormatInput v-model="postStore.form.price_water" placeholder="VD: 25.000"
                        suffix="đ/m³" />
                </div>
                <div class="flex flex-col flex-1">
                    <label class="text-sm text-gray-700 mb-1.5">Mức giá internet</label>
                    <NumberFormatInput v-model="postStore.form.price_internet" placeholder="VD: 200.000"
                        suffix="đ/tháng" />
                </div>
            </div>

            <!-- Tiện ích -->
            <div>
                <label class="text-sm text-gray-700 block mb-1.5">Tiện ích</label>
                <n-checkbox-group v-model:value="postStore.form.amenities">
                    <div class="flex gap-4">
                        <n-checkbox value="camera">Camera</n-checkbox>
                        <n-checkbox value="bao_ve">Bảo vệ</n-checkbox>
                        <n-checkbox value="pccc">PCCC</n-checkbox>
                    </div>
                </n-checkbox-group>
            </div>
        </div>
    </n-card>
</template>

<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { useCreatePost } from '~/stores/create-post'

const postStore = useCreatePost()
const collapsed = ref(false)

// Các ô nhập giá điện/nước/internet dùng NumberFormatInput (auto-import từ common/):
// gõ là format dấu chấm ngay, emit số thực về store qua v-model.

// Giấy tờ pháp lý
const legalDocOptions = ref<SelectOption[]>([
    { label: 'Sổ đỏ/ Sổ hồng', value: 'so_do' },
    { label: 'Hợp đồng mua bán', value: 'hop_dong_mua_ban' },
    { label: 'Đang chờ sổ', value: 'dang_cho_so' },
])

// Nội thất
const interiorOptions = ref<SelectOption[]>([
    { label: 'Đầy đủ nội thất', value: 'day_du' },
    { label: 'Cơ bản', value: 'co_ban' },
    { label: 'Chưa có nội thất', value: 'chua_co' },
])

// Hướng nhà / ban công
const directionOptions = ref<SelectOption[]>([
    { label: 'Đông', value: 'dong' },
    { label: 'Tây', value: 'tay' },
    { label: 'Nam', value: 'nam' },
    { label: 'Bắc', value: 'bac' },
    { label: 'Đông Bắc', value: 'dong_bac' },
    { label: 'Tây Bắc', value: 'tay_bac' },
    { label: 'Đông Nam', value: 'dong_nam' },
    { label: 'Tây Nam', value: 'tay_nam' },
])

// Thời gian vào ở
const moveInTimeOptions = ref<SelectOption[]>([
    { label: 'Ngay', value: 'ngay' },
    { label: 'Trong 1 tháng', value: '1_thang' },
    { label: 'Trong 3 tháng', value: '3_thang' },
    { label: 'Trong 6 tháng', value: '6_thang' },
    { label: 'Thỏa thuận', value: 'thoa_thuan' },
])
</script>
