<template>
  <!-- Lọc theo loại bất động sản (Bang filter — hiển thị mọi trang) -->
  <n-popover placement="bottom-start" trigger="click" :show="showPopover" @update:show="showPopover = $event"
    :style="{ padding: 0 }">
    <template #trigger>
      <n-button ghost size="small">
        <template #icon>
          <IconChevronDownOutline class="h-4 w-4" />
        </template>
        {{ buttonLabel }}
      </n-button>
    </template>

    <div class="w-72 p-3">
      <n-input v-model:value="keyword" placeholder="Tìm loại bất động sản..." clearable size="small" class="mb-2" />
      <n-tree :key="keyword ? 'filtered' : 'all'" :data="treeData" :pattern="keyword" :filter-method="treeFilter"
        :default-expand-all="!!keyword" block-line block-node selectable :selected-keys="selectedKeys"
        class="max-h-72 overflow-y-auto" @update:selected-keys="handleSelect" />
    </div>
  </n-popover>
</template>

<script setup lang="ts">
import { useRealEstateStore } from "~/stores/real_estate";
import type { TreeOption } from "naive-ui";
import type { Category } from "~/types/menu";
import { useFilterStore } from '~/stores/filter';

const showPopover = ref(false);
const keyword = ref("");
const realEstateStore = useRealEstateStore();
const filterStore = useFilterStore();


type CategoryTreeOption = TreeOption & { slug?: string };

const treeData = computed<CategoryTreeOption[]>(() =>
  (window.menu?.settings?.categories ?? []).map(toTreeOption),
);

function toTreeOption(cat: Category): CategoryTreeOption {
  return {
    key: cat.Slug ?? String(cat.ID),
    label: cat.Name,
    slug: cat.Slug,
    children: (cat.children ?? []).map(toTreeOption),
  };
}

function treeFilter(pattern: string, node: CategoryTreeOption): boolean {
  return (node.label ?? "").toLowerCase().includes(pattern.trim().toLowerCase());
}

const buttonLabel = computed(() => {
  const active = findNodeBySlug(treeData.value, realEstateStore.categorySlug);
  return active?.label ?? "Loại BĐS";
});

const selectedKeys = computed<Array<string | number>>(() => {
  const key = findKeyBySlug(treeData.value, realEstateStore.categorySlug);
  return key != null ? [key] : [];
});

function findKeyBySlug(options: CategoryTreeOption[], slug?: string): string | number | null {
  if (!slug) return null;
  for (const opt of options) {
    if (opt.slug === slug && opt.key != null) return opt.key;
    const child = findKeyBySlug((opt.children ?? []) as CategoryTreeOption[], slug);
    if (child != null) return child;
  }
  return null;
}

// ✅ Thêm hàm tìm node theo key
function findNodeByKey(
  options: CategoryTreeOption[],
  key: string | number,
): CategoryTreeOption | null {
  for (const opt of options) {
    if (opt.key === key) return opt;
    const child = findNodeByKey((opt.children ?? []) as CategoryTreeOption[], key);
    if (child) return child;
  }
  return null;
}

// ✅ Thêm hàm tìm node theo slug (dùng cho buttonLabel)
function findNodeBySlug(
  options: CategoryTreeOption[],
  slug?: string,
): CategoryTreeOption | null {
  if (!slug) return null;
  for (const opt of options) {
    if (opt.slug === slug) return opt;
    const child = findNodeBySlug((opt.children ?? []) as CategoryTreeOption[], slug);
    if (child) return child;
  }
  return null;
}

// ✅ Fix: nhận keys thay vì objects
function handleSelect(keys: Array<string | number>) {
  const key = keys[0];
  if (key == null) return;

  const opt = findNodeByKey(treeData.value, key);
  if (!opt?.slug) return;

  realEstateStore.categorySlug = opt.slug;
  showPopover.value = false;
  const url = buildListUrl(
    realEstateStore.categorySlug,
    filterStore.filters,
    filterStore.cityOptions,
  );
  navigateTo(url);
}
</script>