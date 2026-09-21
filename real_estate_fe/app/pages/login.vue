<template>
  <AuthShell title="Đăng nhập" subtitle="Chào mừng bạn quay lại nền tảng bất động sản NhàViệt.">
    <n-form :model="form" label-placement="top" :show-require-mark="false" @submit.prevent="handleLogin">
      <n-form-item label="Email" :feedback="errors.email" :validation-status="errors.email ? 'error' : undefined">
        <n-input v-model:value="form.email" type="text" placeholder="your@email.com" clearable autocomplete="email" />
      </n-form-item>

      <n-form-item label="Mật khẩu" :feedback="errors.password"
        :validation-status="errors.password ? 'error' : undefined">
        <n-input v-model:value="form.password" type="password" placeholder="Nhập mật khẩu" show-password-on="click"
          autocomplete="current-password" />
      </n-form-item>

      <n-button type="primary" attr-type="submit" size="large" block :loading="loading" :disabled="loading"
        class="mt-2">
        Đăng nhập
      </n-button>
    </n-form>

    <n-alert v-if="apiError" type="error" :title="apiError" closable class="mt-4" />

    <template #footer>
      <p class="text-sm text-gray-500">
        Chưa có tài khoản?
        <NuxtLink to="/dang-ky" class="font-semibold text-emerald-600 hover:underline">
          Đăng ký ngay
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
  alias: '/dang-nhap',
})

useHead({
  title: brandTitle('Đăng nhập'),
})

const auth = useAuthStore()
const route = useRoute()

/** Trang đăng ký chuyển sang kèm ?email=... để khỏi phải gõ lại */
const presetEmail = typeof route.query.email === 'string' ? route.query.email : ''

const form = ref({ email: presetEmail, password: '' })
const errors = ref({ email: '', password: '' })
const loading = ref(false)
const apiError = ref('')

/**
 * Lấy đường dẫn cần quay lại sau khi đăng nhập (do middleware route truyền vào
 * khi chặn một trang cần đăng nhập). Chỉ nhận đường dẫn nội bộ để tránh open redirect.
 */
function resolveRedirect(): string {
  const value = route.query.redirect
  if (typeof value !== 'string') return '/'
  if (!value.startsWith('/') || value.startsWith('//')) return '/'
  return value
}

function validate(): boolean {
  let ok = true
  errors.value.email = ''
  errors.value.password = ''

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
  }

  return ok
}

async function handleLogin() {
  if (loading.value) return
  apiError.value = ''
  if (!validate()) return

  loading.value = true
  try {
    await auth.login({ email: form.value.email.trim(), password: form.value.password })
    // Quay lại đúng trang người dùng định vào trước khi bị chặn đăng nhập
    await navigateTo(resolveRedirect())
  } catch (err: unknown) {
    apiError.value = err instanceof Error ? err.message : 'Đăng nhập thất bại'
  } finally {
    loading.value = false
  }
}
</script>
