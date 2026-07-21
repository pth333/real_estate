import { defineStore } from "pinia";
import type { LoginRequest, RegisterRequest } from "@/types/auth";
export const useAuthStore = defineStore("auth", () => {
  // Dùng useCookie để tương thích SSR (Nuxt sẽ hydrate từ request cookie)
  const tokenCookie = useCookie<string | null>("auth_token", {
    maxAge: 60 * 60 * 24 * 7, // 7 ngày
    sameSite: "lax",
  });
  const emailCookie = useCookie<string | null>("auth_email", {
    maxAge: 60 * 60 * 24 * 7,
    sameSite: "lax",
  });
  const nameCookie = useCookie<string | null>("auth_name", {
    maxAge: 60 * 60 * 24 * 7,
    sameSite: "lax",
  });

  const token = ref<string | null>(tokenCookie.value ?? null);
  const email = ref<string | null>(emailCookie.value ?? null);
  const name = ref<string | null>(nameCookie.value ?? null);

  // Đồng bộ ngược lên cookie khi ref thay đổi
  watch(token, (val) => {
    tokenCookie.value = val;
  });
  watch(email, (val) => {
    emailCookie.value = val;
  });
  watch(name, (val) => {
    nameCookie.value = val;
  });

  const isAuthenticated = computed(() => !!token.value);
  const userName = computed(() => name.value ?? "");
  const userEmail = computed(() => email.value ?? "");

  /** Lưu token & user vào cookie */
  function setSession(tok: string, userEmail?: string, userName?: string) {
    token.value = tok;
    email.value = userEmail ?? null;
    name.value = userName ?? null;
  }

  /** Xoá session */
  function clearSession() {
    token.value = null;
    email.value = null;
    name.value = null;
  }

  /** Login */
  async function login(payload: LoginRequest) {
    const res = await $fetch("/api/auth/login", { method: "POST", body: payload });

    if (!res.success || !res.data?.token) {
      throw new Error(res.message || "Đăng nhập thất bại");
    }

    setSession(res.data.token, payload.email);

    return res;
  }

  /** Đăng ký */
  async function register(payload: RegisterRequest) {
    const res = await $fetch("/api/auth/register", { method: "POST", body: payload });

    if (!res.success) {
      throw new Error(res.message || "Đăng ký thất bại");
    }

    return res;
  }

  /** Refresh token */
  async function refreshToken() {
    try {
      const res = await $fetch("/api/auth/refresh", { method: "POST" });
      if (res.success && res.data?.token) {
        token.value = res.data.token;
        return true;
      }
    } catch {
      clearSession();
    }
    return false;
  }

  /** Logout */
  async function logout() {
    try {
      await $fetch("/api/auth/logout", { method: "POST" });
    } catch {
      // ignore
    }
    clearSession();
  }

  return {
    token,
    email,
    name,
    isAuthenticated,
    userName,
    userEmail,
    login,
    register,
    refreshToken,
    logout,
    clearSession,
  };
});
