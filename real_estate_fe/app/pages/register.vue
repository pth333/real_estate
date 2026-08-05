<template>
  <div class="flex h-screen overflow-hidden bg-limestone">
    <!-- Left panel -->
    <div class="flex w-full shrink-0 items-center justify-center bg-foundation px-8 lg:w-2/5">
      <div class="w-full max-w-sm">
        <h1 class="mb-2 font-semibold tracking-tight text-limestone text-3xl">
          Đăng ký
        </h1>
        <p class="mb-8 text-sm text-patina/80">Tạo tài khoản mới</p>

        <n-form :model="form" label-placement="top">
          <n-form-item label="Họ và tên" :feedback="errors.name" :validation-status="errors.name ? 'error' : undefined">
            <n-input v-model:value="form.name" type="text" placeholder="Nguyễn Văn A" clearable />
          </n-form-item>

          <n-form-item label="Email" :feedback="errors.email" :validation-status="errors.email ? 'error' : undefined">
            <n-input v-model:value="form.email" type="text" placeholder="your@email.com" clearable />
          </n-form-item>

          <n-form-item label="Mật khẩu" :feedback="errors.password"
            :validation-status="errors.password ? 'error' : undefined">
            <n-input v-model:value="form.password" type="password" placeholder="••••••••" show-password-on="click"
              clearable />
          </n-form-item>

          <n-form-item label="Xác nhận mật khẩu" :feedback="errors.confirmPassword"
            :validation-status="errors.confirmPassword ? 'error' : undefined">
            <n-input v-model:value="form.confirmPassword" type="password" placeholder="••••••••"
              show-password-on="click" clearable />
          </n-form-item>

          <n-button type="primary" @click="handleRegister" :loading="loading" :disabled="loading" size="large" block>
            Đăng ký
          </n-button>
        </n-form>

        <n-alert v-if="apiError" type="error" :title="apiError" class="mt-4" closable />
        <n-alert v-if="successMsg" type="success" :title="successMsg" class="mt-4" closable />

        <p class="mt-8 text-sm text-limestone/60">
          Đã có tài khoản?
          <NuxtLink to="/login" class="font-medium text-oak hover:underline">Đăng nhập</NuxtLink>
        </p>
      </div>
    </div>

    <!-- Right panel -->
    <div class="hidden lg:block lg:w-3/5">
      <div class="relative h-full w-full">
        <svg class="h-full w-full" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <pattern id="topo-register" x="0" y="0" width="120" height="120" patternUnits="userSpaceOnUse">
              <path d="M20 60c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5"
                opacity="0.15" />
              <path d="M20 80c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5" opacity="0.1" />
              <path d="M20 40c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5" opacity="0.1" />
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#topo-register)" />
        </svg>
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="text-center">
            <div class="mb-4 text-6xl font-bold tracking-tighter text-foundation/10">
              BĐS
            </div>
            <p class="text-sm font-medium tracking-wide text-foundation/30">
              BẮT ĐẦU NGAY
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: false });

import { useAuthStore } from "~/stores/auth";

const auth = useAuthStore();
const form = ref({
  name: "",
  email: "",
  password: "",
  confirmPassword: "",
});
const errors = ref({
  name: "",
  email: "",
  password: "",
  confirmPassword: "",
});
const loading = ref(false);
const apiError = ref("");
const successMsg = ref("");

function validate(): boolean {
  let ok = true;
  errors.value.name = "";
  errors.value.email = "";
  errors.value.password = "";
  errors.value.confirmPassword = "";
  if (!form.value.name.trim()) {
    errors.value.name = "Họ tên không được để trống";
    ok = false;
  }
  if (!form.value.email.trim()) {
    errors.value.email = "Email không được để trống";
    ok = false;
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.value.email)) {
    errors.value.email = "Email không hợp lệ";
    ok = false;
  }
  if (!form.value.password) {
    errors.value.password = "Mật khẩu không được để trống";
    ok = false;
  } else if (form.value.password.length < 6) {
    errors.value.password = "Mật khẩu phải có ít nhất 6 ký tự";
    ok = false;
  }
  if (form.value.password !== form.value.confirmPassword) {
    errors.value.confirmPassword = "Mật khẩu xác nhận không khớp";
    ok = false;
  }
  return ok;
}

async function handleRegister() {
  apiError.value = "";
  successMsg.value = "";
  if (!validate()) return;
  loading.value = true;
  try {
    await auth.register({
      name: form.value.name,
      email: form.value.email,
      password: form.value.password,
    });
    successMsg.value = "Đăng ký thành công! Đang chuyển đến trang đăng nhập…";
    setTimeout(() => navigateTo("/login"), 1500);
  } catch (err: any) {
    apiError.value = err.message || "Đăng ký thất bại";
  } finally {
    loading.value = false;
  }
}
</script>