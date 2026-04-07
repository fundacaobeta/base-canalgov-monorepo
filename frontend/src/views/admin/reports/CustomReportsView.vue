<template>
  <div class="space-y-6 h-full overflow-y-auto p-1">
    <div class="rounded-2xl border bg-gradient-to-br from-slate-50 via-white to-emerald-50 p-6 shadow-sm">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">{{ t('reports.custom.title') }}</h1>
        <p class="text-muted-foreground italic text-sm">
          Crie dashboards dinâmicos para monitorar métricas específicas da prefeitura.
        </p>
      </div>
      <Button @click="openCreateDialog">
        <Plus class="mr-2 h-4 w-4" />
        {{ t('reports.custom.new') }}
      </Button>
    </div>

    <div v-if="loading && !reports?.length" class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      <Card v-for="i in 3" :key="i" class="flex flex-col">
        <CardHeader><Skeleton class="h-6 w-3/4" /><Skeleton class="h-4 w-1/2" /></CardHeader>
        <CardContent><Skeleton class="h-20 w-full" /></CardContent>
      </Card>
    </div>

    <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      <Card v-for="report in reports" :key="report.id" class="flex flex-col group hover:border-primary/50 transition-colors">
        <CardHeader>
          <div class="flex items-center justify-between">
            <CardTitle class="text-lg">{{ report.name }}</CardTitle>
            <BarChart3 class="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
          </div>
          <CardDescription class="line-clamp-2 min-h-[40px]">{{ report.description || t('globals.messages.noData') }}</CardDescription>
        </CardHeader>
        <CardContent class="flex-1 space-y-4">
          <div class="flex items-center gap-2">
            <Badge variant="outline" class="capitalize">{{ getCustomChartTypeLabel(report.chart_type) }}</Badge>
            <Badge variant="secondary" class="capitalize">{{ getMetricTypeLabel(report.metric_type) }}</Badge>
          </div>
          <div class="text-[10px] text-muted-foreground uppercase font-bold tracking-tighter">
            {{ report.filters?.length || 0 }} {{ t('globals.terms.filter', report.filters?.length || 0).toLowerCase() }} ativos
          </div>
        </CardContent>
        <CardFooter class="flex justify-end gap-2 border-t bg-muted/5 pt-4">
          <Button variant="ghost" size="sm" @click="openEditDialog(report)">
            <Edit class="h-4 w-4" />
          </Button>
          <Button variant="ghost" size="sm" class="text-destructive hover:bg-destructive/10" @click="handleDelete(report.id)">
            <Trash2 class="h-4 w-4" />
          </Button>
        </CardFooter>
      </Card>
    </div>

    <div v-if="!loading && !reports?.length" class="flex flex-col items-center justify-center py-20 border-2 border-dashed rounded-xl text-center">
      <BarChart3 class="h-12 w-12 text-muted-foreground/20 mb-4" />
      <p class="text-muted-foreground font-medium text-sm">{{ t('reports.custom.empty') }}</p>
      <Button variant="link" class="mt-2" @click="openCreateDialog">{{ t('reports.custom.createFirst') }}</Button>
    </div>

    <!-- Report Form Dialog -->
    <Dialog v-model:open="isDialogOpen">
      <DialogContent class="sm:max-w-[750px] max-h-[90vh] flex flex-col p-0">
        <DialogHeader class="p-6 pb-0">
          <DialogTitle>{{ editMode ? t('globals.messages.edit') : t('reports.custom.new') }}</DialogTitle>
          <DialogDescription>
            Configure as métricas e filtros para gerar seu relatório personalizado.
          </DialogDescription>
        </DialogHeader>
        
        <div class="flex-1 overflow-y-auto p-6 space-y-6">
          <div class="grid gap-4 md:grid-cols-2">
            <div class="grid gap-2">
              <Label>{{ t('globals.terms.name') }}</Label>
              <Input v-model="form.name" :placeholder="t('globals.terms.name')" />
            </div>
            <div class="grid gap-2">
              <Label>{{ t('globals.terms.chartType') }}</Label>
              <Select v-model="form.chart_type">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="bar">{{ getCustomChartTypeLabel('bar') }}</SelectItem>
                  <SelectItem value="pie">{{ getCustomChartTypeLabel('pie') }}</SelectItem>
                  <SelectItem value="line">{{ getCustomChartTypeLabel('line') }}</SelectItem>
                  <SelectItem value="metric">{{ getCustomChartTypeLabel('metric') }}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div class="grid gap-2">
            <Label>{{ t('globals.terms.description') }}</Label>
            <Textarea v-model="form.description" class="resize-none h-20" />
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

          <div class="space-y-4 pt-4 border-t">
            <Label class="text-sm font-bold uppercase tracking-widest text-primary flex items-center gap-2">
              <Filter class="h-4 w-4" />
              {{ t('globals.terms.filterRules') }}
            </Label>
            
            <FilterBuilder 
              :fields="filterFields" 
              v-model="form.filters"
              class="border rounded-lg p-4 bg-muted/10"
            />
          </div>
        </div>

        <DialogFooter class="p-6 pt-0 border-t mt-0">
          <Button variant="ghost" @click="isDialogOpen = false" :disabled="saving">{{ t('globals.terms.cancel') }}</Button>
          <Button @click="handleSave" :disabled="saving || !form.name">
            <Loader2 v-if="saving" class="mr-2 h-4 w-4 animate-spin" />
            {{ t('globals.terms.save') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { Plus, BarChart3, Edit, Trash2, Filter, Loader2 } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Textarea } from '@/components/ui/textarea'
import { Skeleton } from '@/components/ui/skeleton'
import { 
  Select, SelectContent, SelectItem, 
  SelectTrigger, SelectValue 
} from '@/components/ui/select'
import { 
  Dialog, DialogContent, DialogDescription, 
  DialogFooter, DialogHeader, DialogTitle 
} from '@/components/ui/dialog'
import { useCustomReports } from '@/composables/useCustomReports'
import { useConversationFilters } from '@/composables/useConversationFilters'
import FilterBuilder from '@/components/filter/FilterBuilder.vue'

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

const getCustomChartTypeLabel = (chartType) => {
  const labels = {
    bar: 'Grafico de Barras',
    pie: 'Grafico de Pizza',
    line: 'Grafico de Linha',
    metric: 'Metrica Unica (KPI)'
  }

  return labels[chartType] || chartType
}

const getMetricTypeLabel = (metricType) => {
  const normalizedMetricType = metricType === 'conversations_count'
    ? 'conversations_by_status'
    : metricType

  return metricOptions.find(option => option.value === normalizedMetricType)?.label || normalizedMetricType
}

const normalizeMetricType = (metricType) => (
  metricType === 'conversations_count' ? 'conversations_by_status' : metricType
)

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

const selectedMetricOption = computed(() =>
  metricOptions.find(option => option.value === form.metric_type) || metricOptions[0]
)

const form = reactive({
  name: '',
  description: '',
  chart_type: 'bar',
  metric_type: 'conversations_by_status',
  filters: []
})

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
  form.metric_type = normalizeMetricType(report.metric_type) || 'conversations_by_status'
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
