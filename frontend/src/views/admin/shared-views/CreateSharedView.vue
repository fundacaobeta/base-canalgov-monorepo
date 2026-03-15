<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <SharedViewForm :submitForm="onSubmit" :isLoading="formLoading" />
</template>

<script setup>
import { ref } from 'vue'
import SharedViewForm from '@/features/admin/shared-views/SharedViewForm.vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { useRouter } from 'vue-router'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import { useSharedViewStore } from '@/stores/sharedView'
import api from '@/api'

const router = useRouter()
const { t } = useI18n()
const sharedViewStore = useSharedViewStore()
const formLoading = ref(false)
const { showErrorToast, showSuccessToast } = useAdminErrorToast()
const breadcrumbLinks = [
  { path: 'shared-view-list', label: t('globals.terms.sharedView', 2) },
  { path: '', label: t('globals.messages.new', { name: t('globals.terms.sharedView').toLowerCase() }) }
]

const onSubmit = async (values) => {
  formLoading.value = true
  try {
    await api.createSharedView(values)
    await sharedViewStore.refresh()
    showSuccessToast(t('globals.messages.createdSuccessfully', { name: t('globals.terms.sharedView') }))
    router.push({ name: 'shared-view-list' })
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
}
</script>
