import { defineNuxtConfig } from "nuxt/config";
import tailwindcss from "@tailwindcss/vite";
import Components from "unplugin-vue-components/vite";
import { NaiveUiResolver } from "unplugin-vue-components/resolvers";
export default defineNuxtConfig({
  compatibilityDate: "2026-07-21",
  devtools: { enabled: true },

  modules: ["@pinia/nuxt", "nuxtjs-naive-ui"],

  ssr: false,

  css: ["~/assets/css/main.css"],

  imports: {
    autoImport: true,
  },

  runtimeConfig: {
    public: {
      apiBaseUrl:
        process.env.NUXT_PUBLIC_API_BASE_URL ||
        "http://localhost:8000/api/2026",
      // Geoapify API key — dùng cho bản đồ + geocoding + places (NUXT_PUBLIC_GEOAPIFY_API_KEY)
      geoapifyApiKey: process.env.NUXT_PUBLIC_GEOAPIFY_API_KEY || "",
    },
  },

  components: [
    { path: "~/components", pathPrefix: false },
    { path: "~/icons", pathPrefix: false },
  ],

  build: {
    transpile: ["vueuc", "naive-ui"],
  },

  vite: {
    plugins: [
      tailwindcss(),
      Components({
        resolvers: [NaiveUiResolver()],
        dts: false,
      }),
    ],
    optimizeDeps: {
      include: ["vueuc"],
    },
  },
});
