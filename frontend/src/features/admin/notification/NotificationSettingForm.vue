<template>
  <form @submit="onSmtpSubmit" class="space-y-6">
    <div class="box p-5">
      <FormField name="enabled" v-slot="{ value, handleChange }">
        <FormItem>
          <FormControl>
            <div class="flex items-center space-x-2">
              <Checkbox :checked="value" @update:checked="handleChange" />
              <Label>Habilitar notificações por e-mail</Label>
            </div>
          </FormControl>
          <FormDescription>
            Ative o envio de alertas internos por SMTP para atribuições, SLA e eventos importantes.
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <div class="box p-5 space-y-5">
      <div class="space-y-1">
        <h3 class="font-medium">Conexão SMTP</h3>
        <p class="text-sm text-muted-foreground">
          Dados de autenticação e transporte usados para conexão com o provedor de e-mail.
        </p>
      </div>

      <div class="grid gap-5 md:grid-cols-2">
        <FormField v-slot="{ componentField }" name="host">
          <FormItem>
            <FormLabel>{{ $t('globals.terms.smtpHost') }}</FormLabel>
            <FormControl>
              <Input type="text" placeholder="smtp.gmail.com" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="port">
          <FormItem>
            <FormLabel>{{ $t('globals.terms.smtpPort') }}</FormLabel>
            <FormControl>
              <Input type="number" placeholder="587" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="username">
          <FormItem>
            <FormLabel>{{ $t('globals.terms.username') }}</FormLabel>
            <FormControl>
              <Input type="text" placeholder="admin@seu-orgao.gov.br" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="password">
          <FormItem>
            <FormLabel>{{ $t('globals.terms.password') }}</FormLabel>
            <FormControl>
              <Input type="password" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="auth_protocol">
          <FormItem>
            <FormLabel>{{ $t('admin.inbox.authProtocol') }}</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue :placeholder="t('admin.inbox.authProtocol.description')" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="plain">Plain</SelectItem>
                    <SelectItem value="login">Login</SelectItem>
                    <SelectItem value="cram">CRAM-MD5</SelectItem>
                    <SelectItem value="none">Nenhum</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="tls_type">
          <FormItem>
            <FormLabel>TLS</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue :placeholder="t('globals.messages.selectTLS')" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="none">Off</SelectItem>
                    <SelectItem value="tls">SSL/TLS</SelectItem>
                    <SelectItem value="starttls">STARTTLS</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
    </div>

    <div class="box p-5 space-y-5">
      <div class="space-y-1">
        <h3 class="font-medium">Identidade do remetente</h3>
        <p class="text-sm text-muted-foreground">
          Endereço exibido nas notificações e hostname usado na negociação SMTP.
        </p>
      </div>

      <div class="grid gap-5 md:grid-cols-2">
        <FormField v-slot="{ componentField }" name="email_address">
          <FormItem>
            <FormLabel>{{ $t('globals.terms.fromEmailAddress') }}</FormLabel>
            <FormControl>
              <Input
                type="text"
                :placeholder="t('admin.inbox.fromEmailAddress.placeholder')"
                v-bind="componentField"
              />
            </FormControl>
            <FormDescription>{{ $t('admin.inbox.fromEmailAddress.description') }}</FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="hello_hostname">
          <FormItem>
            <FormLabel>{{ $t('admin.inbox.heloHostname') }}</FormLabel>
            <FormControl>
              <Input type="text" placeholder="mail.seu-orgao.gov.br" v-bind="componentField" />
            </FormControl>
            <FormDescription>{{ $t('admin.inbox.heloHostname.description') }}</FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
    </div>

    <div class="box p-5 space-y-5">
      <div class="space-y-1">
        <h3 class="font-medium">Desempenho e resiliência</h3>
        <p class="text-sm text-muted-foreground">
          Limites de conexão, tempo de espera e política de retentativas do envio.
        </p>
      </div>

      <div class="grid gap-5 md:grid-cols-2">
        <FormField v-slot="{ componentField }" name="max_conns">
          <FormItem>
            <FormLabel>{{ $t('admin.inbox.maxConnections') }}</FormLabel>
            <FormControl>
              <Input type="number" placeholder="2" v-bind="componentField" />
            </FormControl>
            <FormDescription>{{ $t('admin.inbox.maxConnections.description') }}</FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="max_msg_retries">
          <FormItem>
            <FormLabel>{{ $t('admin.inbox.maxRetries') }}</FormLabel>
            <FormControl>
              <Input type="number" placeholder="3" v-bind="componentField" />
            </FormControl>
            <FormDescription>{{ $t('admin.inbox.maxRetries.description') }}</FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="idle_timeout">
          <FormItem>
            <FormLabel>{{ $t('admin.inbox.idleTimeout') }}</FormLabel>
            <FormControl>
              <Input type="text" placeholder="15s" v-bind="componentField" />
            </FormControl>
            <FormDescription>{{ $t('admin.inbox.idleTimeout.description') }}</FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="wait_timeout">
          <FormItem>
            <FormLabel>{{ $t('admin.inbox.waitTimeout') }}</FormLabel>
            <FormControl>
              <Input type="text" placeholder="5s" v-bind="componentField" />
            </FormControl>
            <FormDescription>{{ $t('admin.inbox.waitTimeout.description') }}</FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
    </div>

    <FormField v-slot="{ componentField, handleChange }" name="tls_skip_verify">
      <FormItem class="flex flex-row items-center justify-between box p-4">
        <div class="space-y-0.5">
          <FormLabel class="text-base">{{ $t('admin.inbox.skipTLSVerification') }}</FormLabel>
          <FormDescription>{{ $t('admin.inbox.skipTLSVerification.description') }}</FormDescription>
        </div>
        <FormControl>
          <Switch :checked="componentField.modelValue" @update:checked="handleChange" />
        </FormControl>
      </FormItem>
    </FormField>

    <Button type="submit" :isLoading="isLoading">{{ submitLabel }}</Button>
  </form>
</template>

<script setup>
import { watch, ref, computed } from 'vue'
import { Button } from '@/components/ui/button'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createFormSchema } from './formSchema.js'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription
} from '@/components/ui/form'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { useI18n } from 'vue-i18n'

const isLoading = ref(false)
const { t } = useI18n()
const props = defineProps({
  initialValues: {
    type: Object,
    required: false
  },
  submitForm: {
    type: Function,
    required: true
  },
  submitLabel: {
    type: String,
    required: false,
    default: () => ''
  }
})

const submitLabel = computed(() => {
  if (props.submitLabel) {
    return props.submitLabel
  }
  return t('globals.messages.save')
})

const smtpForm = useForm({
  validationSchema: toTypedSchema(createFormSchema(t))
})

const onSmtpSubmit = smtpForm.handleSubmit(async (values) => {
  isLoading.value = true
  try {
    await props.submitForm(values)
  } finally {
    isLoading.value = false
  }
})

// Watch for changes in initialValues and update the form.
watch(
  () => props.initialValues,
  (newValues) => {
    smtpForm.setValues(newValues)
  },
  { deep: true, immediate: true }
)
</script>
