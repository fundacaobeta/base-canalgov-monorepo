import { h } from 'vue'
import { RouterLink } from 'vue-router'
import TeamDataTableDropdown from '@/features/admin/teams/TeamDataTableDropdown.vue'
import DateTimeMeta from '@/components/datetime/DateTimeMeta.vue'

export const columns = [
  {
    accessorKey: 'name',
    header: function () {
      return h('div', { class: 'text-center' }, 'Nome')
    },
    cell: function ({ row }) {
      return h('div', { class: 'text-center' },
        h(RouterLink,
          {
            to: { name: 'edit-team', params: { id: row.original.id } },
            class: 'text-primary hover:underline'
          },
          () => row.getValue('name')
        )
      )
    }
  },
  {
    accessorKey: 'created_at',
    header: function () {
      return h('div', { class: 'text-center' }, 'Criado em')
    },
    cell: function ({ row }) {
      return h(DateTimeMeta, { value: row.getValue('created_at'), centered: true, compact: true })
    }
  },
  {
    accessorKey: 'updated_at',
    header: function () {
      return h('div', { class: 'text-center' }, 'Atualizado em')
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
      const team = row.original
      return h(
        'div',
        { class: 'relative' },
        h(TeamDataTableDropdown, {
          team
        })
      )
    }
  }
]
