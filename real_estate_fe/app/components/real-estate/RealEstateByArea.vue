<template>
    <section class="py-8">
        <div class="container mx-auto px-24">
        <h1 class="text-lg font-bold text-gray-600 mb-4">Bất động sản theo địa điểm</h1>

        <div class="flex gap-3" style="height: 360px;">
            <div class="relative overflow-hidden cursor-pointer group shrink-0" style="flex: 0 0 45%;"
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
                    class="relative  overflow-hidden cursor-pointer group">
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
interface ListTopCity {
    id: number;
    name: string;
    count: number;
    image: string;
    category_slug: string;
    city_slug: string;
}

function goToCity(city?: ListTopCity | null) {
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

const listTopCity = ref<ListTopCity[]>([])

const { $api } = useNuxtApp()
const fetchListTopCity = async () => {
    try {
        const result = await $api.get<{ data: ListTopCity[] }>('/real-estate/list/top-city')
        listTopCity.value = result.data
    } catch (e) {
        console.log(e)
    }
}

onMounted(() => {
    fetchListTopCity()
})
</script>