<template>
  <Dialog :open="openAIKeyPrompt" @update:open="openAIKeyPrompt = false">
    <DialogContent class="sm:max-w-lg">
      <DialogHeader class="space-y-2">
        <DialogTitle>{{ $t('ai.enterOpenAIAPIKey') }}</DialogTitle>
        <DialogDescription>
          {{
            $t('ai.apiKey.description', {
              provider: 'OpenAI'
            })
          }}
        </DialogDescription>
      </DialogHeader>
      <Form v-slot="{ handleSubmit }" as="" keep-values :validation-schema="formSchema">
        <form id="apiKeyForm" @submit="handleSubmit($event, updateProvider)">
          <FormField v-slot="{ componentField }" name="apiKey">
            <FormItem>
              <FormLabel>{{ $t('globals.terms.apiKey') }}</FormLabel>
              <FormControl>
                <Input type="text" placeholder="sk-am1RLw7XUWGX.." v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
        </form>
        <DialogFooter>
          <Button
            type="submit"
            form="apiKeyForm"
            :is-loading="isOpenAIKeyUpdating"
            :disabled="isOpenAIKeyUpdating"
          >
            {{ $t('globals.messages.save') }}
          </Button>
        </DialogFooter>
      </Form>
    </DialogContent>
  </Dialog>

  <div class="text-foreground bg-background">
    <!-- Fullscreen editor -->
    <Dialog :open="isEditorFullscreen" @update:open="isEditorFullscreen = false">
      <DialogContent
        class="max-w-[60%] max-h-[75%] h-[70%] bg-card text-card-foreground p-4 flex flex-col"
        :class="{ '!bg-private': messageType === 'private_note' }"
        @escapeKeyDown="isEditorFullscreen = false"
        :hide-close-button="true"
      >
        <ReplyBoxContent
          v-if="isEditorFullscreen"
          :isFullscreen="true"
          :aiPrompts="aiPrompts"
          :responseTemplates="responseTemplates"
          :isSending="isSending"
          :isDraftLoading="isDraftLoading"
          :uploadingFiles="uploadingFiles"
          :uploadedFiles="mediaFiles"
          v-model:htmlContent="htmlContent"
          v-model:textContent="textContent"
          v-model:to="to"
          v-model:cc="cc"
          v-model:bcc="bcc"
          v-model:emailErrors="emailErrors"
          v-model:messageType="messageType"
          v-model:selectedResponseChannel="selectedResponseChannel"
          v-model:enabledResponseChannels="enabledResponseChannels"
          v-model:showBcc="showBcc"
          v-model:mentions="mentions"
          @toggleFullscreen="isEditorFullscreen = !isEditorFullscreen"
          @send="processSend"
          @fileUpload="handleFileUpload"
          @fileDelete="handleFileDelete"
          @aiPromptSelected="handleAiPromptSelected"
          class="h-full flex-grow"
        />
      </DialogContent>
    </Dialog>

    <!-- Main Editor non-fullscreen -->
    <div
      class="bg-background text-card-foreground box m-2 px-2 pt-2 flex flex-col"
      :class="{ '!bg-private': messageType === 'private_note' }"
      v-if="!isEditorFullscreen"
    >
      <ReplyBoxContent
        ref="replyBoxContentRef"
        :isFullscreen="false"
        :aiPrompts="aiPrompts"
        :responseTemplates="responseTemplates"
        :isSending="isSending"
        :isDraftLoading="isDraftLoading"
        :uploadingFiles="uploadingFiles"
        :uploadedFiles="mediaFiles"
        v-model:htmlContent="htmlContent"
        v-model:textContent="textContent"
        v-model:to="to"
        v-model:cc="cc"
        v-model:bcc="bcc"
        v-model:emailErrors="emailErrors"
        v-model:messageType="messageType"
        v-model:selectedResponseChannel="selectedResponseChannel"
        v-model:enabledResponseChannels="enabledResponseChannels"
        v-model:showBcc="showBcc"
        v-model:mentions="mentions"
        @toggleFullscreen="isEditorFullscreen = !isEditorFullscreen"
        @send="processSend"
        @fileUpload="handleFileUpload"
        @fileDelete="handleFileDelete"
        @aiPromptSelected="handleAiPromptSelected"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed, toRaw } from 'vue'
import { useStorage } from '@vueuse/core'
import { handleHTTPError } from '@/utils/http'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { MACRO_CONTEXT } from '@/constants/conversation'
import { useUserStore } from '@/stores/user'
import { useDraftManager } from '@/composables/useDraftManager'
import api from '@/api'
import { useI18n } from 'vue-i18n'
import { useConversationStore } from '@/stores/conversation'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { useEmitter } from '@/composables/useEmitter'
import { useFileUpload } from '@/composables/useFileUpload'
import ReplyBoxContent from '@/features/conversation/ReplyBoxContent.vue'
import { UserTypeAgent } from '@/constants/user'
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage
} from '@/components/ui/form'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

const formSchema = toTypedSchema(
  z.object({
    apiKey: z.string().min(1, 'API key is required')
  })
)

const { t } = useI18n()
const conversationStore = useConversationStore()
const emitter = useEmitter()
const userStore = useUserStore()

// Setup file upload composable
const {
  uploadingFiles,
  handleFileUpload,
  handleFileDelete,
  mediaFiles,
  clearMediaFiles,
  setMediaFiles
} = useFileUpload({
  linkedModel: 'messages'
})

// Setup draft management composable
const currentDraftKey = computed(() => conversationStore.current?.uuid || null)
const {
  htmlContent,
  textContent,
  isLoading: isDraftLoading,
  clearDraft,
  loadedAttachments,
  loadedMacroActions
} = useDraftManager(currentDraftKey, mediaFiles)

// Rest of existing state
const openAIKeyPrompt = ref(false)
const isOpenAIKeyUpdating = ref(false)
const isEditorFullscreen = ref(false)
const isSending = ref(false)
const messageType = useStorage('replyBoxMessageType', 'reply')
const selectedResponseChannel = useStorage('replyBoxSelectedResponseChannel', 'email')
const enabledResponseChannels = useStorage('replyBoxEnabledResponseChannels', [
  'email',
  'whatsapp',
  'telegram',
  'sms',
  'official_communication'
])
const to = ref('')
const cc = ref('')
const bcc = ref('')
const showBcc = ref(false)
const emailErrors = ref([])
const aiPrompts = ref([])
const responseTemplates = ref([])
const replyBoxContentRef = ref(null)
const mentions = ref([])

/**
 * Fetches AI prompts from the server.
 */
const fetchAiPrompts = async () => {
  try {
    const resp = await api.getAiPrompts()
    aiPrompts.value = resp.data.data
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  }
}

fetchAiPrompts()

const fetchResponseTemplates = async () => {
  try {
    const teamID = conversationStore.current?.assigned_team_id || null
    const resp = await api.getTemplates('response', teamID ? {
      team_id: teamID,
      include_global: true
    } : {})
    responseTemplates.value = resp.data.data || []
  } catch (error) {
    responseTemplates.value = []
  }
}

fetchResponseTemplates()

/**
 * Handles the AI prompt selection event.
 * Sends the selected prompt key and the current text content to the server for completion.
 * Sets the response as the new content in the editor.
 * @param {String} key - The key of the selected AI prompt
 */
const handleAiPromptSelected = async (key) => {
  try {
    const resp = await api.aiCompletion({
      prompt_key: key,
      content: textContent.value
    })
    htmlContent.value = resp.data.data.replace(/\n/g, '<br>')
  } catch (error) {
    // Check if user needs to enter OpenAI API key and has permission to do so.
    if (error.response?.status === 400 && userStore.can('ai:manage')) {
      openAIKeyPrompt.value = true
    }
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  }
}

/**
 * updateProvider updates the OpenAI API key.
 * @param {Object} values - The form values containing the API key
 */
const updateProvider = async (values) => {
  try {
    isOpenAIKeyUpdating.value = true
    await api.updateAIProvider({ api_key: values.apiKey, provider: 'openai' })
    openAIKeyPrompt.value = false
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      description: t('globals.messages.savedSuccessfully', {
        name: t('globals.terms.apiKey')
      })
    })
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    isOpenAIKeyUpdating.value = false
  }
}

/**
 * Returns true if the editor has text content.
 */
const hasTextContent = computed(() => {
  return textContent.value.trim().length > 0
})

/**
 * Processes the send action.
 */
const processSend = async () => {
  let hasMessageSendingErrored = false
  let pendingUUID = null
  let previousHTMLContent = ''
  let previousTextContent = ''
  let previousConversationPreview = null
  isEditorFullscreen.value = false
  try {
    // Send message if there is text content in the editor or media files are attached.
    if (hasTextContent.value > 0 || mediaFiles.value.length > 0) {
      const conversationUUID = conversationStore.current.uuid
      const message = htmlContent.value
      const plainTextMessage = textContent.value
      previousHTMLContent = message
      previousTextContent = plainTextMessage
      const currentConversation = conversationStore.conversations.data.find(
        conversation => conversation.uuid === conversationUUID
      )
      previousConversationPreview = currentConversation
        ? {
            last_message: currentConversation.last_message,
            last_message_at: currentConversation.last_message_at,
            last_message_sender: currentConversation.last_message_sender
          }
        : null
      const isPrivateMessage = messageType.value === 'private_note'
      const fallbackTo = conversationStore.current?.contact?.email
        ? [conversationStore.current.contact.email.trim()].filter((email) => email)
        : []
      const parsedCC = cc.value
        ? cc.value
            .split(',')
            .map((email) => email.trim())
            .filter((email) => email)
        : []
      const parsedBCC = bcc.value
        ? bcc.value
            .split(',')
            .map((email) => email.trim())
            .filter((email) => email)
        : []
      const parsedTo = to.value
        ? to.value
            .split(',')
            .map((email) => email.trim())
            .filter((email) => email)
        : fallbackTo
      if (!isPrivateMessage && selectedResponseChannel.value === 'email' && parsedTo.length === 0) {
        emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
          variant: 'destructive',
          description: t('globals.messages.required', {
            name: t('globals.emails.to')
          })
        })
        return
      }
      const meta = {}
      if (!isPrivateMessage && parsedTo.length > 0) meta.to = parsedTo
      if (!isPrivateMessage && parsedCC.length > 0) meta.cc = parsedCC
      if (!isPrivateMessage && parsedBCC.length > 0) meta.bcc = parsedBCC

      const author = {
        id: userStore.userID,
        first_name: userStore.firstName,
        last_name: userStore.lastName,
        avatar_url: userStore.avatar,
        type: 'agent'
      }

      pendingUUID = conversationStore.addPendingMessage(
        conversationUUID,
        message,
        isPrivateMessage,
        author,
        [...mediaFiles.value],
        plainTextMessage,
        meta
      )

      htmlContent.value = ''
      textContent.value = ''
      isSending.value = true

      conversationStore.updateConversationLastMessage(
        conversationUUID,
        { text_content: plainTextMessage, created_at: new Date().toISOString(), sender_type: 'agent' }
      )

      const response = await api.sendMessage(conversationUUID, {
        sender_type: UserTypeAgent,
        private: isPrivateMessage,
        message: message,
        attachments: mediaFiles.value.map((file) => file.id),
        // Include mentions only for private notes
        mentions: isPrivateMessage ? mentions.value : [],
        cc: parsedCC,
        bcc: parsedBCC,
        to: parsedTo
      })

      if (response?.data?.data) {
        conversationStore.replacePendingMessage(conversationUUID, pendingUUID, response.data.data)
        conversationStore.updateConversationLastMessage(conversationUUID, response.data.data)
      }
    }

    // Apply macro actions if any, for macro errors just show toast and clear the editor.
    const macroID = conversationStore.getMacro(MACRO_CONTEXT.REPLY)?.id
    const macroActions = conversationStore.getMacro(MACRO_CONTEXT.REPLY)?.actions || []
    if (macroID > 0 && macroActions.length > 0) {
      try {
        await api.applyMacro(conversationStore.current.uuid, macroID, macroActions)
      } catch (error) {
        emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
          variant: 'destructive',
          description: handleHTTPError(error).message
        })
      }
    }
  } catch (error) {
    hasMessageSendingErrored = true
    if (conversationStore.current?.uuid && pendingUUID) {
      conversationStore.removePendingMessage(conversationStore.current.uuid, pendingUUID)
    }
    if (conversationStore.current?.uuid && previousConversationPreview) {
      conversationStore.updateConversationProp({
        uuid: conversationStore.current.uuid,
        prop: 'last_message',
        value: previousConversationPreview.last_message
      })
      conversationStore.updateConversationProp({
        uuid: conversationStore.current.uuid,
        prop: 'last_message_at',
        value: previousConversationPreview.last_message_at
      })
      conversationStore.updateConversationProp({
        uuid: conversationStore.current.uuid,
        prop: 'last_message_sender',
        value: previousConversationPreview.last_message_sender
      })
    }
    htmlContent.value = previousHTMLContent
    textContent.value = previousTextContent
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    // If API has NOT errored clear state.
    if (hasMessageSendingErrored === false) {
      // Clear draft from backend.
      clearDraft(currentDraftKey.value)

      // Clear macro for this conversation reply.
      conversationStore.resetMacro(MACRO_CONTEXT.REPLY)

      // Clear media files.
      clearMediaFiles()

      // Clear any email errors.
      emailErrors.value = []

      // Clear mentions.
      mentions.value = []
    }
    isSending.value = false
  }
}

/**
 * Watches for changes in the conversation's macro id and update message content.
 */
watch(
  () => conversationStore.getMacro('reply').id,
  (newId) => {
    // No macro set.
    if (!newId) return

    // If macro has message content, set it in the editor.
    if (conversationStore.getMacro('reply').message_content) {
      htmlContent.value = conversationStore.getMacro('reply').message_content
    }
  },
  { deep: true }
)

/**
 * Watch loaded macro actions from draft and update conversation store.
 */
watch(
  loadedMacroActions,
  (actions) => {
    if (actions.length > 0) {
      conversationStore.setMacroActions([...toRaw(actions)], MACRO_CONTEXT.REPLY)
    }
  },
  { deep: true }
)

/**
 * Watch for loaded attachments from draft and restore them to mediaFiles.
 */
watch(
  loadedAttachments,
  (attachments) => {
    if (attachments.length > 0) {
      setMediaFiles([...attachments])
    }
  },
  { deep: true }
)

// Initialize to, cc, and bcc fields with the current conversation's values.
watch(
  () => conversationStore.currentCC,
  (newVal) => {
    cc.value = newVal?.join(', ') || ''
  },
  { deep: true, immediate: true }
)

watch(
  () => conversationStore.currentTo,
  (newVal) => {
    to.value = newVal?.join(', ') || ''
  },
  { immediate: true }
)

watch(
  () => conversationStore.currentBCC,
  (newVal) => {
    const newBcc = newVal?.join(', ') || ''
    bcc.value = newBcc
    // Only show BCC field if it has content
    if (newBcc.length > 0) {
      showBcc.value = true
    }
  },
  { deep: true, immediate: true }
)

// Clear media files and reset macro when conversation changes.
watch(
  () => conversationStore.current?.uuid,
  () => {
    clearMediaFiles()
    conversationStore.resetMacro(MACRO_CONTEXT.REPLY)
    // Focus editor on conversation change
    setTimeout(() => {
      replyBoxContentRef.value?.focus()
    }, 100)
  }
)

watch(
  () => conversationStore.current?.assigned_team_id,
  () => {
    fetchResponseTemplates()
  },
  { immediate: false }
)

watch(
  enabledResponseChannels,
  (channels) => {
    if (!channels.includes('official_communication')) {
      enabledResponseChannels.value = [...channels, 'official_communication']
    }
  },
  { deep: true, immediate: true }
)
</script>
