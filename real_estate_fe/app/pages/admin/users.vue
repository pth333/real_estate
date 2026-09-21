<template>
  <div class="flex flex-col gap-5">
    <n-result v-if="!isAdmin" status="403" title="Không có quyền truy cập"
      description="Chỉ tài khoản quản trị viên mới quản lý được người dùng và phân quyền" />

    <template v-else>
      <div>
        <p class="text-sm text-gray-500">
          Gán role cho tài khoản. Quyền của user được suy ra từ các role đang gán.
        </p>
      </div>

      <div class="flex flex-col gap-4 rounded-xl border border-gray-200 bg-white p-5">
        <!-- Tìm kiếm theo tên / email / số điện thoại (debounce 300ms) -->
        <div class="w-full md:w-80">
          <n-input v-model:value="searchQuery" placeholder="Tìm theo tên, email hoặc số điện thoại..." clearable
            @input="handleSearch">
            <template #prefix>
              <n-icon>
                <IconSearch />
              </n-icon>
            </template>
          </n-input>
        </div>

        <n-spin :show="loading">
          <n-data-table :columns="columns" :data="users" :bordered="false" :single-line="false" :scroll-x="900" />

          <n-empty v-if="!users.length && !loading" description="Không tìm thấy người dùng nào" class="py-16" />
        </n-spin>

        <div class="flex justify-end">
          <Pagination :current-page="page" :total-pages="totalPages" @page-change="goToPage" />
        </div>
      </div>
    </template>

    <!-- Modal gán role -->
    <n-modal v-model:show="showAssign" :mask-closable="false">
      <div class="w-[520px] max-w-[94vw] bg-white rounded-xl p-6 flex flex-col gap-4">
        <span class="font-semibold text-gray-800">Phân quyền người dùng</span>

        <div class="rounded-lg bg-gray-50 p-3 flex flex-col">
          <span class="text-sm font-medium text-gray-800">{{ editingUser?.name }}</span>
          <span class="text-xs text-gray-500">{{ editingUser?.email }}</span>
        </div>

        <n-form-item label="Role được gán" :show-feedback="false">
          <n-select v-model:value="selectedRoles" multiple :options="roleOptions" :loading="roleLoading"
            placeholder="Chọn role cho người dùng" />
        </n-form-item>

        <div class="flex justify-end gap-2">
          <n-button @click="showAssign = false">Huỷ</n-button>
          <n-button type="primary" :loading="saving" @click="submitRoles">Lưu</n-button>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { NButton, NIcon, NSpace, NTag, type DataTableColumns, type SelectOption } from 'naive-ui'
import type { AdminUserResponse, RoleResponse } from '~/types/rbac'
import { useRbacService } from '~/services/rbac.service'
import { useAuthStore } from '~/stores/auth'
import { formatDate } from '~/utils/format'
import IconSearch from '~/icons/IconSearch.vue'

definePageMeta({
  layout: 'admin',
  requiresPermission: ['admin.user.list'],
})

const authStore = useAuthStore()
const rbacService = useRbacService()

const isAdmin = computed(() => authStore.user?.roles?.includes('ADMIN') ?? false)

const users = ref<AdminUserResponse[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const loading = ref(false)
const searchQuery = ref('')

const roles = ref<RoleResponse[]>([])
const roleLoading = ref(false)

const showAssign = ref(false)
const editingUser = ref<AdminUserResponse | null>(null)
const selectedRoles = ref<string[]>([])
const saving = ref(false)

const totalPages = computed(() => Math.ceil(total.value / size.value) || 1)

// Options cho n-select: hiển thị tên role, giá trị là code role
const roleOptions = computed<SelectOption[]>(() =>
  roles.value.map((role) => ({ label: role.name, value: role.code })),
)

// Màu badge theo role để bảng dễ đọc
function roleTagType(code: string): 'default' | 'info' | 'success' | 'warning' {
  if (code === 'ADMIN') return 'warning'
  if (code === 'BROKER') return 'info'
  if (code === 'CUSTOMER') return 'success'
  return 'default'
}

const columns: DataTableColumns<AdminUserResponse> = [
  {
    title: 'Người dùng',
    key: 'name',
    minWidth: 220,
    render(row) {
      return h('div', { class: 'flex flex-col gap-0.5 py-1' }, [
        h('span', { class: 'font-medium text-gray-800 text-sm truncate' }, row.name),
        h('span', { class: 'text-xs text-gray-400 truncate' }, row.email),
      ])
    },
  },
  {
    title: 'Số điện thoại',
    key: 'phone',
    width: 150,
    render(row) {
      return h('span', { class: 'text-gray-600 text-sm whitespace-nowrap' }, row.phone || '—')
    },
  },
  {
    title: 'Role',
    key: 'roles',
    minWidth: 220,
    render(row) {
      // User chưa được gán role nào
      if (!row.roles.length) {
        return h(NTag, { size: 'small', bordered: false }, { default: () => 'Chưa gán role' })
      }
      return h(
        NSpace,
        { size: 'small', align: 'center' },
        {
          default: () =>
            row.roles.map((code) =>
              h(NTag, { size: 'small', bordered: false, type: roleTagType(code) }, { default: () => code }),
            ),
        },
      )
    },
  },
  {
    title: 'Ngày tạo',
    key: 'created_at',
    width: 130,
    render(row) {
      return h('span', { class: 'text-gray-500 text-sm whitespace-nowrap' }, formatDate(row.created_at))
    },
  },
  {
    title: 'Hành động',
    key: 'actions',
    width: 130,
    align: 'center',
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'primary', onClick: () => openAssign(row) },
        { default: () => 'Phân quyền' },
      )
    },
  },
]

async function fetchUsers() {
  loading.value = true
  try {
    const result = await rbacService.getUsers({
      search: searchQuery.value || undefined,
      page: page.value,
      size: size.value,
    })
    users.value = result.items
    total.value = result.total
  } catch {
    users.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

async function fetchRoles() {
  roleLoading.value = true
  try {
    roles.value = await rbacService.getRoles()
  } catch {
    roles.value = []
  } finally {
    roleLoading.value = false
  }
}

function goToPage(nextPage: number) {
  if (nextPage < 1 || nextPage > totalPages.value) return
  page.value = nextPage
  fetchUsers()
}

// Debounce 300ms cho ô tìm kiếm
let searchTimeout: ReturnType<typeof setTimeout> | null = null
function handleSearch() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchUsers()
  }, 300)
}

function openAssign(user: AdminUserResponse) {
  editingUser.value = user
  selectedRoles.value = [...user.roles]
  showAssign.value = true
}

async function submitRoles() {
  if (!editingUser.value) return
  saving.value = true
  try {
    await rbacService.setUserRoles(editingUser.value.id, selectedRoles.value)
    window.message?.success('Đã cập nhật role cho người dùng')
    showAssign.value = false
    await fetchUsers()
  } catch {
    // Lỗi 400 (role không tồn tại / mảng rỗng) đã được $api hiện message — giữ modal để admin sửa lại
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  if (!isAdmin.value) return
  fetchRoles()
  fetchUsers()
})
</script>
