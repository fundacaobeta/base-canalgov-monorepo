<template>
  <div class="h-full space-y-6 overflow-y-auto p-1">
    <div class="rounded-2xl border bg-gradient-to-br from-slate-50 via-white to-emerald-50 p-6 shadow-sm">
      <div class="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
        <div class="space-y-3">
          <div class="inline-flex items-center gap-2 rounded-full border bg-white/80 px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] text-slate-600">
            <BarChart3 class="h-3.5 w-3.5" />
            Relatorios personalizados
          </div>
          <div>
            <h1 class="text-3xl font-bold tracking-tight text-slate-900">{{ t('reports.custom.title') }}</h1>
            <p class="mt-2 max-w-3xl text-sm leading-6 text-slate-600">
              Monte relatorios com os dados que ja existem no sistema, como status, prioridade,
              caixas, equipes e agentes. A ideia aqui e permitir que qualquer pessoa organize a
              leitura dos atendimentos sem depender de ajuste tecnico.
            </p>
          </div>
        </div>

        <div class="flex flex-col gap-3 sm:flex-row">
          <Button variant="outline" class="bg-white">
            {{ reports?.length || 0 }} relatorios criados
          </Button>
          <Button @click="openCreateDialog">
            <Plus class="mr-2 h-4 w-4" />
            {{ t('reports.custom.new') }}
          </Button>
        </div>
      </div>
    </div>

    <div class="grid gap-4 xl:grid-cols-[minmax(0,1.45fr)_minmax(320px,1fr)]">
      <Card class="border-dashed">
        <CardHeader>
          <CardTitle class="text-base">Como usar esta pagina</CardTitle>
          <CardDescription>
            Escolha o tema do relatorio, defina como quer enxergar os dados e depois aplique filtros
            com base nas informacoes do sistema.
          </CardDescription>
        </CardHeader>
        <CardContent class="grid gap-3 md:grid-cols-3">
          <div class="rounded-xl border bg-muted/20 p-4">
            <p class="text-xs font-semibold uppercase tracking-[0.18em] text-muted-foreground">1. Escolha o foco</p>
            <p class="mt-2 text-sm font-medium">Status, prioridade, caixa, equipe ou agente</p>
          </div>
          <div class="rounded-xl border bg-muted/20 p-4">
            <p class="text-xs font-semibold uppercase tracking-[0.18em] text-muted-foreground">2. Defina o formato</p>
            <p class="mt-2 text-sm font-medium">Linha, barras, pizza ou indicador rapido</p>
          </div>
          <div class="rounded-xl border bg-muted/20 p-4">
            <p class="text-xs font-semibold uppercase tracking-[0.18em] text-muted-foreground">3. Aplique filtros</p>
            <p class="mt-2 text-sm font-medium">Refine por caixa, agente, prioridade, status e mais</p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-base">Dados que podem virar relatorio</CardTitle>
          <CardDescription>
            Esses recortes usam informacoes que ja existem nos atendimentos do sistema.
          </CardDescription>
        </CardHeader>
        <CardContent class="grid gap-2 sm:grid-cols-2">
          <div
            v-for="option in metricOptions"
            :key="option.value"
            class="rounded-xl border px-3 py-3"
          >
            <p class="text-sm font-semibold text-slate-900">{{ option.label }}</p>
            <p class="mt-1 text-xs leading-5 text-muted-foreground">{{ option.description }}</p>
          </div>
        </CardContent>
      </Card>
    </div>

    <div v-if="loading && !reports?.length" class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <Card v-for="i in 3" :key="i" class="flex flex-col">
        <CardHeader><Skeleton class="h-6 w-3/4" /><Skeleton class="h-4 w-1/2" /></CardHeader>
        <CardContent><Skeleton class="h-20 w-full" /></CardContent>
      </Card>
    </div>

    <div v-else class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <Card
        v-for="report in reports"
        :key="report.id"
        class="group flex flex-col transition-colors hover:border-primary/50"
      >
        <CardHeader>
          <div class="flex items-center justify-between gap-3">
            <CardTitle class="text-lg">{{ report.name }}</CardTitle>
            <BarChart3 class="h-4 w-4 text-muted-foreground transition-colors group-hover:text-primary" />
          </div>
          <CardDescription class="min-h-[40px] line-clamp-2">
            {{ report.description || 'Relatorio personalizado baseado nos dados do sistema.' }}
          </CardDescription>
        </CardHeader>
        <CardContent class="flex-1 space-y-4">
          <div class="flex flex-wrap items-center gap-2">
            <Badge variant="outline" class="capitalize">{{ t('reports.custom.chartType.' + report.chart_type) }}</Badge>
            <Badge variant="secondary">{{ getMetricLabel(report.metric_type) }}</Badge>
          </div>
          <p class="text-xs leading-5 text-muted-foreground">
            {{ getMetricDescription(report.metric_type) }}
          </p>
          <div class="text-[10px] font-bold uppercase tracking-tighter text-muted-foreground">
            {{ report.filters?.length || 0 }} {{ t('globals.terms.filter', report.filters?.length || 0).toLowerCase() }} ativos
          </div>
        </CardContent>
        <CardFooter class="flex justify-end gap-2 border-t bg-muted/5 pt-4">
          <Button variant="ghost" size="sm" @click="openEditDialog(report)">
            <Edit class="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            class="text-destructive hover:bg-destructive/10"
            @click="handleDelete(report.id)"
          >
            <Trash2 class="h-4 w-4" />
          </Button>
        </CardFooter>
      </Card>
    </div>

    <div
      v-if="!loading && !reports?.length"
      class="flex flex-col items-center justify-center rounded-xl border-2 border-dashed py-20 text-center"
    >
      <BarChart3 class="mb-4 h-12 w-12 text-muted-foreground/20" />
      <p class="text-sm font-medium text-muted-foreground">{{ t('reports.custom.empty') }}</p>
      <Button variant="link" class="mt-2" @click="openCreateDialog">{{ t('reports.custom.createFirst') }}</Button>
    </div>

    <Dialog v-model:open="isDialogOpen">
      <DialogContent class="flex max-h-[90vh] flex-col p-0 sm:max-w-[780px]">
        <DialogHeader class="p-6 pb-0">
          <DialogTitle>{{ editMode ? t('globals.messages.edit') : t('reports.custom.new') }}</DialogTitle>
          <DialogDescription>
            Crie um relatorio com base nos dados do sistema e organize a visualizacao do jeito que fizer mais sentido para a equipe.
          </DialogDescription>
        </DialogHeader>

        <div class="flex-1 space-y-6 overflow-y-auto p-6">
          <div class="grid gap-4 md:grid-cols-2">
            <div class="grid gap-2">
              <Label>{{ t('globals.terms.name') }}</Label>
              <Input v-model="form.name" placeholder="Ex.: Atendimentos por equipe" />
            </div>
            <div class="grid gap-2">
              <Label>{{ t('globals.terms.chartType') }}</Label>
              <Select v-model="form.chart_type">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="bar">{{ t('reports.custom.chartType.bar') }}</SelectItem>
                  <SelectItem value="pie">{{ t('reports.custom.chartType.pie') }}</SelectItem>
                  <SelectItem value="line">{{ t('reports.custom.chartType.line') }}</SelectItem>
                  <SelectItem value="metric">{{ t('reports.custom.chartType.metric') }}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div class="grid gap-2">
            <Label>{{ t('globals.terms.description') }}</Label>
            <Textarea
              v-model="form.description"
              class="h-20 resize-none"
              placeholder="Descreva rapidamente o objetivo do relatorio."
            />
          </div>

          <div class="grid gap-2">
            <Label>O que voce quer analisar?</Label>
            <Select v-model="form.metric_type">
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="option in metricOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </SelectItem>
              </SelectContent>
            </Select>
            <p class="text-xs text-muted-foreground">{{ selectedMetricOption?.description }}</p>
          </div>

          <div class="space-y-4 border-t pt-4">
            <Label class="flex items-center gap-2 text-sm font-bold uppercase tracking-widest text-primary">
              <Filter class="h-4 w-4" />
              Regras de filtro
            </Label>
            <p class="text-sm text-muted-foreground">
              Use os filtros abaixo para limitar o relatorio aos dados que realmente importam.
            </p>

            <FilterBuilder
              v-model="form.filters"
              :fields="filterFields"
              class="rounded-lg border bg-muted/10 p-4"
            />
          </div>
        </div>

        <DialogFooter class="mt-0 border-t p-6 pt-0">
          <Button variant="ghost" :disabled="saving" @click="isDialogOpen = false">{{ t('globals.terms.cancel') }}</Button>
          <Button :disabled="saving || !form.name" @click="handleSave">
            <Loader2 v-if="saving" class="mr-2 h-4 w-4 animate-spin" />
            {{ t('globals.terms.save') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { BarChart3, Edit, Filter, Loader2, Plus, Trash2 } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import FilterBuilder from '@/components/filter/FilterBuilder.vue'
import { useConversationFilters } from '@/composables/useConversationFilters'
import { useCustomReports } from '@/composables/useCustomReports'

const { t } = useI18n()
const { reports, loading, fetchReports, createReport, updateReport, deleteReport } = useCustomReports()
const { conversationsListFilters } = useConversationFilters()

const isDialogOpen = ref(false)
const editMode = ref(false)
const selectedId = ref(null)
const saving = ref(false)

const metricOptions = [
  {
    value: 'conversations_by_status',
    label: 'Atendimentos por status',
    description: 'Mostra quantos atendimentos estao abertos, adiados, resolvidos ou fechados.'
  },
  {
    value: 'conversations_by_priority',
    label: 'Atendimentos por prioridade',
    description: 'Ajuda a enxergar a distribuicao entre baixa, media, alta e urgente.'
  },
  {
    value: 'conversations_by_inbox',
    label: 'Atendimentos por caixa',
    description: 'Compara o volume entre caixas de entrada e canais do sistema.'
  },
  {
    value: 'conversations_by_team',
    label: 'Atendimentos por equipe',
    description: 'Mostra como os atendimentos estao distribuidos entre as equipes.'
  },
  {
    value: 'conversations_by_agent',
    label: 'Atendimentos por agente',
    description: 'Permite acompanhar a distribuicao entre os agentes responsaveis.'
  }
]

const filterFields = computed(() =>
  Object.entries(conversationsListFilters.value).map(([field, value]) => ({
    model: 'conversations',
    label: value.label,
    field,
    type: value.type,
    operators: value.operators,
    options: value.options ?? []
  }))
)

const form = reactive({
  name: '',
  description: '',
  chart_type: 'bar',
  metric_type: 'conversations_by_status',
  filters: []
})

const selectedMetricOption = computed(() =>
  metricOptions.find(option => option.value === form.metric_type) || metricOptions[0]
)

const getMetricLabel = (metricType) =>
  metricOptions.find(option => option.value === metricType)?.label || metricType

const getMetricDescription = (metricType) =>
  metricOptions.find(option => option.value === metricType)?.description || 'Relatorio personalizado baseado nos dados do sistema.'

const openCreateDialog = () => {
  editMode.value = false
  selectedId.value = null
  form.name = ''
  form.description = ''
  form.chart_type = 'bar'
  form.metric_type = 'conversations_by_status'
  form.filters = []
  isDialogOpen.value = true
}

const openEditDialog = (report) => {
  editMode.value = true
  selectedId.value = report.id
  form.name = report.name
  form.description = report.description
  form.chart_type = report.chart_type
  form.metric_type = report.metric_type || 'conversations_by_status'
  form.filters = Array.isArray(report.filters) ? report.filters : []
  isDialogOpen.value = true
}

const handleSave = async () => {
  try {
    saving.value = true
    const data = { ...form }
    if (editMode.value) {
      await updateReport(selectedId.value, data)
    } else {
      await createReport(data)
    }
    isDialogOpen.value = false
    await fetchReports()
  } catch (err) {
    console.error(err)
  } finally {
    saving.value = false
  }
}

const handleDelete = async (id) => {
  if (!confirm(t('globals.messages.confirmDelete'))) return
  await deleteReport(id)
  await fetchReports()
}

onMounted(fetchReports)
</script>
