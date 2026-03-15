<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <AgentForm v-if="!isLoading" :initialValues="user" :submitForm="submitForm" :isLoading="formLoading" />
</template>

<script setup>
import { onMounted, ref } from 'vue'
import api from '@/api'
import AgentForm from '@/features/admin/agents/AgentForm.vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'

const user = ref({})
const { t } = useI18n()
const isLoading = ref(false)
const formLoading = ref(false)
const { showErrorToast, showSuccessToast } = useAdminErrorToast()

const props = defineProps({ id: { type: String, required: true } })

const breadcrumbLinks = [
  { path: 'agent-list', label: t('globals.terms.agent', 2) },
  { path: '', label: t('globals.messages.edit', { name: t('globals.terms.agent', 1).toLowerCase() }) }
]

const submitForm = async (values) => {
  formLoading.value = true
  try {
    await api.updateUser(user.value.id, values)
    showSuccessToast(t('globals.messages.updatedSuccessfully', { name: t('globals.terms.agent', 1) }))
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
}

onMounted(async () => {
  isLoading.value = true
  try {
    const resp = await api.getUser(props.id)
    user.value = resp.data.data
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
})
</script>
