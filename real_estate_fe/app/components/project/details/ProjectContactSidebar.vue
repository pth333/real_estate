<template>
  <div class="sticky top-6">
    <n-card class="border border-gray-100 shadow-sm rounded-xl">
      <n-space vertical :size="16">
        <!-- Avatar và thông tin liên hệ -->
        <n-space align="center" :size="12">
          <n-avatar
            round
            :size="48"
            class="bg-emerald-50 text-emerald-600 font-bold border border-emerald-100"
          >
            S
          </n-avatar>
          <n-space vertical :size="2">
            <n-text class="font-bold text-gray-900 text-sm">Phòng Kinh Doanh</n-text>
            <n-text depth="3" class="text-xs">Hỗ trợ tư vấn dự án 24/7</n-text>
          </n-space>
        </n-space>

        <!-- Các nút hành động Naive UI chuyên nghiệp -->
        <n-space vertical :size="10">
          <n-button
            type="primary"
            size="large"
            class="w-full bg-emerald-600! hover:bg-emerald-700! font-semibold! flex items-center justify-center gap-2"
          >
            <IconPhone class="h-4 w-4 shrink-0" />
            Yêu cầu gọi lại tư vấn
          </n-button>
          <n-button
            ghost
            size="large"
            class="w-full text-gray-700 border-gray-300 font-medium"
          >
            Tải bảng giá & pháp lý
          </n-button>
        </n-space>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
</script>
