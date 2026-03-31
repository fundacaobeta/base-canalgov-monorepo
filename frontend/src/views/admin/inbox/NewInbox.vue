<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">{{ t('admin.inbox.title') }}</h1>
      <p class="text-muted-foreground">{{ t('admin.inbox.new.stepConfigureDescription') }}</p>
    </div>

    <div class="max-w-2xl">
      <Card>
        <CardHeader>
          <div class="flex items-center gap-4">
            <div class="rounded-lg bg-primary/10 p-2 text-primary">
              <Mail class="h-6 w-6" />
            </div>
            <div>
              <CardTitle>{{ t('globals.terms.configure') }} {{ t('globals.terms.email') }}</CardTitle>
              <CardDescription>{{ t('admin.inbox.form.managedEmailDescription') }}</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="handleCreate" class="space-y-4">
            <div class="grid gap-2">
              <Label>{{ t('admin.inbox.form.mailboxName') }}</Label>
              <Input v-model="form.name" :placeholder="t('globals.terms.name')" required />
            </div>

            <div class="space-y-4 border-t pt-4">
              <h3 class="font-medium">{{ t('admin.inbox.form.managedEmail') }}</h3>

              <div class="grid gap-4 md:grid-cols-2">
                <div class="grid gap-2">
                  <Label>{{ t('admin.inbox.form.emailPart') }}</Label>
                  <Input v-model="form.emailPart" placeholder="suporte" required />
                  <p class="text-[10px] text-muted-foreground">{{ t('admin.inbox.form.emailPartHelp') }}</p>
                </div>

                <div class="grid gap-2">
                  <Label>{{ t('admin.inbox.form.domain') }}</Label>
                  <Select v-model="form.domainId">
                    <SelectTrigger>
                      <SelectValue :placeholder="t('admin.inbox.form.selectDomain')" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem v-for="domain in domains" :key="domain.id" :value="domain.id.toString()">
                        {{ domain.domain }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div class="rounded-xl border border-border/70 p-4">
                <div class="text-sm font-medium">{{ t('admin.inbox.form.generatedAddress') }}</div>
                <div class="mt-2 font-mono text-sm">
                  {{ managedEmailAddress || t('admin.inbox.form.fillNameAndSelectDomain') }}
                </div>
              </div>
            </div>

            <div class="flex items-center justify-between rounded-lg bg-muted/30 p-4">
              <div class="space-y-0.5">
                <Label class="text-base">{{ t('admin.inbox.form.enabled') }}</Label>
                <p class="text-xs text-muted-foreground">{{ t('admin.inbox.form.enabledHelp') }}</p>
              </div>
              <Switch v-model:checked="form.enabled" />
            </div>

            <div class="flex items-center justify-between rounded-lg bg-muted/30 p-4">
              <div class="space-y-0.5">
                <Label class="text-base">{{ t('admin.inbox.form.csat') }}</Label>
                <p class="text-xs text-muted-foreground">{{ t('admin.inbox.form.csatHelp') }}</p>
              </div>
              <Switch v-model:checked="form.csat_enabled" />
            </div>

            <div class="flex justify-end gap-2 border-t pt-4">
              <Button type="button" variant="ghost" @click="router.push({ name: 'inbox-list' })">
                {{ t('globals.terms.cancel') }}
              </Button>
              <Button type="submit" :disabled="loading || !isManagedEmailValid">
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
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Loader2, Mail } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import api from '@/api'

const { t } = useI18n()
const router = useRouter()
const { showErrorToast } = useAdminErrorToast()
const loading = ref(false)
const domains = ref([])
const emailLocalPartPattern = /^[A-Za-z0-9.!#$%&'*+/=?^_`{|}~-]+$/

const form = reactive({
  name: '',
  enabled: true,
  csat_enabled: false,
  emailPart: '',
  domainId: ''
})

const selectedDomain = computed(() =>
  domains.value.find((item) => item.id.toString() === form.domainId)
)

const managedEmailAddress = computed(() => {
  if (!form.emailPart || !selectedDomain.value?.domain) return ''
  return `${form.emailPart.trim()}@${selectedDomain.value.domain}`
})

const isManagedEmailValid = computed(() =>
  emailLocalPartPattern.test(form.emailPart.trim()) && Boolean(selectedDomain.value?.domain)
)

const handleCreate = async () => {
  if (!managedEmailAddress.value || !isManagedEmailValid.value) {
    showErrorToast(new Error('Preencha um endereço de e-mail válido para a caixa de entrada.'))
    return
  }

  try {
    loading.value = true

    const payload = {
      name: form.name.trim(),
      channel: 'email',
      from: managedEmailAddress.value,
      enabled: form.enabled,
      csat_enabled: form.csat_enabled,
      config: {
        receive_mode: 'managed',
        delivery_provider: 'canalgov_managed',
        managed_domain_id: form.domainId,
        managed_domain: selectedDomain.value?.domain || '',
        managed_local_part: form.emailPart.trim(),
        managed_email_address: managedEmailAddress.value
      }
    }

    await api.createInbox(payload)
    router.push({ name: 'inbox-list' })
  } catch (error) {
    showErrorToast(error)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const resp = await api.getMailDomainsSettings()
    domains.value = (resp.data.data?.['mail.domains'] || []).filter((item) => item.enabled)

    if (!form.domainId && domains.value.length > 0) {
      const defaultDomain = domains.value.find((item) => item.is_default) || domains.value[0]
      form.domainId = defaultDomain.id.toString()
    }
  } catch (error) {
    showErrorToast(error)
  }
})
</script>
