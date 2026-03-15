<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <RoleForm v-if="!isLoading" :initial-values="role" :submitForm="submitForm" :isLoading="formLoading" />
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import RoleForm from '@/features/admin/roles/RoleForm.vue'
import api from '@/api'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const role = ref({})
const isLoading = ref(false)
const formLoading = ref(false)
const { showErrorToast, showSuccessToast } = useAdminErrorToast()

const props = defineProps({ id: { type: String, required: true } })

const breadcrumbLinks = [
  { path: 'role-list', label: t('admin.roles.title') },
  { path: '', label: t('globals.messages.edit', { name: t('globals.terms.role') }) }
]

onMounted(async () => {
  isLoading.value = true
  try {
    const resp = await api.getRole(props.id)
    role.value = resp.data.data
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
})

const submitForm = async (values) => {
  formLoading.value = true
  try {
    await api.updateRole(props.id, values)
    showSuccessToast(t('globals.messages.updatedSuccessfully', { name: t('globals.terms.role') }))
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
}
</script>
