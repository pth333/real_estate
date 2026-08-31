<template>
  <div class="flex h-screen overflow-hidden bg-limestone">
    <!-- Left panel -->
    <div class="flex w-full shrink-0 items-center justify-center bg-foundation px-8 lg:w-2/5">
      <div class="w-full max-w-sm">
        <h1 class="mb-2 font-semibold tracking-tight text-limestone text-3xl">
          Đăng nhập
        </h1>
        <p class="mb-8 text-sm text-patina/80">
          Truy cập vào nền tảng bất động sản
        </p>

        <n-form :model="form" label-placement="top">
          <n-form-item label="Email" :feedback="errors.email" :validation-status="errors.email ? 'error' : undefined">
            <n-input v-model:value="form.email" type="text" placeholder="your@email.com" clearable />
          </n-form-item>

          <n-form-item label="Mật khẩu" :feedback="errors.password"
            :validation-status="errors.password ? 'error' : undefined">
            <n-input v-model:value="form.password" type="password" placeholder="••••••••" show-password-on="click"
              clearable />
          </n-form-item>

          <n-button type="primary" attr-type="submit" :loading="loading" :disabled="loading" size="large" block
            @click="handleLogin">
            Đăng nhập
          </n-button>
        </n-form>

        <n-alert v-if="apiError" type="error" :title="apiError" class="mt-4" closable />

        <p class="mt-8 text-sm text-limestone/60">
          Chưa có tài khoản?
          <NuxtLink to="/dang-ky" class="font-medium text-oak hover:underline">Đăng ký</NuxtLink>
        </p>
      </div>
    </div>

    <!-- Right panel -->
    <div class="hidden lg:block lg:w-3/5">
      <div class="relative h-full w-full">
        <svg class="h-full w-full" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <pattern id="topo" x="0" y="0" width="120" height="120" patternUnits="userSpaceOnUse">
              <path d="M20 60c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5"
                opacity="0.15" />
              <path d="M20 80c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5" opacity="0.1" />
              <path d="M20 40c15-8 35-8 50 0s35 8 50 0" fill="none" stroke="#2a3439" stroke-width="0.5" opacity="0.1" />
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#topo)" />
        </svg>
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="text-center">
            <div class="mb-4 text-6xl font-bold tracking-tighter text-foundation/10">
              BĐS
            </div>
            <p class="text-sm font-medium tracking-wide text-foundation/30">
              NỀN TẢNG QUẢN LÝ
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
  alias: "/dang-nhap",
});

useHead({
  title: "Đăng nhập",
});

import { useAuthStore } from "~/stores/auth";

const auth = useAuthStore();
const form = ref({ email: "", password: "" });
const errors = ref({ email: "", password: "" });
const loading = ref(false);
const apiError = ref("");

function validate(): boolean {
  let ok = true;
  errors.value.email = "";
  errors.value.password = "";
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
  return ok;
}

async function handleLogin() {
  apiError.value = "";
  if (!validate()) return;
  loading.value = true;
  try {
    await auth.login({ email: form.value.email, password: form.value.password });
    await navigateTo("/");
  } catch (err: any) {
    apiError.value = err.message || "Đăng nhập thất bại";
  } finally {
    loading.value = false;
  }
}
</script>
