<template>
  <DataTableRowActions
    :entity-name="t('globals.terms.businessHour').toLowerCase()"
    @edit="edit(props.role.id)"
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
const { emitRefresh } = useAdminListRefresh('business_hours', () => {})
const { showErrorToast } = useAdminErrorToast()

const props = defineProps({
  role: { type: Object, required: true, default: () => ({ id: '' }) }
})

function edit(id) {
  router.push({ name: 'edit-business-hours', params: { id } })
}

async function handleDelete() {
  try {
    await api.deleteBusinessHours(props.role.id)
    emitRefresh()
  } catch (error) {
    showErrorToast(error)
  }
}
</script>
