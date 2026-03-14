<template>
  <div :class="{ 'opacity-50 transition-opacity duration-300': isLoading }">
    <NotificationConfigShell
      title="Comunicações oficiais"
      description="Centralize ofícios, cartas, notificações e intimações como chamados internos com regras claras de triagem e encaminhamento."
      status-label="Fluxo de comunicações oficiais"
      :status-description="statusDescription"
      :help-items="helpItems"
    >
      <OfficialCommunicationsForm
        :initial-values="initialValues"
        :inbox-options="inboxOptions"
        :priority-options="priorityOptions"
        :status-options="statusOptions"
        :team-options="teamOptions"
        :submit-form="submitForm"
      />
      <Spinner v-if="isLoading" />
    </NotificationConfigShell>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import api from '@/api'
import { Spinner } from '@/components/ui/spinner'
import OfficialCommunicationsForm from '@/features/admin/notification/OfficialCommunicationsForm.vue'
import NotificationConfigShell from '@/features/admin/notification/NotificationConfigShell.vue'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents'
import { useConversationStore } from '@/stores/conversation'
import { useInboxStore } from '@/stores/inbox'
import { useTeamStore } from '@/stores/team'
import { handleHTTPError } from '@/utils/http'

const emitter = useEmitter()
const teamStore = useTeamStore()
const inboxStore = useInboxStore()
const conversationStore = useConversationStore()
const isLoading = ref(false)

const defaults = {
  enabled: false,
  auto_create_conversation: true,
  inbox_id: '',
  priority_id: '',
  status_id: '',
  subject_prefix: '[Oficial] Comunicação recebida',
  target_sla_hours: '24',
  default_types: ['Ofício', 'Carta', 'Notificação', 'Intimação'],
  internal_note: '',
  routing_rules: []
}

const initialValues = ref({ ...defaults })

const teamOptions = computed(() => teamStore.options)
const inboxOptions = computed(() => inboxStore.options)
const priorityOptions = computed(() => conversationStore.priorityOptions)
const statusOptions = computed(() => conversationStore.statusOptionsNoSnooze)
const statusDescription = computed(() =>
  initialValues.value.enabled
    ? 'Fluxo habilitado para gerar e encaminhar chamados conforme as regras cadastradas.'
    : 'Fluxo desabilitado ou ainda sem política de encaminhamento definida.'
)

const helpItems = [
  {
    title: 'Classificação',
    description: 'Cadastre os tipos de comunicação mais comuns para padronizar a triagem administrativa e jurídica.'
  },
  {
    title: 'Encaminhamento',
    description: 'Cada regra pode enviar o mesmo chamado para uma ou várias equipes, conforme a natureza do expediente.'
  },
  {
    title: 'Operação',
    description: 'Defina prioridade, status inicial e caixa de entrada padrão para reduzir trabalho manual.'
  },
  {
    title: 'Evolução',
    description: 'A configuração salva nesta fase já prepara o modelo que pode ser conectado depois ao backend de abertura automática.'
  }
]

onMounted(async () => {
  isLoading.value = true
  try {
    await Promise.all([
      teamStore.fetchTeams(),
      inboxStore.fetchInboxes(),
      conversationStore.fetchPriorities(),
      conversationStore.fetchStatuses()
    ])

    const response = await api.getOfficialCommunicationsNotificationSettings()
    initialValues.value = {
      ...defaults,
      ...normalizeResponse(response?.data?.data)
    }
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    isLoading.value = false
  }
})

async function submitForm(values) {
  try {
    const payload = Object.fromEntries(
      Object.entries(values).map(([key, value]) => [
        `notification.official_communications.${key}`,
        value
      ])
    )
    await api.updateOfficialCommunicationsNotificationSettings(payload)
    initialValues.value = { ...values }
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      description: 'Configuração de comunicações oficiais salva com sucesso.'
    })
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
    throw error
  }
}

function normalizeResponse(data = {}) {
  const normalized = {}
  Object.entries(data).forEach(([key, value]) => {
    normalized[key.replace('notification.official_communications.', '')] = value
  })
  return normalized
}
</script>
