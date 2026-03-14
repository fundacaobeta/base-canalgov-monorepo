<template>
  <div class="mb-5">
    <CustomBreadcrumb :links="breadcrumbLinks" />
  </div>
  <div class="space-y-10">
    <div class="mt-10">
      <Stepper class="flex w-full items-start gap-2" v-model="currentStep">
        <StepperItem
          v-for="step in steps"
          :key="step.step"
          v-slot="{ state }"
          class="relative flex w-full flex-col items-center justify-center"
          :step="step.step"
        >
          <StepperSeparator
            v-if="step.step !== steps[steps.length - 1].step"
            class="absolute left-[calc(50%+20px)] right-[calc(-50%+10px)] top-5 block h-0.5 shrink-0 rounded-full bg-muted group-data-[state=completed]:bg-primary"
          />
          <div>
            <Button
              :variant="state === 'completed' || state === 'active' ? 'default' : 'outline'"
              size="icon"
              class="z-10 rounded-full shrink-0"
            >
              <Check v-if="state === 'completed'" class="size-5" />
              <span v-else>{{ step.step }}</span>
            </Button>
          </div>

          <div class="mt-5 flex flex-col items-center text-center">
            <StepperTitle class="text-sm font-semibold transition lg:text-base">
              {{ step.title }}
            </StepperTitle>
            <StepperDescription class="sr-only text-xs text-muted-foreground transition md:not-sr-only lg:text-sm">
              {{ step.description }}
            </StepperDescription>
          </div>
        </StepperItem>
      </Stepper>
    </div>

    <div>
      <div v-if="currentStep === 1" class="grid gap-6 lg:grid-cols-2">
        <MenuCard
          v-for="channel in channels"
          :key="channel.value"
          :onClick="() => selectChannel(channel)"
          :title="channel.title"
          :subTitle="channel.subTitle"
          :icon="channel.icon"
        />
      </div>

      <div v-else-if="currentStep === 2" class="space-y-6">
        <Button @click="goBack" variant="link" size="xs">← {{ $t('globals.messages.back') }}</Button>

        <EmailInboxForm
          v-if="selectedChannel?.value === 'email' && emailMode === 'custom'"
          :initial-values="{}"
          :submitForm="submitEmailInboxForm"
          :isLoading="isLoading"
        />

        <ManagedEmailInboxForm
          v-else-if="selectedChannel?.value === 'email' && emailMode === 'managed'"
          :initial-values="{}"
          :submitForm="submitManagedEmailInboxForm"
          :isLoading="isLoading"
        />

        <div v-else-if="selectedChannel?.value === 'email'" class="space-y-6">
          <div class="grid gap-6 lg:grid-cols-2">
            <MenuCard
              :onClick="() => (emailMode = 'managed')"
              title="E-mail gerenciado pelo CanalGov"
              subTitle="Escolha o nome da caixa e um domínio cadastrado. O CanalGov provisiona o endereço."
              :icon="Send"
            />
            <MenuCard
              :onClick="() => (emailMode = 'custom')"
              title="SMTP/IMAP próprio"
              subTitle="Use AWS SES, SNS, provedor próprio ou servidor de e-mail via Docker."
              :icon="Mail"
            />
          </div>
        </div>

        <GenericInboxForm
          v-else-if="selectedChannel"
          :initial-values="{}"
          :submitForm="submitGenericInboxForm"
          :isLoading="isLoading"
          :channel="selectedChannel.value"
          :channel-label="selectedChannel.title"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { useRouter } from 'vue-router'
import { CustomBreadcrumb } from '@/components/ui/breadcrumb/index.js'
import { BellRing, Check, Mail, MessageSquare, Phone, Send, Webhook } from 'lucide-vue-next'
import MenuCard from '@/components/layout/MenuCard.vue'
import {
  Stepper,
  StepperDescription,
  StepperItem,
  StepperSeparator,
  StepperTitle
} from '@/components/ui/stepper'
import EmailInboxForm from '@/features/admin/inbox/EmailInboxForm.vue'
import GenericInboxForm from '@/features/admin/inbox/GenericInboxForm.vue'
import ManagedEmailInboxForm from '@/features/admin/inbox/ManagedEmailInboxForm.vue'
import api from '@/api'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { useEmitter } from '@/composables/useEmitter'
import { handleHTTPError } from '@/utils/http'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const emitter = useEmitter()
const isLoading = ref(false)
const currentStep = ref(1)
const selectedChannel = ref(null)
const emailMode = ref(null)
const router = useRouter()

const breadcrumbLinks = [
  { path: 'inbox-list', label: t('globals.terms.inbox', 2) },
  { path: '', label: t('globals.messages.new', { name: t('globals.terms.inbox') }) }
]

const steps = [
  {
    step: 1,
    title: t('globals.terms.channel'),
    description: 'Escolha o meio principal desta caixa de entrada.'
  },
  {
    step: 2,
    title: t('globals.terms.configure'),
    description: 'Defina os dados operacionais do canal.'
  }
]

const channels = [
  {
    value: 'none',
    title: 'Nenhum',
    subTitle: 'Caixa manual, sem conexão automática.',
    icon: MessageSquare
  },
  {
    value: 'email',
    title: t('globals.terms.email'),
    subTitle: 'Fluxo completo com IMAP e SMTP.',
    icon: Mail
  },
  {
    value: 'whatsapp',
    title: 'WhatsApp',
    subTitle: 'Caixa preparada para atendimento por WhatsApp.',
    icon: Phone
  },
  {
    value: 'telegram',
    title: 'Telegram',
    subTitle: 'Caixa preparada para bot ou atendimento no Telegram.',
    icon: Send
  },
  {
    value: 'sms',
    title: 'SMS',
    subTitle: 'Caixa para alertas e respostas curtas.',
    icon: MessageSquare
  },
  {
    value: 'push',
    title: 'Notificação push',
    subTitle: 'Caixa para disparos e recepção operacional de eventos push.',
    icon: BellRing
  },
  {
    value: 'webhook',
    title: 'Webhook',
    subTitle: 'Caixa para integrações baseadas em eventos HTTP.',
    icon: Webhook
  }
]

const selectChannel = (channel) => {
  selectedChannel.value = channel
  emailMode.value = null
  currentStep.value = 2
}

const goBack = () => {
  currentStep.value = 1
  selectedChannel.value = null
  emailMode.value = null
}

const createInbox = async (payload) => {
  try {
    isLoading.value = true
    await api.createInbox(payload)
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      description: t('globals.messages.createdSuccessfully', {
        name: t('globals.terms.inbox')
      })
    })
    router.push({ name: 'inbox-list' })
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  } finally {
    isLoading.value = false
  }
}

const submitEmailInboxForm = (values) => {
  createInbox({
    name: values.name,
    from: values.from,
    enabled: values.enabled,
    csat_enabled: values.csat_enabled,
    channel: 'email',
    config: {
      enable_plus_addressing: values.enable_plus_addressing,
      auth_type: values.auth_type,
      oauth: values.oauth,
      imap: [values.imap],
      smtp: [values.smtp]
    }
  })
}

const submitManagedEmailInboxForm = (values) => {
  createInbox({
    name: values.name,
    from: values.managed_email_address,
    enabled: values.enabled,
    csat_enabled: values.csat_enabled,
    channel: 'email',
    config: {
      receive_mode: 'managed',
      delivery_provider: values.delivery_provider,
      managed_domain_id: values.managed_domain_id,
      managed_domain: values.managed_domain,
      managed_local_part: values.managed_local_part,
      managed_email_address: values.managed_email_address,
      provider_config: JSON.parse(values.provider_config || '{}'),
      smtp: [],
      imap: []
    }
  })
}

const submitGenericInboxForm = (values) => {
  createInbox({
    name: values.name,
    from: values.from || '',
    enabled: values.enabled,
    csat_enabled: values.csat_enabled,
    channel: selectedChannel.value.value,
    config: JSON.parse(values.extra_config || '{}')
  })
}
</script>
