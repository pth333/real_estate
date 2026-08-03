import { InformationRealestate } from "~/types/real_estate";

export const useCreatePost = defineStore("create-post", () => {
  // Dữ liệu tin đăng tổ chức trong class InformationRealestate
  const form = ref(new InformationRealestate());

  // Lỗi hiển thị dưới từng ô input, tách riêng theo từng section.
  // Mỗi biến là 1 ref object — mutate field bên trong, không gán lại toàn bộ
  // để giữ reference (tránh bug mất hiển thị lỗi như trước).
  const errorsMainInfo = ref({
    real_estate_type: "",
    area: "",
    price: "",
    unit: "",
  });

  const errorsAddress = ref({
    province: "",
    ward: "",
    detail_address: "",
  });

  const errorsContact = ref({
    contact_name: "",
    contact_email: "",
    contact_phone: "",
  });

  const errorsDescription = ref({
    title: "",
    description: "",
  });

  // Tách id + tên từ chuỗi "12-CanHo" (value của n-select). KHÔNG mutate form
  // để payload gửi đi tách lấy id riêng, còn form giữ chuỗi gốc cho hiển thị select.
  const customValueRealEstateType = computed(() => {
    if (!form.value.real_estate_type) return null;
    const [id, name] = form.value.real_estate_type.split("-");
    return { id: Number(id), name };
  });
  // ── Payload: copy form, tách real_estate_type ra chỉ còn id để gửi backend ──
  const payload = computed(() => {
    const { real_estate_type, ...rest } = form.value;
    const [id] = (real_estate_type ?? "").split("-");
    return {
      ...rest,
      real_estate_type: id || null,
    };
  });

  const stepLabels: Record<string, string> = {
    information: "Bước 1. Thông tin BĐS",
    upload: "Bước 2. Hình ảnh & video",
    review: "Bước 3. Kiểm tra & đăng tin",
  };

  const stepProgress: Record<string, number> = {
    information: 33,
    upload: 66,
    review: 100,
  };

  const currentStepLabel = computed(() => stepLabels[form.value.tab]);
  const currentStepProgress = computed(() => stepProgress[form.value.tab]);

  // ── Validate chung (dùng cho cả nút Tiếp tục và nút Tạo với AI) ──
  const validateMainInfo = (): boolean => {
    errorsMainInfo.value.real_estate_type = form.value.real_estate_type
      ? ""
      : "Vui lòng chọn loại bất động sản";
    errorsMainInfo.value.area = form.value.area
      ? ""
      : "Vui lòng nhập diện tích";
    errorsMainInfo.value.price = form.value.price
      ? ""
      : "Vui lòng nhập mức giá";
    errorsMainInfo.value.unit = form.value.unit
      ? ""
      : "Vui lòng chọn đơn vị giá";
    return Object.values(errorsMainInfo.value).every((msg) => msg === "");
  };

  const validateAddress = (): boolean => {
    errorsAddress.value.province = form.value.province
      ? ""
      : "Vui lòng chọn tỉnh/thành phố";
    errorsAddress.value.ward = form.value.ward ? "" : "Vui lòng chọn phường/xã";
    errorsAddress.value.detail_address = form.value.detail_address?.trim()
      ? ""
      : "Vui lòng nhập địa chỉ chi tiết";
    return Object.values(errorsAddress.value).every((msg) => msg === "");
  };

  const validateContact = (): boolean => {
    errorsContact.value.contact_name = form.value.contact_name
      ? ""
      : "Vui lòng nhập tên liên hệ";
    errorsContact.value.contact_email = form.value.contact_email
      ? ""
      : "Vui lòng nhập email";
    errorsContact.value.contact_phone = form.value.contact_phone
      ? ""
      : "Vui lòng nhập số điện thoại";
    return Object.values(errorsContact.value).every((msg) => msg === "");
  };

  const validateDescription = (): boolean => {
    const titleLen = form.value.title?.length ?? 0;
    const descLen = form.value.description?.length ?? 0;
    errorsDescription.value.title = titleLen
      ? titleLen < 30
        ? `Tiêu đề tối thiểu 30 ký tự (hiện ${titleLen})`
        : ""
      : "Vui lòng nhập tiêu đề";
    errorsDescription.value.description = descLen
      ? descLen < 30
        ? `Mô tả tối thiểu 30 ký tự (hiện ${descLen})`
        : ""
      : "Vui lòng nhập mô tả";
    return Object.values(errorsDescription.value).every((msg) => msg === "");
  };

  // Validate toàn bộ thông tin bắt buộc ở bước 1
  const validateInformation = (): boolean => {
    const validations = [
      validateAddress(),
      validateMainInfo(),
      validateContact(),
      validateDescription(),
    ];
    return validations.every(Boolean);
  };

  const validateForAI = (): boolean => {
    const validations = [
      validateAddress(),
      validateMainInfo(),
      validateContact(),
    ];
    return validations.every(Boolean);
  };

  return {
    form,
    payload,
    errorsMainInfo,
    errorsAddress,
    errorsContact,
    errorsDescription,
    currentStepLabel,
    currentStepProgress,
    customValueRealEstateType,
    validateMainInfo,
    validateAddress,
    validateContact,
    validateDescription,
    validateInformation,
    validateForAI,
  };
});
