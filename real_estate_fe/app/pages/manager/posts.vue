<template>
  <!-- Chỉ hiển thị duy nhất nội dung của bảng danh sách bài viết -->
  <PostListTable />
</template>

<script setup lang="ts">
import PostListTable from "~/components/manager/PostListTable.vue";

definePageMeta({
  alias: [
    "/nguoi-ban/quan-ly-tin-dang"
  ],
});
</script>
