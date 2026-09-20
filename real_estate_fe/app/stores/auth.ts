import { defineStore } from "pinia";
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  UserInfo,
} from "~/types/auth";
import { useSession } from "~/composables/useSession";
import { useAuthService } from "~/services/auth.service";
import { useTrackingService } from "~/services/tracking.service";

/**
 * Store auth — cũng là STATE GLOBAL của user hiện tại:
 * - user               : thông tin user (kèm roles[] + permissions[])
 * - permissions/roles  : bản sao dạng mảng để component đọc trực tiếp
 * - hasPermission()/hasRole() : dùng để ẩn/hiện UI
 * - fetchCurrentUser(): gọi GET /auth/user-current-info, dùng ở middleware route và lúc khởi động
 */
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
  // Cờ cho biết đã gọi API user-current-info ít nhất 1 lần trong phiên này
  const accessLoaded = ref(false);
  // Promise của request đang bay — dùng để chống gọi trùng khi nhiều nơi cùng yêu cầu
  let accessPromise: Promise<UserInfo | null> | null = null;

  watch(token, (val) => {
    tokenCookie.value = val;
  });

  watch(user, (val) => {
    userCookie.value = val;
    // Đồng bộ ra biến global trên window để nơi không dùng được store vẫn đọc được
    if (import.meta.client) {
      window.currentUser = val ?? undefined;
    }
  });

  const isAuthenticated = computed(() => !!token.value);
  const roles = computed<string[]>(() => user.value?.roles ?? []);
  const permissions = computed<string[]>(() => user.value?.permissions ?? []);
  const isAdmin = computed(() => roles.value.includes("ADMIN"));

  /** User có ít nhất 1 trong các role này không */
  function hasAnyRole(expected: string[]): boolean {
    if (!expected.length) return true;
    return expected.some((code) => roles.value.includes(code));
  }

  /** User có ít nhất 1 trong các permission này không (dùng để ẩn/hiện UI) */
  function hasAnyPermission(expected: string[]): boolean {
    if (!expected.length) return true;
    return expected.some((code) => permissions.value.includes(code));
  }

  /** Tiện dụng cho template: `v-if="auth.can('admin.escrow.view')"` */
  function can(permission: string): boolean {
    return permissions.value.includes(permission);
  }

  function setSession(tok: string, usr?: UserInfo) {
    token.value = tok;
    if (usr) {
      user.value = usr;
      accessLoaded.value = true;
    }
  }

  function clearSession() {
    token.value = null;
    user.value = null;
    accessLoaded.value = false;
    if (import.meta.client) {
      window.currentUser = undefined;
    }
  }

  /**
   * Lấy thông tin user hiện tại từ backend và ghi vào state global.
   * Trả về user (hoặc user cũ đọc từ cookie nếu chỉ là lỗi mạng).
   *
   * CHỐNG GỌI TRÙNG: middleware route và GlobalInit có thể cùng gọi hàm này
   * lúc tải trang → dùng chung 1 promise đang bay, chỉ phát 1 request.
   */
  function fetchCurrentUser(force = false): Promise<UserInfo | null> {
    if (!token.value) return Promise.resolve(null);

    // Đã có quyền rồi và không ép gọi lại → dùng luôn state hiện tại
    if (!force && accessLoaded.value) return Promise.resolve(user.value);

    // Đang có request bay → trả về chính promise đó, không tạo request thứ 2
    if (accessPromise) return accessPromise;

    // Ghi nhớ token lúc bắt đầu: nếu user đăng xuất/đổi phiên trong lúc chờ
    // thì bỏ qua kết quả trả về, tránh ghi đè lại state của phiên đã kết thúc.
    const tokenAtRequest = token.value;

    accessPromise = (async () => {
      try {
        const current = await useAuthService().getUserCurrentInfo();
        if (token.value !== tokenAtRequest) return null;
        user.value = current;
        accessLoaded.value = true;
        return current;
      } catch (error: unknown) {
        if (token.value !== tokenAtRequest) return null;
        const status = (error as { response?: { status?: number } })?.response?.status;
        if (status === 401 || status === 403) {
          // Token hỏng/hết hạn hoặc tài khoản bị khoá → xoá phiên,
          // middleware route sẽ đá về trang đăng nhập.
          clearSession();
          return null;
        }
        // Lỗi mạng/server tạm thời: giữ user cũ để không văng đăng nhập oan
        return user.value;
      } finally {
        accessPromise = null;
      }
    })();

    return accessPromise;
  }

  /** Đảm bảo đã có quyền mới nhất — gọi nhiều lần cũng chỉ phát 1 request */
  async function ensureAccessLoaded(): Promise<void> {
    await fetchCurrentUser();
  }

  async function login(payload: LoginRequest) {
    const res = await useAuthService().login(payload);

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
    const res = await useAuthService().register(payload);

    if (!res.success) {
      throw new Error(res.message || "Đăng ký thất bại");
    }

    return res;
  }

  async function refreshToken() {
    try {
      const res = await useAuthService().refresh();
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
    try {
      await useAuthService().logout();
    } catch {
      // ignore
    }
    clearSession();
  }

  return {
    token,
    user,
    accessLoaded,
    roles,
    permissions,
    isAdmin,
    isAuthenticated,

    hasAnyRole,
    hasAnyPermission,
    can,
    fetchCurrentUser,
    ensureAccessLoaded,

    login,
    register,
    refreshToken,
    logout,
    clearSession,
  };
});
