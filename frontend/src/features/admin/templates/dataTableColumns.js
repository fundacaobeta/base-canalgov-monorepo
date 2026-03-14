import { h } from 'vue'
import { RouterLink } from 'vue-router'
import dropdown from './dataTableDropdown.vue'
import DateTimeMeta from '@/components/datetime/DateTimeMeta.vue'

const createNameColumn = (t) => ({
  accessorKey: 'name',
  header: function () {
    return h('div', { class: 'text-center' }, t('globals.terms.name'))
  },
  cell: function ({ row }) {
    return h('div', { class: 'text-center' },
      h(RouterLink,
        {
          to: { name: 'edit-template', params: { id: row.original.id } },
          class: 'text-primary hover:underline'
        },
        () => row.getValue('name')
      )
    )
  }
})

const createScopeColumn = () => ({
  accessorKey: 'team_name',
  header: function () {
    return h('div', { class: 'text-center' }, 'Escopo')
  },
  cell: function ({ row }) {
    return h(
      'div',
      { class: 'text-center text-sm' },
      row.original.team_name || 'Global'
    )
  }
})

const createDefaultColumn = (t) => ({
  accessorKey: 'is_default',
  header: () => h('div', { class: 'text-center' }, t('globals.terms.default')),
  cell: ({ row }) => {
    const isDefault = row.getValue('is_default')

    return h('div', { class: 'text-center' }, [
      h('input', {
        type: 'checkbox',
        checked: isDefault,
        disabled: true
      })
    ])
  }
})

const createUpdatedAtColumn = (t) => ({
  accessorKey: 'updated_at',
  header: function () {
    return h('div', { class: 'text-center' }, t('globals.terms.updatedAt'))
  },
  cell: function ({ row }) {
    return h(DateTimeMeta, { value: row.getValue('updated_at'), centered: true, compact: true })
  }
})

const createActionsColumn = () => ({
  id: 'actions',
  enableHiding: false,
  enableSorting: false,
  cell: ({ row }) => {
    const template = row.original
    return h(
      'div',
      { class: 'relative' },
      h(dropdown, {
        template
      })
    )
  }
})

export const createResponseTemplateColumns = (t) => [
  createNameColumn(t),
  createScopeColumn(),
  createDefaultColumn(t),
  createUpdatedAtColumn(t),
  createActionsColumn()
]

export const createOutgoingEmailTableColumns = (t) => [
  createNameColumn(t),
  createDefaultColumn(t),
  createUpdatedAtColumn(t),
  createActionsColumn()
]

export const createEmailNotificationTableColumns = (t) => [
  createNameColumn(t),
  createUpdatedAtColumn(t),
  createActionsColumn()
]
