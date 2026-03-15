<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <SharedViewForm
    v-if="!isLoading"
    :initialValues="sharedView"
    :submitForm="submitForm"
    :isLoading="formLoading"
  />
</template>

<script setup>
import { onMounted, ref } from 'vue'
import api from '@/api'
import SharedViewForm from '@/features/admin/shared-views/SharedViewForm.vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import { useSharedViewStore } from '@/stores/sharedView'

const sharedView = ref({})
const { t } = useI18n()
const isLoading = ref(false)
const formLoading = ref(false)
const sharedViewStore = useSharedViewStore()
const { showErrorToast, showSuccessToast } = useAdminErrorToast()

const props = defineProps({ id: { type: String, required: true } })

const breadcrumbLinks = [
  { path: 'shared-view-list', label: t('globals.terms.sharedView', 2) },
  { path: '', label: t('globals.messages.edit', { name: t('globals.terms.sharedView') }) }
]

const submitForm = async (values) => {
  formLoading.value = true
  try {
    await api.updateSharedView(sharedView.value.id, values)
    await sharedViewStore.refresh()
    showSuccessToast(t('globals.messages.updatedSuccessfully', { name: t('globals.terms.sharedView') }))
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
}

onMounted(async () => {
  isLoading.value = true
  try {
    const resp = await api.getSharedView(props.id)
    sharedView.value = resp.data.data
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
})
</script>
