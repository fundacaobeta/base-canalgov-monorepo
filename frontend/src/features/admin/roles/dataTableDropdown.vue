<template>
  <DataTableRowActions
    :entity-name="t('globals.terms.role').toLowerCase()"
    :show-delete="!Roles.includes(props.role.name)"
    @edit="editRole(props.role.id)"
    @delete="handleDelete"
  />
</template>

<script setup>
import { useRouter } from 'vue-router'
import DataTableRowActions from '@/components/admin/DataTableRowActions.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { Roles } from '@/constants/user'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const router = useRouter()
const { emitRefresh } = useAdminListRefresh('team', () => {})
const { showErrorToast } = useAdminErrorToast()

const props = defineProps({
  role: { type: Object, required: true, default: () => ({ id: '' }) }
})

function editRole(id) {
  router.push({ name: 'edit-role', params: { id } })
}

async function handleDelete() {
  try {
    await api.deleteRole(props.role.id)
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}
</script>
