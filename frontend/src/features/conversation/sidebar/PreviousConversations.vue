<template>
  <div
    v-if="
      conversationStore.current?.previous_conversations?.length === 0 ||
      conversationStore.conversation?.loading
    "
    class="text-center text-sm text-muted-foreground py-4"
  >
    {{ $t('conversation.sidebar.noPreviousConvo') }}
  </div>
  <div v-else class="space-y-1">
    <router-link
      v-for="conversation in conversationStore.current.previous_conversations"
      :key="conversation.uuid"
      :to="{
        name: 'inbox-conversation',
        params: {
          uuid: conversation.uuid,
          type: 'assigned'
        }
      }"
      class="block p-2 rounded hover:bg-muted"
    >
      <div class="flex flex-wrap items-start justify-between gap-1">
        <div class="flex flex-col flex-1 min-w-[120px]">
          <Tooltip>
            <TooltipTrigger asChild>
              <span class="font-medium text-sm truncate max-w-[300px]">
                {{ conversation.subject }}
              </span>
            </TooltipTrigger>
            <TooltipContent>
              {{ conversation.subject }}
            </TooltipContent>
          </Tooltip>
          <span class="text-xs text-muted-foreground truncate max-w-[400px]">
            {{ conversation.last_message }}
          </span>
        </div>
        <Tooltip>
          <TooltipTrigger asChild>
            <div class="flex flex-col items-end text-xs text-muted-foreground flex-shrink-0">
              <DateTimeMeta
                v-if="conversation.created_at"
                :value="conversation.created_at"
                :show-absolute="false"
                compact
              />
              <DateTimeMeta
                v-if="conversation.last_message_at"
                :value="conversation.last_message_at"
                :show-absolute="false"
                compact
              />
            </div>
          </TooltipTrigger>
          <TooltipContent>
            <div class="space-y-1 text-xs">
              <p>
                {{ $t('globals.terms.createdAt') }}:
                {{ formatFullTimestamp(new Date(conversation.created_at)) }}
              </p>
              <p v-if="conversation.last_message_at">
                {{ $t('globals.terms.lastMessageAt') }}:
                {{ formatFullTimestamp(new Date(conversation.last_message_at)) }}
              </p>
            </div>
          </TooltipContent>
        </Tooltip>
      </div>
    </router-link>
  </div>
</template>

<script setup>
import { useConversationStore } from '@/stores/conversation'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { formatFullTimestamp } from '@/utils/datetime'
import DateTimeMeta from '@/components/datetime/DateTimeMeta.vue'

const conversationStore = useConversationStore()
</script>
