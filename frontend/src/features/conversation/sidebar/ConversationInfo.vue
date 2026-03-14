<template>
  <div class="space-y-4">
    <div class="flex flex-col" v-if="conversation.subject">
      <p class="font-medium">{{ $t('globals.terms.subject') }}</p>
      <Skeleton v-if="conversationStore.conversation.loading" class="w-32 h-4" />
      <p v-else>
        {{ conversation.subject }}
      </p>
    </div>

    <div class="flex flex-col">
      <p class="font-medium">{{ $t('globals.terms.referenceNumber') }}</p>
      <Skeleton v-if="conversationStore.conversation.loading" class="w-32 h-4" />
      <p v-else>
        {{ conversation.reference_number }}
      </p>
    </div>
    <div class="flex flex-col">
      <p class="font-medium">{{ $t('globals.terms.initiatedAt') }}</p>
      <Skeleton v-if="conversationStore.conversation.loading" class="w-32 h-4" />
      <DateTimeMeta v-if="conversation.created_at" :value="conversation.created_at" />
      <p v-else>-</p>
    </div>

    <div class="flex flex-col">
      <div class="flex justify-start items-center space-x-2">
        <p class="font-medium">{{ $t('globals.terms.firstReplyAt') }}</p>
        <SlaBadge
          v-if="conversation.first_response_deadline_at"
          :dueAt="conversation.first_response_deadline_at"
          :actualAt="conversation.first_reply_at"
          :key="`${conversation.uuid}-${conversation.first_response_deadline_at}-${conversation.first_reply_at}`"
        />
      </div>
      <Skeleton v-if="conversationStore.conversation.loading" class="w-32 h-4" />
      <div v-else>
        <DateTimeMeta v-if="conversation.first_reply_at" :value="conversation.first_reply_at" />
        <p v-else>-</p>
      </div>
    </div>

    <div class="flex flex-col">
      <div class="flex justify-start items-center space-x-2">
        <p class="font-medium">{{ $t('globals.terms.resolvedAt') }}</p>
        <SlaBadge
          v-if="conversation.resolution_deadline_at"
          :dueAt="conversation.resolution_deadline_at"
          :actualAt="conversation.resolved_at"
          :key="`${conversation.uuid}-${conversation.resolution_deadline_at}-${conversation.resolved_at}`"
        />
      </div>
      <Skeleton v-if="conversationStore.conversation.loading" class="w-32 h-4" />
      <div v-else>
        <DateTimeMeta v-if="conversation.resolved_at" :value="conversation.resolved_at" />
        <p v-else>-</p>
      </div>
    </div>

    <div class="flex flex-col">
      <div class="flex justify-start items-center space-x-2">
        <p class="font-medium">{{ $t('globals.terms.lastReplyAt') }}</p>
        <SlaBadge
          v-if="conversation.next_response_deadline_at"
          :dueAt="conversation.next_response_deadline_at"
          :actualAt="conversation.next_response_met_at"
          :key="`${conversation.uuid}-${conversation.next_response_deadline_at}-${conversation.next_response_met_at}`"
        />
      </div>
      <Skeleton v-if="conversationStore.conversation.loading" class="w-32 h-4" />
      <DateTimeMeta v-if="conversation.last_reply_at" :value="conversation.last_reply_at" />
      <p v-else>-</p>
    </div>

    <div class="flex flex-col" v-if="conversation.closed_at">
      <p class="font-medium">{{ $t('globals.terms.closedAt') }}</p>
      <Skeleton v-if="conversationStore.conversation.loading" class="w-32 h-4" />
      <DateTimeMeta v-else :value="conversation.closed_at" />
    </div>

    <div class="flex flex-col" v-if="conversation.sla_policy_name">
      <p class="font-medium">{{ $t('globals.terms.slaPolicy') }}</p>
      <div>
        <p>
          {{ conversation.sla_policy_name }}
        </p>
      </div>
    </div>

    <CustomAttributes
      v-if="customAttributeStore.conversationAttributeOptions.length > 0"
      :loading="conversationStore.conversation.loading"
      :attributes="customAttributeStore.conversationAttributeOptions"
      :custom-attributes="conversation.custom_attributes || {}"
      @update:setattributes="updateCustomAttributes"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import SlaBadge from '@/features/sla/SlaBadge.vue'
import { useConversationStore } from '@/stores/conversation'
import { Skeleton } from '@/components/ui/skeleton'
import CustomAttributes from '@/features/conversation/sidebar/CustomAttributes.vue'
import { useCustomAttributeStore } from '@/stores/customAttributes'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useEmitter } from '@/composables/useEmitter'
import { handleHTTPError } from '@/utils/http'
import api from '@/api'
import { useI18n } from 'vue-i18n'
import DateTimeMeta from '@/components/datetime/DateTimeMeta.vue'

const emitter = useEmitter()
const { t } = useI18n()
const customAttributeStore = useCustomAttributeStore()
const conversationStore = useConversationStore()
const conversation = computed(() => conversationStore.current)
customAttributeStore.fetchCustomAttributes()

const updateCustomAttributes = async (attributes) => {
  let previousAttributes = conversationStore.current.custom_attributes
  try {
    conversationStore.current.custom_attributes = attributes
    await api.updateConversationCustomAttribute(conversation.value.uuid, attributes)
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
    conversationStore.current.custom_attributes = previousAttributes
  }
}
</script>
