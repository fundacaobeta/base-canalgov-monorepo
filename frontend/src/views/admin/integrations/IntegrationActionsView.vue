<template>
  <div class="space-y-6">
    <div class="box p-6">
      <h2 class="text-xl font-semibold">Ações e gatilhos de integrações</h2>
      <p class="mt-2 text-sm text-muted-foreground">
        Cadastre múltiplas ações por integração, defina se elas são de entrada ou saída e escolha em quais gatilhos devem ser executadas.
      </p>
    </div>

    <div class="box space-y-5 p-6">
      <div class="grid gap-4 lg:grid-cols-2">
        <div>
          <label class="text-sm font-medium">Integração</label>
          <Select v-model="draft.integration">
            <SelectTrigger class="mt-2">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="integration in integrationCatalog"
                :key="integration.value"
                :value="integration.value"
              >
                {{ integration.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div>
          <label class="text-sm font-medium">Direção</label>
          <Select v-model="draft.direction">
            <SelectTrigger class="mt-2">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="direction in integrationDirections"
                :key="direction.value"
                :value="direction.value"
              >
                {{ direction.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div class="grid gap-4 lg:grid-cols-2">
        <div>
          <label class="text-sm font-medium">Nome da ação</label>
          <Input v-model="draft.name" class="mt-2" placeholder="Ex.: Abrir protocolo externo" />
        </div>
        <div>
          <label class="text-sm font-medium">Método HTTP</label>
          <Select v-model="draft.method">
            <SelectTrigger class="mt-2">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="POST">POST</SelectItem>
              <SelectItem value="PUT">PUT</SelectItem>
              <SelectItem value="PATCH">PATCH</SelectItem>
              <SelectItem value="GET">GET</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div>
        <label class="text-sm font-medium">Descrição</label>
        <Input
          v-model="draft.description"
          class="mt-2"
          placeholder="Explique rapidamente quando esta ação deve ser usada."
        />
      </div>

      <div>
        <label class="text-sm font-medium">URL do request</label>
        <Input
          v-model="draft.url"
          class="mt-2"
          placeholder="https://seu-endpoint.exemplo.com/api/integracoes"
        />
      </div>

      <div class="space-y-3">
        <div>
          <label class="text-sm font-medium">Gatilhos</label>
          <p class="mt-1 text-xs text-muted-foreground">
            Para aparecer no menu da conversa, marque `Manual na conversa` em ações de saída.
          </p>
        </div>

        <div class="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
          <label
            v-for="trigger in integrationTriggers"
            :key="trigger.value"
            class="flex items-center gap-3 rounded-xl border border-border/70 px-3 py-3"
          >
            <Checkbox
              :checked="draft.triggers.includes(trigger.value)"
              @update:checked="(checked) => toggleTrigger(trigger.value, checked)"
            />
            <span class="text-sm">{{ trigger.label }}</span>
          </label>
        </div>
      </div>

      <div>
        <label class="text-sm font-medium">Body template</label>
        <Textarea
          v-model="draft.bodyTemplate"
          class="mt-2 min-h-40 font-mono text-xs"
          placeholder='{"conversation_uuid":"{{conversation_uuid}}"}'
        />
        <p class="mt-2 text-xs text-muted-foreground">
          Variáveis disponíveis: `{{conversation_uuid}}`, `{{reference_number}}`, `{{contact_name}}`, `{{contact_email}}`, `{{contact_phone}}`, `{{assigned_user_name}}`, `{{assigned_team_name}}`.
        </p>
      </div>

      <div class="flex items-center justify-between rounded-xl border border-border/70 px-4 py-3">
        <div>
          <p class="text-sm font-medium">Ação habilitada</p>
          <p class="text-xs text-muted-foreground">Somente ações habilitadas podem ser executadas.</p>
        </div>
        <Switch v-model:checked="draft.enabled" />
      </div>

      <div class="flex gap-3">
        <Button @click="saveAction">{{ isEditing ? 'Salvar ação' : 'Adicionar ação' }}</Button>
        <Button variant="outline" @click="resetDraft">Limpar</Button>
      </div>
    </div>

    <div class="space-y-5">
      <div v-for="group in groupedActions" :key="group.value" class="box p-6">
        <div class="mb-4 flex items-center justify-between gap-3">
          <div>
            <h3 class="text-lg font-semibold">{{ group.label }}</h3>
            <p class="text-sm text-muted-foreground">
              {{ group.actions.length }} {{ group.actions.length === 1 ? 'ação configurada' : 'ações configuradas' }}
            </p>
          </div>
        </div>

        <div class="space-y-4">
          <div
            v-for="action in group.actions"
            :key="action.id"
            class="rounded-2xl border border-border/70 p-4"
          >
            <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
              <div class="space-y-2">
                <div class="flex flex-wrap items-center gap-2">
                  <span class="font-medium">{{ action.name }}</span>
                  <span class="rounded-full bg-muted px-2 py-0.5 text-xs">{{ action.method }}</span>
                  <span class="rounded-full bg-muted px-2 py-0.5 text-xs">
                    {{ labelForDirection(action.direction) }}
                  </span>
                  <span
                    class="rounded-full px-2 py-0.5 text-xs"
                    :class="action.enabled ? 'bg-emerald-100 text-emerald-700' : 'bg-zinc-100 text-zinc-600'"
                  >
                    {{ action.enabled ? 'Ativa' : 'Inativa' }}
                  </span>
                </div>
                <p v-if="action.description" class="text-sm text-muted-foreground">
                  {{ action.description }}
                </p>
                <p class="break-all text-sm text-muted-foreground">
                  {{ action.url || 'URL não configurada' }}
                </p>
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="trigger in action.triggers"
                    :key="`${action.id}-${trigger}`"
                    class="rounded-full border border-border px-2 py-0.5 text-xs text-muted-foreground"
                  >
                    {{ labelForTrigger(trigger) }}
                  </span>
                </div>
              </div>

              <div class="flex gap-2">
                <Button variant="outline" @click="editAction(action)">Editar</Button>
                <Button variant="destructive" @click="removeAction(action.id)">Remover</Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import {
  useIntegrationActions,
  integrationCatalog,
  integrationDirections,
  integrationTriggers
} from '@/composables/useIntegrationActions'

const { actions, groupedActions } = useIntegrationActions()

const blankDraft = () => ({
  id: '',
  integration: 'whatsapp',
  direction: 'outgoing',
  triggers: ['manual_conversation'],
  name: '',
  description: '',
  method: 'POST',
  url: '',
  enabled: true,
  bodyTemplate: '{\n  "conversation_uuid": "{{conversation_uuid}}"\n}'
})

const draft = ref(blankDraft())
const isEditing = computed(() => Boolean(draft.value.id))

function resetDraft() {
  draft.value = blankDraft()
}

function toggleTrigger(triggerValue, checked) {
  if (checked) {
    draft.value.triggers = Array.from(new Set([...draft.value.triggers, triggerValue]))
    return
  }

  draft.value.triggers = draft.value.triggers.filter((item) => item !== triggerValue)
}

function saveAction() {
  const payload = {
    ...draft.value,
    id: draft.value.id || `action-${Date.now()}-${Math.round(Math.random() * 10000)}`
  }

  const existingIndex = actions.value.findIndex((item) => item.id === payload.id)
  if (existingIndex >= 0) {
    actions.value[existingIndex] = payload
  } else {
    actions.value.unshift(payload)
  }

  resetDraft()
}

function editAction(action) {
  draft.value = {
    ...action,
    triggers: [...action.triggers]
  }
}

function removeAction(id) {
  actions.value = actions.value.filter((item) => item.id !== id)
  if (draft.value.id === id) resetDraft()
}

function labelForTrigger(triggerValue) {
  return integrationTriggers.find((item) => item.value === triggerValue)?.label || triggerValue
}

function labelForDirection(directionValue) {
  return integrationDirections.find((item) => item.value === directionValue)?.label || directionValue
}
</script>
