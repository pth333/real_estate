<template>
    <div class="min-h-screen bg-white flex flex-col font-sans">

        <!-- Header -->
        <div class="flex justify-between items-center px-6 py-3.5 border-b border-gray-200">
            <span class="text-lg font-semibold text-gray-900">Tạo tin đăng</span>
            <div class="flex items-center gap-2.5">
                <n-button text>
                    <template #icon><n-icon>
                            <EyeOutline />
                        </n-icon></template>
                    Xem trước
                </n-button>
                <n-button>Thoát</n-button>
            </div>
        </div>

        <!-- Progress -->
        <div class="px-6 pb-3 border-b border-gray-200">
            <span class="text-xs text-red-500 font-medium block mb-1.5 mt-1.5">Bước 1. Thông tin BĐS</span>
            <n-progress type="line" :percentage="33" :show-indicator="false" color="#ef4444" rail-color="#f3f4f6" />
        </div>

        <!-- Form body -->
        <div class="flex-1 px-6 py-6 w-full max-w-2xl mx-auto flex flex-col gap-4">

            <!-- Section: Nhu cầu -->P
            <div class="border border-gray-200 rounded-xl overflow-hidden">
                <div class="flex justify-between items-center px-4 py-3.5 cursor-pointer select-none hover:bg-gray-50">
                    <span class="text-sm font-semibold text-gray-900">
                        Nhu cầu <span class="text-red-500">*</span>
                    </span>
                    <n-icon class="text-gray-500">
                        <IconChevronDown />
                    </n-icon>
                </div>

                <div class="px-4 pt-3 pb-4 border-t border-gray-100">
                    <div class="grid grid-cols-2 gap-3">
                        <!-- Bán: active -->
                        <div
                            class="flex items-center gap-2.5 px-5 py-3.5 rounded-lg cursor-pointer text-sm border border-gray-900 text-gray-900 font-semibold">
                            <n-icon size="20">
                                <PricetagOutline />
                            </n-icon>
                            <span>Bán</span>
                        </div>
                        <!-- Cho thuê -->
                        <div
                            class="flex items-center gap-2.5 px-5 py-3.5 rounded-lg cursor-pointer text-sm border border-gray-300 text-gray-600 hover:border-red-500 hover:text-red-500 transition-colors">
                            <n-icon size="20">
                                <KeyOutline />
                            </n-icon>
                            <span>Cho thuê</span>
                        </div>
                    </div>
                </div>
            </div>
            <!-- Section: Địa chỉ -->
            <AddressSection />
            <!-- Section: Thông tin chính -->
            <MainInformationSection />

        </div>

        <!-- Footer -->
        <div class="flex justify-end px-6 py-4 border-t border-gray-200">
            <n-button size="large" class="!bg-red-500 !border-red-500 !text-white !rounded-full !px-8 !font-semibold">
                Tiếp tục
            </n-button>
        </div>

        <ModalOTPAuthentication />

    </div>
</template>
<script setup lang="ts">

</script>