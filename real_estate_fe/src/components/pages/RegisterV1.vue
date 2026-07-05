<template>
  <div class="flex min-h-screen bg-limestone">
    <!-- Form bên trái - 40% -->
    <div class="flex w-full items-center justify-center bg-foundation px-8 py-12 lg:w-2/5">
      <div class="w-full max-w-sm">
        <h1 class="mb-2 font-semibold tracking-tight text-limestone text-3xl">Đăng ký</h1>
        <p class="mb-8 text-sm text-patina/80">Tạo tài khoản mới</p>

        <form @submit.prevent="handleRegister" class="space-y-5">
          <!-- Họ tên -->
          <div>
            <label class="mb-2 block text-sm font-medium text-limestone/90">Họ và tên</label>
            <div class="plot-input-wrapper">
              <input
                v-model="form.name"
                type="text"
                placeholder="Nguyễn Văn A"
                class="w-full border border-limestone/20 bg-foundation px-4 py-3 text-sm text-limestone placeholder:text-limestone/40 focus:border-patina focus:outline-none"
                :class="{ '!border-red-400': errors.name }"
              />
            </div>
            <p v-if="errors.name" class="mt-1.5 text-xs text-red-400">{{ errors.name }}</p>
          </div>

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

          <!-- Mật khẩu -->
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

          <!-- Xác nhận mật khẩu -->
          <div>
            <label class="mb-2 block text-sm font-medium text-limestone/90">Xác nhận mật khẩu</label>
            <div class="plot-input-wrapper">
              <input
                v-model="form.confirmPassword"
                type="password"
                placeholder="••••••••"
                class="w-full border border-limestone/20 bg-foundation px-4 py-3 text-sm text-limestone placeholder:text-limestone/40 focus:border-patina focus:outline-none"
                :class="{ '!border-red-400': errors.confirmPassword }"
              />
            </div>
            <p v-if="errors.confirmPassword" class="mt-1.5 text-xs text-red-400">{{ errors.confirmPassword }}</p>
          </div>

          <!-- Submit -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full bg-patina px-4 py-3 text-sm font-semibold tracking-tight text-foundation transition hover:bg-patina/90 disabled:cursor-not-allowed disabled:opacity-50"
          >
            <span v-if="loading">Đang xử lý…</span>
            <span v-else>Đăng ký</span>
          </button>
        </form>

        <!-- Thông báo lỗi / thành công -->
        <p v-if="apiError" class="mt-4 text-sm text-red-400">{{ apiError }}</p>
        <p v-if="successMsg" class="mt-4 text-sm text-success">{{ successMsg }}</p>

        <!-- Link sang đăng nhập -->
        <p class="mt-8 text-sm text-limestone/60">
          Đã có tài khoản?
          <router-link :to="{ name: 'Login' }" class="font-medium text-oak hover:underline">
            Đăng nhập
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
            <pattern id="topo-register" x="0" y="0" width="120" height="120" patternUnits="userSpaceOnUse">
              <path d="M20 60c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5" opacity="0.15"/>
              <path d="M20 80c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5" opacity="0.1"/>
              <path d="M20 40c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5" opacity="0.1"/>
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#topo-register)"/>
        </svg>
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="text-center">
            <div class="mb-4 text-6xl font-bold tracking-tighter text-foundation/10">BĐS</div>
            <p class="text-sm font-medium tracking-wide text-foundation/30">BẮT ĐẦU NGAY</p>
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
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const errors = reactive({
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
  errors.name = ''
  errors.email = ''
  errors.password = ''
  errors.confirmPassword = ''

  if (!form.name.trim()) {
    errors.name = 'Họ tên không được để trống'
    ok = false
  }

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

  if (form.password !== form.confirmPassword) {
    errors.confirmPassword = 'Mật khẩu xác nhận không khớp'
    ok = false
  }

  return ok
}

async function handleRegister() {
  apiError.value = ''
  successMsg.value = ''
  if (!validate()) return

  loading.value = true
  try {
    await auth.register({ name: form.name, email: form.email, password: form.password })
    successMsg.value = 'Đăng ký thành công! Đang chuyển đến trang đăng nhập…'
    setTimeout(() => {
      router.push({ name: 'Login' })
    }, 1500)
  } catch (err: any) {
    apiError.value = err.message || 'Đăng ký thất bại'
  } finally {
    loading.value = false
  }
}
</script>
