import naive from "naive-ui";
import { setup } from "@css-render/vue3-ssr";

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(naive);

  if (import.meta.server) {
    const { collect } = setup(nuxtApp.vueApp);

    // Cache styles để tránh mất dữ liệu khi collect() nhiều lần trong cùng request
    let collectedStyles = "";
    const getStyles = () => {
      if (!collectedStyles) {
        collectedStyles = collect();
      }
      return collectedStyles;
    };

    // Cách 1: Hook vào 'app:rendered' (Hỗ trợ Unhead/Nuxt 3 hiện đại)
    // nuxtApp.hooks.hook("app:rendered", ({ ssrContext }) => {
    //   if (ssrContext && ssrContext.head) {
    //     const css = getStyles();
    //     if (css) {
    //       ssrContext.head.push({
    //         style: [
    //           {
    //             children: css,
    //             "cssr-id": "naive-ui-ssr"
    //           }
    //         ]
    //       });
    //     }
    //   }
    // });

    // Cách 2: Fallback qua ssrContext.renderMeta (Cho các bản Nuxt cũ/truyền thống)
    const { ssrContext } = nuxtApp;
    if (ssrContext) {
      const originalRenderMeta = ssrContext.renderMeta;
      ssrContext.renderMeta = async () => {
        const meta = await (originalRenderMeta as any)?.() ?? {};
        const css = getStyles();
        return {
          ...meta,
          headTags: (meta.headTags ?? "") + css,
        };
      };
    }
  }
});