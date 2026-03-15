<template>
  <DataTableRowActions
    :entity-name="t('globals.terms.agent', 1)"
    :delete-description="t('admin.agent.deleteConfirmation')"
    @edit="editUser(props.user.id)"
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
const { emitRefresh } = useAdminListRefresh('agent', () => {})
const { showErrorToast } = useAdminErrorToast()

const props = defineProps({
  user: { type: Object, required: true, default: () => ({ id: '' }) }
})

function editUser(id) {
  router.push({ name: 'edit-agent', params: { id } })
}

async function handleDelete() {
  try {
    await api.deleteUser(props.user.id)
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}
</script>
