<template>
  <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false"
    flex-height class="flex-1 min-h-0" :scroll-x="980" :row-props="rowProps" />
</template>

<script setup lang="ts">
import { computed, h } from 'vue'
import { NButton, NIcon, NImage, NSpace, NTag, type DataTableColumns } from 'naive-ui'
import type { Deposit, DepositActorRole } from '~/types/deposit'
import { formatSlot, formatVnd } from '~/utils/deposit'
import { fromNow } from '~/utils/format'
import DepositStatusTag from '~/components/deposit/DepositStatusTag.vue'
import IconEyeOutline from '~/icons/IconEyeOutline.vue'

const props = defineProps<{
  items: Deposit[]
  loading: boolean
  role: DepositActorRole
}>()

const emit = defineEmits<{
  view: [id: number]
}>()

// Cột hiển thị đối tác: khách nhìn thấy môi giới, môi giới/admin nhìn thấy khách
const partnerTitle = computed(() => (props.role === 'CUSTOMER' ? 'Môi giới' : 'Khách hàng'))

const columns = computed<DataTableColumns<Deposit>>(() => [
  {
    title: 'Mã đơn',
    key: 'id',
    width: 90,
    render: (row) => h('span', { class: 'text-sm font-medium text-gray-700' }, `#${row.id}`),
  },
  {
    title: 'Bất động sản',
    key: 'real_estate',
    width: 280,
    render(row) {
      return h('div', { class: 'flex items-center gap-3 py-1' }, [
        h(NImage, {
          src: row.real_estate_thumbnail || 'https://picsum.photos/200/150?random=deposit',
          width: 56,
          height: 42,
          class: 'rounded object-cover border border-gray-100 flex-shrink-0',
          previewDisabled: true,
        }),
        h('div', { class: 'flex flex-col gap-0.5 min-w-0' }, [
          h('span', { class: 'font-medium text-gray-800 text-sm truncate' }, row.real_estate_title || 'Bất động sản'),
          h('span', { class: 'text-xs text-gray-400 truncate' }, row.real_estate_address || '—'),
        ]),
      ])
    },
  },
  {
    title: partnerTitle.value,
    key: 'partner',
    width: 170,
    render(row) {
      const name = props.role === 'CUSTOMER' ? row.broker_name : row.customer_name
      const phone = props.role === 'CUSTOMER' ? row.broker_phone : row.customer_phone
      return h('div', { class: 'flex flex-col gap-0.5' }, [
        h('span', { class: 'text-sm text-gray-700' }, name || '—'),
        h('span', { class: 'text-xs text-gray-400' }, phone || '—'),
      ])
    },
  },
  {
    title: 'Lịch xem',
    key: 'viewing',
    width: 180,
    render: (row) => h('span', { class: 'text-sm text-gray-600 whitespace-nowrap' }, formatSlot(row)),
  },
  {
    title: 'Tiền cọc',
    key: 'amount',
    width: 130,
    render: (row) =>
      h('div', { class: 'flex flex-col gap-0.5' }, [
        h('span', { class: 'text-sm font-semibold text-gray-800 whitespace-nowrap' }, formatVnd(row.amount)),
        h('span', { class: 'text-xs text-gray-400 whitespace-nowrap' }, `Phí: ${formatVnd(row.broker_fee)}`),
      ]),
  },
  {
    title: 'Trạng thái',
    key: 'status',
    width: 170,
    render: (row) =>
      h('div', { class: 'flex flex-col gap-1 items-start' }, [
        h(DepositStatusTag, { status: row.status }),
        row.has_dispute ? h(NTag, { size: 'tiny', type: 'warning', bordered: false }, { default: () => 'Có tranh chấp' }) : null,
      ]),
  },
  {
    title: 'Cập nhật',
    key: 'updated_at',
    width: 120,
    render: (row) => h('span', { class: 'text-xs text-gray-400 whitespace-nowrap' }, fromNow(row.updated_at)),
  },
  {
    title: 'Hành động',
    key: 'actions',
    width: 100,
    align: 'center',
    render: (row) =>
      h(
        NSpace,
        { justify: 'center', size: 'small' },
        {
          default: () => [
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                type: 'info',
                onClick: (event: MouseEvent) => {
                  // Chặn bubble để không mở chi tiết 2 lần
                  event.stopPropagation()
                  emit('view', row.id)
                },
              },
              { icon: () => h(NIcon, null, { default: () => h(IconEyeOutline) }) },
            ),
          ],
        },
      ),
  },
])

// Click vào dòng cũng mở chi tiết
function rowProps(row: Deposit) {
  return {
    style: 'cursor: pointer',
    onClick: () => emit('view', row.id),
  }
}
</script>
