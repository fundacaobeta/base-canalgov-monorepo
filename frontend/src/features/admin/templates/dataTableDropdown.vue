<template>
  <DataTableRowActions
    :entity-name="t('globals.terms.template').toLowerCase()"
    :show-delete="props.template.type !== 'email_notification'"
    @edit="editTemplate(props.template.id)"
    @delete="handleDelete"
  />
</template>

<script setup>
import { useRouter } from 'vue-router'
import DataTableRowActions from '@/components/admin/DataTableRowActions.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const router = useRouter()
const { emitRefresh } = useAdminListRefresh('templates', () => {})
const { showErrorToast } = useAdminErrorToast()

const props = defineProps({
  template: { type: Object, required: true, default: () => ({ id: '' }) }
})

const editTemplate = (id) => {
  router.push({ name: 'edit-template', params: { id } })
}

const handleDelete = async () => {
  try {
    await api.deleteTemplate(props.template.id)
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}
</script>
