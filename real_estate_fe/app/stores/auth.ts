import { defineStore } from "pinia";
import type { LoginRequest, RegisterRequest, AuthResponse } from "~/types/auth";

export const useAuthStore = defineStore("auth", () => {
  const tokenCookie = useCookie<string | null>("auth_token", {
    maxAge: 60 * 60 * 24 * 7,
    sameSite: "lax",
    path: "/",
  });
  // const emailCookie = useCookie<string | null>("auth_email", {
  //   maxAge: 60 * 60 * 24 * 7,
  //   sameSite: "lax",
  //   path: "/",
  // });
  // const nameCookie = useCookie<string | null>("auth_name", {
  //   maxAge: 60 * 60 * 24 * 7,
  //   sameSite: "lax",
  //   path: "/",
  // });

  const token = ref<string | null>(tokenCookie.value ?? null);
  // const email = ref<string | null>(emailCookie.value ?? null);
  // const name = ref<string | null>(nameCookie.value ?? null);

  watch(token, (val) => {
    tokenCookie.value = val;
  });
  // watch(email, (val) => {
  //   emailCookie.value = val;
  // });
  // watch(name, (val) => {
  //   nameCookie.value = val;
  // });

  const isAuthenticated = computed(() => !!token.value);

  function setSession(tok: string, userEmail?: string, userName?: string) {
    token.value = tok;
  }

  function clearSession() {
    token.value = null;
  }

  async function login(payload: LoginRequest) {
    const { $api } = useNuxtApp();
    const res = await $api.post<AuthResponse>("/auth/login", payload);

    if (!res.success || !res.data?.token) {
      throw new Error(res.message || "Đăng nhập thất bại");
    }

    setSession(res.data.token);
    return res;
  }

  async function register(payload: RegisterRequest) {
    const { $api } = useNuxtApp();
    const res = await $api.post<AuthResponse>("/auth/register", payload);

    if (!res.success) {
      throw new Error(res.message || "Đăng ký thất bại");
    }

    return res;
  }

  async function refreshToken() {
    try {
      const { $api } = useNuxtApp();
      const res = await $api.post<AuthResponse>("/auth/refresh");
      if (res.success && res.data?.token) {
        token.value = res.data.token;
        return true;
      }
    } catch {
      clearSession();
    }
    return false;
  }

  async function logout() {
    try {
      const { $api } = useNuxtApp();
      await $api.post("/auth/logout");
    } catch {
      // ignore
    }
    clearSession();
  }

  return {
    token,
    isAuthenticated,

    login,
    register,
    refreshToken,
    logout,
    clearSession,
  };
});
