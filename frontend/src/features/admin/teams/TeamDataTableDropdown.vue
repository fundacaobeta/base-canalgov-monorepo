<template>
  <DataTableRowActions
    :entity-name="t('globals.terms.team')"
    :delete-title="t('globals.messages.deleteTeam')"
    :delete-description="t('admin.team.deleteConfirmation')"
    @edit="editTeam(props.team.id)"
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
const { emitRefresh } = useAdminListRefresh('team', () => {})
const { showErrorToast } = useAdminErrorToast()

const props = defineProps({
  team: { type: Object, required: true, default: () => ({ id: '' }) }
})

function editTeam(id) {
  router.push({ name: 'edit-team', params: { id } })
}

async function handleDelete() {
  try {
    await api.deleteTeam(props.team.id)
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}
</script>
