<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <Spinner v-if="formLoading" />

  <EmailInboxForm
    v-else-if="isEmailInbox && emailMode === 'custom'"
    :initialValues="inbox"
    :submitForm="submitEmailForm"
    :isLoading="isLoading"
  />

  <ManagedEmailInboxForm
    v-else-if="isEmailInbox && emailMode === 'managed'"
    :initial-values="managedInbox"
    :submitForm="submitManagedEmailForm"
    :isLoading="isLoading"
  />

  <GenericInboxForm
    v-else
    :initial-values="genericInbox"
    :submitForm="submitGenericForm"
    :isLoading="isLoading"
    :channel="inbox.channel"
    :channel-label="channelLabel"
  />
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import api from '@/api'
import EmailInboxForm from '@/features/admin/inbox/EmailInboxForm.vue'
import GenericInboxForm from '@/features/admin/inbox/GenericInboxForm.vue'
import ManagedEmailInboxForm from '@/features/admin/inbox/ManagedEmailInboxForm.vue'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb/index.js'
import { Spinner } from '@/components/ui/spinner'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { AUTH_TYPE_PASSWORD, AUTH_TYPE_OAUTH2 } from '@/constants/auth.js'
import { useEmitter } from '@/composables/useEmitter'
import { handleHTTPError } from '@/utils/http'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  id: {
    type: String,
    required: true
  }
})

const emitter = useEmitter()
const { t } = useI18n()
const formLoading = ref(false)
const isLoading = ref(false)
const inbox = ref({})

const breadcrumbLinks = [
  { path: 'inbox-list', label: t('globals.terms.inbox', 2) },
  { path: '', label: t('globals.messages.edit') }
]

const isEmailInbox = computed(() => inbox.value.channel === 'email')
const emailMode = computed(() => inbox.value?.config?.receive_mode === 'managed' ? 'managed' : 'custom')
const genericInbox = computed(() => ({
  ...inbox.value,
  extra_config: JSON.stringify(inbox.value.config || {}, null, 2)
}))
const managedInbox = computed(() => ({
  name: inbox.value.name || '',
  managed_local_part: inbox.value?.config?.managed_local_part || '',
  managed_domain_id: inbox.value?.config?.managed_domain_id || '',
  delivery_provider: inbox.value?.config?.delivery_provider || 'canalgov_managed',
  provider_config: JSON.stringify(inbox.value?.config?.provider_config || {}, null, 2),
  enabled: typeof inbox.value.enabled === 'boolean' ? inbox.value.enabled : true,
  csat_enabled: typeof inbox.value.csat_enabled === 'boolean' ? inbox.value.csat_enabled : false
}))
const channelLabel = computed(() => {
  const labels = {
    none: 'Nenhum',
    whatsapp: 'WhatsApp',
    telegram: 'Telegram',
    sms: 'SMS',
    push: 'Notificação push',
    webhook: 'Webhook'
  }
  return labels[inbox.value.channel] || 'Canal'
})

const updateInbox = async (payload) => {
  try {
    isLoading.value = true
    await api.updateInbox(inbox.value.id, payload)
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      description: t('globals.messages.updatedSuccessfully', {
        name: t('globals.terms.inbox')
      })
    })
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    isLoading.value = false
  }
}

const submitEmailForm = (values) => {
  const config = {
    auth_type: values.auth_type,
    enable_plus_addressing: values.enable_plus_addressing,
    imap: [{ ...values.imap }],
    smtp: [{ ...values.smtp }]
  }

  if (values.auth_type === AUTH_TYPE_OAUTH2) {
    config.oauth = values.oauth
  }

  const payload = {
    ...values,
    channel: inbox.value.channel,
    config
  }

  if (payload.config.imap[0].password?.includes('•')) {
    payload.config.imap[0].password = ''
  }

  if (payload.config.auth_type === AUTH_TYPE_OAUTH2) {
    if (payload.config.oauth.access_token?.includes('•')) payload.config.oauth.access_token = ''
    if (payload.config.oauth.client_secret?.includes('•')) payload.config.oauth.client_secret = ''
    if (payload.config.oauth.refresh_token?.includes('•')) payload.config.oauth.refresh_token = ''
  }

  payload.config.smtp.forEach((smtp) => {
    if (smtp.password?.includes('•')) smtp.password = ''
  })

  updateInbox(payload)
}

const submitGenericForm = (values) => {
  updateInbox({
    name: values.name,
    from: values.from || '',
    enabled: values.enabled,
    csat_enabled: values.csat_enabled,
    channel: inbox.value.channel,
    config: JSON.parse(values.extra_config || '{}')
  })
}

const submitManagedEmailForm = (values) => {
  updateInbox({
    name: values.name,
    from: values.managed_email_address,
    enabled: values.enabled,
    csat_enabled: values.csat_enabled,
    channel: 'email',
    config: {
      receive_mode: 'managed',
      delivery_provider: values.delivery_provider,
      managed_domain_id: values.managed_domain_id,
      managed_domain: values.managed_domain,
      managed_local_part: values.managed_local_part,
      managed_email_address: values.managed_email_address,
      provider_config: JSON.parse(values.provider_config || '{}'),
      smtp: [],
      imap: []
    }
  })
}

onMounted(async () => {
  try {
    formLoading.value = true
    const resp = await api.getInbox(props.id)
    const inboxData = resp.data.data

    if (inboxData?.config?.imap) {
      inboxData.imap = inboxData.config.imap[0]
    }
    if (inboxData?.config?.smtp) {
      inboxData.smtp = inboxData.config.smtp[0]
    }
    inboxData.auth_type = inboxData?.config?.auth_type || AUTH_TYPE_PASSWORD
    inboxData.oauth = inboxData?.config?.oauth || {}
    inboxData.enable_plus_addressing = inboxData?.config?.enable_plus_addressing || false
    inbox.value = inboxData
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    formLoading.value = false
  }
})
</script>
