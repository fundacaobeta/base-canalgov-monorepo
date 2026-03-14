<template>
  <form @submit.prevent="onSubmit" class="space-y-6">
    <div v-if="isResponseTemplate" class="rounded-2xl border border-border/70 bg-muted/20 p-4">
      <div class="space-y-1">
        <h3 class="text-base font-medium">Modelo de resposta</h3>
        <p class="text-sm text-muted-foreground">
          Use este modelo no composer de resposta. Ele pode ser global ou associado a uma equipe específica.
        </p>
      </div>
    </div>

    <FormField v-slot="{ componentField }" name="name">
      <FormItem v-auto-animate>
        <FormLabel>{{ isResponseTemplate ? 'Nome do modelo' : $t('globals.terms.name') }}</FormLabel>
        <FormControl>
          <Input type="text" v-bind="componentField" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField, handleChange }" name="team_id" v-if="isResponseTemplate">
      <FormItem>
        <FormLabel>Equipe</FormLabel>
        <FormControl>
          <Select
            :model-value="componentField.modelValue == null ? 'global' : String(componentField.modelValue)"
            @update:model-value="(value) => handleChange(value === 'global' ? null : Number(value))"
          >
            <SelectTrigger>
              <SelectValue placeholder="Modelo global" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="global">Global</SelectItem>
              <SelectItem v-for="team in teamOptions" :key="team.value" :value="team.value">
                {{ team.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </FormControl>
        <FormDescription>
          Se escolher uma equipe, o modelo aparece primeiro para conversas encaminhadas a ela.
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="subject" v-if="showSubjectField">
      <FormItem>
        <FormLabel>{{ $t('globals.terms.subject') }}</FormLabel>
        <FormControl>
          <Input type="text" placeholder="" v-bind="componentField" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField, handleChange }" name="body">
      <FormItem>
        <FormLabel>{{ isResponseTemplate ? 'Conteúdo da resposta' : $t('globals.terms.body') }}</FormLabel>
        <FormControl>
          <CodeEditor v-model="componentField.modelValue" @update:modelValue="handleChange" />
        </FormControl>
        <FormDescription v-if="isOutgoingTemplate">
          {{
            $t('admin.template.makeSureTemplateHasContent', {
              content: '\u007b\u007b template "content" . \u007d\u007d'
            })
          }}
        </FormDescription>
        <FormDescription v-else-if="isResponseTemplate">
          Escreva o texto-base da resposta. O atendente ainda pode ajustar o conteúdo antes de enviar.
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="is_default" v-slot="{ value, handleChange }">
      <FormItem>
        <FormControl>
          <div class="flex items-center space-x-2">
            <Checkbox :checked="value" @update:checked="handleChange" />
            <Label>{{ defaultLabel }}</Label>
          </div>
        </FormControl>
        <FormDescription>{{ defaultDescription }}</FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <Button type="submit" :isLoading="isLoading"> {{ submitLabel }} </Button>
  </form>
</template>

<script setup>
import { watch, computed, onMounted } from 'vue'
import { Button } from '@/components/ui/button'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createFormSchema } from './formSchema.js'
import { vAutoAnimate } from '@formkit/auto-animate/vue'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import CodeEditor from '@/components/editor/CodeEditor.vue'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import { useI18n } from 'vue-i18n'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { useTeamStore } from '@/stores/team'

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
  },
  isLoading: {
    type: Boolean,
    required: false
  }
})
const { t } = useI18n()
const teamStore = useTeamStore()

const submitLabel = computed(() => {
  return props.submitLabel || t('globals.messages.save')
})

const form = useForm({
  validationSchema: toTypedSchema(createFormSchema(t)),
  initialValues: props.initialValues
})

const onSubmit = form.handleSubmit((values) => {
  props.submitForm({
    ...values,
    team_id: values.team_id ? Number(values.team_id) : null
  })
})

const templateType = computed(() => form.values.type || props.initialValues?.type || 'response')
const isResponseTemplate = computed(() => templateType.value === 'response')
const isOutgoingTemplate = computed(() => templateType.value === 'email_outgoing')
const showSubjectField = computed(() => templateType.value === 'email_notification')
const teamOptions = computed(() => teamStore.options)
const defaultLabel = computed(() =>
  isResponseTemplate.value ? 'Usar como padrão nas respostas' : t('globals.terms.isDefault')
)
const defaultDescription = computed(() => {
  if (isResponseTemplate.value) {
    return 'Se for global, vale para todo o sistema. Se estiver vinculada a uma equipe, vira o padrão daquela equipe.'
  }
  if (isOutgoingTemplate.value) {
    return t('admin.template.onlyOneDefaultOutgoingTemplate')
  }
  return 'Marque apenas quando este modelo deve ser priorizado dentro do seu tipo.'
})

onMounted(() => {
  teamStore.fetchTeams()
})

// Watch for changes in initialValues and update the form.
watch(
  () => props.initialValues,
  (newValues) => {
    form.setValues(newValues)
  },
  { deep: true }
)
</script>
