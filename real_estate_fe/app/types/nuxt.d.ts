// types/nuxt.d.ts
import type { api } from '~/plugins/api.client'

declare module '#app' {
  interface NuxtApp {
    $api: typeof api
  }
}

declare module 'vue' {
  interface ComponentCustomProperties {
    $api: typeof api
  }
}

export {}