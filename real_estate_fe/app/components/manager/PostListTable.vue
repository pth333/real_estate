<!-- Component danh sách bài viết dành riêng cho Manager -->
<template>
  <div class="flex-1 bg-white p-6 rounded-xl border border-gray-200 flex flex-col gap-4 h-full min-h-0">
    <!-- Header: Bộ lọc và tìm kiếm -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 flex-shrink-0">
      <div class="flex flex-wrap gap-2">
        <n-radio-group>
          <n-radio-button value="all">Tất cả ({{ totalPosts }})</n-radio-button>
        </n-radio-group>
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
      <n-pagination v-model:page="currentPage" :page-count="totalPages" :page-size="pageSize"
        @update:page="handlePageChange" />
    </div>

    <!-- Modal confirm xóa -->
    <n-modal v-model:show="showDeleteModal" :mask-closable="false">
      <div class="bg-white rounded-xl p-6 w-[420px] flex flex-col gap-4">
        <!-- Header -->
        <div class="flex items-center gap-3">
          <n-icon size="22" color="#ef4444">
            <IconCloseOutline />
          </n-icon>
          <span class="font-semibold text-gray-800 text-base">Xác nhận xóa bài đăng</span>
        </div>

        <!-- Content -->
        <p class="text-gray-500 text-sm">
          Bạn có chắc muốn xóa bài đăng này không? Hành động này <strong>không thể hoàn tác</strong>.
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
  useMessage, type DataTableColumns
} from "naive-ui";
import type { IManagerPostItem, IManagerPostListResponse } from "~/types/manager";
import IconSearch from "~/icons/IconSearch.vue";
import IconEyeOutline from "~/icons/IconEyeOutline.vue";
import IconCreateOutline from "~/icons/IconCreateOutline.vue";
import IconCloseOutline from "~/icons/IconCloseOutline.vue";

const { $api } = useNuxtApp();

// State quản lý bài viết từ API
const loading = ref<boolean>(false);
const posts = ref<IManagerPostItem[]>([]);
const totalPosts = ref<number>(0);
const pageSize = ref<number>(10);
const currentPage = ref<number>(1);
const searchQuery = ref<string>("");

// State quản lý modal xóa
const showDeleteModal = ref<boolean>(false);
const deleteLoading = ref<boolean>(false);
const pendingDeleteId = ref<number | null>(null);

// Tính tổng số trang dựa trên API total
const totalPages = computed(() => {
  return Math.ceil(totalPosts.value / pageSize.value) || 1;
});

// Hàm trực tiếp gọi API lấy danh sách bài viết bằng $api
const fetchPostsData = async () => {
  loading.value = true;
  try {
    const res = await $api.get<{ data: IManagerPostListResponse }>("/manager/posts", {
      params: {
        search: searchQuery.value,
        page: currentPage.value,
        size: pageSize.value,
      },
    });

    if (res) {
      posts.value = res.data.posts;
      totalPosts.value = res.data.total;
    }
  } catch (error: any) {
    window.message?.error("Lỗi khi tải danh sách bài viết: " + (error?.message || "Lỗi máy chủ"));
  } finally {
    loading.value = false;
  }
};

// Gọi dữ liệu khi mount
onMounted(() => {
  fetchPostsData();
});

// Kích hoạt khi tìm kiếm
let searchTimeout: any = null;
const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    currentPage.value = 1;
    fetchPostsData();
  }, 300);
};

// Kích hoạt khi đổi trang
const handlePageChange = (page: number) => {
  currentPage.value = page;
  fetchPostsData();
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
        h("div", { class: "flex flex-col gap-1 min-w-0" }, [
          // Tooltip hiện full title khi hover
          h(
            NTooltip,
            { trigger: "hover", placement: "top" },
            {
              trigger: () =>
                h(
                  "span",
                  { class: "font-medium text-gray-800 text-sm truncate block max-w-[200px] cursor-default" },
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
            h(NButton, { size: "small", quaternary: true, onClick: () => handleViewPost(row.slug) },
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