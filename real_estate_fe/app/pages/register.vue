<template>
  <AuthShell title="Tạo tài khoản" subtitle="Đăng ký miễn phí để đặt cọc giữ lịch xem nhà và lưu tin yêu thích.">
    <n-form :model="form" label-placement="top" :show-require-mark="false" @submit.prevent="handleRegister">
      <n-form-item label="Họ và tên" :feedback="errors.name" :validation-status="errors.name ? 'error' : undefined">
        <n-input v-model:value="form.name" type="text" placeholder="Nguyễn Văn A" clearable autocomplete="name" />
      </n-form-item>

      <n-form-item label="Email" :feedback="errors.email" :validation-status="errors.email ? 'error' : undefined">
        <n-input v-model:value="form.email" type="text" placeholder="your@email.com" clearable autocomplete="email" />
      </n-form-item>

      <n-form-item label="Mật khẩu" :feedback="errors.password"
        :validation-status="errors.password ? 'error' : undefined">
        <n-input v-model:value="form.password" type="password" placeholder="Ít nhất 6 ký tự" show-password-on="click"
          autocomplete="new-password" />
      </n-form-item>

      <n-form-item label="Xác nhận mật khẩu" :feedback="errors.confirmPassword"
        :validation-status="errors.confirmPassword ? 'error' : undefined">
        <n-input v-model:value="form.confirmPassword" type="password" placeholder="Nhập lại mật khẩu"
          show-password-on="click" autocomplete="new-password" />
      </n-form-item>

      <n-button type="primary" attr-type="submit" size="large" block :loading="loading" :disabled="loading"
        class="mt-2">
        Đăng ký
      </n-button>
    </n-form>

    <n-alert v-if="apiError" type="error" :title="apiError" closable class="mt-4" />
    <n-alert v-if="successMsg" type="success" :title="successMsg" class="mt-4" />

    <template #footer>
      <p class="text-sm text-gray-500">
        Đã có tài khoản?
        <NuxtLink to="/dang-nhap" class="font-semibold text-emerald-600 hover:underline">
          Đăng nhập
        </NuxtLink>
      </p>
    </template>
  </AuthShell>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { brandTitle } from '~/utils/brand'
import AuthShell from '~/components/auth/AuthShell.vue'

definePageMeta({
  layout: false,
  alias: '/dang-ky',
})

useHead({
  title: brandTitle('Đăng ký'),
})

const auth = useAuthStore()
const form = ref({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
})
const errors = ref({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
})
const loading = ref(false)
const apiError = ref('')
const successMsg = ref('')

function validate(): boolean {
  let ok = true
  errors.value.name = ''
  errors.value.email = ''
  errors.value.password = ''
  errors.value.confirmPassword = ''

  if (!form.value.name.trim()) {
    errors.value.name = 'Họ tên không được để trống'
    ok = false
  }

  const email = form.value.email.trim()
  if (!email) {
    errors.value.email = 'Email không được để trống'
    ok = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    errors.value.email = 'Email không hợp lệ'
    ok = false
  }

  if (!form.value.password) {
    errors.value.password = 'Mật khẩu không được để trống'
    ok = false
  } else if (form.value.password.length < 6) {
    errors.value.password = 'Mật khẩu phải có ít nhất 6 ký tự'
    ok = false
  }

  if (form.value.password !== form.value.confirmPassword) {
    errors.value.confirmPassword = 'Mật khẩu xác nhận không khớp'
    ok = false
  }

  return ok
}

async function handleRegister() {
  if (loading.value) return
  apiError.value = ''
  successMsg.value = ''
  if (!validate()) return

  loading.value = true
  try {
    await auth.register({
      name: form.value.name.trim(),
      email: form.value.email.trim(),
      password: form.value.password,
    })
    successMsg.value = 'Đăng ký thành công! Đang chuyển sang trang đăng nhập…'
    // Tài khoản mới cần đăng nhập lại nên đưa về trang đăng nhập kèm email vừa dùng
    setTimeout(() => {
      navigateTo({ path: '/dang-nhap', query: { email: form.value.email.trim() } })
    }, 1200)
  } catch (err: unknown) {
    apiError.value = err instanceof Error ? err.message : 'Đăng ký thất bại'
  } finally {
    loading.value = false
  }
}
</script>
