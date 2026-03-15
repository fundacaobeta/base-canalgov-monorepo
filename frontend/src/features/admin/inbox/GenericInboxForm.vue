<template>
  <form class="space-y-6" @submit.prevent="onSubmit">
    <div class="rounded-2xl border border-border/70 bg-muted/20 p-4">
      <h3 class="font-medium">{{ channelLabel }}</h3>
      <p class="mt-1 text-sm text-muted-foreground">
        {{ t('admin.inbox.form.genericChannelDescription') }}
      </p>
    </div>

    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>{{ t('globals.terms.name') }}</FormLabel>
        <FormControl>
          <Input v-bind="componentField" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="from">
      <FormItem>
        <FormLabel>{{ t('admin.inbox.form.sourceIdentifier') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            :placeholder="channel === 'none' ? t('globals.terms.optional') : t('admin.inbox.form.sourceIdentifierPlaceholder')"
          />
        </FormControl>
        <FormDescription>
          {{ t('admin.inbox.form.sourceIdentifierDescription') }}
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField, handleChange }" name="enabled">
      <FormItem class="box flex flex-row items-center justify-between p-4">
        <div class="space-y-0.5">
          <FormLabel class="text-base">{{ t('globals.terms.enabled') }}</FormLabel>
          <FormDescription>{{ t('admin.inbox.form.enabledDescription') }}</FormDescription>
        </div>
        <FormControl>
          <Switch :checked="componentField.modelValue" @update:checked="handleChange" />
        </FormControl>
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField, handleChange }" name="csat_enabled">
      <FormItem class="box flex flex-row items-center justify-between p-4">
        <div class="space-y-0.5">
          <FormLabel class="text-base">{{ t('admin.inbox.csatSurveys') }}</FormLabel>
          <FormDescription>{{ t('admin.inbox.form.csatDescription') }}</FormDescription>
        </div>
        <FormControl>
          <Switch :checked="componentField.modelValue" @update:checked="handleChange" />
        </FormControl>
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="extra_config">
      <FormItem>
        <FormLabel>{{ t('admin.inbox.form.extraConfig') }}</FormLabel>
        <FormControl>
          <Textarea
            v-bind="componentField"
            class="min-h-40 font-mono text-xs"
            :placeholder="t('admin.inbox.form.extraConfigPlaceholder')"
          />
        </FormControl>
        <FormDescription>
          {{ t('admin.inbox.form.extraConfigDescription') }}
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
import { useI18n } from 'vue-i18n'
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

const { t } = useI18n()

const schema = toTypedSchema(
  z.object({
    name: z.string().min(1, t('admin.inbox.form.validation.required')),
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
      }, t('admin.inbox.form.validation.invalidJson'))
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

const submitLabel = computed(() => props.submitLabel || t('globals.messages.save'))

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
