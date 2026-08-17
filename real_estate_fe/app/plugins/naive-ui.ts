import naive from "naive-ui";
import { setup } from "@css-render/vue3-ssr";

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(naive);

  if (import.meta.server) {
    const { collect } = setup(nuxtApp.vueApp);
    const { ssrContext } = nuxtApp;
    if (ssrContext) {
      const originalRenderMeta = ssrContext.renderMeta;
      ssrContext.renderMeta = async () => {
        const meta = await (originalRenderMeta as any)?.() ?? {};
        return {
          ...meta,
          headTags: (meta.headTags ?? "") + collect(),
        };
      };
    }
  }
});