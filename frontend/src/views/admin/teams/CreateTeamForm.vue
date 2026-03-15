<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <TeamForm :initial-values="{}" :submitForm="submitForm" :isLoading="formLoading" />
</template>

<script setup>
import { ref } from 'vue'
import TeamForm from '@/features/admin/teams/TeamForm.vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { useRouter } from 'vue-router'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'
import api from '@/api'

const formLoading = ref(false)
const router = useRouter()
const { t } = useI18n()
const { showErrorToast, showSuccessToast } = useAdminErrorToast()
const breadcrumbLinks = [
  { path: 'team-list', label: t('admin.teams.title') },
  { path: '', label: t('globals.messages.new', { name: t('globals.terms.team') }) }
]

const submitForm = async (values) => {
  formLoading.value = true
  try {
    await api.createTeam(values)
    showSuccessToast(t('globals.messages.teamCreated'))
    router.push({ name: 'team-list' })
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
}
</script>
