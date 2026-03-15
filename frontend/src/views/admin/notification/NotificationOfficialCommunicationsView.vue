<template>
  <div :class="{ 'opacity-50 transition-opacity duration-300': isLoading }">
    <NotificationConfigShell
      :title="$t('admin.notification.officialCommunications.title')"
      :description="$t('admin.notification.officialCommunications.description')"
      :status-label="$t('admin.notification.officialCommunications.statusLabel')"
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
import { useI18n } from 'vue-i18n'
import api from '@/api'
import { Spinner } from '@/components/ui/spinner'
import OfficialCommunicationsForm from '@/features/admin/notification/OfficialCommunicationsForm.vue'
import NotificationConfigShell from '@/features/admin/notification/NotificationConfigShell.vue'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useConversationStore } from '@/stores/conversation'
import { useInboxStore } from '@/stores/inbox'
import { useTeamStore } from '@/stores/team'

const { t } = useI18n()
const { showErrorToast, showSuccessToast } = useAdminErrorToast()
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
    ? t('admin.notification.officialCommunications.statusEnabled')
    : t('admin.notification.officialCommunications.statusDisabled')
)

const helpItems = computed(() => [
  {
    title: t('admin.notification.officialCommunications.help1.title'),
    description: t('admin.notification.officialCommunications.help1.description')
  },
  {
    title: t('admin.notification.officialCommunications.help2.title'),
    description: t('admin.notification.officialCommunications.help2.description')
  },
  {
    title: t('admin.notification.officialCommunications.help3.title'),
    description: t('admin.notification.officialCommunications.help3.description')
  },
  {
    title: t('admin.notification.officialCommunications.help4.title'),
    description: t('admin.notification.officialCommunications.help4.description')
  }
])

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
    showErrorToast(error)
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
    showSuccessToast(t('admin.notification.officialCommunications.savedSuccess'))
  } catch (error) {
    showErrorToast(error)
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
