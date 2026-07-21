export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const body = await readBody(event);
  return $fetch(`${config.public.apiBaseUrl}/real-estate/list`, {
    method: "POST",
    body,
  });
});
