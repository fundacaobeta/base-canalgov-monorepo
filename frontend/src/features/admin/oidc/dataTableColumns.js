import { h } from 'vue'
import dropdown from './dataTableDropdown.vue'
import DateTimeMeta from '@/components/datetime/DateTimeMeta.vue'

export const createColumns = (t) => [
  {
    accessorKey: 'name',
    header: function () {
      return h('div', { class: 'text-center' }, t('globals.terms.name'))
    },
    cell: function ({ row }) {
      return h('div', { class: 'text-center' }, row.getValue('name'))
    }
  },
  {
    accessorKey: 'provider',
    header: function () {
      return h('div', { class: 'text-center' }, t('globals.terms.provider'))
    },
    cell: function ({ row }) {
      return h('div', { class: 'text-center' }, row.getValue('provider'))
    }
  },
  {
    accessorKey: 'enabled',
    header: () => h('div', { class: 'text-center' }, t('globals.terms.enabled')),
    cell: ({ row }) => {
      const enabled = row.getValue('enabled')
      return h('div', { class: 'text-center' }, enabled ? t('globals.messages.yes') : t('globals.messages.no'))
    }
  },
  {
    accessorKey: 'created_at',
    header: function () {
      return h('div', { class: 'text-center' }, t('globals.terms.createdAt'))
    },
    cell: function ({ row }) {
      return h(DateTimeMeta, { value: row.getValue('created_at'), centered: true, compact: true })
    }
  },
  {
    accessorKey: 'updated_at',
    header: function () {
      return h('div', { class: 'text-center' }, t('globals.terms.updatedAt'))
    },
    cell: function ({ row }) {
      return h(DateTimeMeta, { value: row.getValue('updated_at'), centered: true, compact: true })
    }
  },
  {
    id: 'actions',
    enableHiding: false,
    enableSorting: false,
    cell: ({ row }) => {
      const role = row.original
      return h(
        'div',
        { class: 'relative' },
        h(dropdown, {
          role
        })
      )
    }
  }
]
