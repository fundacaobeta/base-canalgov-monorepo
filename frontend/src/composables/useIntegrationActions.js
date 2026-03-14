import { computed } from 'vue'
import { useStorage } from '@vueuse/core'

export const integrationCatalog = [
  { value: 'whatsapp', label: 'WhatsApp' },
  { value: 'telegram', label: 'Telegram' },
  { value: 'sms', label: 'SMS' },
  { value: 'push', label: 'Notificação push' },
  { value: 'webhook', label: 'Webhook' },
  { value: 'custom', label: 'Integração personalizada' }
]

export const integrationDirections = [
  { value: 'outgoing', label: 'Saída' },
  { value: 'incoming', label: 'Entrada' }
]

export const integrationTriggers = [
  { value: 'manual_conversation', label: 'Manual na conversa' },
  { value: 'conversation_created', label: 'Novo chamado' },
  { value: 'conversation_updated', label: 'Atualização de chamado' },
  { value: 'conversation_closed', label: 'Encerramento de chamado' },
  { value: 'message_created', label: 'Nova mensagem' }
]

const legacyDefaultActions = [
  {
    id: 'acao-whatsapp',
    name: 'Notificar no WhatsApp',
    method: 'POST',
    url: '',
    enabled: true,
    bodyTemplate:
      '{\n  "canal": "whatsapp",\n  "conversation_uuid": "{{conversation_uuid}}",\n  "reference_number": "{{reference_number}}",\n  "contact_name": "{{contact_name}}",\n  "contact_email": "{{contact_email}}"\n}'
  },
  {
    id: 'acao-telegram',
    name: 'Notificar no Telegram',
    method: 'POST',
    url: '',
    enabled: false,
    bodyTemplate:
      '{\n  "canal": "telegram",\n  "conversation_uuid": "{{conversation_uuid}}",\n  "reference_number": "{{reference_number}}"\n}'
  }
]

const defaultActions = legacyDefaultActions.map((action) => ({
  ...action,
  integration: action.bodyTemplate.includes('"canal": "telegram"') ? 'telegram' : 'whatsapp',
  direction: 'outgoing',
  triggers: ['manual_conversation'],
  description: ''
}))

const normalizeAction = (action) => ({
  id: action.id || `action-${Date.now()}`,
  name: action.name || 'Nova ação',
  description: action.description || '',
  integration: action.integration || 'custom',
  direction: action.direction || 'outgoing',
  triggers: Array.isArray(action.triggers) && action.triggers.length
    ? action.triggers
    : ['manual_conversation'],
  method: action.method || 'POST',
  url: action.url || '',
  enabled: typeof action.enabled === 'boolean' ? action.enabled : true,
  bodyTemplate:
    action.bodyTemplate || '{\n  "conversation_uuid": "{{conversation_uuid}}"\n}'
})

export function useIntegrationActions() {
  const actions = useStorage('canalgov_integration_actions', defaultActions)

  const normalizedActions = computed({
    get: () => (actions.value || []).map(normalizeAction),
    set: (value) => {
      actions.value = value.map(normalizeAction)
    }
  })

  const groupedActions = computed(() => {
    return integrationCatalog
      .map((integration) => ({
        ...integration,
        actions: normalizedActions.value.filter((action) => action.integration === integration.value)
      }))
      .filter((group) => group.actions.length > 0)
  })

  return {
    actions: normalizedActions,
    groupedActions,
    integrationCatalog,
    integrationDirections,
    integrationTriggers
  }
}
