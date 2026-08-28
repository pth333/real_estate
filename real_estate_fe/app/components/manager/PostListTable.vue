<!-- Component danh sách bài viết dành riêng cho Manager -->
<template>
  <div class="flex-1 bg-white p-6 rounded-lg border border-gray-200 flex flex-col gap-4 h-full min-h-0">
    <!-- Header: Thống kê + Bộ lọc và tìm kiếm -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 flex-shrink-0">
      <div class="flex items-center gap-3">
        <n-radio-group>
          <n-radio-button value="all">Tất cả</n-radio-button>
        </n-radio-group>
        <span class="inline-flex items-center gap-1.5 rounded-full bg-red-50 px-3 py-1 text-xs font-medium text-red-500">
          {{ total }} bài đăng
        </span>
      </div>

      <div class="w-full md:w-80">
        <n-input v-model:value="searchQuery" placeholder="Tìm kiếm theo tiêu đề..." clearable @input="handleSearch">
          <template #prefix>
            <n-icon>
              <IconSearch />
            </n-icon>
          </template>
        </n-input>
      </div>
    </div>

    <!-- Data Table hiển thị danh sách bài đăng từ API - Dùng flex-height để lấp đầy khoảng trống -->
    <div class="flex-1 min-h-0 flex flex-col">
      <n-data-table :columns="columns" :data="posts" :loading="loading" :bordered="false" :single-line="false"
        flex-height class="flex-1 min-h-0" :scroll-x="740" />
    </div>

    <!-- Pagination điều hướng phân trang từ API -->
    <div class="flex justify-end flex-shrink-0">
      <Pagination :current-page="managerStore.postsPage" :total-pages="totalPages" @page-change="goToPage" />
    </div>

    <!-- Modal confirm xóa -->
    <n-modal v-model:show="showDeleteModal" :mask-closable="false">
      <div class="bg-white rounded-xl p-6 w-[420px] flex flex-col gap-4">
        <!-- Header -->
        <div class="flex items-center gap-3">
          <div class="flex h-9 w-9 items-center justify-center rounded-full bg-red-50 text-red-500">
            <n-icon size="20">
              <IconCloseOutline />
            </n-icon>
          </div>
          <span class="font-semibold text-gray-800 text-base">Xác nhận xóa bài đăng</span>
        </div>

        <!-- Content -->
        <p class="text-gray-500 text-sm">
          Bạn có chắc muốn xóa bài đăng này không? Hành động này <strong class="text-red-500">không thể hoàn tác</strong>.
        </p>

        <!-- Actions -->
        <div class="flex justify-end gap-2 mt-2">
          <n-button @click="cancelDelete">Hủy</n-button>
          <n-button type="error" :loading="deleteLoading" @click="confirmDelete">Xóa</n-button>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import {
  NButton, NIcon, NSpace, NImage, NTooltip, NModal,
  type DataTableColumns
} from "naive-ui";
import type { IManagerPostItem } from "~/types/manager";
import { useManagerStore } from "~/stores/manager";
import IconSearch from "~/icons/IconSearch.vue";
import IconEyeOutline from "~/icons/IconEyeOutline.vue";
import IconCreateOutline from "~/icons/IconCreateOutline.vue";
import IconCloseOutline from "~/icons/IconCloseOutline.vue";

const { $api } = useNuxtApp();
const managerStore = useManagerStore();

// State loading của bảng (data được cache trong managerStore)
const loading = ref<boolean>(false);

const posts = computed(() => managerStore.posts);
const total = computed(() => managerStore.postsTotal);
const totalPages = computed(() => Math.ceil(managerStore.postsTotal / managerStore.postsSize) || 1);
const searchQuery = computed({
  get: () => managerStore.postsSearch,
  set: (val: string) => { managerStore.postsSearch = val },
});

// Gọi qua store — store tự bỏ qua nếu đã có dữ liệu cho đúng trang/từ khoá
async function fetchPostsData() {
  loading.value = true;
  try {
    await managerStore.fetchPosts({
      page: managerStore.postsPage,
      size: managerStore.postsSize,
      search: managerStore.postsSearch,
    });
  } catch {
    // store đã hiện message lỗi
  } finally {
    loading.value = false;
  }
}

function goToPage(page: number) {
  if (page < 1 || page > totalPages.value) return;
  managerStore.postsPage = page;
  fetchPostsData();
}

// State quản lý modal xóa
const showDeleteModal = ref<boolean>(false);
const deleteLoading = ref<boolean>(false);
const pendingDeleteId = ref<number | null>(null);

// Gọi dữ liệu khi mount (store cache → navigate quay lại không gọi lại API)
onMounted(() => {
  fetchPostsData();
});

// Kích hoạt khi tìm kiếm
let searchTimeout: any = null;
const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    managerStore.postsPage = 1;
    fetchPostsData();
  }, 300);
};

const handleDeletePost = (id: number) => {
  pendingDeleteId.value = id;
  showDeleteModal.value = true;
};

const confirmDelete = async () => {
  if (!pendingDeleteId.value) return;
  deleteLoading.value = true;
  try {
    await $api.delete<{ message: string }>(`/manager/posts/${pendingDeleteId.value}`);
    window.message?.success("Đã xóa bài viết thành công!");
    showDeleteModal.value = false;
    pendingDeleteId.value = null;
    // Đánh dấu cache cũ → fetch lại
    managerStore.invalidatePosts();
    fetchPostsData();
  } catch (error: any) {
    window.message?.error("Lỗi khi xóa bài viết: " + (error?.message || "Lỗi máy chủ"));
  } finally {
    deleteLoading.value = false;
  }
};

// Hủy xóa — đóng modal, clear id
const cancelDelete = () => {
  showDeleteModal.value = false;
  pendingDeleteId.value = null;
};

// Định nghĩa các cột hiển thị trong bảng Naive UI
const columns: DataTableColumns<IManagerPostItem> = [
  {
    title: "Bất động sản",
    key: "property",
    width: 300,
    render(row) {
      return h("div", { class: "flex items-center gap-3 py-1" }, [
        h(NImage, {
          src: row.thumbnail || "https://picsum.photos/200/150?random=default",
          alt: "BĐS Thumbnail",
          width: 60,
          height: 45,
          class: "rounded object-cover border border-gray-100 flex-shrink-0",
          previewDisabled: false,
        }),
        h("div", { class: "flex flex-col gap-1 min-w-0 flex-1" }, [
          // Tooltip hiện full title khi hover
          h(
            NTooltip,
            { trigger: "hover", placement: "top" },
            {
              trigger: () =>
                h(
                  "span",
                  { class: "font-medium text-gray-800 text-sm truncate block w-full cursor-default" },
                  row.title
                ),
              default: () => row.title,
            }
          ),
          h("span", { class: "text-xs text-gray-400 truncate" }, row.type),
        ]),
      ]);
    },
  },
  {
    title: "Giá",
    key: "price",
    width: 100,
    render(row) {
      return h("span", { class: "font-semibold text-red-500 text-sm whitespace-nowrap" }, formatPrice(row.price));
    },
  },
  {
    title: "Diện tích",
    key: "area",
    width: 100,
    render(row) {
      return h("span", { class: "text-gray-700 text-sm whitespace-nowrap" }, `${row.area} m²`);
    },
  },
  {
    title: "Ngày đăng",
    key: "created_at",
    width: 120,
    render(row) {
      return h("span", { class: "text-gray-500 text-sm whitespace-nowrap" }, row.created_at);
    },
  },
  {
    title: "Hành động",
    key: "actions",
    width: 120,
    align: "center",
    render(row) {
      return h(
        NSpace,
        { justify: "center", size: "small", align: "center" },
        {
          default: () => [
            h(NButton, { size: "small", quaternary: true, type: "info", onClick: () => handleViewPost(row.slug) },
              { icon: () => h(NIcon, null, { default: () => h(IconEyeOutline) }) }
            ),
            h(NButton, { size: "small", quaternary: true, type: "info", onClick: () => handleEditPost(row.id) },
              { icon: () => h(NIcon, null, { default: () => h(IconCreateOutline) }) }
            ),
            // Nút xóa — mở modal confirm thay vì xóa thẳng
            h(NButton, { size: "small", quaternary: true, type: "error", onClick: () => handleDeletePost(row.id) },
              { icon: () => h(NIcon, null, { default: () => h(IconCloseOutline) }) }
            ),
          ],
        }
      );
    },
  },
];

// Điều hướng xem chi tiết
const handleViewPost = (slug: string) => {
  navigateTo(`${slug}`);
};

// Điều hướng chỉnh sửa
const handleEditPost = (id: number) => {
  navigateTo(`/nguoi-ban/dang-tin?id=${id}`);
};
</script>
