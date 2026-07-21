export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  return $fetch(`${config.public.apiBaseUrl}/category`);
});
