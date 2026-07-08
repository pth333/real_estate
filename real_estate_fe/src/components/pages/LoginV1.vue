<template>
  <div class="flex min-h-screen bg-limestone">
    <!-- Form bên trái - 40% -->
    <div class="flex w-full items-center justify-center bg-foundation px-8 py-12 lg:w-2/5">
      <div class="w-full max-w-sm">
        <h1 class="mb-2 font-semibold tracking-tight text-limestone text-3xl">Đăng nhập</h1>
        <p class="mb-8 text-sm text-patina/80">Truy cập vào nền tảng bất động sản</p>

        <form @submit.prevent="handleLogin" class="space-y-6">
          <!-- Email -->
          <div>
            <label class="mb-2 block text-sm font-medium text-limestone/90">Email</label>
            <div class="plot-input-wrapper">
              <input
                v-model="form.email"
                type="email"
                placeholder="your@email.com"
                class="w-full border border-limestone/20 bg-foundation px-4 py-3 text-sm text-limestone placeholder:text-limestone/40 focus:border-patina focus:outline-none"
                :class="{ '!border-red-400': errors.email }"
              />
            </div>
            <p v-if="errors.email" class="mt-1.5 text-xs text-red-400">{{ errors.email }}</p>
          </div>

          <!-- Password -->
          <div>
            <label class="mb-2 block text-sm font-medium text-limestone/90">Mật khẩu</label>
            <div class="plot-input-wrapper">
              <input
                v-model="form.password"
                type="password"
                placeholder="••••••••"
                class="w-full border border-limestone/20 bg-foundation px-4 py-3 text-sm text-limestone placeholder:text-limestone/40 focus:border-patina focus:outline-none"
                :class="{ '!border-red-400': errors.password }"
              />
            </div>
            <p v-if="errors.password" class="mt-1.5 text-xs text-red-400">{{ errors.password }}</p>
          </div>

          <!-- Submit -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full bg-patina px-4 py-3 text-sm font-semibold tracking-tight text-foundation transition hover:bg-patina/90 disabled:cursor-not-allowed disabled:opacity-50"
          >
            <span v-if="loading">Đang xử lý…</span>
            <span v-else>Đăng nhập</span>
          </button>
        </form>

        <!-- Thông báo lỗi -->
        <p v-if="apiError" class="mt-4 text-sm text-red-400">{{ apiError }}</p>

        <!-- Link sang đăng ký -->
        <p class="mt-8 text-sm text-limestone/60">
          Chưa có tài khoản?
          <router-link :to="{ name: 'Register' }" class="font-medium text-oak hover:underline">
            Đăng ký
          </router-link>
        </p>
      </div>
    </div>

    <!-- Pattern bên phải - 60% -->
    <div class="hidden lg:block lg:w-3/5">
      <div class="relative h-full w-full">
        <!-- Topographic pattern -->
        <svg class="h-full w-full" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <pattern id="topo" x="0" y="0" width="120" height="120" patternUnits="userSpaceOnUse">
              <path d="M20 60c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5" opacity="0.15"/>
              <path d="M20 80c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5" opacity="0.1"/>
              <path d="M20 40c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5" opacity="0.1"/>
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#topo)"/>
        </svg>
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="text-center">
            <div class="mb-4 text-6xl font-bold tracking-tighter text-foundation/10">BĐS</div>
            <p class="text-sm font-medium tracking-wide text-foundation/30">NỀN TẢNG QUẢN LÝ</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.plot-input-wrapper {
  position: relative;
}

.plot-input-wrapper::after {
  content: '';
  position: absolute;
  right: -12px;
  top: 50%;
  transform: translateY(-50%);
  width: 8px;
  height: 1px;
  background: #f8f5f0;
  opacity: 0.3;
}
</style>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const form = reactive({
  email: '',
  password: '',
})

const errors = reactive({
  email: '',
  password: '',
})

const loading = ref(false)
const apiError = ref('')

function validate(): boolean {
  let ok = true
  errors.email = ''
  errors.password = ''

  if (!form.email.trim()) {
    errors.email = 'Email không được để trống'
    ok = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = 'Email không hợp lệ'
    ok = false
  }

  if (!form.password) {
    errors.password = 'Mật khẩu không được để trống'
    ok = false
  } else if (form.password.length < 6) {
    errors.password = 'Mật khẩu phải có ít nhất 6 ký tự'
    ok = false
  }

  return ok
}

async function handleLogin() {
  apiError.value = ''
  if (!validate()) return

  loading.value = true
  try {
    await auth.login({ email: form.email, password: form.password })
    router.push({ name: 'Dashboard' })
  } catch (err: any) {
    apiError.value = err.message || 'Đăng nhập thất bại'
  } finally {
    loading.value = false
  }
}
</script>
