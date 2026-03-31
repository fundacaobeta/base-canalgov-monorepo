<template>
  <div class="h-full min-w-0 overflow-hidden">
    <LineChart
      class="h-full min-w-0"
      :data="data"
      index="date"
      :categories="[t('report.chart.newConversations'), t('report.chart.resolvedConversations')]"
      :x-formatter="xFormatter"
      :y-formatter="yFormatter"
      :margin="{ left: 24, right: 24, top: 12, bottom: 52 }"
    />
  </div>
</template>

<script setup>
import { LineChart } from '@/components/ui/chart-line'
import { useI18n } from 'vue-i18n'
const props = defineProps({
  data: {
    type: Array,
    default: () => []
  }
})
const { t } = useI18n()

const formatShortDate = (value) => {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return String(value ?? '')
  return date.toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit'
  })
}

const xFormatter = (tick) => {
  const total = props.data.length
  const step = total > 12 ? Math.ceil(total / 6) : 2

  if (typeof tick === 'number') {
    const isEdge = tick === 0 || tick === total - 1
    if (!isEdge && tick % step !== 0) return ''
    return formatShortDate(props.data[tick]?.date)
  }

  return formatShortDate(tick)
}

const yFormatter = (tick) => {
  return Number.isInteger(tick) ? tick : ''
}
</script>
