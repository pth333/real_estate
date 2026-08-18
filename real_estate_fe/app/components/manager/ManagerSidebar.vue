<template>
  <div class="w-64 border-r border-gray-200 bg-white h-[calc(100vh-64px)] flex flex-col justify-between py-4 select-none">
    <div class="flex flex-col gap-6 px-4">
      <!-- Nút Đăng Tin Mới -->
      <n-button type="error" size="large" block class="rounded-lg shadow-sm" @click="goToCreatePost">
        <template #icon>
          <n-icon>
            <IconAddOutline />
          </n-icon>
        </template>
        Đăng tin mới
      </n-button>

      <!-- Menu Điều Hướng -->
      <n-menu
        v-model:value="activeKey"
        :options="menuOptions"
        @update:value="handleMenuSelect"
      />
    </div>

    <!-- Thông tin phiên làm việc hoặc chân trang Sidebar -->
    <div class="px-6 py-4 border-t border-gray-100 flex items-center gap-3">
      <n-avatar round size="medium" class="bg-gray-100">
        <template #icon>
          <n-icon><IconUser /></n-icon>
        </template>
      </n-avatar>
      <div class="flex flex-col">
        <span class="text-sm font-semibold text-gray-800">Quản lý</span>
        <span class="text-xs text-gray-400">Trang quản trị viên</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { h, ref, type Component } from "vue";
import { NIcon, type MenuOption } from "naive-ui";
import IconBuilding from "~/icons/IconBuilding.vue";
import IconUser from "~/icons/IconUser.vue";
import IconAddOutline from "~/icons/IconAddOutline.vue";

const props = defineProps<{
  activeKey: "posts" | "customers";
}>();

const activeKey = ref<string>(props.activeKey);

// Render icon cho Naive UI Menu
function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

// Danh sách các mục quản lý trong Sidebar
const menuOptions: MenuOption[] = [
  {
    label: "Quản lý bài viết",
    key: "posts",
    icon: renderIcon(IconBuilding),
  },
  {
    label: "Quản lý khách hàng",
    key: "customers",
    icon: renderIcon(IconUser),
  },
];

// Chuyển hướng tới trang đăng tin
const goToCreatePost = () => {
  navigateTo("/nguoi-ban/dang-tin");
};

// Xử lý khi click vào menu điều hướng sang trang tương ứng
const handleMenuSelect = (key: string) => {
  activeKey.value = key;
  if (key === "posts") {
    navigateTo("/nguoi-ban/quan-ly-tin-dang");
  } else if (key === "customers") {
    navigateTo("/nguoi-ban/quan-ly-khach-hang");
  }
};
</script>
