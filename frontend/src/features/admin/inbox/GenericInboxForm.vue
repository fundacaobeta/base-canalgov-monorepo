<template>
  <form class="space-y-6" @submit.prevent="onSubmit">
    <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
      <h3 class="font-medium">{{ channelLabel }}</h3>
      <p class="mt-1 text-sm text-muted-foreground">
        Canal operacional sem configuração técnica obrigatória. Use o JSON adicional para registrar credenciais, IDs externos ou metadados.
      </p>
    </div>

    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>Nome</FormLabel>
        <FormControl>
          <Input v-bind="componentField" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="from">
      <FormItem>
        <FormLabel>Identificador de origem</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            :placeholder="channel === 'none' ? 'Opcional' : 'Ex.: número, bot, sender id ou endereço de referência'"
          />
        </FormControl>
        <FormDescription>
          Campo opcional para identificar a origem visível deste canal.
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField, handleChange }" name="enabled">
      <FormItem class="box flex flex-row items-center justify-between p-4">
        <div class="space-y-0.5">
          <FormLabel class="text-base">Habilitada</FormLabel>
          <FormDescription>A caixa de entrada fica disponível para uso operacional.</FormDescription>
        </div>
        <FormControl>
          <Switch :checked="componentField.modelValue" @update:checked="handleChange" />
        </FormControl>
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField, handleChange }" name="csat_enabled">
      <FormItem class="box flex flex-row items-center justify-between p-4">
        <div class="space-y-0.5">
          <FormLabel class="text-base">CSAT</FormLabel>
          <FormDescription>Ativa pesquisa de satisfação quando este fluxo fizer sentido.</FormDescription>
        </div>
        <FormControl>
          <Switch :checked="componentField.modelValue" @update:checked="handleChange" />
        </FormControl>
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="extra_config">
      <FormItem>
        <FormLabel>Configurações adicionais</FormLabel>
        <FormControl>
          <Textarea
            v-bind="componentField"
            class="min-h-40 font-mono text-xs"
            placeholder='{"provider":"meta","phone_number_id":"","token":""}'
          />
        </FormControl>
        <FormDescription>
          Informe um JSON válido. Isso permite adicionar mais dados de configuração sem depender de um formulário fixo.
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <Button type="submit" :is-loading="isLoading" :disabled="isLoading">
      {{ submitLabel }}
    </Button>
  </form>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'
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

const props = defineProps({
  initialValues: {
    type: Object,
    default: () => ({})
  },
  submitForm: {
    type: Function,
    required: true
  },
  submitLabel: {
    type: String,
    default: ''
  },
  isLoading: {
    type: Boolean,
    default: false
  },
  channel: {
    type: String,
    required: true
  },
  channelLabel: {
    type: String,
    required: true
  }
})

const schema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Obrigatório'),
    from: z.string().optional().default(''),
    enabled: z.boolean().optional().default(true),
    csat_enabled: z.boolean().optional().default(false),
    extra_config: z
      .string()
      .optional()
      .default('{}')
      .refine((value) => {
        try {
          JSON.parse(value || '{}')
          return true
        } catch {
          return false
        }
      }, 'Informe um JSON válido')
  })
)

const form = useForm({
  validationSchema: schema,
  initialValues: {
    name: '',
    from: '',
    enabled: true,
    csat_enabled: false,
    extra_config: '{}'
  }
})

const submitLabel = computed(() => props.submitLabel || 'Salvar')

const onSubmit = form.handleSubmit((values) => {
  props.submitForm({
    ...values,
    extra_config: values.extra_config || '{}'
  })
})

watch(
  () => props.initialValues,
  (newValues) => {
    if (!newValues || !Object.keys(newValues).length) return

    form.setValues({
      name: newValues.name || '',
      from: newValues.from || '',
      enabled: typeof newValues.enabled === 'boolean' ? newValues.enabled : true,
      csat_enabled: typeof newValues.csat_enabled === 'boolean' ? newValues.csat_enabled : false,
      extra_config: newValues.extra_config || '{}'
    })
  },
  { immediate: true, deep: true }
)
</script>
