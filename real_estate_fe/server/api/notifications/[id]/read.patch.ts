export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const { id } = getRouterParams(event);
  return $fetch(`${config.public.apiBaseUrl}/notifications/${id}/read`, {
    method: "PATCH",
  });
});
