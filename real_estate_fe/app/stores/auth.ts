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
 * - init()             : NẠP PHIÊN duy nhất 1 lần — tự refresh token trước (nếu cần) rồi
 *                        mới gọi GET /auth/user-current-info. Dùng ở GlobalInit + middleware.
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
  // Promise của lần nạp phiên đang chạy — mọi nơi gọi init() đều chờ CHUNG 1 lần
  let initPromise: Promise<void> | null = null;

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
    initPromise = null;
    if (import.meta.client) {
      window.currentUser = undefined;
    }
  }

  /**
   * Nạp phiên đăng nhập — MỘT hàm duy nhất, chạy tuần tự, có chống gọi trùng.
   *
   * Thứ tự cố định (đây là chỗ trước đây bị race):
   *   1. Có dấu hiệu đã đăng nhập (còn token HOẶC còn user cookie) mà chưa có access token
   *      → gọi /auth/refresh. BE đọc cookie httpOnly `refresh_token` nên KHÔNG cần token cũ.
   *   2. Có token → GET /auth/user-current-info (silent: request nền không hiện toast lỗi).
   *   3. Chỉ khi 401/403 (kể cả sau khi đã refresh) mới xoá phiên; lỗi mạng thì giữ user đọc
   *      từ cookie để không văng đăng nhập oan và không mất menu của user.
   */
  function init(force = false): Promise<void> {
    if (!force && initPromise) return initPromise;

    const promise = (async () => {
      // Khách vãng lai hoàn toàn: không có gì để nạp, cũng KHÔNG gọi refresh
      if (!token.value && !user.value) return;

      // (1) Chưa có access token (hết hạn/bị xoá) nhưng còn phiên → lấy lại bằng refresh token
      if (!token.value) {
        const refreshed = await refreshToken();
        if (!refreshed) {
          clearSession();
          return;
        }
      }

      // (2) Lấy quyền MỚI NHẤT (admin vừa đổi role thì không cần đăng nhập lại)
      const tokenAtRequest = token.value;
      try {
        const current = await useAuthService().getUserCurrentInfo();
        // Phiên bị đổi/đăng xuất trong lúc chờ → bỏ kết quả, không ghi đè state mới
        if (!tokenAtRequest || token.value !== tokenAtRequest) return;
        user.value = current;
        accessLoaded.value = true;
      } catch (error: unknown) {
        if (token.value !== tokenAtRequest) return;
        const status = (error as { response?: { status?: number } })?.response?.status;
        // (3) Token vừa refresh vẫn bị từ chối → phiên thật sự hết hạn
        if (status === 401 || status === 403) {
          clearSession();
          return;
        }
        // Lỗi mạng/server tạm thời: giữ user hiện có, KHÔNG toast, KHÔNG xoá phiên
      }
    })();

    // Giữ promise cục bộ để trả về (clearSession ở trong có thể đặt initPromise = null
    // để lần nạp sau được phép chạy lại)
    initPromise = promise;
    return promise;
  }

  /** Giữ tên cũ cho nơi đang gọi (middleware/route) — nay chỉ là alias của init() */
  async function ensureAccessLoaded(): Promise<void> {
    await init();
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
    init,
    ensureAccessLoaded,

    login,
    register,
    refreshToken,
    logout,
    clearSession,
  };
});
