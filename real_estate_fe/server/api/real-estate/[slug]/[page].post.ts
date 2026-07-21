export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const { slug, page } = getRouterParams(event);
  const body = await readBody(event);
  return $fetch(`${config.public.apiBaseUrl}/real-estate/${slug}/${page}`, {
    method: "POST",
    body,
  });
});
