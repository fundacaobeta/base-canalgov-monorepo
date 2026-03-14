<template>
  <div class="space-y-6">
    <div class="box p-6">
      <h2 class="text-xl font-semibold">Domínios</h2>
      <p class="mt-2 text-sm text-muted-foreground">
        Cadastre os domínios que poderão ser usados nas caixas de e-mail gerenciadas pelo CanalGov.
      </p>
    </div>

    <div class="box space-y-5 p-6">
      <div class="grid gap-4 lg:grid-cols-2">
        <div>
          <label class="text-sm font-medium">Nome</label>
          <Input v-model="draft.name" class="mt-2" placeholder="Domínio principal" />
        </div>
        <div>
          <label class="text-sm font-medium">Domínio</label>
          <Input v-model="draft.domain" class="mt-2" placeholder="atendimento.gov.br" />
        </div>
      </div>

      <div class="grid gap-4 lg:grid-cols-2">
        <div>
          <label class="text-sm font-medium">Provedor</label>
          <Select v-model="draft.provider">
            <SelectTrigger class="mt-2">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ses">AWS SES</SelectItem>
              <SelectItem value="sns">AWS SNS</SelectItem>
              <SelectItem value="mailgun">Mailgun</SelectItem>
              <SelectItem value="sendgrid">SendGrid</SelectItem>
              <SelectItem value="self_hosted">Servidor próprio</SelectItem>
              <SelectItem value="custom">Customizado</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div>
          <label class="text-sm font-medium">Estratégia de entrada</label>
          <Select v-model="draft.inbound_strategy">
            <SelectTrigger class="mt-2">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="managed_mailbox">Caixa gerenciada</SelectItem>
              <SelectItem value="smtp_imap">SMTP/IMAP próprio</SelectItem>
              <SelectItem value="webhook">Webhook inbound</SelectItem>
              <SelectItem value="docker_mailserver">Servidor via Docker</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div>
        <label class="text-sm font-medium">Notas</label>
        <Textarea
          v-model="draft.notes"
          class="mt-2 min-h-28"
          placeholder="Observações sobre DNS, MX, provedor ou operação."
        />
      </div>

      <div class="grid gap-4 lg:grid-cols-2">
        <div class="flex items-center justify-between rounded-xl border border-border/70 p-4">
          <div>
            <p class="text-sm font-medium">Habilitado</p>
            <p class="text-xs text-muted-foreground">Disponível para novas caixas.</p>
          </div>
          <Switch v-model:checked="draft.enabled" />
        </div>

        <div class="flex items-center justify-between rounded-xl border border-border/70 p-4">
          <div>
            <p class="text-sm font-medium">Padrão</p>
            <p class="text-xs text-muted-foreground">Usado como sugestão inicial nas caixas gerenciadas.</p>
          </div>
          <Switch v-model:checked="draft.is_default" />
        </div>
      </div>

      <div class="flex gap-3">
        <Button @click="saveDomain" :isLoading="isSaving">{{ isEditing ? 'Salvar domínio' : 'Adicionar domínio' }}</Button>
        <Button variant="outline" @click="resetDraft">Limpar</Button>
      </div>
    </div>

    <div class="space-y-4">
      <div
        v-for="domain in domains"
        :key="domain.id"
        class="box flex flex-col gap-4 p-5 lg:flex-row lg:items-start lg:justify-between"
      >
        <div class="space-y-2">
          <div class="flex flex-wrap items-center gap-2">
            <span class="font-medium">{{ domain.name }}</span>
            <span class="rounded-full bg-muted px-2 py-0.5 text-xs">{{ domain.domain }}</span>
            <span class="rounded-full bg-muted px-2 py-0.5 text-xs">{{ providerLabel(domain.provider) }}</span>
            <span v-if="domain.is_default" class="rounded-full bg-foreground px-2 py-0.5 text-xs text-background">Padrão</span>
            <span
              class="rounded-full px-2 py-0.5 text-xs"
              :class="domain.enabled ? 'bg-emerald-100 text-emerald-700' : 'bg-zinc-100 text-zinc-600'"
            >
              {{ domain.enabled ? 'Ativo' : 'Inativo' }}
            </span>
          </div>
          <p class="text-sm text-muted-foreground">{{ inboundLabel(domain.inbound_strategy) }}</p>
          <p v-if="domain.notes" class="text-sm text-muted-foreground">{{ domain.notes }}</p>
        </div>

        <div class="flex gap-2">
          <Button variant="outline" @click="editDomain(domain)">Editar</Button>
          <Button variant="destructive" @click="removeDomain(domain.id)">Remover</Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import api from '@/api'
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
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents'
import { handleHTTPError } from '@/utils/http'

const emitter = useEmitter()
const isSaving = ref(false)
const domains = ref([])

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
    domains.value = resp.data.data?.domains || []
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  }
}

onMounted(fetchDomains)

const persistDomains = async (nextDomains) => {
  await api.updateMailDomainsSettings({ domains: nextDomains })
  domains.value = nextDomains
}

const resetDraft = () => {
  draft.value = blankDraft()
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
    resetDraft()
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, { description: 'Domínio salvo.' })
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    isSaving.value = false
  }
}

const editDomain = (domain) => {
  draft.value = { ...domain }
}

const removeDomain = async (id) => {
  try {
    isSaving.value = true
    await persistDomains(domains.value.filter((item) => item.id !== id))
    if (draft.value.id === id) resetDraft()
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, { description: 'Domínio removido.' })
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    isSaving.value = false
  }
}

const providerLabel = (value) => ({
  ses: 'AWS SES',
  sns: 'AWS SNS',
  mailgun: 'Mailgun',
  sendgrid: 'SendGrid',
  self_hosted: 'Servidor próprio',
  custom: 'Customizado'
}[value] || value)

const inboundLabel = (value) => ({
  managed_mailbox: 'Caixa gerenciada pelo CanalGov',
  smtp_imap: 'SMTP/IMAP próprio',
  webhook: 'Webhook inbound',
  docker_mailserver: 'Servidor de e-mail via Docker'
}[value] || value)
</script>
