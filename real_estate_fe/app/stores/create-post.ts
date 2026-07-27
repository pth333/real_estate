export const useCreatePost = defineStore("create-post", () => {
  type Tab = "information" | "upload" | "review";
  const province = ref<string | null>(null);
  const ward = ref<string | null>(null);
  const detail_address = ref<string | null>(null);
  const tab = ref<Tab>("information");
  const listingType = ref<"sell" | "rent">("sell");

  const payload = computed(() => {});

  const submitCreatePost = () => {};

  const stepLabels: Record<Tab, string> = {
    information: "Bước 1. Thông tin BĐS",
    upload: "Bước 2. Hình ảnh & video",
    review: "Bước 3. Kiểm tra & đăng tin",
  };

  const stepProgress: Record<Tab, number> = {
    information: 33,
    upload: 66,
    review: 100,
  };

  const currentStepLabel = computed(() => stepLabels[tab.value]);
  const currentStepProgress = computed(() => stepProgress[tab.value]);

  return {
    province,
    ward,
    detail_address,
    tab,
    currentStepLabel,
    currentStepProgress,
    listingType,
  };
});
