<template>
  <form class="space-y-6" @submit.prevent="onSubmit">
    <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
      <h3 class="font-medium">E-mail gerenciado pelo CanalGov</h3>
      <p class="mt-1 text-sm text-muted-foreground">
        O CanalGov provisiona o endereço da caixa. Mensagens recebidas nesse endereço entram no fluxo de atendimento.
      </p>
    </div>

    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>Nome da caixa</FormLabel>
        <FormControl>
          <Input v-bind="componentField" placeholder="Ouvidoria" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(280px,1fr)]">
      <FormField v-slot="{ componentField }" name="managed_local_part">
        <FormItem>
          <FormLabel>Nome do e-mail</FormLabel>
          <FormControl>
            <Input v-bind="componentField" placeholder="ouvidoria" />
          </FormControl>
          <FormDescription>Parte antes do `@`.</FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField, handleChange }" name="managed_domain_id">
        <FormItem>
          <FormLabel>Domínio</FormLabel>
          <FormControl>
            <Select v-bind="componentField" @update:model-value="handleChange">
              <SelectTrigger>
                <SelectValue placeholder="Selecione um domínio" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="domain in domainOptions" :key="domain.id" :value="domain.id">
                  {{ domain.domain }}
                </SelectItem>
              </SelectContent>
            </Select>
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <div class="rounded-xl border border-border/70 p-4">
      <div class="text-sm font-medium">Endereço gerado</div>
      <div class="mt-2 font-mono text-sm">{{ managedEmailAddress || 'Preencha o nome e selecione um domínio' }}</div>
    </div>

    <FormField v-slot="{ componentField, handleChange }" name="delivery_provider">
      <FormItem>
        <FormLabel>Provedor de entrega</FormLabel>
        <FormControl>
          <Select v-bind="componentField" @update:model-value="handleChange">
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="canalgov_managed">CanalGov gerenciado</SelectItem>
              <SelectItem value="ses">AWS SES</SelectItem>
              <SelectItem value="sns">AWS SNS</SelectItem>
              <SelectItem value="self_hosted">Servidor próprio</SelectItem>
              <SelectItem value="docker_mailserver">Servidor via Docker</SelectItem>
              <SelectItem value="custom">Customizado</SelectItem>
            </SelectContent>
          </Select>
        </FormControl>
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="provider_config">
      <FormItem>
        <FormLabel>Configuração adicional do provedor</FormLabel>
        <FormControl>
          <Textarea
            v-bind="componentField"
            class="min-h-32 font-mono text-xs"
            placeholder='{"region":"us-east-1","topic_arn":"","bucket":""}'
          />
        </FormControl>
        <FormDescription>JSON opcional para integrações como AWS SES, SNS ou infraestrutura própria.</FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <div class="grid gap-4 lg:grid-cols-2">
      <FormField v-slot="{ componentField, handleChange }" name="enabled">
        <FormItem class="box flex flex-row items-center justify-between p-4">
          <div>
            <FormLabel class="text-base">Habilitada</FormLabel>
          </div>
          <FormControl>
            <Switch :checked="componentField.modelValue" @update:checked="handleChange" />
          </FormControl>
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField, handleChange }" name="csat_enabled">
        <FormItem class="box flex flex-row items-center justify-between p-4">
          <div>
            <FormLabel class="text-base">CSAT</FormLabel>
          </div>
          <FormControl>
            <Switch :checked="componentField.modelValue" @update:checked="handleChange" />
          </FormControl>
        </FormItem>
      </FormField>
    </div>

    <Button type="submit" :is-loading="isLoading" :disabled="isLoading">Salvar</Button>
  </form>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'
import api from '@/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form'
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

const props = defineProps({
  initialValues: {
    type: Object,
    default: () => ({})
  },
  submitForm: {
    type: Function,
    required: true
  },
  isLoading: {
    type: Boolean,
    default: false
  }
})

const emitter = useEmitter()
const domains = ref([])

const form = useForm({
  validationSchema: toTypedSchema(z.object({
    name: z.string().min(1, 'Obrigatório'),
    managed_local_part: z.string().min(1, 'Obrigatório'),
    managed_domain_id: z.string().min(1, 'Obrigatório'),
    delivery_provider: z.string().min(1, 'Obrigatório'),
    provider_config: z.string().default('{}').refine((value) => {
      try { JSON.parse(value || '{}'); return true } catch { return false }
    }, 'Informe um JSON válido'),
    enabled: z.boolean().default(true),
    csat_enabled: z.boolean().default(false)
  })),
  initialValues: {
    name: '',
    managed_local_part: '',
    managed_domain_id: '',
    delivery_provider: 'canalgov_managed',
    provider_config: '{}',
    enabled: true,
    csat_enabled: false
  }
})

const fetchDomains = async () => {
  try {
    const resp = await api.getMailDomainsSettings()
    domains.value = (resp.data.data?.domains || []).filter((item) => item.enabled)
    if (!form.values.managed_domain_id) {
      const defaultDomain = domains.value.find((item) => item.is_default) || domains.value[0]
      if (defaultDomain) form.setFieldValue('managed_domain_id', defaultDomain.id)
    }
  } catch (error) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: handleHTTPError(error).message
    })
  }
}

onMounted(fetchDomains)

const domainOptions = computed(() => domains.value)
const selectedDomain = computed(() =>
  domains.value.find((item) => item.id === form.values.managed_domain_id)
)
const managedEmailAddress = computed(() => {
  if (!form.values.managed_local_part || !selectedDomain.value?.domain) return ''
  return `${form.values.managed_local_part}@${selectedDomain.value.domain}`
})

const onSubmit = form.handleSubmit((values) => {
  props.submitForm({
    ...values,
    managed_domain: selectedDomain.value?.domain || '',
    managed_email_address: managedEmailAddress.value
  })
})

watch(
  () => props.initialValues,
  (newValues) => {
    if (!newValues || !Object.keys(newValues).length) return
    form.setValues({
      name: newValues.name || '',
      managed_local_part: newValues.managed_local_part || '',
      managed_domain_id: newValues.managed_domain_id || '',
      delivery_provider: newValues.delivery_provider || 'canalgov_managed',
      provider_config: newValues.provider_config || '{}',
      enabled: typeof newValues.enabled === 'boolean' ? newValues.enabled : true,
      csat_enabled: typeof newValues.csat_enabled === 'boolean' ? newValues.csat_enabled : false
    })
  },
  { immediate: true, deep: true }
)
</script>
