export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const query = getQuery(event);
  return $fetch(`${config.public.apiBaseUrl}/dashboard/summary`, {
    params: query,
  });
});
