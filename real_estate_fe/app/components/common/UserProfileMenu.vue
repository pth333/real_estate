<template>
  <n-dropdown
    trigger="click"
    :options="dropdownOptions"
    @select="handleSelect"
  >
    <n-button circle quaternary size="large">
      <template #icon>
        <n-icon class="text-emerald-600">
          <IconUser />
        </n-icon>
      </template>
    </n-button>
  </n-dropdown>
</template>

<script setup lang="ts">
import {resolveComponent } from "vue";
import { NIcon } from "naive-ui";
import { useAuthStore } from "~/stores/auth";
import { UserMenu } from "~/types/window";

// Khởi tạo các store và class quản lý menu người dùng
const auth = useAuthStore();
const userMenu = new UserMenu();

// Map danh sách tùy chọn từ class sang cấu trúc của Naive UI dropdown dựa trên role, kèm icon và divider
const dropdownOptions = computed(() => {
  const options = userMenu.getFilteredOptions();
  const list: any[] = [];

  options.forEach((option) => {
    // Thêm đường gạch phân cách (divider) trước mục "Đăng xuất" để bám sát giao diện chuẩn
    if (option.key === "logout" && list.length > 0) {
      list.push({
        type: "divider",
        key: "divider-logout",
      });
    }

    list.push({
      label: option.label,
      key: option.key,
      icon: () => {
        let iconComp = null;
        if (option.key === "manage-posts") {
          iconComp = resolveComponent("IconBuilding");
        } else if (option.key === "manage-customers") {
          iconComp = resolveComponent("IconPhone");
        } else if (option.key === "logout") {
          iconComp = resolveComponent("IconLock");
        }

        // resolveComponent trả về component object nếu được tìm thấy trong hệ thống auto-import của Nuxt
        if (iconComp && typeof iconComp !== "string") {
          return h(NIcon, null, { default: () => h(iconComp) });
        }
        return null;
      },
    });
  });

  return list;
});

/**
 * Xử lý đăng xuất tài khoản và chuyển hướng về trang đăng nhập
 */
const handleLogout = async () => {
  await auth.logout();
  navigateTo("/login");
};

/**
 * Xử lý khi người dùng chọn một tùy chọn trong menu dropdown
 * @param key Key của tùy chọn được chọn
 */
const handleSelect = async (key: string) => {
  const option = userMenu.getOptionByKey(key);
  if (option?.path) {
    // Đảm bảo đường dẫn luôn bắt đầu bằng '/' để không gặp lỗi relative route khi chuyển trang
    const path = option.path.startsWith("/") ? option.path : `/${option.path}`;
    navigateTo(path);
  } else if (key === "logout") {
    await handleLogout();
  }
};
</script>
