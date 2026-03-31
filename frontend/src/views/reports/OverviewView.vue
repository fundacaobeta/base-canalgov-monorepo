<template>
  <div class="overflow-y-auto h-full">
    <div
      class="mx-auto w-full max-w-7xl p-6"
      :class="{ 'opacity-50 transition-opacity duration-300': isLoading }"
    >
      <div v-if="isLoading" class="fixed inset-0 flex items-center justify-center z-50 pointer-events-none">
        <Loader2 class="h-8 w-8 animate-spin text-primary" />
      </div>

      <div class="space-y-8">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <h1 class="text-3xl font-bold tracking-tight">{{ t('reports.overview.title') }}</h1>
            <p class="text-sm text-muted-foreground">
              {{ $t('globals.terms.lastUpdated') }}: {{ lastUpdateFormatted }}
            </p>
          </div>
          <div class="flex items-center gap-2">
            <Button variant="outline" size="sm" @click="loadDashboardData" :disabled="isLoading">
              <RefreshCw class="h-4 w-4 mr-2" :class="{ 'animate-spin': isLoading }" />
              {{ t('globals.messages.update', { name: '' }) }}
            </Button>
          </div>
        </div>

        <!-- Row 1: Open Conversations and Agent Status -->
        <div class="grid gap-6 xl:grid-cols-2">
          <Card
            :title="$t('report.openConversations')"
            :counts="cardCounts"
            :labels="conversationCountLabels"
            size="large"
          />
          <Card
            :title="$t('report.agentStatus')"
            :counts="agentStatusCounts"
            :labels="agentStatusLabels"
            size="large"
          />
        </div>

        <!-- Row 2: CSAT and Message Volume -->
        <div class="grid gap-6 xl:grid-cols-2">
          <!-- CSAT Card -->
          <div class="box p-5 flex flex-col">
            <div class="flex justify-between items-center mb-6">
              <p class="card-title">{{ $t('report.csat.cardTitle', { days: csatDays }) }}</p>
              <DateFilter @filter-change="handleCSATFilterChange" :label="''" />
            </div>
            <div class="grid gap-6 sm:grid-cols-3 flex-1 items-center">
              <div class="metric-item">
                <span class="metric-value">{{ formatRating(csatData.average_rating) }}</span>
                <span class="metric-label">{{ $t('report.csat.avgRating') }}</span>
              </div>
              <div class="metric-item">
                <span class="metric-value">{{ formatPercent(csatData.response_rate) }}</span>
                <span class="metric-label">{{ $t('report.csat.responseRate') }}</span>
              </div>
              <div class="metric-item">
                <span class="metric-value">{{
                  formatCompactNumber(csatData.total_responses || 0)
                }}</span>
                <span class="metric-label">{{ $t('report.csat.responses') }}</span>
              </div>
            </div>
          </div>

          <!-- Message Volume Card -->
          <div class="box p-5 flex flex-col">
            <div class="flex justify-between items-center mb-6">
              <p class="card-title">
                {{ $t('report.messages.cardTitle', { days: messageVolumeDays }) }}
              </p>
              <DateFilter @filter-change="handleMessageVolumeFilterChange" :label="''" />
            </div>
            <div class="grid gap-6 sm:grid-cols-2 xl:grid-cols-4 flex-1 items-center">
              <div class="metric-item">
                <span class="metric-value">{{
                  formatCompactNumber(messageVolumeData.total_messages || 0)
                }}</span>
                <span class="metric-label">{{ $t('report.messages.total') }}</span>
              </div>
              <div class="metric-item">
                <span class="metric-value">{{
                  formatCompactNumber(messageVolumeData.incoming_messages || 0)
                }}</span>
                <span class="metric-label">{{ $t('report.messages.incoming') }}</span>
              </div>
              <div class="metric-item">
                <span class="metric-value">{{
                  formatCompactNumber(messageVolumeData.outgoing_messages || 0)
                }}</span>
                <span class="metric-label">{{ $t('report.messages.outgoing') }}</span>
              </div>
              <div class="metric-item">
                <span class="metric-value">{{
                  messageVolumeData.messages_per_conversation || 0
                }}</span>
                <span class="metric-label">{{ $t('report.messages.perConversation') }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Row 3: SLA Card -->
        <div class="w-full box p-6">
          <div class="flex justify-between items-center mb-8">
            <p class="card-title">{{ slaCardTitle }}</p>
            <DateFilter @filter-change="handleSlaFilterChange" :label="''" />
          </div>

          <div class="grid grid-cols-1 gap-8 xl:grid-cols-3 xl:gap-12">
            <!-- First Response -->
            <div class="space-y-6">
              <p class="section-title">{{ $t('report.sla.firstResponse') }}</p>
              <div class="flex justify-center">
                <div class="relative flex flex-col items-center">
                  <span class="text-4xl font-black text-green-600"
                    >{{ slaCounts.first_response_compliance_percent || 0 }}%</span
                  >
                  <span class="metric-label">{{ $t('report.sla.compliance') }}</span>
                </div>
              </div>
              <div class="grid grid-cols-2 gap-4 text-center border-t border-dashed pt-4">
                <div>
                  <span class="text-xl font-bold text-green-600">{{
                    slaCounts.first_response_met_count || 0
                  }}</span>
                  <p class="text-[10px] uppercase text-muted-foreground">{{ $t('report.sla.met') }}</p>
                </div>
                <div>
                  <span class="text-xl font-bold text-red-600">{{
                    slaCounts.first_response_breached_count || 0
                  }}</span>
                  <p class="text-[10px] uppercase text-muted-foreground">{{ $t('report.sla.breached') }}</p>
                </div>
              </div>
              <div class="text-center bg-muted/30 p-2 rounded">
                <span class="text-sm font-bold">{{
                  formattedSlaCounts.avg_first_response_time_sec
                }}</span>
                <p class="text-[10px] text-muted-foreground uppercase">{{ $t('report.sla.avgFirstResp') }}</p>
              </div>
            </div>

            <!-- Next Response -->
            <div class="space-y-6 xl:border-x xl:border-dashed xl:px-8">
              <p class="section-title">{{ $t('report.sla.nextResponse') }}</p>
              <div class="flex justify-center">
                <div class="relative flex flex-col items-center">
                  <span class="text-4xl font-black text-green-600"
                    >{{ slaCounts.next_response_compliance_percent || 0 }}%</span
                  >
                  <span class="metric-label">{{ $t('report.sla.compliance') }}</span>
                </div>
              </div>
              <div class="grid grid-cols-2 gap-4 text-center border-t border-dashed pt-4">
                <div>
                  <span class="text-xl font-bold text-green-600">{{
                    slaCounts.next_response_met_count || 0
                  }}</span>
                  <p class="text-[10px] uppercase text-muted-foreground">{{ $t('report.sla.met') }}</p>
                </div>
                <div>
                  <span class="text-xl font-bold text-red-600">{{
                    slaCounts.next_response_breached_count || 0
                  }}</span>
                  <p class="text-[10px] uppercase text-muted-foreground">{{ $t('report.sla.breached') }}</p>
                </div>
              </div>
              <div class="text-center bg-muted/30 p-2 rounded">
                <span class="text-sm font-bold">{{
                  formattedSlaCounts.avg_next_response_time_sec
                }}</span>
                <p class="text-[10px] text-muted-foreground uppercase">{{ $t('report.sla.avgNextResp') }}</p>
              </div>
            </div>

            <!-- Resolution -->
            <div class="space-y-6">
              <p class="section-title">{{ $t('report.sla.resolution') }}</p>
              <div class="flex justify-center">
                <div class="relative flex flex-col items-center">
                  <span class="text-4xl font-black text-green-600"
                    >{{ slaCounts.resolution_compliance_percent || 0 }}%</span
                  >
                  <span class="metric-label">{{ $t('report.sla.compliance') }}</span>
                </div>
              </div>
              <div class="grid grid-cols-2 gap-4 text-center border-t border-dashed pt-4">
                <div>
                  <span class="text-xl font-bold text-green-600">{{
                    slaCounts.resolution_met_count || 0
                  }}</span>
                  <p class="text-[10px] uppercase text-muted-foreground">{{ $t('report.sla.met') }}</p>
                </div>
                <div>
                  <span class="text-xl font-bold text-red-600">{{
                    slaCounts.resolution_breached_count || 0
                  }}</span>
                  <p class="text-[10px] uppercase text-muted-foreground">{{ $t('report.sla.breached') }}</p>
                </div>
              </div>
              <div class="text-center bg-muted/30 p-2 rounded">
                <span class="text-sm font-bold">{{
                  formattedSlaCounts.avg_resolution_time_sec
                }}</span>
                <p class="text-[10px] text-muted-foreground uppercase">{{ $t('report.sla.avgResolution') }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Row 4: Tag Distribution & History -->
        <div class="grid gap-6 xl:grid-cols-3">
          <div class="box p-5 col-span-1">
            <div class="flex justify-between items-center mb-6">
              <p class="card-title">{{ t('reports.overview.tagDistribution') }}</p>
              <DateFilter @filter-change="handleTagDistributionFilterChange" :label="''" />
            </div>
            <div class="space-y-4">
              <div v-for="tag in (tagDistributionData.top_tags || []).slice(0, 6)" :key="tag.tag_id" class="space-y-1">
                <div class="flex justify-between text-xs font-medium">
                  <span>{{ tag.tag_name }}</span>
                  <span>{{ tag.count }}</span>
                </div>
                <div class="w-full bg-muted rounded-full h-1.5 overflow-hidden">
                  <div class="bg-primary h-full rounded-full" :style="{ width: `${(tag.count / (tagDistributionData.tagged_conversations || 1) * 100) || 0}%` }"></div>
                </div>
              </div>
              <p v-if="!tagDistributionData.top_tags?.length" class="text-center py-10 text-muted-foreground text-sm italic">
                {{ t('globals.messages.noData') }}
              </p>
            </div>
          </div>

          <div class="box p-5 xl:col-span-2">
            <div class="flex justify-between items-center mb-6">
              <p class="card-title">{{ t('reports.overview.history') }}</p>
              <DateFilter @filter-change="handleChartFilterChange" :label="''" />
            </div>
            <div class="h-[320px] w-full xl:h-[360px]">
              <LineChart :data="processedLineData" />
            </div>
          </div>
        </div>

        <!-- Row 5: Custom Reports (Widgets) -->
        <div v-if="customReportsData?.length > 0" class="space-y-6 pt-8 border-t border-dashed">
          <div class="flex items-center justify-between">
            <div>
              <h2 class="text-2xl font-bold tracking-tight">{{ t('globals.terms.customReports') }}</h2>
              <p class="text-xs text-muted-foreground">{{ t('reports.overview.customDescription') }}</p>
            </div>
            <router-link :to="{ name: 'custom-reports' }">
              <Button variant="outline" size="sm" class="h-8">
                <Settings class="h-3.5 w-3.5 mr-2" />
                {{ t('globals.terms.manage') }}
              </Button>
            </router-link>
          </div>
          
          <div class="grid gap-6 xl:grid-cols-3">
            <div v-for="report in customReportsData" :key="report.id" class="box p-5 min-h-[250px] flex flex-col group hover:border-primary/30 transition-colors">
              <div class="flex items-start justify-between mb-6">
                <div>
                  <p class="font-bold text-sm uppercase tracking-wider">{{ report.name }}</p>
                  <p class="text-[10px] text-muted-foreground line-clamp-1" :title="report.description">{{ report.description }}</p>
                </div>
                <Badge variant="secondary" class="uppercase text-[9px] font-black tracking-tighter">{{ t('reports.custom.chartType.' + report.chart_type) }}</Badge>
              </div>

              <!-- Metric type (KPI) -->
              <div v-if="report.chart_type === 'metric'" class="flex-1 flex items-center justify-center">
                <div class="text-center">
                  <div class="text-6xl font-black text-primary tracking-tighter">
                    {{ report.results?.[0]?.value || 0 }}
                  </div>
                  <div class="text-[10px] font-bold uppercase text-muted-foreground mt-2 tracking-widest">
                    {{ report.results?.[0]?.label || t('globals.terms.total') }}
                  </div>
                </div>
              </div>

              <!-- Bar/List -->
              <div v-else class="flex-1 space-y-3 justify-center flex flex-col">
                <div v-for="(res, idx) in (report.results || []).slice(0, 5)" :key="idx" class="space-y-1">
                  <div class="flex justify-between text-[11px] font-bold uppercase">
                    <span class="text-muted-foreground truncate mr-2">{{ res.label }}</span>
                    <span>{{ res.value }}</span>
                  </div>
                  <div class="w-full bg-muted/50 rounded-full h-1.5 overflow-hidden">
                    <div class="bg-primary h-full rounded-full transition-all duration-500" :style="{ width: `${(res.value / (report.total || 1) * 100) || 0}%` }"></div>
                  </div>
                </div>
                <div v-if="!report.results?.length" class="flex-1 flex items-center justify-center italic text-muted-foreground text-xs">
                  {{ t('globals.messages.noData') }}
                </div>
                <div v-if="report.results?.length > 5" class="text-[9px] text-center text-muted-foreground font-bold uppercase mt-2">
                  + {{ report.results.length - 5 }} {{ t('globals.terms.others') }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useEmitter } from '@/composables/useEmitter'
import { EMITTER_EVENTS } from '@/constants/emitterEvents.js'
import { handleHTTPError } from '@/utils/http'
import { formatDuration } from '@/utils/datetime'
import Card from '@/features/reports/OverviewCard.vue'
import LineChart from '@/features/reports/OverviewLineChart.vue'
import { DateFilter } from '@/components/ui/date-filter'
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Loader2, RefreshCw, Settings } from 'lucide-vue-next'
import { useCustomReports } from '@/composables/useCustomReports'
import { translateConversationStatus } from '@/utils/conversationStatus'
import api from '@/api'

const emitter = useEmitter()
const { t } = useI18n()
const { reports: customReports, fetchReports, executeReport } = useCustomReports()

const isLoading = ref(false)
const lastUpdate = ref(new Date())
const cardCounts = ref({})
const chartData = ref({ new_conversations: [], resolved_conversations: [] })
const customReportsData = ref([])
let updateInterval = null

const agentStatusCounts = ref({
  agents_online: 0,
  agents_offline: 0,
  agents_away: 0,
  agents_reassigning: 0
})

const slaCounts = ref({
  first_response_compliance_percent: 0,
  next_response_compliance_percent: 0,
  resolution_compliance_percent: 0,
  avg_first_response_time_sec: 0,
  avg_next_response_time_sec: 0,
  avg_resolution_time_sec: 0,
  first_response_met_count: 0,
  first_response_breached_count: 0,
  next_response_met_count: 0,
  next_response_breached_count: 0,
  resolution_met_count: 0,
  resolution_breached_count: 0
})

const csatData = ref({ average_rating: 0, response_rate: 0, total_responses: 0 })
const messageVolumeData = ref({ total_messages: 0, incoming_messages: 0, outgoing_messages: 0, messages_per_conversation: 0 })
const tagDistributionData = ref({ top_tags: [], tagged_conversations: 0, untagged_conversations: 0, tagged_percentage: 0 })

// Date filters
const slaDays = ref(30)
const chartDays = ref(90)
const csatDays = ref(30)
const messageVolumeDays = ref(30)
const tagDistributionDays = ref(30)

const formatRating = (v) => Number(v || 0).toFixed(1)
const formatPercent = (v) => `${Math.round(v || 0)}%`
const formatCompactNumber = (v) => new Intl.NumberFormat('pt-BR', { notation: 'compact' }).format(v || 0)

const formattedSlaCounts = computed(() => ({
  avg_first_response_time_sec: formatDuration(slaCounts.value.avg_first_response_time_sec, false),
  avg_next_response_time_sec: formatDuration(slaCounts.value.avg_next_response_time_sec, false),
  avg_resolution_time_sec: formatDuration(slaCounts.value.avg_resolution_time_sec, false)
}))

const slaCardTitle = computed(() => t('report.sla.cardTitle', { days: slaDays.value }))
const lastUpdateFormatted = computed(() => lastUpdate.value.toLocaleTimeString())

const conversationCountLabels = computed(() => ({
  open: translateConversationStatus('Open', t),
  snoozed: translateConversationStatus('Snoozed', t),
  resolved: translateConversationStatus('Resolved', t),
  closed: translateConversationStatus('Closed', t)
}))

const agentStatusLabels = computed(() => ({
  agents_online: t('globals.terms.online'),
  agents_offline: t('globals.terms.offline'),
  agents_away: t('globals.terms.away'),
  agents_reassigning: t('globals.messages.reassigning')
}))

const processedLineData = computed(() => {
  const { new_conversations = [], resolved_conversations = [] } = chartData.value
  const dateMap = new Map()
  new_conversations.forEach(i => dateMap.set(i.date, { date: i.date, [t('report.chart.newConversations')]: i.count, [t('report.chart.resolvedConversations')]: 0 }))
  resolved_conversations.forEach(i => {
    if (dateMap.has(i.date)) dateMap.get(i.date)[t('report.chart.resolvedConversations')] = i.count
    else dateMap.set(i.date, { date: i.date, [t('report.chart.newConversations')]: 0, [t('report.chart.resolvedConversations')]: i.count })
  })
  return Array.from(dateMap.values()).sort((a, b) => new Date(a.date) - new Date(b.date))
})

const fetchCardStats = async () => {
  const { data } = await api.getOverviewCounts()
  cardCounts.value = data.data
  agentStatusCounts.value = {
    agents_online: data.data.agents_online || 0,
    agents_offline: data.data.agents_offline || 0,
    agents_away: data.data.agents_away || 0,
    agents_reassigning: data.data.agents_reassigning || 0
  }
}

const fetchSLAStats = async (days = slaDays.value) => {
  const { data } = await api.getOverviewSLA({ days })
  slaCounts.value = { ...slaCounts.value, ...data.data }
}

const fetchChartData = async (days = chartDays.value) => {
  const { data } = await api.getOverviewCharts({ days })
  chartData.value = data.data
}

const fetchCSATStats = async (days = csatDays.value) => {
  const { data } = await api.getOverviewCSAT({ days })
  csatData.value = data.data
}

const fetchMessageVolumeStats = async (days = messageVolumeDays.value) => {
  const { data } = await api.getOverviewMessageVolume({ days })
  messageVolumeData.value = data.data
}

const fetchTagDistributionStats = async (days = tagDistributionDays.value) => {
  const { data } = await api.getOverviewTagDistribution({ days })
  tagDistributionData.value = data.data
}

const handleSlaFilterChange = (d) => { slaDays.value = d; fetchSLAStats(d) }
const handleChartFilterChange = (d) => { chartDays.value = d; fetchChartData(d) }
const handleCSATFilterChange = (d) => { csatDays.value = d; fetchCSATStats(d) }
const handleMessageVolumeFilterChange = (d) => { messageVolumeDays.value = d; fetchMessageVolumeStats(d) }
const handleTagDistributionFilterChange = (d) => { tagDistributionDays.value = d; fetchTagDistributionStats(d) }

const loadDashboardData = async () => {
  isLoading.value = true
  try {
    await Promise.allSettled([
      fetchCardStats(),
      fetchSLAStats(),
      fetchChartData(),
      fetchCSATStats(),
      fetchMessageVolumeStats(),
      fetchTagDistributionStats(),
      fetchReports()
    ])

    if (customReports.value?.length > 0) {
      const results = await Promise.all(
        customReports.value.map(async (r) => {
          const data = await executeReport(r.id)
          const total = data.reduce((acc, curr) => acc + curr.value, 0)
          return { ...r, results: data, total }
        })
      )
      customReportsData.value = results
    }
  } catch (err) {
    emitter.emit(EMITTER_EVENTS.SHOW_TOAST, { variant: 'destructive', description: handleHTTPError(err).message })
  } finally {
    isLoading.value = false
    lastUpdate.value = new Date()
  }
}

onMounted(() => {
  loadDashboardData()
  updateInterval = setInterval(loadDashboardData, 60000)
})

onUnmounted(() => {
  if (updateInterval) clearInterval(updateInterval)
})
</script>

<style scoped>
.metric-value { @apply text-3xl font-bold tracking-tight; }
.metric-label { @apply text-[10px] text-muted-foreground uppercase font-bold tracking-widest; }
.card-title { @apply text-lg font-bold tracking-tight; }
.metric-item { @apply flex flex-col items-center gap-1 text-center; }
.section-title { @apply text-[11px] font-black text-center text-muted-foreground uppercase tracking-widest; }
.box { @apply rounded-xl border bg-card text-card-foreground shadow-sm; }
</style>
