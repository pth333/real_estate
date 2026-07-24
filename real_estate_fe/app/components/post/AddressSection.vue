<template>
    <div class="border border-gray-200 rounded-xl overflow-hidden">
        <div class="flex justify-between items-center px-4 py-3.5 cursor-pointer select-none hover:bg-gray-50">
            <span class="text-sm font-semibold text-gray-900">
                Địa chỉ <span class="text-red-500">*</span>
            </span>
            <n-icon class="text-gray-500 rotate-180">
                <ChevronDownOutline />
            </n-icon>
        </div>

        <div class="px-4 pt-3 pb-4 border-t border-gray-100 flex flex-col">
            <!-- Địa chỉ cũ -->
            <div class="flex items-start gap-2.5 py-1">
                <n-icon class="mt-0.5 shrink-0 text-gray-400">
                    <LocationOutline />
                </n-icon>
                <div>
                    <div class="text-xs text-gray-400 mb-0.5">Địa chỉ cũ</div>
                    <div class="text-sm text-gray-400 line-through">Thăng Long Capital, Xã An Khánh, Huyện Hoài
                        Đức, Hà
                        Nội</div>
                </div>
            </div>

            <n-divider class="!my-3" />

            <!-- Địa chỉ mới -->
            <div class="flex items-start gap-2.5 py-1">
                <n-icon class="mt-0.5 shrink-0 text-red-500">
                    <LocationOutline />
                </n-icon>
                <div class="flex-1">
                    <div class="text-xs text-gray-400 mb-0.5">Địa chỉ mới</div>
                    <div class="text-sm text-gray-700 leading-relaxed">
                        Thăng Long Capital, <strong>Đường 10A</strong>, <strong>Xã An Khánh</strong>, Hà Nội
                    </div>
                </div>
                <button
                    class="shrink-0 ml-auto p-1.5 rounded-full text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors">
                    <n-icon size="18">
                        <CreateOutline />
                    </n-icon>
                </button>
            </div>
        </div>
    </div>
</template>