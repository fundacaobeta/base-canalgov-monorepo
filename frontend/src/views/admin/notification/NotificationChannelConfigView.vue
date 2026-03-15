<template>
  <NotificationConfigShell
    :title="translatedConfig.title"
    :description="translatedConfig.description"
    :status-label="translatedConfig.statusLabel"
    :status-description="form.enabled ? translatedConfig.statusEnabled : translatedConfig.statusDisabled"
    :help-items="translatedConfig.helpItems"
  >
    <form class="space-y-6" @submit.prevent="saveConfig">
          <div class="box p-5">
            <div class="flex items-center justify-between gap-6">
              <div class="space-y-1">
                <h3 class="font-medium">{{ $t('admin.notification.channel.enableChannel') }}</h3>
                <p class="text-sm text-muted-foreground">
                  {{ $t('admin.notification.channel.enableChannelDescription') }}
                </p>
              </div>
              <Switch v-model:checked="form.enabled" />
            </div>
          </div>

          <div
            v-for="section in translatedConfig.sections"
            :key="section.title"
            class="box p-6 space-y-5"
          >
            <div class="space-y-1">
              <h3 class="font-medium">{{ section.title }}</h3>
              <p class="text-sm text-muted-foreground">{{ section.description }}</p>
            </div>

            <div class="grid gap-5 md:grid-cols-2">
              <div v-for="field in section.fields" :key="field.key" :class="field.fullWidth ? 'md:col-span-2' : ''">
                <label class="text-sm font-medium">{{ field.label }}</label>

                <Select
                  v-if="field.type === 'select'"
                  v-model="form[field.key]"
                >
                  <SelectTrigger class="mt-2">
                    <SelectValue :placeholder="field.placeholder" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="option in field.options"
                      :key="option.value"
                      :value="option.value"
                    >
                      {{ option.label }}
                    </SelectItem>
                  </SelectContent>
                </Select>

                <Textarea
                  v-else-if="field.type === 'textarea'"
                  v-model="form[field.key]"
                  class="mt-2 min-h-32 font-mono text-xs"
                  :placeholder="field.placeholder"
                />

                <Input
                  v-else
                  v-model="form[field.key]"
                  class="mt-2"
                  :type="field.type || 'text'"
                  :placeholder="field.placeholder"
                />

                <p v-if="field.help" class="mt-2 text-xs text-muted-foreground">
                  {{ field.help }}
                </p>
              </div>
            </div>
          </div>

          <div class="box p-5">
            <div class="space-y-3">
              <h3 class="font-medium">{{ $t('admin.notification.channel.coveredEvents') }}</h3>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="event in translatedConfig.events"
                  :key="event"
                  class="rounded-full border bg-background px-3 py-1 text-xs"
                >
                  {{ event }}
                </span>
              </div>
            </div>
          </div>

      <Button type="submit" :isLoading="isSaving">{{ $t('globals.messages.save') }}</Button>
    </form>
  </NotificationConfigShell>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import NotificationConfigShell from '@/features/admin/notification/NotificationConfigShell.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'

const props = defineProps({
  channel: {
    type: String,
    required: true
  }
})

const { t } = useI18n()
const { showErrorToast, showSuccessToast } = useAdminErrorToast()
const isSaving = ref(false)
const isLoading = ref(false)

const channelConfigs = {
  whatsapp: {
    title: 'WhatsApp',
    descriptionKey: 'admin.notification.whatsapp.description',
    statusLabelKey: 'admin.notification.whatsapp.statusLabel',
    statusEnabledKey: 'admin.notification.whatsapp.statusEnabled',
    statusDisabledKey: 'admin.notification.whatsapp.statusDisabled',
    helpItems: [
      {
        titleKey: 'admin.notification.whatsapp.help1.title',
        descriptionKey: 'admin.notification.whatsapp.help1.description'
      },
      {
        titleKey: 'admin.notification.whatsapp.help2.title',
        descriptionKey: 'admin.notification.whatsapp.help2.description'
      },
      {
        titleKey: 'admin.notification.whatsapp.help3.title',
        descriptionKey: 'admin.notification.whatsapp.help3.description'
      },
      {
        titleKey: 'admin.notification.whatsapp.help4.title',
        descriptionKey: 'admin.notification.whatsapp.help4.description'
      }
    ],
    defaults: {
      enabled: false,
      provider: 'meta',
      display_name: '',
      phone_number_id: '',
      access_token: '',
      base_url: '',
      webhook_verify_token: '',
      default_template: '',
      department_hint: ''
    },
    secretFields: ['access_token', 'webhook_verify_token'],
    events: [
      'admin.notification.whatsapp.event1',
      'admin.notification.whatsapp.event2',
      'admin.notification.whatsapp.event3',
      'admin.notification.whatsapp.event4'
    ],
    sections: [
      {
        titleKey: 'admin.notification.whatsapp.section1.title',
        descriptionKey: 'admin.notification.whatsapp.section1.description',
        fields: [
          {
            key: 'provider',
            labelKey: 'admin.notification.channel.provider',
            type: 'select',
            placeholderKey: 'admin.notification.channel.selectProvider',
            options: [
              { value: 'meta', label: 'Meta Cloud API' },
              { value: 'zapi', label: 'Z-API' },
              { value: 'gupshup', label: 'Gupshup' },
              { value: 'custom', label: 'globals.terms.custom' }
            ]
          },
          { key: 'display_name', labelKey: 'admin.notification.whatsapp.fields.displayName', placeholder: 'CanalGov Atendimento' },
          { key: 'phone_number_id', labelKey: 'admin.notification.whatsapp.fields.phoneNumberId', placeholder: '1234567890' },
          { key: 'department_hint', labelKey: 'admin.notification.whatsapp.fields.departmentHint', placeholder: 'Atendimento digital' }
        ]
      },
      {
        titleKey: 'admin.notification.whatsapp.section2.title',
        descriptionKey: 'admin.notification.whatsapp.section2.description',
        fields: [
          { key: 'access_token', labelKey: 'admin.notification.whatsapp.fields.accessToken', type: 'password', placeholder: '••••••••••••••••' },
          { key: 'base_url', labelKey: 'admin.notification.whatsapp.fields.baseUrl', placeholder: 'https://graph.facebook.com/v20.0' },
          { key: 'webhook_verify_token', labelKey: 'admin.notification.whatsapp.fields.webhookVerifyToken', placeholder: 'canalgov-whatsapp-verify' },
          { key: 'default_template', labelKey: 'admin.notification.whatsapp.fields.defaultTemplate', type: 'textarea', fullWidth: true, placeholder: 'Mensagem padrão de confirmação ou recepção.' }
        ]
      }
    ]
  },
  telegram: {
    title: 'Telegram',
    descriptionKey: 'admin.notification.telegram.description',
    statusLabelKey: 'admin.notification.telegram.statusLabel',
    statusEnabledKey: 'admin.notification.channel.statusEnabled',
    statusDisabledKey: 'admin.notification.channel.statusDisabled',
    helpItems: [
      {
        titleKey: 'admin.notification.telegram.help1.title',
        descriptionKey: 'admin.notification.telegram.help1.description'
      },
      {
        titleKey: 'admin.notification.telegram.help2.title',
        descriptionKey: 'admin.notification.telegram.help2.description'
      },
      {
        titleKey: 'admin.notification.telegram.help3.title',
        descriptionKey: 'admin.notification.telegram.help3.description'
      },
      {
        titleKey: 'admin.notification.telegram.help4.title',
        descriptionKey: 'admin.notification.telegram.help4.description'
      }
    ],
    defaults: {
      enabled: false,
      bot_name: '',
      bot_token: '',
      webhook_url: '',
      default_chat_id: '',
      allowed_updates: 'message,callback_query',
      default_message: ''
    },
    secretFields: ['bot_token'],
    events: [
      'admin.notification.telegram.event1',
      'admin.notification.telegram.event2',
      'admin.notification.telegram.event3',
      'admin.notification.telegram.event4'
    ],
    sections: [
      {
        titleKey: 'admin.notification.telegram.section1.title',
        descriptionKey: 'admin.notification.telegram.section1.description',
        fields: [
          { key: 'bot_name', labelKey: 'admin.notification.telegram.fields.botName', placeholder: 'CanalGov Bot' },
          { key: 'default_chat_id', labelKey: 'admin.notification.telegram.fields.defaultChatId', placeholder: '-1001234567890' }
        ]
      },
      {
        titleKey: 'admin.notification.telegram.section2.title',
        descriptionKey: 'admin.notification.telegram.section2.description',
        fields: [
          { key: 'bot_token', labelKey: 'admin.notification.telegram.fields.botToken', type: 'password', placeholder: '••••••••••••••••' },
          { key: 'webhook_url', labelKey: 'admin.notification.telegram.fields.webhookUrl', placeholder: 'https://seu-dominio.exemplo.com/telegram/webhook' },
          { key: 'allowed_updates', labelKey: 'admin.notification.telegram.fields.allowedUpdates', placeholder: 'message,callback_query' },
          { key: 'default_message', labelKey: 'admin.notification.telegram.fields.defaultMessage', type: 'textarea', fullWidth: true, placeholder: 'Mensagem inicial enviada pelo bot.' }
        ]
      }
    ]
  },
  sms: {
    title: 'SMS',
    descriptionKey: 'admin.notification.sms.description',
    statusLabelKey: 'admin.notification.sms.statusLabel',
    statusEnabledKey: 'admin.notification.channel.statusEnabled',
    statusDisabledKey: 'admin.notification.channel.statusDisabled',
    helpItems: [
      {
        titleKey: 'admin.notification.sms.help1.title',
        descriptionKey: 'admin.notification.sms.help1.description'
      },
      {
        titleKey: 'admin.notification.sms.help2.title',
        descriptionKey: 'admin.notification.sms.help2.description'
      },
      {
        titleKey: 'admin.notification.sms.help3.title',
        descriptionKey: 'admin.notification.sms.help3.description'
      },
      {
        titleKey: 'admin.notification.sms.help4.title',
        descriptionKey: 'admin.notification.sms.help4.description'
      }
    ],
    defaults: {
      enabled: false,
      provider: 'twilio',
      sender_id: '',
      api_key: '',
      api_secret: '',
      base_url: '',
      fallback_country_code: '55',
      default_message: ''
    },
    secretFields: ['api_key', 'api_secret'],
    events: [
      'admin.notification.sms.event1',
      'admin.notification.sms.event2',
      'admin.notification.sms.event3',
      'admin.notification.sms.event4'
    ],
    sections: [
      {
        titleKey: 'admin.notification.sms.section1.title',
        descriptionKey: 'admin.notification.sms.section1.description',
        fields: [
          {
            key: 'provider',
            labelKey: 'admin.notification.channel.provider',
            type: 'select',
            placeholderKey: 'admin.notification.channel.selectProvider',
            options: [
              { value: 'twilio', label: 'Twilio' },
              { value: 'zenvia', label: 'Zenvia' },
              { value: 'totalvoice', label: 'TotalVoice' },
              { value: 'custom', label: 'globals.terms.custom' }
            ]
          },
          { key: 'sender_id', labelKey: 'admin.notification.sms.fields.senderId', placeholder: 'CANALGOV' }
        ]
      },
      {
        titleKey: 'admin.notification.sms.section2.title',
        descriptionKey: 'admin.notification.sms.section2.description',
        fields: [
          { key: 'api_key', labelKey: 'admin.notification.sms.fields.apiKey', type: 'password', placeholder: '••••••••••••••••' },
          { key: 'api_secret', labelKey: 'admin.notification.sms.fields.apiSecret', type: 'password', placeholder: '••••••••••••••••' },
          { key: 'base_url', labelKey: 'admin.notification.sms.fields.baseUrl', placeholder: 'https://api.provedor.exemplo.com' },
          { key: 'fallback_country_code', labelKey: 'admin.notification.sms.fields.fallbackCountryCode', placeholder: '55' },
          { key: 'default_message', labelKey: 'admin.notification.sms.fields.defaultMessage', type: 'textarea', fullWidth: true, placeholder: 'Texto base usado para alertas e confirmações.' }
        ]
      }
    ]
  },
  push: {
    title: 'Notificação push',
    descriptionKey: 'admin.notification.push.description',
    statusLabelKey: 'admin.notification.push.statusLabel',
    statusEnabledKey: 'admin.notification.channel.statusEnabled',
    statusDisabledKey: 'admin.notification.channel.statusDisabled',
    helpItems: [
      {
        titleKey: 'admin.notification.push.help1.title',
        descriptionKey: 'admin.notification.push.help1.description'
      },
      {
        titleKey: 'admin.notification.push.help2.title',
        descriptionKey: 'admin.notification.push.help2.description'
      },
      {
        titleKey: 'admin.notification.push.help3.title',
        descriptionKey: 'admin.notification.push.help3.description'
      },
      {
        titleKey: 'admin.notification.push.help4.title',
        descriptionKey: 'admin.notification.push.help4.description'
      }
    ],
    defaults: {
      enabled: false,
      provider: 'firebase',
      app_id: '',
      project_id: '',
      api_key: '',
      topic_default: '',
      click_action_url: '',
      payload_template: '{\n  "title": "Novo alerta",\n  "body": "Você recebeu uma nova atualização."\n}'
    },
    secretFields: ['api_key'],
    events: [
      'admin.notification.push.event1',
      'admin.notification.push.event2',
      'admin.notification.push.event3',
      'admin.notification.push.event4'
    ],
    sections: [
      {
        titleKey: 'admin.notification.push.section1.title',
        descriptionKey: 'admin.notification.push.section1.description',
        fields: [
          {
            key: 'provider',
            labelKey: 'admin.notification.channel.provider',
            type: 'select',
            placeholderKey: 'admin.notification.channel.selectProvider',
            options: [
              { value: 'firebase', label: 'Firebase Cloud Messaging' },
              { value: 'onesignal', label: 'OneSignal' },
              { value: 'custom', label: 'globals.terms.custom' }
            ]
          },
          { key: 'app_id', labelKey: 'admin.notification.push.fields.appId', placeholder: 'canalgov.app' },
          { key: 'project_id', labelKey: 'admin.notification.push.fields.projectId', placeholder: 'canalgov-prod' },
          { key: 'topic_default', labelKey: 'admin.notification.push.fields.topicDefault', placeholder: 'atendimento-geral' }
        ]
      },
      {
        titleKey: 'admin.notification.push.section2.title',
        descriptionKey: 'admin.notification.push.section2.description',
        fields: [
          { key: 'api_key', labelKey: 'admin.notification.push.fields.apiKey', type: 'password', placeholder: '••••••••••••••••' },
          { key: 'click_action_url', labelKey: 'admin.notification.push.fields.clickActionUrl', placeholder: 'https://portal.exemplo.com/chamados' },
          { key: 'payload_template', labelKey: 'admin.notification.push.fields.payloadTemplate', type: 'textarea', fullWidth: true, placeholder: '{ "title": "Novo alerta" }' }
        ]
      }
    ]
  }
}

const config = computed(() => channelConfigs[props.channel])
const form = reactive({ ...config.value.defaults })

const translatedConfig = computed(() => {
  const c = config.value
  return {
    title: c.title,
    description: t(c.descriptionKey),
    statusLabel: t(c.statusLabelKey),
    statusEnabled: t(c.statusEnabledKey),
    statusDisabled: t(c.statusDisabledKey),
    helpItems: c.helpItems.map(item => ({
      title: t(item.titleKey),
      description: t(item.descriptionKey)
    })),
    events: c.events.map(key => t(key)),
    sections: c.sections.map(section => ({
      title: t(section.titleKey),
      description: t(section.descriptionKey),
      fields: section.fields.map(field => ({
        ...field,
        label: t(field.labelKey),
        placeholder: field.placeholderKey ? t(field.placeholderKey) : field.placeholder,
        help: field.helpKey ? t(field.helpKey) : undefined,
        options: field.options
          ? field.options.map(opt =>
              opt.label.startsWith('globals.')
                ? { ...opt, label: t(opt.label) }
                : opt
            )
          : undefined
      }))
    }))
  }
})

const channelAPI = {
  whatsapp: {
    get: api.getWhatsAppNotificationSettings,
    update: api.updateWhatsAppNotificationSettings
  },
  telegram: {
    get: api.getTelegramNotificationSettings,
    update: api.updateTelegramNotificationSettings
  },
  sms: {
    get: api.getSMSNotificationSettings,
    update: api.updateSMSNotificationSettings
  },
  push: {
    get: api.getPushNotificationSettings,
    update: api.updatePushNotificationSettings
  }
}

onMounted(async () => {
  await loadConfig()
})

async function loadConfig() {
  isLoading.value = true
  try {
    const response = await channelAPI[props.channel].get()
    Object.assign(form, config.value.defaults, normalizeResponse(response?.data?.data))
  } catch (error) {
    showErrorToast(error)
  } finally {
    isLoading.value = false
  }
}

function normalizeResponse(data = {}) {
  const normalized = {}
  Object.entries(data).forEach(([key, value]) => {
    normalized[key.replace(`notification.${props.channel}.`, '')] = value
  })
  return normalized
}

async function saveConfig() {
  isSaving.value = true
  try {
    const payload = Object.fromEntries(
      Object.entries(form).map(([key, value]) => {
        if (config.value.secretFields?.includes(key) && typeof value === 'string' && value.includes('•')) {
          return [`notification.${props.channel}.${key}`, '']
        }
        return [`notification.${props.channel}.${key}`, value]
      })
    )

    await channelAPI[props.channel].update(payload)
    showSuccessToast(t('admin.notification.channel.savedSuccess', { title: config.value.title }))
    await loadConfig()
  } catch (error) {
    showErrorToast(error)
  } finally {
    isSaving.value = false
  }
}
</script>
