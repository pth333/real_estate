export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const { action } = getRouterParams(event);
  const body = await readBody(event);
  return $fetch(`${config.public.apiBaseUrl}/auth/${action}`, {
    method: "POST",
    body,
  });
});
