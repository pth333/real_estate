import { refElement } from 'vuetify/lib/util/helpers.mjs';
import { InformationRealestate } from '~/types/real_estate'

export const useCreatePost = defineStore("create-post", () => {
  type Tab = "information" | "upload" | "review";

  // Dữ liệu tin đăng tổ chức trong class InformationRealestate
  const form = ref(new InformationRealestate())

  // Lỗi hiển thị dưới từng ô input, nhóm theo section.
  // Mỗi section component tự cập nhật lỗi khi validate.
  const errors = ref({
    address: { province: "", ward: "", detail_address: "" },
    mainInfo: { real_estate_type: "", area: "", price: "", unit: "" },
    contact: { contact_name: "", contact_email: "", contact_phone: "" },
    description: { title: "", description: "" },
  })

  const tab = ref<Tab>("information");
  const next = ref(false);

  // ── Payload: chính là class InformationRealestate ──
  const payload = computed(() => form.value)

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
    form,
    payload,
    errors,
    errorsMainInfo: errors.value.mainInfo,
    errorsAddress: errors.value.address,
    errorsContact: errors.value.contact,
    errorsDescription: errors.value.description,
    tab,
    next,
    currentStepLabel,
    currentStepProgress,
  };
});
