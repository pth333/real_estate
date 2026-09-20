import { defineStore } from "pinia";
import type {
  LoginRequest,
  RegisterRequest,
  UserInfo,
} from "~/types/auth";
import { useSession } from "~/composables/useSession";
import { useAuthService } from "~/services/auth.service";
import { useTrackingService } from "~/services/tracking.service";

export const useAuthStore = defineStore("auth", () => {
  const tokenCookie = useCookie<string | null>("auth_token", {
    maxAge: 60 * 60 * 24 * 7,
    sameSite: "lax",
    path: "/",
  });

  const userCookie = useCookie<UserInfo | null>("auth_user", {
    maxAge: 60 * 60 * 24 * 7,
    sameSite: "lax",
    path: "/",
  });

  const { sessionId } = useSession();

  const token = ref<string | null>(tokenCookie.value ?? null);
  const user = ref<UserInfo | null>(userCookie.value ?? null);

  watch(token, (val) => {
    tokenCookie.value = val;
  });

  watch(user, (val) => {
    userCookie.value = val;
  });

  const isAuthenticated = computed(() => !!token.value);

  function setSession(tok: string, usr?: UserInfo) {
    token.value = tok;
    if (usr) {
      user.value = usr;
    }
  }

  function clearSession() {
    token.value = null;
    user.value = null;
  }

  async function login(payload: LoginRequest) {
    const authService = useAuthService();
    const res = await authService.login(payload);
    
    if (!res.success || !res.data?.token) {
      throw new Error(res.message || "Đăng nhập thất bại");
    }

    setSession(res.data.token, res.data.user);

    // Tích hợp Session Merging sau khi đăng nhập thành công
    try {
      if (sessionId) {
        await useTrackingService().mergeSession(sessionId.value);
      }
    } catch (err) {
      console.error("Failed to merge session on login", err);
    }

    return res;
  }

  async function register(payload: RegisterRequest) {
    const authService = useAuthService();
    const res = await authService.register(payload);

    if (!res.success) {
      throw new Error(res.message || "Đăng ký thất bại");
    }

    return res;
  }

  async function refreshToken() {
    const authService = useAuthService();
    try {
      const res = await authService.refresh();
      if (!res.success) {
        throw new Error(res.message || "Refresh token failed");
      }
      if (res.data?.token) {
        token.value = res.data.token;
        return true;
      }
      throw new Error("No token in response");
    } catch {
      clearSession();
      return false;
    }
  }

  async function logout() {
    const authService = useAuthService();
    try {
      await authService.logout();
    } catch {
      // ignore
    }
    clearSession();
  }

  return {
    token,
    user,
    isAuthenticated,

    login,
    register,
    refreshToken,
    logout,
    clearSession,
  };
});
