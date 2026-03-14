<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <Spinner v-if="isLoading"></Spinner>
  <TeamForm :initial-values="team" :submitForm="submitForm" :isLoading="formLoading" v-else />
</template>

<script setup>
import { onMounted, ref } from 'vue'
import api from '@/api'
import TeamForm from '@/features/admin/teams/TeamForm.vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useEmitter } from '@/composables/useEmitter'
import { handleHTTPError } from '@/utils/http'
import { Spinner } from '@/components/ui/spinner'
import { useI18n } from 'vue-i18n'

const team = ref({})
const emitter = useEmitter()
const formLoading = ref(false)
const isLoading = ref(false)
const { t } = useI18n()

const breadcrumbLinks = [
  { path: 'team-list', label: 'Equipes' },
  { path: '', label: 'Editar equipe' }
]

const props = defineProps({
  id: {
    type: String,
    required: true
  }
})

const submitForm = (values) => {
  updateTeam(values)
}

const updateTeam = async (payload) => {
  try {
    formLoading.value = true
    await api.updateTeam(team.value.id, payload)
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      title: t('globals.terms.success'),
      description: t('globals.messages.teamUpdated')
    })
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      title: t('globals.terms.error'),
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    formLoading.value = false
  }
}

onMounted(async () => {
  try {
    isLoading.value = true
    const resp = await api.getTeam(props.id)
    team.value = resp.data.data
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      title: t('globals.terms.error'),
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    isLoading.value = false
  }
})
</script>
