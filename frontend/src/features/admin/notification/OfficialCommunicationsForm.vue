<template>
  <form @submit="onSubmit" class="space-y-6">
    <div class="box p-5">
      <FormField name="enabled" v-slot="{ value, handleChange }">
        <FormItem>
          <FormControl>
            <div class="flex items-center space-x-2">
              <Checkbox :checked="value" @update:checked="handleChange" />
              <Label>Habilitar comunicações oficiais</Label>
            </div>
          </FormControl>
          <FormDescription>
            Ative este fluxo para transformar ofícios, cartas, notificações e intimações em chamados internos.
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <div class="box p-5 space-y-5">
      <div class="space-y-1">
        <h3 class="font-medium">Configuração padrão do chamado</h3>
        <p class="text-sm text-muted-foreground">
          Defina os parâmetros usados quando uma comunicação oficial gerar atendimento interno.
        </p>
      </div>

      <div class="grid gap-5 md:grid-cols-2">
        <FormField v-slot="{ componentField }" name="inbox_id">
          <FormItem>
            <FormLabel>Caixa de entrada</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue placeholder="Selecione a caixa de entrada" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem
                    v-for="option in inboxOptions"
                    :key="option.value"
                    :value="option.value"
                  >
                    {{ option.label }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </FormControl>
            <FormDescription>
              Caixa usada para concentrar os chamados gerados por esse fluxo.
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="priority_id">
          <FormItem>
            <FormLabel>Prioridade padrão</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue placeholder="Selecione a prioridade" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem
                    v-for="option in priorityOptions"
                    :key="option.value"
                    :value="option.value"
                  >
                    {{ option.label }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="status_id">
          <FormItem>
            <FormLabel>Status inicial</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue placeholder="Selecione o status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem
                    v-for="option in statusOptions"
                    :key="option.value"
                    :value="option.value"
                  >
                    {{ option.label }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="target_sla_hours">
          <FormItem>
            <FormLabel>SLA alvo em horas</FormLabel>
            <FormControl>
              <Input type="number" min="1" placeholder="24" v-bind="componentField" />
            </FormControl>
            <FormDescription>
              Referência operacional para primeira análise do caso.
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="subject_prefix">
          <FormItem class="md:col-span-2">
            <FormLabel>Prefixo do assunto</FormLabel>
            <FormControl>
              <Input
                type="text"
                placeholder="[Oficial] Comunicação recebida"
                v-bind="componentField"
              />
            </FormControl>
            <FormDescription>
              Esse texto ajuda a identificar rapidamente chamados oriundos de comunicações oficiais.
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField, handleChange }" name="default_types">
          <FormItem class="md:col-span-2">
            <FormLabel>Tipos monitorados</FormLabel>
            <FormControl>
              <TagsInput :modelValue="componentField.modelValue" @update:modelValue="handleChange">
                <TagsInputItem v-for="item in componentField.modelValue" :key="item" :value="item">
                  <TagsInputItemText />
                  <TagsInputItemDelete />
                </TagsInputItem>
                <TagsInputInput placeholder="Ofício" />
              </TagsInput>
            </FormControl>
            <FormDescription>
              Cadastre termos usados para classificar esse tipo de comunicação.
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="internal_note">
          <FormItem class="md:col-span-2">
            <FormLabel>Instruções internas</FormLabel>
            <FormControl>
              <Textarea
                class="min-h-32"
                placeholder="Descreva como a equipe deve tratar esse tipo de expediente."
                v-bind="componentField"
              />
            </FormControl>
            <FormDescription>
              Essas instruções servem como referência para a triagem e resposta.
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
    </div>

    <div class="box p-5">
      <FormField name="auto_create_conversation" v-slot="{ componentField, handleChange }">
        <FormItem class="flex flex-row items-center justify-between">
          <div class="space-y-0.5">
            <FormLabel class="text-base">Abrir chamado automaticamente</FormLabel>
            <FormDescription>
              Ao receber uma comunicação oficial, o sistema poderá abrir um chamado usando as regras abaixo.
            </FormDescription>
          </div>
          <FormControl>
            <Switch :checked="componentField.modelValue" @update:checked="handleChange" />
          </FormControl>
        </FormItem>
      </FormField>
    </div>

    <div class="box p-5 space-y-5">
      <div class="space-y-1">
        <h3 class="font-medium">Regras de encaminhamento</h3>
        <p class="text-sm text-muted-foreground">
          Mapeie quais tipos de comunicação devem gerar chamado para uma ou mais equipes.
        </p>
      </div>

      <div class="grid gap-5 md:grid-cols-2">
        <div>
          <Label>Nome da regra</Label>
          <Input
            v-model="routingDraft.name"
            class="mt-2"
            placeholder="Ex.: Jurídico e Ouvidoria"
          />
        </div>

        <div>
          <Label>Tipos da regra</Label>
          <TagsInput
            :modelValue="routingDraft.types"
            class="mt-2"
            @update:modelValue="(value) => (routingDraft.types = value)"
          >
            <TagsInputItem v-for="item in routingDraft.types" :key="item" :value="item">
              <TagsInputItemText />
              <TagsInputItemDelete />
            </TagsInputItem>
            <TagsInputInput placeholder="Intimação" />
          </TagsInput>
          <p class="mt-2 text-xs text-muted-foreground">
            Exemplo: Ofício, Carta, Notificação extrajudicial, Intimação.
          </p>
        </div>
      </div>

      <div>
        <Label>Equipes destinatárias</Label>
        <div class="mt-3 grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
          <label
            v-for="option in teamOptions"
            :key="option.value"
            class="flex items-start space-x-3 rounded border border-border bg-card p-3"
          >
            <Checkbox
              :checked="routingDraft.team_ids.includes(option.value)"
              @update:checked="(checked) => toggleDraftTeam(option.value, checked)"
            />
            <div class="space-y-1">
              <div class="text-sm font-medium">{{ option.label }}</div>
              <p class="text-xs text-muted-foreground">
                Recebe os chamados gerados por esta regra.
              </p>
            </div>
          </label>
        </div>
      </div>

      <div class="flex gap-3">
        <Button type="button" @click="saveRoutingRule">
          {{ isEditingRule ? 'Salvar regra' : 'Adicionar regra' }}
        </Button>
        <Button type="button" variant="outline" @click="resetRoutingDraft">Limpar</Button>
      </div>

      <div class="space-y-4" v-if="routingRules.length">
        <div
          v-for="rule in routingRules"
          :key="rule.id"
          class="rounded border border-border bg-card p-4 space-y-4"
        >
          <div class="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
            <div class="space-y-2">
              <div class="font-medium">{{ rule.name }}</div>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="type in rule.types"
                  :key="`${rule.id}-${type}`"
                  class="rounded bg-muted px-2 py-1 text-xs"
                >
                  {{ type }}
                </span>
              </div>
              <p class="text-sm text-muted-foreground">
                Equipes: {{ getRuleTeamsLabel(rule.team_ids) }}
              </p>
            </div>
            <div class="flex gap-2">
              <Button type="button" variant="outline" @click="editRoutingRule(rule)">Editar</Button>
              <Button type="button" variant="destructive" @click="removeRoutingRule(rule.id)">
                Remover
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="rounded border border-dashed p-4 text-sm text-muted-foreground">
        Nenhuma regra cadastrada ainda. Adicione pelo menos uma regra para mapear as equipes que receberão os chamados.
      </div>
    </div>

    <Button type="submit" :isLoading="isLoading">Salvar configuração</Button>
  </form>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { useForm } from 'vee-validate'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import {
  TagsInput,
  TagsInputInput,
  TagsInputItem,
  TagsInputItemDelete,
  TagsInputItemText
} from '@/components/ui/tags-input'
import { Textarea } from '@/components/ui/textarea'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents'

const props = defineProps({
  initialValues: {
    type: Object,
    required: false,
    default: () => ({})
  },
  inboxOptions: {
    type: Array,
    required: true
  },
  priorityOptions: {
    type: Array,
    required: true
  },
  statusOptions: {
    type: Array,
    required: true
  },
  teamOptions: {
    type: Array,
    required: true
  },
  submitForm: {
    type: Function,
    required: true
  }
})

const emitter = useEmitter()
const isLoading = ref(false)

const defaultValues = {
  enabled: false,
  auto_create_conversation: true,
  inbox_id: '',
  priority_id: '',
  status_id: '',
  subject_prefix: '[Oficial] Comunicação recebida',
  target_sla_hours: '24',
  default_types: ['Ofício', 'Carta', 'Notificação', 'Intimação'],
  internal_note: ''
}

const form = useForm({
  initialValues: defaultValues
})

const routingRules = ref([])
const blankRoutingDraft = () => ({
  id: '',
  name: '',
  types: [],
  team_ids: []
})
const routingDraft = reactive(blankRoutingDraft())
const isEditingRule = computed(() => Boolean(routingDraft.id))

const onSubmit = form.handleSubmit(async (values) => {
  isLoading.value = true
  try {
    await props.submitForm({
      ...values,
      routing_rules: routingRules.value
    })
  } finally {
    isLoading.value = false
  }
})

watch(
  () => props.initialValues,
  (newValues) => {
    const mergedValues = {
      ...defaultValues,
      ...(newValues || {})
    }
    routingRules.value = Array.isArray(mergedValues.routing_rules) ? mergedValues.routing_rules : []
    form.setValues(mergedValues)
  },
  { deep: true, immediate: true }
)

function resetRoutingDraft() {
  Object.assign(routingDraft, blankRoutingDraft())
}

function toggleDraftTeam(teamId, checked) {
  if (checked) {
    routingDraft.team_ids = Array.from(new Set([...routingDraft.team_ids, teamId]))
    return
  }

  routingDraft.team_ids = routingDraft.team_ids.filter((value) => value !== teamId)
}

function saveRoutingRule() {
  if (!routingDraft.name.trim()) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: 'Informe um nome para a regra de encaminhamento.'
    })
    return
  }

  if (!routingDraft.types.length) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: 'Informe ao menos um tipo de comunicação para a regra.'
    })
    return
  }

  if (!routingDraft.team_ids.length) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: 'Selecione ao menos uma equipe destinatária.'
    })
    return
  }

  const payload = {
    id: routingDraft.id || `routing-${Date.now()}-${Math.round(Math.random() * 10000)}`,
    name: routingDraft.name.trim(),
    types: [...routingDraft.types],
    team_ids: [...routingDraft.team_ids]
  }

  const existingIndex = routingRules.value.findIndex((rule) => rule.id === payload.id)
  if (existingIndex >= 0) {
    routingRules.value[existingIndex] = payload
  } else {
    routingRules.value.unshift(payload)
  }

  resetRoutingDraft()
}

function editRoutingRule(rule) {
  Object.assign(routingDraft, {
    id: rule.id,
    name: rule.name,
    types: [...rule.types],
    team_ids: [...rule.team_ids]
  })
}

function removeRoutingRule(ruleId) {
  routingRules.value = routingRules.value.filter((rule) => rule.id !== ruleId)
  if (routingDraft.id === ruleId) {
    resetRoutingDraft()
  }
}

function getRuleTeamsLabel(teamIds) {
  return teamIds
    .map((teamId) => props.teamOptions.find((option) => option.value === teamId)?.label || teamId)
    .join(', ')
}
</script>
