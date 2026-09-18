// Quản lý trạng thái xác thực số điện thoại (dùng cho đăng tin)
import { useState } from "#app";

const KEY = "phone_verified";
const PHONE_KEY = "verified_phone";
const PHONES_LIST_KEY = "verified_phones_list";

export const usePhoneVerification = () => {
  const phoneVerified = useState(KEY, () => false);
  const verifiedPhone = useState(PHONE_KEY, () => "");
  const verifiedPhones = useState<string[]>(PHONES_LIST_KEY, () => []);

  onMounted(() => {
    phoneVerified.value = localStorage.getItem(KEY) === "true";
    verifiedPhone.value = localStorage.getItem(PHONE_KEY) || "";

    try {
      const listStr = localStorage.getItem(PHONES_LIST_KEY);
      if (listStr) {
        verifiedPhones.value = JSON.parse(listStr);
      } else {
        const currentPhone = localStorage.getItem(PHONE_KEY);
        if (currentPhone) {
          verifiedPhones.value = [currentPhone];
          localStorage.setItem(PHONES_LIST_KEY, JSON.stringify([currentPhone]));
        }
      }
    } catch (e) {
      console.error("Lỗi parse verified_phones_list:", e);
    }
  });

  function setPhoneVerified(phone: string) {
    if (import.meta.client) {
      localStorage.setItem(KEY, "true");
      localStorage.setItem(PHONE_KEY, phone);

      // Thêm vào danh sách các số điện thoại đã xác thực nếu chưa có
      if (!verifiedPhones.value.includes(phone)) {
        verifiedPhones.value.push(phone);
        localStorage.setItem(
          PHONES_LIST_KEY,
          JSON.stringify(verifiedPhones.value),
        );
      }
    }
    phoneVerified.value = true;
    verifiedPhone.value = phone;
  }

  function clearPhoneVerified() {
    if (import.meta.client) {
      localStorage.removeItem(KEY);
      localStorage.removeItem(PHONE_KEY);
      localStorage.removeItem(PHONES_LIST_KEY);
    }
    phoneVerified.value = false;
    verifiedPhone.value = "";
    verifiedPhones.value = [];
  }

  return {
    phoneVerified,
    verifiedPhone,
    verifiedPhones,
    setPhoneVerified,
    clearPhoneVerified,
  };
};
