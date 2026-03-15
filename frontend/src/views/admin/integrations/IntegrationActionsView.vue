<template>
  <div>
    <AdminPageHeader
      :title="$t('admin.integration.actions.title')"
      :description="$t('admin.integration.actions.description')"
      :breadcrumbs="[{ label: $t('globals.terms.administration'), to: '/admin' }, { label: $t('globals.terms.integration', 2) }, { label: $t('globals.terms.action', 2) }]"
    >
      <template #actions>
        <Button @click="openCreate">
          <Plus class="h-4 w-4 mr-1.5" aria-hidden="true" />
          {{ $t('globals.messages.new', { name: $t('globals.terms.action') }) }}
        </Button>
      </template>
    </AdminPageHeader>

    <AdminPageWithHelp>
      <template #content>
        <div class="space-y-5">
          <div
            v-if="groupedActions.length === 0"
            class="box flex flex-col items-center justify-center gap-2 p-10 text-center text-muted-foreground"
          >
            <Zap class="h-8 w-8 opacity-40" />
            <p class="text-sm">{{ $t('admin.integration.actions.empty') }}</p>
          </div>

          <div v-for="group in groupedActions" :key="group.value" class="box p-6">
            <div class="mb-4 flex items-center justify-between gap-3">
              <div>
                <h3 class="text-lg font-semibold">{{ group.label }}</h3>
                <p class="text-sm text-muted-foreground">
                  {{ $t('admin.integration.actions.actionsConfigured', group.actions.length) }}
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
                        {{ action.enabled ? $t('globals.terms.active') : $t('globals.terms.inactive', 1) }}
                      </span>
                    </div>
                    <p v-if="action.description" class="text-sm text-muted-foreground">
                      {{ action.description }}
                    </p>
                    <p class="break-all text-sm text-muted-foreground font-mono text-xs">
                      {{ action.url || $t('admin.integration.actions.urlNotConfigured') }}
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
                    <Button variant="outline" size="sm" @click="openEdit(action)">{{ $t('globals.terms.edit') }}</Button>
                    <Button variant="destructive" size="sm" @click="confirmRemove(action)">{{ $t('globals.terms.remove') }}</Button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <template #help>
        <p>{{ $t('admin.integration.actions.help') }}</p>
        <p v-html="$t('admin.integration.actions.help2')"></p>
      </template>
    </AdminPageWithHelp>

    <!-- Create / Edit Dialog -->
    <Dialog v-model:open="dialogOpen">
      <DialogContent class="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{{ isEditing ? $t('globals.messages.edit', { name: $t('globals.terms.action') }) : $t('globals.messages.new', { name: $t('globals.terms.action') }) }}</DialogTitle>
        </DialogHeader>

        <div class="space-y-4 py-2">
          <div class="grid gap-4 lg:grid-cols-2">
            <div>
              <label class="text-sm font-medium">{{ $t('admin.integration.actions.fields.integration') }}</label>
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
              <label class="text-sm font-medium">{{ $t('admin.integration.actions.fields.direction') }}</label>
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
              <label class="text-sm font-medium">{{ $t('admin.integration.actions.fields.actionName') }}</label>
              <Input v-model="draft.name" class="mt-2" :placeholder="$t('admin.integration.actions.fields.actionNamePlaceholder')" />
            </div>
            <div>
              <label class="text-sm font-medium">{{ $t('admin.integration.actions.fields.httpMethod') }}</label>
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
            <label class="text-sm font-medium">{{ $t('admin.integration.actions.fields.description') }}</label>
            <Input
              v-model="draft.description"
              class="mt-2"
              :placeholder="$t('admin.integration.actions.fields.descriptionPlaceholder')"
            />
          </div>

          <div>
            <label class="text-sm font-medium">{{ $t('admin.integration.actions.fields.requestUrl') }}</label>
            <Input
              v-model="draft.url"
              class="mt-2"
              :placeholder="$t('admin.integration.actions.fields.requestUrlPlaceholder')"
            />
          </div>

          <div class="space-y-3">
            <div>
              <label class="text-sm font-medium">{{ $t('admin.integration.actions.fields.triggers') }}</label>
              <p class="mt-1 text-xs text-muted-foreground" v-html="$t('admin.integration.actions.fields.triggersHelp')"></p>
            </div>

            <div class="grid gap-3 md:grid-cols-2">
              <label
                v-for="trigger in integrationTriggers"
                :key="trigger.value"
                class="flex items-center gap-3 rounded-xl border border-border/70 px-3 py-3 cursor-pointer"
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
            <label class="text-sm font-medium">{{ $t('admin.integration.actions.fields.bodyTemplate') }}</label>
            <Textarea
              v-model="draft.bodyTemplate"
              class="mt-2 min-h-32 font-mono text-xs"
              placeholder='{"conversation_uuid":"{{conversation_uuid}}"}'
            />
            <p class="mt-2 text-xs text-muted-foreground">
              {{ $t('admin.integration.actions.fields.bodyTemplateHelp') }}
            </p>
          </div>

          <div class="flex items-center justify-between rounded-xl border border-border/70 px-4 py-3">
            <div>
              <p class="text-sm font-medium">{{ $t('admin.integration.actions.fields.enabled') }}</p>
              <p class="text-xs text-muted-foreground">{{ $t('admin.integration.actions.fields.enabledDescription') }}</p>
            </div>
            <Switch v-model:checked="draft.enabled" />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" @click="dialogOpen = false">{{ $t('globals.messages.cancel') }}</Button>
          <Button @click="saveAction">{{ isEditing ? $t('globals.messages.save') : $t('globals.messages.add', { name: '' }).trim() }}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Delete Confirmation -->
    <AlertDialog v-model:open="deleteDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ $t('admin.integration.actions.removeTitle') }}</AlertDialogTitle>
          <AlertDialogDescription>
            <span v-html="$t('admin.integration.actions.removeDescription', { name: actionToDelete?.name })"></span>
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{{ $t('globals.messages.cancel') }}</AlertDialogCancel>
          <AlertDialogAction @click="removeAction" class="bg-destructive text-destructive-foreground hover:bg-destructive/90">
            {{ $t('globals.terms.remove') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Plus, Zap } from 'lucide-vue-next'
import AdminPageWithHelp from '@/layouts/admin/AdminPageWithHelp.vue'
import AdminPageHeader from '@/components/layout/AdminPageHeader.vue'
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
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle
} from '@/components/ui/alert-dialog'
import {
  useIntegrationActions,
  integrationCatalog,
  integrationDirections,
  integrationTriggers
} from '@/composables/useIntegrationActions'

const { t } = useI18n()
const { actions, groupedActions } = useIntegrationActions()

const dialogOpen = ref(false)
const deleteDialogOpen = ref(false)
const actionToDelete = ref(null)

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

const openCreate = () => {
  draft.value = blankDraft()
  dialogOpen.value = true
}

const openEdit = (action) => {
  draft.value = { ...action, triggers: [...action.triggers] }
  dialogOpen.value = true
}

const confirmRemove = (action) => {
  actionToDelete.value = action
  deleteDialogOpen.value = true
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

  dialogOpen.value = false
  draft.value = blankDraft()
}

function removeAction() {
  if (!actionToDelete.value) return
  actions.value = actions.value.filter((item) => item.id !== actionToDelete.value.id)
  actionToDelete.value = null
}

function labelForTrigger(triggerValue) {
  return integrationTriggers.find((item) => item.value === triggerValue)?.label || triggerValue
}

function labelForDirection(directionValue) {
  return integrationDirections.find((item) => item.value === directionValue)?.label || directionValue
}
</script>
