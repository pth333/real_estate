<template>
    <section class="py-8">
        <div class="container mx-auto px-24">
            <h2 class="text-xl font-bold text-gray-900 mb-5">Bất động sản theo địa điểm</h2>

            <SkeletonCard v-if="loading" type="area" />

            <div v-else class="flex gap-3" style="height: 360px;">
                <div class="relative overflow-hidden rounded-lg cursor-pointer group shrink-0" style="flex: 0 0 45%;"
                    @click="goToCity(featured)">
                    <img :src="featured?.image" :alt="featured?.name"
                        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105" />
                    <div class="absolute inset-0 bg-linear-to-t from-black/60 via-black/20 to-transparent" />
                    <div class="absolute bottom-0 left-0 p-4 text-white">
                        <p class="font-bold text-lg leading-tight"> {{ featured?.name }} </p>
                        <p class="text-sm text-white/80 mt-0.5"> {{ featured?.count }} tin đăng </p>
                    </div>
                </div>

                <div class="flex-1 grid grid-cols-2 grid-rows-2 gap-3">
                    <div v-for="location in restLocations" :key="location.id"
                        class="relative overflow-hidden rounded-lg cursor-pointer group">
                        <div @click="goToCity(location)" class="w-full h-full">
                            <img :src="location.image" :alt="location.name"
                                class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105" />
                            <div class="absolute inset-0 bg-linear-to-t from-black/60 via-black/20 to-transparent" />
                            <div class="absolute bottom-0 left-0 p-3 text-white">
                                <p class="font-bold text-sm leading-tight"> {{ location.name }} </p>
                                <p class="text-xs text-white/80 mt-0.5"> {{ location.count }} tin đăng </p>
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