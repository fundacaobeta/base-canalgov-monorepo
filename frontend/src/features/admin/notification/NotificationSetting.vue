<template>
  <div :class="{ 'opacity-50 transition-opacity duration-300': isLoading }">
    <NotificationConfigShell
      title="E-mail"
      description="Configure o envio de notificações por SMTP para alertas operacionais, atribuições, violações de SLA e outros eventos internos."
      status-label="Canal de notificação por e-mail"
      :status-description="statusDescription"
      :help-items="helpItems"
    >
      <Spinner v-if="isLoading" />
      <NotificationsForm :initial-values="initialValues" :submit-form="submitForm" />
    </NotificationConfigShell>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import api from '@/api'
import { useI18n } from 'vue-i18n'
import NotificationsForm from './NotificationSettingForm.vue'
import NotificationConfigShell from './NotificationConfigShell.vue'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useEmitter } from '@/composables/useEmitter'
import { handleHTTPError } from '@/utils/http'
import { Spinner } from '@/components/ui/spinner'
import { useAppSettingsStore } from '@/stores/appSettings'

const initialValues = ref({})
const { t } = useI18n()
const isLoading = ref(false)
const emitter = useEmitter()
const appSettingsStore = useAppSettingsStore()

const statusDescription = computed(() =>
  initialValues.value.enabled
    ? 'Canal ativo para envio de notificações internas por SMTP.'
    : 'Canal desabilitado ou ainda não configurado.'
)

const helpItems = [
  {
    title: 'Segurança',
    description: 'Prefira SMTP autenticado com TLS ou STARTTLS e evite armazenar credenciais de teste em produção.'
  },
  {
    title: 'Entregabilidade',
    description: 'Configure corretamente o endereço de remetente e o hostname HELO para reduzir bloqueios e spam.'
  },
  {
    title: 'Capacidade',
    description: 'Ajuste conexões, timeouts e retentativas conforme o volume de notificações da operação.'
  },
  {
    title: 'Aplicação',
    description: 'Depois de salvar, reinicie a aplicação para que as alterações de SMTP tenham efeito.'
  }
]

onMounted(() => {
  getNotificationSettings()
})

const getNotificationSettings = async () => {
  try {
    isLoading.value = true
    const resp = await api.getEmailNotificationSettings()
    initialValues.value = Object.fromEntries(
      Object.entries(resp.data.data).map(([key, value]) => [
        key.replace('notification.email.', ''),
        value
      ])
    )
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    isLoading.value = false
  }
}

const submitForm = async (values) => {
  try {
    const updatedValues = Object.fromEntries(
      Object.entries(values).map(([key, value]) => {
        if (key === 'password' && value.includes('•')) {
          return [`notification.email.${key}`, '']
        }
        return [`notification.email.${key}`, value]
      })
    )
    await api.updateEmailNotificationSettings(updatedValues)
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      description: t('admin.notification.restartApp')
    })
    await getNotificationSettings()
    appSettingsStore.fetchSettings()
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  }
}
</script>
