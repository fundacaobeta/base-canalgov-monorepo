<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <AgentForm :submitForm="onSubmit" :initialValues="{}" :isNewForm="true" :isLoading="formLoading" />
</template>

<script setup>
import { ref } from 'vue'
import AgentForm from '@/features/admin/agents/AgentForm.vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { useRouter } from 'vue-router'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const { t } = useI18n()
const router = useRouter()
const formLoading = ref(false)
const { showErrorToast, showSuccessToast } = useAdminErrorToast()
const breadcrumbLinks = [
  { path: 'agent-list', label: t('globals.terms.agent', 2) },
  { path: '', label: t('globals.messages.new', { name: t('globals.terms.agent', 1).toLowerCase() }) }
]

const onSubmit = async (values) => {
  formLoading.value = true
  try {
    await api.createUser(values)
    showSuccessToast(t('globals.messages.createdSuccessfully', { name: t('globals.terms.agent', 1) }))
    router.push({ name: 'agent-list' })
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
}
</script>
