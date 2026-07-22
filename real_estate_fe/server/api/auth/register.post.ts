export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const body = await readBody(event);
  return $fetch(`${config.public.apiBaseUrl}/auth/register`, {
    method: "POST",
    body,
  });
});
