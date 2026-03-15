<template>
  <DataTableRowActions
    :entity-name="t('globals.terms.macro')"
    @edit="editMacro"
    @delete="handleDelete"
  />
</template>

<script setup>
import { useRouter } from 'vue-router'
import DataTableRowActions from '@/components/admin/DataTableRowActions.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useMacroStore } from '@/stores/macro'
import { useI18n } from 'vue-i18n'
import api from '@/api/index.js'

const { t } = useI18n()
const router = useRouter()
const macroStore = useMacroStore()
const { emitRefresh } = useAdminListRefresh('macros', () => {})
const { showErrorToast } = useAdminErrorToast()

const props = defineProps({
  macro: { type: Object, required: true }
})

const handleDelete = async () => {
  try {
    await api.deleteMacro(props.macro.id)
    await macroStore.loadMacros(true)
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}

const editMacro = () => {
  router.push({ name: 'edit-macro', params: { id: props.macro.id } })
}
</script>
