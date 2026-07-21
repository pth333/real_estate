import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
} from "~/app/types/auth";

import useApi from "~/app/composables/useApi";

const api = useApi();

export const authApi = {
  login: async (payload: LoginRequest) => {
    const res = await api.request<AuthResponse>({
      url: "/auth/login",
      data: payload,
      method: "post",
    });
    return res.data;
  },
  register: async (payload: RegisterRequest) => {
    const res = await api.request<AuthResponse>({
      url: "/auth/register",
      data: payload,
      method: "post",
    });
    return res.data;
  },
  refreshToken: async () => {
    const res = await api.request<AuthResponse>({
      url: "/auth/refresh",
      data: {},
      method: "post",
    });
    return res.data;
  },
  logout: async () => {
    const res = await api.request<AuthResponse>({
      url: "/auth/logout",
      data: {},
      method: "post",
    });
    return res.data;
  },
};
