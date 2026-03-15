<template>
  <DataTableRowActions
    :entity-name="t('globals.terms.customAttribute').toLowerCase()"
    @edit="editCustomAttribute"
    @delete="handleDelete"
  />
</template>

<script setup>
import DataTableRowActions from '@/components/admin/DataTableRowActions.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const emitter = useEmitter()
const { emitRefresh } = useAdminListRefresh('custom-attributes', () => {})
const { showErrorToast } = useAdminErrorToast()

const props = defineProps({
  customAttribute: { type: Object, required: true, default: () => ({ id: '' }) }
})

async function handleDelete() {
  try {
    await api.deleteCustomAttribute(props.customAttribute.id)
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}

const editCustomAttribute = () => {
  emitter.emit(EMITTER_EVENTS.EDIT_MODEL, {
    model: 'custom-attributes',
    data: props.customAttribute
  })
}
</script>
