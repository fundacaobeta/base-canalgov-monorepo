import { h } from 'vue'
import dataTableDropdown from '@/features/admin/custom-attributes/dataTableDropdown.vue'
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
        accessorKey: 'key',
        header: function () {
            return h('div', { class: 'text-center' }, t('globals.terms.key'))
        },
        cell: function ({ row }) {
            return h('div', { class: 'text-center' }, row.getValue('key'))
        }
    },
    {
        accessorKey: 'data_type',
        header: function () {
            return h('div', { class: 'text-center' }, t('globals.terms.type'))
        },
        cell: function ({ row }) {
            return h('div', { class: 'text-center' }, row.getValue('data_type'))
        }
    },
    {
        accessorKey: 'applies_to',
        header: function () {
            return h('div', { class: 'text-center' }, t('globals.terms.appliesTo'))
        },
        cell: function ({ row }) {
            return h('div', { class: 'text-center' }, row.getValue('applies_to'))
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
            const customAttribute = row.original
            return h(
                'div',
                { class: 'relative' },
                h(dataTableDropdown, {
                    customAttribute
                })
            )
        }
    }
]
