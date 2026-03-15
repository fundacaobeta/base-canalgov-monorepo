<template>
  <div class="mx-auto max-w-5xl space-y-6">
    <Card class="border-border/70 shadow-sm">
      <CardHeader>
        <CardTitle>{{ t('account.appPreferences') }}</CardTitle>
        <CardDescription>
          {{ t('account.preferencesDescription') }}
        </CardDescription>
      </CardHeader>

      <CardContent class="space-y-4">
        <div class="grid gap-4 md:grid-cols-2">
          <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
            <div class="flex items-start justify-between gap-4">
              <div>
                <p class="text-sm font-medium">{{ t('account.preferences.mainSidebar') }}</p>
                <p class="mt-1 text-sm text-muted-foreground">
                  {{ t('account.preferences.mainSidebarDescription') }}
                </p>
              </div>
              <Switch v-model:checked="mainSidebarOpen" />
            </div>
          </div>

          <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
            <div class="flex items-start justify-between gap-4">
              <div>
                <p class="text-sm font-medium">{{ t('account.preferences.conversationSidebar') }}</p>
                <p class="mt-1 text-sm text-muted-foreground">
                  {{ t('account.preferences.conversationSidebarDescription') }}
                </p>
              </div>
              <Switch v-model:checked="conversationSidebarOpen" />
            </div>
          </div>

          <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
            <div class="space-y-3">
              <div>
                <p class="text-sm font-medium">{{ t('account.preferences.defaultReplyType') }}</p>
                <p class="mt-1 text-sm text-muted-foreground">
                  {{ t('account.preferences.defaultReplyTypeDescription') }}
                </p>
              </div>
              <Select v-model="replyBoxMessageType">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="note">{{ t('account.preferences.replyType.note') }}</SelectItem>
                  <SelectItem value="reply">{{ t('account.preferences.replyType.reply') }}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
            <div class="space-y-3">
              <div>
                <p class="text-sm font-medium">{{ t('account.preferences.defaultReplyChannel') }}</p>
                <p class="mt-1 text-sm text-muted-foreground">
                  {{ t('account.preferences.defaultReplyChannelDescription') }}
                </p>
              </div>
              <Select v-model="selectedResponseChannel">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="email">{{ $t('globals.terms.email') }}</SelectItem>
                  <SelectItem value="official">{{ $t('globals.terms.officialCommunication') }}</SelectItem>
                  <SelectItem value="whatsapp">{{ $t('globals.terms.whatsapp') }}</SelectItem>
                  <SelectItem value="telegram">{{ $t('globals.terms.telegram') }}</SelectItem>
                  <SelectItem value="sms">{{ $t('globals.terms.sms') }}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>

    <Card class="border-border/70 shadow-sm">
      <CardHeader>
        <CardTitle>{{ t('account.preferences.enabledChannels') }}</CardTitle>
        <CardDescription>
          {{ t('account.preferences.enabledChannelsDescription') }}
        </CardDescription>
      </CardHeader>

      <CardContent class="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
        <div
          v-for="channel in channelOptions"
          :key="channel.value"
          class="rounded-2xl border border-border/70 bg-muted/20 p-4"
        >
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="text-sm font-medium">{{ channel.label }}</p>
              <p class="mt-1 text-sm text-muted-foreground">{{ channel.description }}</p>
            </div>
            <Switch
              :checked="enabledResponseChannels.includes(channel.value)"
              @update:checked="toggleChannel(channel.value, $event)"
            />
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup>
import { useStorage } from '@vueuse/core'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'

const { t } = useI18n()

const mainSidebarOpen = useStorage('mainSidebarOpen', true)
const conversationSidebarOpen = useStorage('conversationSidebarOpen', true)
const replyBoxMessageType = useStorage('replyBoxMessageType', 'reply')
const selectedResponseChannel = useStorage('replyBoxSelectedResponseChannel', 'email')
const enabledResponseChannels = useStorage('replyBoxEnabledResponseChannels', [
  'email',
  'official',
  'whatsapp',
  'telegram',
  'sms'
])

const channelOptions = computed(() => [
  {
    value: 'email',
    label: t('globals.terms.email'),
    description: t('account.preferences.channelDescription.email')
  },
  {
    value: 'official',
    label: t('globals.terms.officialCommunication'),
    description: t('account.preferences.channelDescription.official')
  },
  {
    value: 'whatsapp',
    label: t('globals.terms.whatsapp'),
    description: t('account.preferences.channelDescription.whatsapp')
  },
  {
    value: 'telegram',
    label: t('globals.terms.telegram'),
    description: t('account.preferences.channelDescription.telegram')
  },
  {
    value: 'sms',
    label: t('globals.terms.sms'),
    description: t('account.preferences.channelDescription.sms')
  }
])

const toggleChannel = (value, checked) => {
  const next = new Set(enabledResponseChannels.value)
  if (checked) {
    next.add(value)
  } else {
    next.delete(value)
  }

  if (next.size === 0) {
    next.add('email')
  }

  enabledResponseChannels.value = Array.from(next)

  if (!next.has(selectedResponseChannel.value)) {
    selectedResponseChannel.value = enabledResponseChannels.value[0]
  }
}
</script>
