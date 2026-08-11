<template>
    <section class="py-8">
        <div class="container mx-auto px-24">
        <!-- Header -->
        <div class="flex items-center justify-between mb-5">
            <h2 class="text-xl font-bold text-gray-900">Dự án bất động sản nổi bật</h2>
            <a href="#" class="flex items-center gap-1 text-emerald-600 text-sm font-medium hover:underline">
                Xem thêm
                <IconArrowRight class="h-4 w-4" />
            </a>
        </div>

        <!-- Slider wrapper -->
        <div class="relative">
            <!-- Nút prev -->
            <button
                class="absolute -left-5 top-1/2 -translate-y-1/2 z-10 w-9 h-9 bg-white border border-gray-200 shadow flex items-center justify-center hover:shadow-md transition-shadow"
                @click="prev">
                <IconChevronLeft class="h-4 w-4 text-gray-600" />
            </button>

            <!-- Cards -->
            <div class="grid grid-cols-4 gap-4 overflow-hidden">
                <div v-for="item in visibleItems" :key="item.id"
                    class="bg-white overflow-hidden cursor-pointer group">
                    <!-- Ảnh -->
                    <div class="relative aspect-[4/3] overflow-hidden bg-gray-100">
                        <img :src="item.thumbnail" :alt="item.name"
                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                        <!-- Số ảnh -->
                        <div
                            class="absolute bottom-2 right-2 flex items-center gap-1 bg-black/50 text-white text-xs px-1.5 py-0.5">
                            <IconImage class="h-3 w-3" />
                            {{ item.imageCount }}
                        </div>
                    </div>

                    <!-- Nội dung -->
                    <div class="pt-2.5 flex flex-col gap-1">
                        <!-- Badge trạng thái -->
                        <div class="flex items-center gap-2 flex-wrap">
                            <span class="text-xs font-medium px-2 py-0.5 border"
                                :class="statusClass(item.status)">
                                {{ item.status }}
                            </span>
                            <span v-if="item.openDate" class="text-xs text-gray-400">
                                · Ngày {{ item.openDate }}
                            </span>
                        </div>

                        <!-- Tên dự án -->
                        <div
                            class="font-semibold text-gray-900 text-sm truncate group-hover:text-emerald-600 transition-colors">
                            {{ item.name }}
                        </div>

                        <!-- Diện tích -->
                        <div class="text-sm text-gray-500">{{ item.area }}</div>

                        <!-- Địa chỉ -->
                        <div class="text-sm text-gray-500">{{ item.location }}</div>
                    </div>
                </div>
            </div>

            <!-- Nút next -->
            <button
                class="absolute -right-5 top-1/2 -translate-y-1/2 z-10 w-9 h-9 bg-white border border-gray-200 shadow flex items-center justify-center hover:shadow-md transition-shadow"
                @click="next">
                <IconChevronRight class="h-4 w-4 text-gray-600" />
            </button>
        </div>
        </div>
    </section>
</template>

<script setup lang="ts">
const projects = ref([
    {
        id: 1,
        name: 'Gladia Heights',
        status: 'Đang mở bán',
        area: '1,29 ha',
        location: 'Quận 2, Hồ Chí Minh',
        imageCount: 4,
        openDate: null,
        thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=Project',
    },
    {
        id: 2,
        name: 'Legacy89 Signature',
        status: 'Sắp mở bán',
        area: 'Mỹ Hào, Hưng Yên',
        location: 'Mỹ Hào, Hưng Yên',
        imageCount: 12,
        openDate: null,
        thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=Project',
    },
    {
        id: 3,
        name: 'Park City Phú Quốc',
        status: 'Đang mở bán',
        area: '4,41 ha',
        location: 'Phú Quốc, Kiên Giang',
        imageCount: 9,
        openDate: null,
        thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=Project',
    },
    {
        id: 4,
        name: 'Vinhomes Global Gate Hạ...',
        status: 'Đang mở bán',
        area: '4.109,64 ha',
        location: 'Quảng Yên, Quảng Ninh',
        imageCount: 9,
        openDate: '19/12/2025:...',
        thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=Project',
    },
    {
        id: 5,
        name: 'The Opus One Residences',
        status: 'Đang mở bán',
        area: '2,8 ha',
        location: 'Quận 9, Hồ Chí Minh',
        imageCount: 7,
        openDate: null,
        thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=Project',
    },
    {
        id: 6,
        name: 'Eaton Park',
        status: 'Sắp mở bán',
        area: '3,5 ha',
        location: 'Bình Thạnh, Hồ Chí Minh',
        imageCount: 5,
        openDate: null,
        thumbnail: 'https://placehold.co/400x300/e2e8f0/94a3b8?text=Project',
    },
])

const currentIndex = ref(0)
const pageSize = 4

const visibleItems = computed(() =>
    projects.value.slice(currentIndex.value, currentIndex.value + pageSize)
)

function prev() {
    if (currentIndex.value > 0) currentIndex.value -= pageSize
}

function next() {
    if (currentIndex.value + pageSize < projects.value.length) currentIndex.value += pageSize
}

function statusClass(status: string) {
    if (status === 'Đang mở bán') return 'border-green-400 text-green-600 bg-green-50'
    if (status === 'Sắp mở bán') return 'border-red-300 text-red-500 bg-red-50'
    return 'border-gray-300 text-gray-500 bg-gray-50'
}
</script>