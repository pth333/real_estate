<template>
    <section class="py-8">
        <div class="container mx-auto px-4 md:px-12 lg:px-24">
            <h2 class="mb-5 text-lg font-bold text-gray-900 md:text-xl">Bất động sản theo địa điểm</h2>

            <SkeletonCard v-if="loading" type="area" />

            <!-- Mobile: địa điểm nổi bật trên đầu, lưới 2x2 bên dưới.
                 Từ tablet trở lên: ảnh lớn bên trái 45%, lưới 2x2 bên phải. -->
            <div v-else class="flex flex-col gap-3 md:h-[360px] md:flex-row">
                <div class="group relative h-52 shrink-0 cursor-pointer overflow-hidden rounded-lg md:h-full md:w-[45%]"
                    @click="goToCity(featured)">
                    <img :src="featured?.image" :alt="featured?.name"
                        class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105" />
                    <div class="absolute inset-0 bg-linear-to-t from-black/60 via-black/20 to-transparent" />
                    <div class="absolute bottom-0 left-0 p-4 text-white">
                        <p class="text-lg font-bold leading-tight"> {{ featured?.name }} </p>
                        <p class="mt-0.5 text-sm text-white/80"> {{ featured?.count }} tin đăng </p>
                    </div>
                </div>

                <div class="grid flex-1 grid-cols-2 grid-rows-2 gap-3">
                    <div v-for="location in restLocations" :key="location.id"
                        class="relative h-32 cursor-pointer overflow-hidden rounded-lg group md:h-auto">
                        <div @click="goToCity(location)" class="h-full w-full">
                            <img :src="location.image" :alt="location.name"
                                class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105" />
                            <div class="absolute inset-0 bg-linear-to-t from-black/60 via-black/20 to-transparent" />
                            <div class="absolute bottom-0 left-0 p-3 text-white">
                                <p class="text-sm font-bold leading-tight"> {{ location.name }} </p>
                                <p class="mt-0.5 text-xs text-white/80"> {{ location.count }} tin đăng </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>

<script setup lang="ts">
import type { TopCityOption } from '~/types/real_estate';
import { useRealEstateService } from '~/services/real-estate.service';

function goToCity(city?: TopCityOption | null) {
    if (!city) return;
    navigateTo(`/${city.category_slug}-${city.city_slug}`);
}

// Các địa điểm còn lại (lưới 2x2 bên phải)
const restLocations = computed(() => {
    if (listTopCity.value) {
        const feat = listTopCity.value.slice(1)
        return feat
    }
})
const featured = computed(() => {
    return listTopCity.value[0]
})


const listTopCity = ref<TopCityOption[]>([])
const loading = ref(true)

const realEstateService = useRealEstateService()
const fetchListTopCity = async () => {
    loading.value = true
    try {
        listTopCity.value = await realEstateService.getTopCities()
    } catch (e) {
        console.log(e)
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    fetchListTopCity()
})
</script>