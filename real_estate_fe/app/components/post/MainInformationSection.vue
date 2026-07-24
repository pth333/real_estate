<template>
    <div class="border border-gray-200 rounded-xl overflow-hidden">
        <div class="flex justify-between items-center px-4 py-3.5 cursor-pointer select-none hover:bg-gray-50">
            <span class="text-sm font-semibold text-gray-900">
                Thông tin chính <span class="text-red-500">*</span>
            </span>
            <n-icon class="text-gray-500">
                <ChevronDownOutline />
            </n-icon>
        </div>

        <div class="px-4 pt-3 pb-4 border-t border-gray-100">

            <!-- Loại BĐS -->
            <div class="mb-4">
                <label class="text-sm text-gray-700 block mb-1.5">
                    Loại bất động sản <span class="text-red-500">*</span>
                </label>
                <n-select placeholder="Chọn loại BĐS" :options="[]" />
            </div>

            <!-- Diện tích -->
            <div>
                <label class="text-sm text-gray-700 block mb-1.5">
                    Diện tích <span class="text-red-500">*</span>
                </label>
                <n-input-number placeholder="Nhập diện tích" :min="0" class="w-full" status="error">
                    <template #suffix>m²</template>
                </n-input-number>
                <span class="text-xs text-red-500 mt-1 block">Thông tin này bắt buộc phải điền</span>
            </div>

        </div>
    </div>
</template>