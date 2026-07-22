import { defineNuxtConfig } from "nuxt/config";
import tailwindcss from "@tailwindcss/vite";
import Components from "unplugin-vue-components/vite";
import { NaiveUiResolver } from "unplugin-vue-components/resolvers";

export default defineNuxtConfig({
  compatibilityDate: "2026-07-21",
  devtools: { enabled: true },

  // Modules
  modules: ["@pinia/nuxt"],

  css: ["~/assets/css/main.css"],

  // Runtime config cho API URL — có thể override bằng biến môi trường
  runtimeConfig: {
    public: {
      apiBaseUrl:
        process.env.NUXT_PUBLIC_API_BASE_URL ||
        "http://localhost:8000/api/2026",
    },
  },

  // Nuxt auto-imports: Vue, Vue Router, Pinia composables
  imports: {
    autoImport: true,
  },

  ssr: false,

  // Components auto-import
  components: [
    { path: "~/components", pathPrefix: false },
    { path: "~/icons", pathPrefix: false },
  ],

  build: {
    transpile: ["naive-ui", "vueuc"],
  },

  // Vite config (Nuxt dùng Vite bên trong)
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
    vue: {
      script: {
        propsDestructure: true,
      },
    },
  },
});
