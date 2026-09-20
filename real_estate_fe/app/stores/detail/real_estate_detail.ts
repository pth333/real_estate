import { defineStore } from "pinia";
import { useRealEstateService } from "~/services/real-estate.service";
import type { RealEstateResponse } from "~/types/real_estate";

export const useRealEstateDetail = defineStore("real_estate_detail", () => {
  const loading = ref(false);
  const listing = ref<RealEstateResponse | null>(null);
  const showPhone = ref(false);

  async function fetchDetail(id: number) {
    const realEstateService = useRealEstateService();
    loading.value = true;
    try {
      const res = await realEstateService.getDetail(id);
      listing.value = res ?? null;
    } catch {
      listing.value = null;
    } finally {
      loading.value = false;
    }
  }

  const maskedPhone = computed(() => {
    const phone = listing.value?.agent_phone || "";
    if (showPhone.value || !phone) return phone;
    return phone.replace(/.{3}$/, "***");
  });
  function handleShare() {
    navigator.clipboard?.writeText(window.location.href);
    window.message?.success("Đã copy link chia sẻ");
  }

  return {
    loading,
    listing,
    showPhone,
    maskedPhone,
    handleShare,
    fetchDetail,
  };
});
