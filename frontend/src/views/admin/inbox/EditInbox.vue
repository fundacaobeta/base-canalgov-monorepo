<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">{{ t('admin.inbox.title') }}</h1>
        <p class="text-muted-foreground">{{ t('globals.messages.edit') }}: {{ inbox?.name }}</p>
      </div>
    </div>

    <div v-if="isLoading" class="flex justify-center py-20">
      <Loader2 class="h-8 w-8 animate-spin text-primary" />
    </div>

    <div v-else class="max-w-3xl space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>{{ t('globals.terms.general') }}</CardTitle>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="handleUpdate" class="space-y-4">
            <div class="grid gap-2">
              <Label>{{ t('globals.terms.name') }}</Label>
              <Input v-model="form.name" required />
            </div>

            <div class="grid gap-2">
              <Label>{{ t('globals.terms.identifier') }}</Label>
              <Input v-model="form.from" readonly class="bg-muted" />
            </div>

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

            <div class="flex justify-end pt-4">
              <Button type="submit" :disabled="isSaving">
                <Loader2 v-if="isSaving" class="mr-2 h-4 w-4 animate-spin" />
                {{ t('globals.messages.save') }}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <!-- Provider Config (WhatsApp/Telegram/Email Custom) -->
      <Card v-if="hasProviderConfig">
        <CardHeader>
          <CardTitle>{{ t('admin.inbox.form.additionalConfig') }}</CardTitle>
          <CardDescription>{{ t('admin.inbox.form.additionalConfigHelp') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="grid gap-2">
            <Textarea 
              v-model="configString" 
              class="font-mono text-xs min-h-[200px]" 
              placeholder='{ "key": "value" }'
            />
          </div>
          <div class="flex justify-end pt-4">
            <Button variant="outline" @click="saveProviderConfig" :disabled="isSaving">
              {{ t('globals.messages.save') }}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Textarea } from '@/components/ui/textarea'
import api from '@/api'

const props = defineProps({
  id: { type: String, required: true }
})

const { t } = useI18n()
const isLoading = ref(true)
const isSaving = ref(false)
const inbox = ref(null)
const configString = ref('{}')

const form = reactive({
  name: '',
  from: '',
  enabled: true,
  csat_enabled: false
})

const hasProviderConfig = computed(() => {
  return ['whatsapp', 'telegram', 'email', 'webhook'].includes(inbox.value?.channel)
})

const loadInbox = async () => {
  try {
    isLoading.value = true
    const resp = await api.getInbox(props.id)
    inbox.value = resp.data.data
    Object.assign(form, {
      name: inbox.value.name,
      from: inbox.value.from,
      enabled: inbox.value.enabled,
      csat_enabled: inbox.value.csat_enabled
    })
    configString.value = JSON.stringify(inbox.value.config || {}, null, 2)
  } catch (err) {
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const handleUpdate = async () => {
  try {
    isSaving.value = true
    await api.updateInbox(props.id, form)
    await loadInbox()
  } catch (err) {
    console.error(err)
  } finally {
    isSaving.value = false
  }
}

const saveProviderConfig = async () => {
  try {
    isSaving.value = true
    const config = JSON.parse(configString.value)
    await api.updateInbox(props.id, { ...form, config })
    await loadInbox()
  } catch (err) {
    alert('JSON inválido')
  } finally {
    isSaving.value = false
  }
}

onMounted(loadInbox)
</script>
