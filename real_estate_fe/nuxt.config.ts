// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2026-07-21",
  devtools: { enabled: true },

  // Modules
  modules: ["@pinia/nuxt"],

  css: ["~/app/assets/css/main.css"],

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

  ssr: true,

  // Components auto-import
  components: [{ path: "~/app/components", pathPrefix: false }],

  // Vite config (Nuxt dùng Vite bên trong)
  vite: {
    plugins: [
      // @ts-ignore
      (await import("@tailwindcss/vite")).default,
    ],
    vue: {
      script: {
        propsDestructure: true,
      },
    },
  },
});
