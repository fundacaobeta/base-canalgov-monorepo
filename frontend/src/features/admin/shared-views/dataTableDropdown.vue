<template>
  <DataTableRowActions
    :entity-name="t('globals.terms.sharedView')"
    @edit="editSharedView"
    @delete="handleDelete"
  />
</template>

<script setup>
import { useRouter } from 'vue-router'
import DataTableRowActions from '@/components/admin/DataTableRowActions.vue'
import { useAdminListRefresh } from '@/composables/useAdminListRefresh'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useSharedViewStore } from '@/stores/sharedView'
import { useI18n } from 'vue-i18n'
import api from '@/api/index.js'

const { t } = useI18n()
const router = useRouter()
const sharedViewStore = useSharedViewStore()
const { emitRefresh } = useAdminListRefresh('shared-views', () => {})
const { showErrorToast } = useAdminErrorToast()

const props = defineProps({
  sharedView: { type: Object, required: true }
})

const handleDelete = async () => {
  try {
    await api.deleteSharedView(props.sharedView.id)
    await sharedViewStore.refresh()
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}

const editSharedView = () => {
  router.push({ name: 'edit-shared-view', params: { id: props.sharedView.id } })
}
</script>
