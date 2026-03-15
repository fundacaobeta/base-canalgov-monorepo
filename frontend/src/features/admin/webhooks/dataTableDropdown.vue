<template>
  <DataTableRowActions
    :entity-name="t('globals.terms.webhook')"
    :show-edit="false"
    @delete="handleDelete"
  >
    <template #extra-items>
      <DropdownMenuItem :as-child="true">
        <RouterLink :to="{ name: 'edit-webhook', params: { id: props.webhook.id } }">
          {{ t('globals.messages.edit') }}
        </RouterLink>
      </DropdownMenuItem>
      <DropdownMenuItem @click="handleToggle">
        {{ props.webhook.is_active ? t('globals.messages.disable') : t('globals.messages.enable') }}
      </DropdownMenuItem>
      <DropdownMenuItem @click="handleTest">
        {{ t('globals.messages.send', { name: t('globals.terms.test').toLowerCase() }) }}
      </DropdownMenuItem>
      <DropdownMenuSeparator />
    </template>
  </DataTableRowActions>
</template>

<script setup>
import { DropdownMenuItem, DropdownMenuSeparator } from '@/components/ui/dropdown-menu'
import DataTableRowActions from '@/components/admin/DataTableRowActions.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const emitter = useEmitter()
const { emitRefresh } = useAdminListRefresh('webhook', () => {})
const { showErrorToast, showSuccessToast } = useAdminErrorToast()

const props = defineProps({
  webhook: { type: Object, required: true, default: () => ({ id: '', name: '', is_active: false }) }
})

async function handleDelete() {
  try {
    await api.deleteWebhook(props.webhook.id)
    showSuccessToast(t('globals.messages.deletedSuccessfully', { name: t('globals.terms.webhook') }))
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}

async function handleToggle() {
  try {
    await api.toggleWebhook(props.webhook.id)
    showSuccessToast(t('globals.messages.updatedSuccessfully', { name: t('globals.terms.webhook') }))
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}

async function handleTest() {
  try {
    await api.testWebhook(props.webhook.id)
    showSuccessToast(t('globals.messages.sentSuccessfully', { name: t('globals.terms.webhook') }))
  } catch (error) {
    showErrorToast(error)
  }
}
</script>
