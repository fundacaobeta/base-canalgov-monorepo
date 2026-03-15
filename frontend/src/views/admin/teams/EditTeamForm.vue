<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <TeamForm v-if="!isLoading" :initial-values="team" :submitForm="submitForm" :isLoading="formLoading" />
</template>

<script setup>
import { onMounted, ref } from 'vue'
import api from '@/api'
import TeamForm from '@/features/admin/teams/TeamForm.vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'

const team = ref({})
const formLoading = ref(false)
const isLoading = ref(false)
const { t } = useI18n()
const { showErrorToast, showSuccessToast } = useAdminErrorToast()

const props = defineProps({ id: { type: String, required: true } })

const breadcrumbLinks = [
  { path: 'team-list', label: t('admin.teams.title') },
  { path: '', label: t('globals.messages.edit', { name: t('globals.terms.team') }) }
]

const submitForm = async (values) => {
  formLoading.value = true
  try {
    await api.updateTeam(team.value.id, values)
    showSuccessToast(t('globals.messages.teamUpdated'))
  } catch (error) {
    showErrorToast(error)
  } finally {
    formLoading.value = false
  }
}

onMounted(async () => {
  isLoading.value = true
  try {
    const resp = await api.getTeam(props.id)
    team.value = resp.data.data
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
})
</script>
