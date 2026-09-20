<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
    <div class="w-full max-w-[560px] bg-white rounded-xl border border-gray-200 p-8 flex flex-col gap-4">
      <n-result status="403" title="Không có quyền truy cập"
        description="Tài khoản của bạn không được phép vào trang này. Nếu bạn cho rằng đây là nhầm lẫn, hãy liên hệ quản trị viên.">
        <template #footer>
          <n-space justify="center">
            <n-button @click="goHome">Về trang chủ</n-button>
            <n-button type="primary" @click="goBack">Quay lại</n-button>
          </n-space>
        </template>
      </n-result>
    </div>
  </div>
</template>

<script setup lang="ts">
// Trang đích khi middleware route chặn vì thiếu role/permission.
// Tên file đặt tiếng Anh; giữ alias tiếng Việt để URL cũ không chết.
definePageMeta({
  layout: 'empty',
  alias: ['/khong-co-quyen'],
})

const router = useRouter()

function goHome() {
  navigateTo('/')
}

function goBack() {
  // Không có lịch sử (mở trực tiếp URL) thì về trang chủ
  if (window.history.length > 1) {
    router.back()
    return
  }
  navigateTo('/')
}
</script>
