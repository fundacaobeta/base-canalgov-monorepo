<template>
  <div class="flex h-full flex-col" :class="{ 'max-h-[600px]': !isFullscreen }">
    <div
      class="flex items-center justify-between gap-3"
      :class="{ 'mb-4': !isFullscreen, 'border-b border-border pb-4': isFullscreen }"
    >
      <Tabs v-model="messageType">
        <TabsList class="h-10 rounded-full border border-border bg-muted/50 p-1">
          <TabsTrigger
            value="private_note"
            class="rounded-full px-4 py-1.5 text-sm transition-colors"
            :class="{ 'bg-background text-foreground shadow-sm': messageType === 'private_note' }"
          >
            {{ $t('globals.terms.privateNote') }}
          </TabsTrigger>
          <TabsTrigger
            value="reply"
            class="rounded-full px-4 py-1.5 text-sm transition-colors"
            :class="{ 'bg-background text-foreground shadow-sm': messageType === 'reply' }"
          >
            {{ $t('globals.terms.reply') }}
          </TabsTrigger>
        </TabsList>
      </Tabs>

      <Button type="button" class="text-muted-foreground" variant="ghost" @click="toggleFullscreen">
        <component :is="isFullscreen ? Minimize2 : Maximize2" />
      </Button>
    </div>

    <div
      v-if="messageType === 'reply'"
      :class="['space-y-3', isFullscreen ? 'border-b border-border pb-4' : 'mb-4']"
    >
      <div class="space-y-3 rounded-2xl border border-border/70 bg-background/70 px-4 py-3">
        <div class="flex flex-col gap-3 xl:flex-row xl:items-start xl:justify-between">
          <div class="min-w-0 flex-1 space-y-2">
            <div class="flex flex-wrap items-center gap-2">
              <span class="text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground">
                {{ t('replyBox.channel') }}
              </span>
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="channel in visibleResponseChannels"
                  :key="channel.value"
                  type="button"
                  class="rounded-full border px-3 py-1 text-sm transition-colors"
                  :class="
                    selectedResponseChannel === channel.value
                      ? 'border-foreground bg-foreground text-background'
                      : 'border-border bg-background text-muted-foreground hover:text-foreground'
                  "
                  @click="selectedResponseChannel = channel.value"
                >
                  {{ channel.label }}
                </button>
              </div>
            </div>

            <div class="flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-muted-foreground">
              <span>{{ selectedChannelDescription }}</span>
              <button
                type="button"
                class="font-medium text-foreground underline-offset-4 hover:underline"
                @click="showChannelSettings = !showChannelSettings"
              >
                {{ showChannelSettings ? t('replyBox.hideChannels') : t('replyBox.showChannels') }}
              </button>
              <button
                v-if="selectedResponseChannel === 'email'"
                type="button"
                class="font-medium text-foreground underline-offset-4 hover:underline"
                @click="showEmailEnvelope = !showEmailEnvelope"
              >
                {{ emailSummary }}
              </button>
            </div>
          </div>

          <div class="w-full xl:w-[320px]">
            <label class="mb-2 block text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground">
              {{ t('replyBox.template') }}
            </label>
            <div class="flex gap-2">
              <Select v-model="selectedTemplateId">
                <SelectTrigger class="bg-background">
                  <SelectValue :placeholder="t('replyBox.newMessage')" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="new-message">{{ t('replyBox.newMessage') }}</SelectItem>
                  <SelectItem
                    v-for="template in responseTemplates"
                    :key="template.id"
                    :value="String(template.id)"
                  >
                    {{ template.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
              <Button
                v-if="selectedTemplateId !== 'new-message'"
                type="button"
                variant="ghost"
                class="shrink-0"
                @click="clearSelectedTemplate"
              >
                {{ t('replyBox.clear') }}
              </Button>
            </div>
            <p class="mt-2 text-xs text-muted-foreground">
              {{ templateSummary }}
            </p>
          </div>
        </div>

        <Collapsible v-model:open="showChannelSettings">
          <CollapsibleContent>
            <div class="grid gap-2 border-t border-border/60 pt-3 sm:grid-cols-2 xl:grid-cols-3">
              <label
                v-for="channel in allResponseChannels"
                :key="`config-${channel.value}`"
                class="flex items-start gap-3 rounded-2xl border border-border/70 bg-muted/20 px-3 py-3"
              >
                <Checkbox
                  :checked="enabledResponseChannels.includes(channel.value)"
                  @update:checked="(checked) => toggleResponseChannel(channel.value, checked)"
                />
                <div class="space-y-1">
                  <div class="text-sm font-medium">{{ channel.label }}</div>
                  <p class="text-xs text-muted-foreground">{{ channel.description }}</p>
                </div>
              </label>
            </div>
          </CollapsibleContent>
        </Collapsible>

        <div class="flex items-start justify-between gap-3 border-t border-border/60 pt-3">
          <div class="space-y-1">
            <div class="text-sm font-medium">{{ t('replyBox.activeChannel') }}: {{ selectedResponseChannelLabel }}</div>
            <p class="text-xs text-muted-foreground">{{ selectedChannelHint }}</p>
          </div>
        </div>

        <Collapsible v-if="selectedResponseChannel === 'email'" v-model:open="showEmailEnvelope">
          <CollapsibleContent>
            <div class="grid gap-3 border-t border-border/60 pt-3">
              <div class="grid gap-2">
                <label class="text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground">{{ t('globals.emails.to') }}</label>
                <Input
                  v-model="to"
                  type="text"
                  :placeholder="t('replyBox.emailAddresess')"
                  @blur="validateEmails"
                />
              </div>

              <div class="grid gap-3 lg:grid-cols-[minmax(0,1fr)_auto]">
                <div class="grid gap-2">
                  <label class="text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground">{{ t('globals.emails.cc') }}</label>
                  <Input
                    v-model="cc"
                    type="text"
                    :placeholder="t('replyBox.emailAddresess')"
                    @blur="validateEmails"
                  />
                </div>

                <div class="flex items-end">
                  <Button type="button" size="sm" variant="ghost" @click="toggleBcc">
                    {{ showBcc ? t('replyBox.hideBcc') : t('replyBox.addBcc') }}
                  </Button>
                </div>
              </div>

              <div v-if="showBcc" class="grid gap-2">
                <label class="text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground">{{ t('globals.emails.bcc') }}</label>
                <Input
                  v-model="bcc"
                  type="text"
                  :placeholder="t('replyBox.emailAddresess')"
                  @blur="validateEmails"
                />
              </div>
            </div>
          </CollapsibleContent>
        </Collapsible>
      </div>
    </div>

    <div
      v-if="emailErrors.length > 0"
      class="mb-4 rounded-lg border border-destructive bg-destructive/10 px-3 py-2 text-destructive"
    >
      <p v-for="error in emailErrors" :key="error" class="text-sm">{{ error }}</p>
    </div>

    <div class="flex min-h-0 flex-1 flex-col overflow-hidden rounded-lg border border-border bg-background">
      <Editor
        ref="editorRef"
        v-model:htmlContent="htmlContent"
        v-model:textContent="textContent"
        :placeholder="editorPlaceholder"
        :aiPrompts="aiPrompts"
        :insertContent="insertContent"
        :autoFocus="true"
        :disabled="isDraftLoading"
        :enableMentions="messageType === 'private_note'"
        :getSuggestions="getSuggestions"
        @aiPromptSelected="handleAiPromptSelected"
        @send="handleSend"
        @mentionsChanged="handleMentionsChanged"
      />
    </div>

    <MacroActionsPreview
      v-if="conversationStore.getMacro(MACRO_CONTEXT.REPLY)?.actions?.length > 0"
      :actions="conversationStore.getMacro(MACRO_CONTEXT.REPLY).actions"
      :onRemove="(action) => conversationStore.removeMacroAction(action, MACRO_CONTEXT.REPLY)"
      class="mt-2"
    />

    <AttachmentsPreview
      v-if="uploadedFiles.length > 0 || uploadingFiles.length > 0"
      :attachments="uploadedFiles"
      :uploadingFiles="uploadingFiles"
      :onDelete="handleOnFileDelete"
      class="mt-2"
    />

    <ReplyBoxMenuBar
      class="mt-2 shrink-0"
      :isFullscreen="isFullscreen"
      :handleFileUpload="handleFileUpload"
      :isSending="isSending"
      :enableSend="enableSend"
      :handleSend="handleSend"
      @emojiSelect="handleEmojiSelect"
    />
  </div>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { Maximize2, Minimize2 } from 'lucide-vue-next'
import Editor from '@/components/editor/TextEditor.vue'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger
} from '@/components/ui/collapsible'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { MACRO_CONTEXT } from '@/constants/conversation'
import { useConversationStore } from '@/stores/conversation'
import AttachmentsPreview from '@/features/conversation/message/attachment/AttachmentsPreview.vue'
import MacroActionsPreview from '@/features/conversation/MacroActionsPreview.vue'
import ReplyBoxMenuBar from '@/features/conversation/ReplyBoxMenuBar.vue'
import { useI18n } from 'vue-i18n'
import { validateEmail } from '@/utils/strings'
import { useMacroStore } from '@/stores/macro'
import { useAgentsStore } from '@/stores/agents'
import { useTeamStore } from '@/stores/team'

const messageType = defineModel('messageType', { default: 'reply' })
const to = defineModel('to', { default: '' })
const cc = defineModel('cc', { default: '' })
const bcc = defineModel('bcc', { default: '' })
const showBcc = defineModel('showBcc', { default: false })
const selectedResponseChannel = defineModel('selectedResponseChannel', { default: 'email' })
const enabledResponseChannels = defineModel('enabledResponseChannels', {
  default: () => ['email', 'whatsapp', 'telegram', 'sms', 'official_communication']
})
const emailErrors = defineModel('emailErrors', { default: () => [] })
const htmlContent = defineModel('htmlContent', { default: '' })
const textContent = defineModel('textContent', { default: '' })
const mentions = defineModel('mentions', { default: () => [] })

const props = defineProps({
  isFullscreen: {
    type: Boolean,
    default: false
  },
  aiPrompts: {
    type: Array,
    required: true
  },
  responseTemplates: {
    type: Array,
    required: false,
    default: () => []
  },
  isSending: {
    type: Boolean,
    required: true
  },
  uploadingFiles: {
    type: Array,
    required: true
  },
  uploadedFiles: {
    type: Array,
    required: false,
    default: () => []
  },
  isDraftLoading: {
    type: Boolean,
    required: false,
    default: false
  }
})

const emit = defineEmits([
  'toggleFullscreen',
  'send',
  'fileUpload',
  'inlineImageUpload',
  'fileDelete',
  'aiPromptSelected'
])

const macroStore = useMacroStore()
const usersStore = useAgentsStore()
const teamStore = useTeamStore()
const conversationStore = useConversationStore()
const emitter = useEmitter()
const { t } = useI18n()
const insertContent = ref(null)
const editorRef = ref(null)
const showChannelSettings = ref(false)
const showEmailEnvelope = ref(false)
const selectedTemplateId = ref('new-message')

const allResponseChannels = computed(() => [
  {
    value: 'email',
    label: t('globals.terms.email'),
    description: t('replyBox.channels.email.description')
  },
  {
    value: 'whatsapp',
    label: t('globals.terms.whatsapp'),
    description: t('replyBox.channels.whatsapp.description')
  },
  {
    value: 'telegram',
    label: t('globals.terms.telegram'),
    description: t('replyBox.channels.telegram.description')
  },
  {
    value: 'sms',
    label: t('globals.terms.sms'),
    description: t('replyBox.channels.sms.description')
  },
  {
    value: 'official_communication',
    label: t('globals.terms.officialCommunication'),
    description: t('replyBox.channels.official.description')
  }
])

const visibleResponseChannels = computed(() =>
  allResponseChannels.value.filter((channel) => enabledResponseChannels.value.includes(channel.value))
)

const selectedChannel = computed(() => {
  return allResponseChannels.value.find((channel) => channel.value === selectedResponseChannel.value) || allResponseChannels.value[0]
})

const selectedResponseChannelLabel = computed(() => selectedChannel.value.label)
const selectedChannelDescription = computed(() => selectedChannel.value.description)

const selectedChannelHint = computed(() => {
  if (selectedResponseChannel.value === 'email') {
    return t('replyBox.hints.email')
  }
  if (selectedResponseChannel.value === 'official_communication') {
    return t('replyBox.hints.official')
  }
  return t('replyBox.hints.generic')
})

const templateSummary = computed(() => {
  if (selectedTemplateId.value === 'new-message') {
    return t('replyBox.templates.newMessageSummary')
  }

  const selectedTemplate = props.responseTemplates.find(
    (template) => String(template.id) === String(selectedTemplateId.value)
  )

  if (!selectedTemplate) {
    return t('replyBox.templates.newMessageSummary')
  }

  return selectedTemplate.team_name
    ? t('replyBox.templates.teamTemplate', { name: selectedTemplate.team_name })
    : t('replyBox.templates.globalTemplate')
})

const emailSummary = computed(() => {
  if (to.value || cc.value || bcc.value) {
    return t('replyBox.revisingRecipients')
  }
  return showEmailEnvelope.value ? t('replyBox.hideRecipients') : t('replyBox.addRecipients')
})

const editorPlaceholder = computed(() => {
  if (messageType.value === 'private_note') {
    return t('replyBox.placeholders.privateNote')
  }
  return t('replyBox.placeholders.reply', { channel: selectedResponseChannelLabel.value.toLowerCase() })
})

const getSuggestions = async (query) => {
  if (messageType.value !== 'private_note') {
    return []
  }

  await Promise.all([usersStore.fetchAgents(), teamStore.fetchTeams()])

  const q = query.toLowerCase()

  const users = usersStore.agents
    .filter((user) => user.enabled)
    .filter((user) => `${user.first_name} ${user.last_name}`.toLowerCase().includes(q))
    .map((user) => ({
      id: user.id,
      type: 'agent',
      label: `${user.first_name} ${user.last_name}`.trim(),
      avatar_url: user.avatar_url
    }))

  const teams = teamStore.teams
    .filter((team) => team.name.toLowerCase().includes(q))
    .map((team) => ({
      id: team.id,
      type: 'team',
      label: team.name,
      emoji: team.emoji
    }))

  return [...users, ...teams].slice(0, 25)
}

const handleMentionsChanged = (newMentions) => {
  mentions.value = newMentions
}

const toggleFullscreen = () => {
  emit('toggleFullscreen')
}

const toggleBcc = async () => {
  showBcc.value = !showBcc.value
  await nextTick()
  if (!showBcc.value) {
    bcc.value = ''
    await nextTick()
    validateEmails()
  }
}

const toggleResponseChannel = (channelValue, checked) => {
  const currentChannels = [...enabledResponseChannels.value]

  if (checked) {
    enabledResponseChannels.value = Array.from(new Set([...currentChannels, channelValue]))
    return
  }

  if (currentChannels.length === 1) {
    return
  }

  enabledResponseChannels.value = currentChannels.filter((value) => value !== channelValue)
  if (!enabledResponseChannels.value.includes(selectedResponseChannel.value)) {
    selectedResponseChannel.value = enabledResponseChannels.value[0] || 'email'
  }
}

const applyTemplate = (templateId) => {
  if (templateId === 'new-message') {
    return
  }

  const selectedTemplate = props.responseTemplates.find(
    (template) => String(template.id) === String(templateId)
  )

  if (!selectedTemplate) {
    return
  }

  htmlContent.value = selectedTemplate.body || ''
  textContent.value = selectedTemplate.body || ''
}

const clearSelectedTemplate = () => {
  selectedTemplateId.value = 'new-message'
}

const enableSend = computed(() => {
  return (
    (textContent.value.trim().length > 0 ||
      conversationStore.getMacro(MACRO_CONTEXT.REPLY)?.actions?.length > 0 ||
      props.uploadedFiles.length > 0) &&
    emailErrors.value.length === 0 &&
    !props.uploadingFiles.length &&
    !props.isDraftLoading
  )
})

const validateEmails = async () => {
  emailErrors.value = []
  await nextTick()

  if (selectedResponseChannel.value !== 'email') {
    return
  }

  const fields = ['to', 'cc', 'bcc']
  const values = { to: to.value, cc: cc.value, bcc: bcc.value }

  fields.forEach((field) => {
    const invalid = values[field]
      .split(',')
      .map((email) => email.trim())
      .filter((email) => email && !validateEmail(email))

    if (invalid.length) {
      emailErrors.value.push(`${t('replyBox.invalidEmailsIn')} '${field}': ${invalid.join(', ')}`)
    }
  })
}

const handleSend = async () => {
  await validateEmails()
  if (emailErrors.value.length > 0) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: t('globals.messages.correctEmailErrors')
    })
    return
  }
  emit('send')
}

const handleFileUpload = (event) => {
  emit('fileUpload', event)
}

const handleOnFileDelete = (uuid) => {
  emit('fileDelete', uuid)
}

const handleEmojiSelect = (emoji) => {
  insertContent.value = undefined
  nextTick(() => (insertContent.value = emoji))
}

const handleAiPromptSelected = (key) => {
  emit('aiPromptSelected', key)
}

watch(
  messageType,
  (newType) => {
    if (newType === 'reply') {
      macroStore.setCurrentView('replying')
    } else {
      macroStore.setCurrentView('adding_private_note')
    }

    setTimeout(() => {
      editorRef.value?.focus()
    }, 50)
  },
  { immediate: true }
)

watch(
  enabledResponseChannels,
  (channels) => {
    if (!channels.includes(selectedResponseChannel.value)) {
      selectedResponseChannel.value = channels[0] || 'email'
    }
  },
  { deep: true, immediate: true }
)

watch(selectedTemplateId, (templateId) => {
  if (templateId !== 'new-message') {
    applyTemplate(templateId)
  }
})

watch(
  () => props.responseTemplates,
  (templates) => {
    if (
      messageType.value !== 'reply' ||
      selectedTemplateId.value !== 'new-message' ||
      htmlContent.value.trim() ||
      textContent.value.trim()
    ) {
      return
    }

    const defaultTemplate = templates.find((template) => template.is_default)
    if (defaultTemplate) {
      selectedTemplateId.value = String(defaultTemplate.id)
    }
  },
  { immediate: true, deep: true }
)

watch(selectedResponseChannel, (channel) => {
  if (channel !== 'email') {
    showEmailEnvelope.value = false
    showBcc.value = false
    bcc.value = ''
    emailErrors.value = []
  }
})

watch(
  () => [to.value, cc.value, bcc.value],
  ([nextTo, nextCC, nextBCC]) => {
    showEmailEnvelope.value = Boolean(nextTo || nextCC || nextBCC)
  },
  { immediate: true }
)

const focus = () => {
  editorRef.value?.focus()
}

defineExpose({ focus })
</script>
