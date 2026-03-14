<template>
  <NotificationConfigShell
    :title="config.title"
    :description="config.description"
    :status-label="config.statusLabel"
    :status-description="form.enabled ? 'Canal habilitado para uso operacional.' : 'Canal desabilitado ou ainda não configurado.'"
    :help-items="config.helpItems"
  >
    <form class="space-y-6" @submit.prevent="saveConfig">
          <div class="box p-5">
            <div class="flex items-center justify-between gap-6">
              <div class="space-y-1">
                <h3 class="font-medium">Habilitar canal</h3>
                <p class="text-sm text-muted-foreground">
                  Ative este canal para preparar o atendimento e as notificações.
                </p>
              </div>
              <Switch v-model:checked="form.enabled" />
            </div>
          </div>

          <div
            v-for="section in config.sections"
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
              <h3 class="font-medium">Eventos cobertos</h3>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="event in config.events"
                  :key="event"
                  class="rounded-full border bg-background px-3 py-1 text-xs"
                >
                  {{ event }}
                </span>
              </div>
            </div>
          </div>

      <Button type="submit" :isLoading="isSaving">Salvar configuração</Button>
    </form>
  </NotificationConfigShell>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
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
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents'
import { handleHTTPError } from '@/utils/http'

const props = defineProps({
  channel: {
    type: String,
    required: true
  }
})

const emitter = useEmitter()
const isSaving = ref(false)
const isLoading = ref(false)

const channelConfigs = {
  whatsapp: {
    title: 'WhatsApp',
    description: 'Configure credenciais, identificação do número e webhook do provedor de WhatsApp.',
    help: 'Use esta tela para preparar integração com provedores como Meta Cloud API, Z-API, Gupshup ou gateways próprios.',
    help2: 'Defina credenciais, URL base e eventos para envio e recebimento de mensagens operacionais.',
    helpItems: [
      {
        title: 'Provedor',
        description: 'Mapeie o provedor oficial do canal e valide o identificador do número ou instância conectada.'
      },
      {
        title: 'Webhook',
        description: 'Mantenha token de verificação e URL base consistentes para recebimento de eventos.'
      },
      {
        title: 'Operação',
        description: 'Defina mensagem ou template inicial para padronizar a primeira resposta do atendimento.'
      },
      {
        title: 'Governança',
        description: 'Use fila padrão para orientar triagem de mensagens quando o canal começar a operar.'
      }
    ],
    statusLabel: 'Canal de atendimento via WhatsApp',
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
    events: ['Mensagem recebida', 'Mensagem enviada', 'Confirmação de leitura', 'Falha de envio'],
    sections: [
      {
        title: 'Identificação do canal',
        description: 'Dados básicos da conta e do número conectado.',
        fields: [
          {
            key: 'provider',
            label: 'Provedor',
            type: 'select',
            placeholder: 'Selecione o provedor',
            options: [
              { value: 'meta', label: 'Meta Cloud API' },
              { value: 'zapi', label: 'Z-API' },
              { value: 'gupshup', label: 'Gupshup' },
              { value: 'custom', label: 'Customizado' }
            ]
          },
          { key: 'display_name', label: 'Nome exibido', placeholder: 'CanalGov Atendimento' },
          { key: 'phone_number_id', label: 'ID do número / instância', placeholder: '1234567890' },
          { key: 'department_hint', label: 'Fila padrão', placeholder: 'Atendimento digital' }
        ]
      },
      {
        title: 'Credenciais e webhook',
        description: 'Informações técnicas para autenticação e processamento.',
        fields: [
          { key: 'access_token', label: 'Token de acesso', type: 'password', placeholder: '••••••••••••••••' },
          { key: 'base_url', label: 'URL base da API', placeholder: 'https://graph.facebook.com/v20.0' },
          { key: 'webhook_verify_token', label: 'Token de verificação do webhook', placeholder: 'canalgov-whatsapp-verify' },
          { key: 'default_template', label: 'Template inicial', type: 'textarea', fullWidth: true, placeholder: 'Mensagem padrão de confirmação ou recepção.' }
        ]
      }
    ]
  },
  telegram: {
    title: 'Telegram',
    description: 'Configure bots, chats e webhook para operação via Telegram.',
    help: 'Use o bot do Telegram para receber mensagens, notificar equipes e automatizar respostas transacionais.',
    help2: 'Mapeie grupos, canais ou chats específicos para o fluxo operacional desejado.',
    helpItems: [
      {
        title: 'Bot',
        description: 'Garanta que o token e o nome do bot correspondam ao ambiente configurado.'
      },
      {
        title: 'Destino',
        description: 'Use chat padrão ou grupos específicos para separar notificações internas e atendimento.'
      },
      {
        title: 'Webhook',
        description: 'Restrinja os updates permitidos ao necessário para reduzir ruído operacional.'
      },
      {
        title: 'Mensagem',
        description: 'Padronize a resposta inicial do bot para orientar a abertura de atendimento.'
      }
    ],
    statusLabel: 'Canal de atendimento via Telegram',
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
    events: ['Mensagem recebida', 'Comando do bot', 'Callback de botão', 'Notificação de fila'],
    sections: [
      {
        title: 'Bot e destino',
        description: 'Dados do bot e do chat padrão.',
        fields: [
          { key: 'bot_name', label: 'Nome do bot', placeholder: 'CanalGov Bot' },
          { key: 'default_chat_id', label: 'Chat ID padrão', placeholder: '-1001234567890' }
        ]
      },
      {
        title: 'Conectividade',
        description: 'Credenciais e eventos do webhook.',
        fields: [
          { key: 'bot_token', label: 'Token do bot', type: 'password', placeholder: '••••••••••••••••' },
          { key: 'webhook_url', label: 'URL do webhook', placeholder: 'https://seu-dominio.exemplo.com/telegram/webhook' },
          { key: 'allowed_updates', label: 'Eventos permitidos', placeholder: 'message,callback_query' },
          { key: 'default_message', label: 'Mensagem padrão', type: 'textarea', fullWidth: true, placeholder: 'Mensagem inicial enviada pelo bot.' }
        ]
      }
    ]
  },
  sms: {
    title: 'SMS',
    description: 'Configure provedores, credenciais e remetente padrão para envio por SMS.',
    help: 'Ideal para alertas críticos, autenticação, cobranças ou comunicação com usuários sem app dedicado.',
    help2: 'Defina URL base, chave da API, remetente e regras de fallback por operadora ou país.',
    helpItems: [
      {
        title: 'Remetente',
        description: 'Cadastre um sender ID reconhecível para facilitar a identificação da mensagem pelo cidadão.'
      },
      {
        title: 'Credenciais',
        description: 'Separe chave e segredo por ambiente para evitar envio acidental em produção.'
      },
      {
        title: 'Fallback',
        description: 'Configure DDI padrão e texto base para suportar alertas transacionais rápidos.'
      },
      {
        title: 'Uso',
        description: 'Priorize SMS para eventos críticos, autenticação e comunicações sem dependência de aplicativo.'
      }
    ],
    statusLabel: 'Canal de atendimento e alerta por SMS',
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
    events: ['Disparo de alerta', 'Confirmação transacional', 'Retentativa de SLA', 'Fallback operacional'],
    sections: [
      {
        title: 'Provedor e remetente',
        description: 'Dados básicos do provedor de SMS.',
        fields: [
          {
            key: 'provider',
            label: 'Provedor',
            type: 'select',
            placeholder: 'Selecione o provedor',
            options: [
              { value: 'twilio', label: 'Twilio' },
              { value: 'zenvia', label: 'Zenvia' },
              { value: 'totalvoice', label: 'TotalVoice' },
              { value: 'custom', label: 'Customizado' }
            ]
          },
          { key: 'sender_id', label: 'Sender ID', placeholder: 'CANALGOV' }
        ]
      },
      {
        title: 'Credenciais e padrão de envio',
        description: 'Autenticação e comportamento padrão do canal.',
        fields: [
          { key: 'api_key', label: 'Chave da API', type: 'password', placeholder: '••••••••••••••••' },
          { key: 'api_secret', label: 'Segredo da API', type: 'password', placeholder: '••••••••••••••••' },
          { key: 'base_url', label: 'URL base da API', placeholder: 'https://api.provedor.exemplo.com' },
          { key: 'fallback_country_code', label: 'DDI padrão', placeholder: '55' },
          { key: 'default_message', label: 'Mensagem padrão', type: 'textarea', fullWidth: true, placeholder: 'Texto base usado para alertas e confirmações.' }
        ]
      }
    ]
  },
  push: {
    title: 'Notificação push',
    description: 'Configure o envio de notificações push para aplicativos móveis e portais web.',
    help: 'Use push para avisos operacionais, atualização de protocolos, menções e alertas de SLA.',
    help2: 'Defina credenciais do provedor, app alvo, segmentos e payload padrão.',
    helpItems: [
      {
        title: 'Aplicação',
        description: 'Associe corretamente app ID e project ID para evitar envio ao app errado.'
      },
      {
        title: 'Entrega',
        description: 'Defina tópico padrão e URL de clique para levar o usuário ao fluxo correto.'
      },
      {
        title: 'Payload',
        description: 'Use um template simples e previsível para padronizar alertas operacionais.'
      },
      {
        title: 'Segmentação',
        description: 'Combine push com menções, SLA e atualizações de protocolo para uma experiência mais útil.'
      }
    ],
    statusLabel: 'Canal de notificação push',
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
    events: ['Nova atribuição', 'Menção em conversa', 'Violação de SLA', 'Atualização de protocolo'],
    sections: [
      {
        title: 'Aplicação e provedor',
        description: 'Dados do app e do projeto de push.',
        fields: [
          {
            key: 'provider',
            label: 'Provedor',
            type: 'select',
            placeholder: 'Selecione o provedor',
            options: [
              { value: 'firebase', label: 'Firebase Cloud Messaging' },
              { value: 'onesignal', label: 'OneSignal' },
              { value: 'custom', label: 'Customizado' }
            ]
          },
          { key: 'app_id', label: 'ID da aplicação', placeholder: 'canalgov.app' },
          { key: 'project_id', label: 'ID do projeto', placeholder: 'canalgov-prod' },
          { key: 'topic_default', label: 'Tópico padrão', placeholder: 'atendimento-geral' }
        ]
      },
      {
        title: 'Entrega e payload',
        description: 'Credenciais e estrutura de notificação.',
        fields: [
          { key: 'api_key', label: 'Chave / token do servidor', type: 'password', placeholder: '••••••••••••••••' },
          { key: 'click_action_url', label: 'URL ao clicar', placeholder: 'https://portal.exemplo.com/chamados' },
          { key: 'payload_template', label: 'Payload padrão', type: 'textarea', fullWidth: true, placeholder: '{ "title": "Novo alerta" }' }
        ]
      }
    ]
  }
}

const config = computed(() => channelConfigs[props.channel])
const form = reactive({ ...config.value.defaults })

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
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
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
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      description: `${config.value.title} configurado com sucesso.`
    })
    await loadConfig()
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    isSaving.value = false
  }
}
</script>
