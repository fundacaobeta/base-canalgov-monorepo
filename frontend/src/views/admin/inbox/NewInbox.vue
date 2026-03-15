<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">{{ t('admin.inbox.title') }}</h1>
        <p class="text-muted-foreground">{{ t('admin.inbox.new.stepChannelDescription') }}</p>
      </div>
    </div>

    <!-- Step 1: Select Channel -->
    <div v-if="step === 1" class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      <Card
        v-for="channel in channels"
        :key="channel.id"
        class="cursor-pointer hover:border-primary/50 transition-colors group"
        @click="selectChannel(channel.id)"
      >
        <CardHeader>
          <div class="flex items-center gap-4">
            <div class="p-2 rounded-lg bg-primary/5 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
              <component :is="channel.icon" class="h-6 w-6" />
            </div>
            <div>
              <CardTitle class="text-lg">{{ t('globals.terms.' + channel.id) }}</CardTitle>
              <CardDescription>{{ t('admin.inbox.new.channel' + capitalize(channel.id) + 'Subtitle') }}</CardDescription>
            </div>
          </div>
        </CardHeader>
      </Card>
    </div>

    <!-- Step 2: Configure Channel -->
    <div v-if="step === 2" class="max-w-2xl">
      <Card>
        <CardHeader>
          <div class="flex items-center gap-4">
            <Button variant="ghost" size="icon" @click="step = 1">
              <ArrowLeft class="h-4 w-4" />
            </Button>
            <div>
              <CardTitle>{{ t('globals.terms.configure') }} {{ t('globals.terms.' + selectedChannel) }}</CardTitle>
              <CardDescription>{{ t('admin.inbox.new.stepConfigureDescription') }}</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="handleCreate" class="space-y-4">
            <div class="grid gap-2">
              <Label>{{ t('admin.inbox.form.mailboxName') }}</Label>
              <Input v-model="form.name" :placeholder="t('globals.terms.name')" required />
            </div>

            <!-- Email Specific Fields -->
            <template v-if="selectedChannel === 'email'">
              <div class="space-y-4 pt-4 border-t">
                <h3 class="font-medium">{{ t('admin.inbox.form.managedEmail') }}</h3>
                <div class="grid gap-4 md:grid-cols-2">
                  <div class="grid gap-2">
                    <Label>{{ t('admin.inbox.form.emailPart') }}</Label>
                    <Input v-model="form.emailPart" placeholder="suporte" />
                    <p class="text-[10px] text-muted-foreground">{{ t('admin.inbox.form.emailPartHelp') }}</p>
                  </div>
                  <div class="grid gap-2">
                    <Label>{{ t('admin.inbox.form.domain') }}</Label>
                    <Select v-model="form.domainId">
                      <SelectTrigger><SelectValue /></SelectTrigger>
                      <SelectContent>
                        <SelectItem v-for="d in domains" :key="d.id" :value="d.id.toString()">{{ d.name }}</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </div>
            </template>

            <!-- Other Channels -->
            <template v-else>
              <div class="grid gap-2">
                <Label>{{ t('globals.terms.identifier') }}</Label>
                <Input v-model="form.from" placeholder="Ex: 5511999998888" required />
              </div>
            </template>

            <div class="flex items-center justify-between p-4 bg-muted/30 rounded-lg">
              <div class="space-y-0.5">
                <Label class="text-base">{{ t('admin.inbox.form.enabled') }}</Label>
                <p class="text-xs text-muted-foreground">{{ t('admin.inbox.form.enabledHelp') }}</p>
              </div>
              <Switch v-model:checked="form.enabled" />
            </div>

            <div class="flex items-center justify-between p-4 bg-muted/30 rounded-lg">
              <div class="space-y-0.5">
                <Label class="text-base">{{ t('admin.inbox.form.csat') }}</Label>
                <p class="text-xs text-muted-foreground">{{ t('admin.inbox.form.csatHelp') }}</p>
              </div>
              <Switch v-model:checked="form.csat_enabled" />
            </div>

            <div class="flex justify-end gap-2 pt-4 border-t">
              <Button type="button" variant="ghost" @click="step = 1">{{ t('globals.terms.cancel') }}</Button>
              <Button type="submit" :disabled="loading">
                <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
                {{ t('globals.messages.create', { name: t('admin.inbox.title').toLowerCase() }) }}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { 
  Mail, MessageSquare, Send, Smartphone, 
  Globe, ArrowLeft, Loader2, Bot, Bell 
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import api from '@/api'

const { t } = useI18n()
const router = useRouter()
const step = ref(1)
const loading = ref(false)
const selectedChannel = ref(null)
const domains = ref([])

const channels = [
  { id: 'email', icon: Mail },
  { id: 'whatsapp', icon: MessageSquare },
  { id: 'telegram', icon: Bot },
  { id: 'sms', icon: Send },
  { id: 'push', icon: Bell },
  { id: 'webhook', icon: Globe }
]

const form = reactive({
  name: '',
  from: '',
  channel: '',
  enabled: true,
  csat_enabled: false,
  emailPart: '',
  domainId: null
})

const capitalize = (s) => s.charAt(0).toUpperCase() + s.slice(1)

const selectChannel = (id) => {
  selectedChannel.value = id
  form.channel = id
  step.value = 2
}

const handleCreate = async () => {
  try {
    loading.value = true
    const payload = { ...form }
    if (selectedChannel.value === 'email') {
      const domain = domains.value.find(d => d.id.toString() === form.domainId)
      payload.from = `${form.emailPart}@${domain.name}`
    }
    await api.createInbox(payload)
    router.push({ name: 'inbox-list' })
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const resp = await api.getSettings() // Assuming domains are in settings or a specific endpoint
    // domains.value = resp.data.data.mail_domains || []
  } catch (err) {
    console.error(err)
  }
})
</script>
