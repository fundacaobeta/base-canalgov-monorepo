<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.domain.title')"
      :description="$t('admin.domain.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('admin.domain.title') }]"
    >
      <template #actions>
        <Button @click="openCreate">
          <Plus class="h-4 w-4 mr-1.5" aria-hidden="true" />
          {{ $t('globals.messages.new', { name: $t('globals.terms.domain') }) }}
        </Button>
      </template>
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <div class="space-y-3">
          <div
            v-if="domains.length === 0"
            class="box flex flex-col items-center justify-center gap-2 p-10 text-center text-muted-foreground"
          >
            <Globe class="h-8 w-8 opacity-40" />
            <p class="text-sm">{{ $t('admin.domain.empty') }}</p>
          </div>

          <div
            v-for="domain in domains"
            :key="domain.id"
            class="box flex flex-col gap-4 p-5 lg:flex-row lg:items-start lg:justify-between"
          >
            <div class="space-y-2">
              <div class="flex flex-wrap items-center gap-2">
                <span class="font-medium">{{ domain.name }}</span>
                <span class="rounded-full bg-muted px-2 py-0.5 text-xs font-mono">{{ domain.domain }}</span>
                <span class="rounded-full bg-muted px-2 py-0.5 text-xs">{{ providerLabel(domain.provider) }}</span>
                <span v-if="domain.is_default" class="rounded-full bg-foreground px-2 py-0.5 text-xs text-background">
                  {{ $t('globals.terms.default') }}
                </span>
                <span
                  class="rounded-full px-2 py-0.5 text-xs"
                  :class="domain.enabled ? 'bg-emerald-100 text-emerald-700' : 'bg-zinc-100 text-zinc-600'"
                >
                  {{ domain.enabled ? $t('globals.terms.active') : $t('globals.terms.inactive') }}
                </span>
              </div>
              <p class="text-sm text-muted-foreground">{{ inboundLabel(domain.inbound_strategy) }}</p>
              <p v-if="domain.notes" class="text-sm text-muted-foreground">{{ domain.notes }}</p>
            </div>

            <div class="flex gap-2">
              <Button variant="outline" size="sm" @click="openEdit(domain)">
                {{ $t('globals.messages.edit', { name: '' }).trim() }}
              </Button>
              <Button variant="destructive" size="sm" @click="confirmRemove(domain)">
                {{ $t('globals.messages.delete', { name: '' }).trim() }}
              </Button>
            </div>
          </div>
        </div>
      </template>

      <template #help>
        <p>{{ $t('admin.domain.help') }}</p>
        <p>{{ $t('admin.domain.help2') }}</p>
      </template>
    </AdminPageWithHelp>

    <!-- Create / Edit Dialog -->
    <Dialog v-model:open="dialogOpen">
      <DialogContent class="max-w-2xl">
        <DialogHeader>
          <DialogTitle>{{ isEditing ? $t('admin.domain.editTitle') : $t('globals.messages.new', { name: $t('globals.terms.domain') }) }}</DialogTitle>
        </DialogHeader>

        <div class="space-y-4 py-2">
          <div class="grid gap-4 lg:grid-cols-2">
            <div>
              <label class="text-sm font-medium">{{ $t('globals.terms.name') }}</label>
              <Input v-model="draft.name" class="mt-2" :placeholder="t('globals.terms.name')" />
            </div>
            <div>
              <label class="text-sm font-medium">{{ $t('globals.terms.domain') }}</label>
              <Input v-model="draft.domain" class="mt-2" placeholder="atendimento.gov.br" />
            </div>
          </div>

          <div class="grid gap-4 lg:grid-cols-2">
            <div>
              <label class="text-sm font-medium">{{ $t('globals.terms.provider') }}</label>
              <Select v-model="draft.provider">
                <SelectTrigger class="mt-2">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="ses">{{ t('admin.inbox.form.provider.ses') }}</SelectItem>
                  <SelectItem value="sns">{{ t('admin.inbox.form.provider.sns') }}</SelectItem>
                  <SelectItem value="mailgun">{{ t('admin.inbox.form.provider.mailgun') }}</SelectItem>
                  <SelectItem value="sendgrid">{{ t('admin.inbox.form.provider.sendgrid') }}</SelectItem>
                  <SelectItem value="self_hosted">{{ t('admin.inbox.form.provider.self_hosted') }}</SelectItem>
                  <SelectItem value="custom">{{ t('globals.terms.custom') }}</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div>
              <label class="text-sm font-medium">{{ $t('admin.domain.inboundStrategy') }}</label>
              <Select v-model="draft.inbound_strategy">
                <SelectTrigger class="mt-2">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="managed_mailbox">{{ t('admin.inbox.form.provider.managed') }}</SelectItem>
                  <SelectItem value="smtp_imap">{{ t('admin.inbox.form.provider.smtp_imap') }}</SelectItem>
                  <SelectItem value="webhook">{{ t('admin.inbox.form.provider.webhook') }}</SelectItem>
                  <SelectItem value="docker_mailserver">{{ t('admin.inbox.form.provider.docker') }}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div>
            <label class="text-sm font-medium">{{ $t('globals.terms.notes') }}</label>
            <Textarea
              v-model="draft.notes"
              class="mt-2 min-h-24"
              :placeholder="t('globals.terms.notes')"
            />
          </div>

          <div class="grid gap-4 lg:grid-cols-2">
            <div class="flex items-center justify-between rounded-xl border border-border/70 p-4">
              <div>
                <p class="text-sm font-medium">{{ $t('globals.terms.enabled') }}</p>
                <p class="text-xs text-muted-foreground">{{ $t('admin.domain.enabledDescription') }}</p>
              </div>
              <Switch v-model:checked="draft.enabled" />
            </div>

            <div class="flex items-center justify-between rounded-xl border border-border/70 p-4">
              <div>
                <p class="text-sm font-medium">{{ $t('globals.terms.default') }}</p>
                <p class="text-xs text-muted-foreground">{{ $t('admin.domain.defaultDescription') }}</p>
              </div>
              <Switch v-model:checked="draft.is_default" />
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" @click="dialogOpen = false">{{ $t('globals.messages.cancel') }}</Button>
          <Button @click="saveDomain" :isLoading="isSaving">{{ $t('globals.messages.save') }}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <AlertDialog v-model:open="deleteDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ $t('globals.messages.areYouAbsolutelySure') }}</AlertDialogTitle>
          <AlertDialogDescription>
            {{ $t('globals.messages.deletionConfirmation', { name: domainToDelete?.name || $t('globals.terms.domain') }) }}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{{ $t('globals.messages.cancel') }}</AlertDialogCancel>
          <AlertDialogAction @click="removeDomain" class="bg-destructive text-destructive-foreground hover:bg-destructive/90">
            {{ $t('globals.messages.delete') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { Plus, Globe } from 'lucide-vue-next'
import api from '@/api'
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
import AdminPageHeader from '@/components/layout/AdminPageHeader.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle
} from '@/components/ui/alert-dialog'
import { useAdminErrorToast } from '@/composables/useAdminErrorToast'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const { showErrorToast, showSuccessToast } = useAdminErrorToast()
const isSaving = ref(false)
const domains = ref([])
const dialogOpen = ref(false)
const deleteDialogOpen = ref(false)
const domainToDelete = ref(null)

const blankDraft = () => ({
  id: '',
  name: '',
  domain: '',
  provider: 'ses',
  inbound_strategy: 'managed_mailbox',
  enabled: true,
  is_default: false,
  notes: ''
})

const draft = ref(blankDraft())
const isEditing = computed(() => Boolean(draft.value.id))

const fetchDomains = async () => {
  try {
    const resp = await api.getMailDomainsSettings()
    domains.value = resp.data.data?.['mail.domains'] || []
  } catch (error) {
    showErrorToast(error)
  }
}

onMounted(fetchDomains)

const persistDomains = async (nextDomains) => {
  await api.updateMailDomainsSettings({ 'mail.domains': nextDomains })
  domains.value = nextDomains
}

const openCreate = () => {
  draft.value = blankDraft()
  dialogOpen.value = true
}

const openEdit = (domain) => {
  draft.value = { ...domain }
  dialogOpen.value = true
}

const confirmRemove = (domain) => {
  domainToDelete.value = domain
  deleteDialogOpen.value = true
}

const saveDomain = async () => {
  try {
    isSaving.value = true
    const payload = {
      ...draft.value,
      id: draft.value.id || `domain-${Date.now()}`
    }
    let nextDomains = domains.value.filter((item) => item.id !== payload.id)
    if (payload.is_default) {
      nextDomains = nextDomains.map((item) => ({ ...item, is_default: false }))
    }
    nextDomains.unshift(payload)
    await persistDomains(nextDomains)
    dialogOpen.value = false
    showSuccessToast(t('globals.messages.savedSuccessfully', { name: t('globals.terms.domain') }))
  } catch (error) {
    showErrorToast(error)
  } finally {
    isSaving.value = false
  }
}

const removeDomain = async () => {
  if (!domainToDelete.value) return
  try {
    isSaving.value = true
    await persistDomains(domains.value.filter((item) => item.id !== domainToDelete.value.id))
    domainToDelete.value = null
    showSuccessToast(t('globals.messages.deletedSuccessfully', { name: t('globals.terms.domain') }))
  } catch (error) {
    showErrorToast(error)
  } finally {
    isSaving.value = false
  }
}

const providerLabel = (value) =>
  ({
    ses: 'AWS SES',
    sns: 'AWS SNS',
    mailgun: 'Mailgun',
    sendgrid: 'SendGrid',
    self_hosted: t('admin.domain.providerSelfHosted'),
    custom: t('globals.terms.custom')
  }[value] || value)

const inboundLabel = (value) =>
  ({
    managed_mailbox: t('admin.domain.inboundManaged'),
    smtp_imap: t('admin.inbox.new.emailCustom'),
    webhook: t('admin.domain.inboundWebhook'),
    docker_mailserver: t('admin.domain.inboundDocker')
  }[value] || value)
</script>
