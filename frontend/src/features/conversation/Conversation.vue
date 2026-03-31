<template>
  <div class="flex flex-col h-screen">
    <!-- Header -->
    <div class="h-12 flex-shrink-0 px-2 border-b flex items-center justify-between">
      <div>
        <span v-if="!conversationStore.conversation.loading">
          {{ conversationStore.currentContactName }}
        </span>
        <Skeleton class="w-[130px] h-6" v-else />
      </div>
      <div>
        <div class="flex items-center gap-2">
          <DropdownMenu v-if="availableActions.length">
            <DropdownMenuTrigger>
              <div
                class="flex items-center space-x-1 cursor-pointer border px-2 py-1 rounded text-sm"
                v-if="!conversationStore.conversation.loading"
              >
                <span class="font-medium inline-block">{{ $t('globals.terms.actions') }}</span>
              </div>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <template v-for="group in availableActionGroups" :key="group.value">
                <DropdownMenuLabel>{{ group.label }}</DropdownMenuLabel>
                <DropdownMenuItem
                  v-for="action in group.actions"
                  :key="action.id"
                  @click="executeIntegrationAction(action)"
                >
                  {{ action.name }}
                </DropdownMenuItem>
                <DropdownMenuSeparator />
              </template>
            </DropdownMenuContent>
          </DropdownMenu>

          <div
            v-if="!conversationStore.conversation.loading"
            class="flex items-center gap-3 rounded-md border px-3 py-2"
          >
            <div class="flex flex-col leading-tight">
              <span class="text-[11px] uppercase tracking-wide text-muted-foreground">
                Status do chamado
              </span>
              <Badge
                :variant="getConversationStatusBadgeVariant(conversationStore.current?.status)"
                class="w-fit"
              >
                {{ translateConversationStatus(conversationStore.current?.status, t) }}
              </Badge>
            </div>

            <div v-if="conversationStore.current?.priority" class="flex flex-col leading-tight">
              <span class="text-[11px] uppercase tracking-wide text-muted-foreground">
                Prioridade
              </span>
              <Badge
                :variant="getConversationPriorityBadgeVariant(conversationStore.current?.priority)"
                class="w-fit"
              >
                {{ translateConversationPriority(conversationStore.current?.priority) }}
              </Badge>
            </div>

            <DropdownMenu>
              <DropdownMenuTrigger>
                <Button size="sm" variant="outline">
                  Alterar status
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem
                  v-if="conversationStore.current?.status !== CONVERSATION_DEFAULT_STATUSES.OPEN"
                  class="text-destructive focus:text-destructive"
                  @click="handleUpdateStatus(CONVERSATION_DEFAULT_STATUSES.OPEN)"
                >
                  Reabrir chamado
                </DropdownMenuItem>
                <DropdownMenuItem
                  v-if="conversationStore.current?.status !== CONVERSATION_DEFAULT_STATUSES.RESOLVED"
                  @click="handleUpdateStatus(CONVERSATION_DEFAULT_STATUSES.RESOLVED)"
                >
                  Marcar como resolvido
                </DropdownMenuItem>
                <DropdownMenuItem
                  v-if="conversationStore.current?.status !== CONVERSATION_DEFAULT_STATUSES.SNOOZED"
                  @click="handleUpdateStatus(CONVERSATION_DEFAULT_STATUSES.SNOOZED)"
                >
                  Adiar chamado
                </DropdownMenuItem>
                <DropdownMenuItem
                  v-if="conversationStore.current?.status !== CONVERSATION_DEFAULT_STATUSES.CLOSED"
                  @click="handleUpdateStatus(CONVERSATION_DEFAULT_STATUSES.CLOSED)"
                >
                  Fechar chamado
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
          <Skeleton class="w-[220px] h-10 rounded-md" v-else />
        </div>
      </div>
    </div>

    <!-- Messages & reply box -->
    <div class="flex flex-col flex-grow overflow-hidden">
      <MessageList class="flex-1 overflow-y-auto" />
      <div class="sticky bottom-0">
        <ReplyBox />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useConversationStore } from '@/stores/conversation'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import MessageList from '@/features/conversation/message/MessageList.vue'
import ReplyBox from './ReplyBox.vue'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { CONVERSATION_DEFAULT_STATUSES } from '@/constants/conversation'
import { useEmitter } from '@/composables/useEmitter'
import { useIntegrationActions } from '@/composables/useIntegrationActions'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { useI18n } from 'vue-i18n'
import { getConversationStatusBadgeVariant, translateConversationStatus } from '@/utils/conversationStatus'
import { getConversationPriorityBadgeVariant, translateConversationPriority } from '@/utils/conversationPriority'

const conversationStore = useConversationStore()
const { t } = useI18n()
const emitter = useEmitter()
const { groupedActions } = useIntegrationActions()
const availableActionGroups = computed(() =>
  groupedActions.value
    .map((group) => ({
      ...group,
      actions: group.actions.filter(
        (item) =>
          item.enabled &&
          item.direction === 'outgoing' &&
          item.triggers.includes('manual_conversation')
      )
    }))
    .filter((group) => group.actions.length > 0)
)
const availableActions = computed(() => availableActionGroups.value.flatMap((group) => group.actions))

const handleUpdateStatus = (status) => {
  if (status === CONVERSATION_DEFAULT_STATUSES.SNOOZED) {
    emitter.emit(EMITTER_EVENTS.SET_NESTED_COMMAND, {
      command: 'snooze',
      open: true
    })
    return
  }
  conversationStore.updateStatus(status)
}

const buildPayload = (action) => {
  const replacements = {
    '{{conversation_uuid}}': conversationStore.current?.uuid || '',
    '{{reference_number}}': conversationStore.current?.reference_number || '',
    '{{contact_name}}': conversationStore.currentContactName || '',
    '{{contact_email}}': conversationStore.current?.contact?.email || '',
    '{{contact_phone}}': conversationStore.current?.contact?.phone_number || '',
    '{{assigned_user_name}}': conversationStore.current?.assigned_user?.first_name || '',
    '{{assigned_team_name}}': conversationStore.current?.assigned_team?.name || ''
  }

  return Object.entries(replacements).reduce((body, [key, value]) => {
    return body.replaceAll(key, value)
  }, action.bodyTemplate || '')
}

const executeIntegrationAction = async (action) => {
  if (!action.url) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: t('admin.integration.actions.notConfiguredWarning')
      })
    return
  }

  try {
    const options = {
      method: action.method || 'POST',
      headers: {
        'Content-Type': 'application/json'
      }
    }

    if (!['GET', 'HEAD'].includes(options.method)) {
      options.body = buildPayload(action)
    }

    await fetch(action.url, options)
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      description: `Ação "${action.name}" executada.`
    })
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: `Falha ao executar a ação "${action.name}".`
    })
  }
}
</script>
