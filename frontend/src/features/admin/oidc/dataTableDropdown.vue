<template>
  <DataTableRowActions
    :entity-name="t('globals.terms.sso')"
    :show-edit="false"
    @delete="handleDelete"
  >
    <template #extra-items>
      <DropdownMenuItem :as-child="true">
        <RouterLink :to="{ name: 'edit-sso', params: { id: props.role.id } }">
          {{ t('globals.messages.edit') }}
        </RouterLink>
      </DropdownMenuItem>
    </template>
  </DataTableRowActions>
</template>

<script setup>
import { DropdownMenuItem } from '@/components/ui/dropdown-menu'
import DataTableRowActions from '@/components/admin/DataTableRowActions.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const { emitRefresh } = useAdminListRefresh('oidc', () => {})
const { showErrorToast } = useAdminErrorToast()

const props = defineProps({
  role: { type: Object, required: true, default: () => ({ id: '' }) }
})

async function handleDelete() {
  try {
    await api.deleteOIDC(props.role.id)
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}
</script>
