<template>
  <div class="w-64 border-r border-gray-200 bg-white h-[calc(100vh-65px)] flex flex-col justify-between py-4 select-none">
    <div class="flex flex-col gap-6 px-4">
      <!-- Nút Đăng Tin Mới -->
      <n-button type="error" size="large" block class="rounded-lg" @click="goToCreatePost">
        <template #icon>
          <n-icon>
            <IconAddOutline />
          </n-icon>
        </template>
        Đăng tin mới
      </n-button>

      <!-- Menu Điều Hướng -->
      <div>
        <p class="px-2 mb-2 text-xs font-semibold uppercase tracking-wide text-gray-400">Quản lý</p>
        <n-menu
          v-model:value="activeKey"
          :options="menuOptions"
          :indent="18"
          @update:value="handleMenuSelect"
        />
      </div>
    </div>

    <!-- Thông tin phiên làm việc hoặc chân trang Sidebar -->
    <div class="px-4 py-4 border-t border-gray-100 flex items-center gap-3">
      <n-avatar round size="medium" class="bg-emerald-50 text-emerald-500">
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
import IconCreateOutline from "~/icons/IconCreateOutline.vue";
import IconUser from "~/icons/IconUser.vue";
import IconHeart from "~/icons/IconHeart.vue";
import IconAddOutline from "~/icons/IconAddOutline.vue";

const props = defineProps<{
  activeKey: "projects" | "posts" | "customers" | "favorites";
}>();

const activeKey = ref<string>(props.activeKey);

// Render icon cho Naive UI Menu
function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

// Danh sách các mục quản lý trong Sidebar
const menuOptions: MenuOption[] = [
  {
    label: "Quản lý dự án",
    key: "projects",
    icon: renderIcon(IconBuilding),
  },
  {
    label: "Quản lý bài viết",
    key: "posts",
    icon: renderIcon(IconCreateOutline),
  },
  {
    label: "Quản lý khách hàng",
    key: "customers",
    icon: renderIcon(IconUser),
  },
  {
    label: "Danh mục yêu thích",
    key: "favorites",
    icon: renderIcon(IconHeart),
  },
];

// Chuyển hướng tới trang đăng tin
const goToCreatePost = () => {
  navigateTo("/nguoi-ban/dang-tin");
};

// Xử lý khi click vào menu điều hướng sang trang tương ứng
const handleMenuSelect = (key: string) => {
  activeKey.value = key;
  if (key === "projects") {
    navigateTo("/nguoi-ban/quan-ly-du-an");
  } else if (key === "posts") {
    navigateTo("/nguoi-ban/quan-ly-tin-dang");
  } else if (key === "customers") {
    navigateTo("/nguoi-ban/quan-ly-khach-hang");
  } else if (key === "favorites") {
    navigateTo("/nguoi-ban/quan-ly-yeu-thich");
  }
};
</script>
