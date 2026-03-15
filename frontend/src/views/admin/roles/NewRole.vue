<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <RoleForm :initial-values="{}" :submitForm="submitForm" :isLoading="formLoading" />
</template>

<script setup>
import { ref } from 'vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import RoleForm from '@/features/admin/roles/RoleForm.vue'
import api from '@/api'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { t } = useI18n()
const router = useRouter()
const formLoading = ref(false)
const { showErrorToast, showSuccessToast } = useAdminErrorToast()
const breadcrumbLinks = [
  { path: 'role-list', label: t('globals.terms.role', 2) },
  { path: '', label: t('globals.messages.new', { name: t('globals.terms.role') }) }
]

const submitForm = async (values) => {
  formLoading.value = true
  try {
    await api.createRole(values)
    showSuccessToast(t('globals.messages.createdSuccessfully', { name: t('globals.terms.role') }))
    router.push({ name: 'role-list' })
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
}
</script>
