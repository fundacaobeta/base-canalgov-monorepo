<template>
  <form @submit="onSubmit" class="space-y-6">
    <div class="box p-5">
      <FormField name="enabled" v-slot="{ value, handleChange }">
        <FormItem>
          <FormControl>
            <div class="flex items-center space-x-2">
              <Checkbox :checked="value" @update:checked="handleChange" />
              <Label>{{ t('admin.notification.official.enable') }}</Label>
            </div>
          </FormControl>
          <FormDescription>
            {{ t('admin.notification.official.enableDescription') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <div class="box p-5 space-y-5">
      <div class="space-y-1">
        <h3 class="font-medium">{{ t('admin.notification.official.defaultConfig') }}</h3>
        <p class="text-sm text-muted-foreground">
          {{ t('admin.notification.official.defaultConfigDescription') }}
        </p>
      </div>

      <div class="grid gap-5 md:grid-cols-2">
        <FormField v-slot="{ componentField }" name="inbox_id">
          <FormItem>
            <FormLabel>{{ t('admin.notification.official.recipientTeams') }}</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue :placeholder="t('globals.messages.select', { name: t('globals.terms.inbox') })" />
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
              {{ t('admin.notification.official.inboxDescription') }}
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="priority_id">
          <FormItem>
            <FormLabel>{{ t('globals.terms.priority') }}</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue :placeholder="t('globals.messages.select', { name: t('globals.terms.priority') })" />
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
            <FormLabel>{{ t('globals.terms.status') }}</FormLabel>
            <FormControl>
              <Select v-bind="componentField" v-model="componentField.modelValue">
                <SelectTrigger>
                  <SelectValue :placeholder="t('globals.messages.select', { name: t('globals.terms.status') })" />
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
            <FormLabel>{{ t('admin.notification.official.targetSlaHours') }}</FormLabel>
            <FormControl>
              <Input type="number" min="1" placeholder="24" v-bind="componentField" />
            </FormControl>
            <FormDescription>
              {{ t('admin.notification.official.targetSlaDescription') }}
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="subject_prefix">
          <FormItem class="md:col-span-2">
            <FormLabel>{{ t('admin.notification.official.subjectPrefix') }}</FormLabel>
            <FormControl>
              <Input
                type="text"
                :placeholder="t('admin.notification.official.subjectPrefix')"
                v-bind="componentField"
              />
            </FormControl>
            <FormDescription>
              {{ t('admin.notification.official.subjectPrefixDescription') }}
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField, handleChange }" name="default_types">
          <FormItem class="md:col-span-2">
            <FormLabel>{{ t('admin.notification.official.monitoredTypes') }}</FormLabel>
            <FormControl>
              <TagsInput :modelValue="componentField.modelValue" @update:modelValue="handleChange">
                <TagsInputItem v-for="item in componentField.modelValue" :key="item" :value="item">
                  <TagsInputItemText />
                  <TagsInputItemDelete />
                </TagsInputItem>
                <TagsInputInput :placeholder="t('admin.notification.official.ruleTypesPlaceholder')" />
              </TagsInput>
            </FormControl>
            <FormDescription>
              {{ t('admin.notification.official.monitoredTypesDescription') }}
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="internal_note">
          <FormItem class="md:col-span-2">
            <FormLabel>{{ t('admin.notification.official.internalInstructions') }}</FormLabel>
            <FormControl>
              <Textarea
                class="min-h-32"
                :placeholder="t('admin.notification.official.internalInstructionsDescription')"
                v-bind="componentField"
              />
            </FormControl>
            <FormDescription>
              {{ t('admin.notification.official.internalInstructionsDescription') }}
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
            <FormLabel class="text-base">{{ t('admin.notification.official.autoOpen') }}</FormLabel>
            <FormDescription>
              {{ t('admin.notification.official.autoOpenDescription') }}
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
        <h3 class="font-medium">{{ t('admin.notification.official.routingRules') }}</h3>
        <p class="text-sm text-muted-foreground">
          {{ t('admin.notification.official.routingRulesDescription') }}
        </p>
      </div>

      <div class="grid gap-5 md:grid-cols-2">
        <div>
          <Label>{{ t('admin.notification.official.ruleName') }}</Label>
          <Input
            v-model="routingDraft.name"
            class="mt-2"
            :placeholder="t('admin.notification.official.ruleNamePlaceholder')"
          />
        </div>

        <div>
          <Label>{{ t('admin.notification.official.ruleTypes') }}</Label>
          <TagsInput
            :modelValue="routingDraft.types"
            class="mt-2"
            @update:modelValue="(value) => (routingDraft.types = value)"
          >
            <TagsInputItem v-for="item in routingDraft.types" :key="item" :value="item">
              <TagsInputItemText />
              <TagsInputItemDelete />
            </TagsInputItem>
            <TagsInputInput :placeholder="t('admin.notification.official.ruleTypesPlaceholder')" />
          </TagsInput>
          <p class="mt-2 text-xs text-muted-foreground">
            {{ t('admin.notification.official.ruleTypesHint') }}
          </p>
        </div>
      </div>

      <div>
        <Label>{{ t('admin.notification.official.recipientTeams') }}</Label>
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
                {{ t('admin.notification.official.recipientTeamsDescription') }}
              </p>
            </div>
          </label>
        </div>
      </div>

      <div class="flex gap-3">
        <Button type="button" @click="saveRoutingRule">
          {{ isEditingRule ? t('admin.notification.official.saveRule') : t('admin.notification.official.addRule') }}
        </Button>
        <Button type="button" variant="outline" @click="resetRoutingDraft">{{ t('admin.notification.official.clear') }}</Button>
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
                {{ t('globals.terms.team', 2) }}: {{ getRuleTeamsLabel(rule.team_ids) }}
              </p>
            </div>
            <div class="flex gap-2">
              <Button type="button" variant="outline" @click="editRoutingRule(rule)">{{ t('admin.notification.official.edit') }}</Button>
              <Button type="button" variant="destructive" @click="removeRoutingRule(rule.id)">
                {{ t('admin.notification.official.remove') }}
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="rounded border border-dashed p-4 text-sm text-muted-foreground">
        {{ t('admin.notification.official.noRules') }}
      </div>
    </div>

    <Button type="submit" :isLoading="isLoading">{{ t('admin.notification.form.saveConfig') }}</Button>
  </form>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { useForm } from 'vee-validate'
import { useI18n } from 'vue-i18n'
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

const { t } = useI18n()
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
      description: t('admin.notification.official.validation.ruleNameRequired')
    })
    return
  }

  if (!routingDraft.types.length) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: t('admin.notification.official.validation.ruleTypesRequired')
    })
    return
  }

  if (!routingDraft.team_ids.length) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, {
      variant: 'destructive',
      description: t('admin.notification.official.validation.ruleTeamsRequired')
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
