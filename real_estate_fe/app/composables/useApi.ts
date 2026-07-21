import type { AxiosInstance } from "axios";

export default function useApi(): AxiosInstance {
  const { api } = useNuxtApp();
  return api as AxiosInstance;
}
