<template>
    <n-card title="Địa chỉ" size="small" :segmented="{ content: true }">
        <template #header-extra>
            <n-icon class="text-gray-500 rotate-180">
                <IconChevronDownOutline />
            </n-icon>
        </template>

        <div class="flex flex-col">
            <div class="flex items-start gap-2.5 py-1">
                <n-icon class="mt-0.5 shrink-0 text-gray-400">
                    <IconLocationOutline />
                </n-icon>
                <div>
                    <div class="text-xs text-gray-400 mb-0.5">Địa chỉ cũ</div>
                    <div class="text-sm text-gray-400 line-through">Thăng Long Capital, Xã An Khánh, Huyện Hoài
                        Đức, Hà
                        Nội</div>
                </div>
            </div>

            <n-divider class="!my-3" />

            <div class="flex items-start gap-2.5 py-1">
                <n-icon class="mt-0.5 shrink-0 text-red-500">
                    <IconLocationOutline />
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
                        <IconCreateOutline />
                    </n-icon>
                </button>
            </div>
        </div>
    </n-card>
</template>
